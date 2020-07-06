package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/XeroAPI/golang-oauth2-example/config"
	"github.com/XeroAPI/golang-oauth2-example/xero"
)

func (s *Server) handleOrganisationPage(w http.ResponseWriter, req *http.Request) {
	if s.preFlightCheck() {
		s.redirectToLogin(w)
		return
	}
	tenantID := req.URL.Query().Get("tenantId")
	if tenantID == "" {
		w.Write([]byte("<p>Missing query string parameter 'tenantId'.</p>"))
		w.Write([]byte(returnToHomepageLink))
		return
	}
	if config.DebugMode {
		log.Println("Attempting to retrieve information for organisation ID:", tenantID)
	}
	orgsRequest, err := http.NewRequest("GET", "https://api.xero.com/api.xro/2.0/organisation", nil)
	if err != nil {
		errMsg := "An error occurred while trying to create a request to send to the Xero API."
		w.Write([]byte("<p>" + errMsg + "</p>"))
		w.Write([]byte("<p>Please check the log for more details.</p>"))
		w.Write([]byte(returnToHomepageLink))
		log.Println(errMsg)
		log.Println(err)
		return
	}
	orgsRequest.Header.Add("xero-tenant-id", tenantID)
	orgsRequest.Header.Add(s.getAuthorisationHeader())
	orgsResponse, err := s.httpClient.Do(orgsRequest)

	if err != nil && orgsResponse.StatusCode != 200 {
		errMsg := "An error occurred while trying to retrieve the organisations connected to this account."
		w.Write([]byte(errMsg + " Please check the log for details"))
		log.Println(errMsg)
		log.Println(err)
	}
	respBody, err := ioutil.ReadAll(orgsResponse.Body)
	if err != nil {
		errMsg := "An error occurred while trying to read the response from the Xero API"
		w.Write([]byte(errMsg + " Please check the log for details"))
		log.Println(errMsg)
		log.Fatalln(err)
		return
	}
	var orgs xero.OrganisationResponse
	err = json.Unmarshal(respBody, &orgs)
	if err != nil {
		errMsg := "There was an error attempting to unmarshal the data from the organisation endpoint."
		w.Write([]byte(errMsg + " Please check the log for details."))
		log.Println(errMsg)
		log.Println(err)
		if config.DebugMode {
			log.Println("Response Status:", orgsResponse.Status)
			log.Println("Response Body:")
			fmt.Println(string(respBody))
		}
		return
	}

	org, err := orgs.GetOrgByID(tenantID)
	if err != nil {
		w.Write([]byte("<p>Sorry, we couldn't find the specified organisation!</p>"))
		w.Write([]byte(returnToHomepageLink))
		log.Println("Unable to find organisation with ID:", tenantID)
		if config.DebugMode {
			fmt.Println(string(respBody))
			log.Println("Available orgs:")
			for _, org := range orgs.Organisations {
				log.Println("-", org.OrganisationID, "-", org.LegalName)
			}
		}
		return
	}
	w.Write([]byte("<h1>" + org.LegalName + "</h1>"))
	if org.IsDemoCompany {
		w.Write([]byte("<p>This is a demo company.</p>"))
	}
	primaryAddress := org.Addresses[0]
	addressString := fmt.Sprintf("%s, %s %s", primaryAddress.AddressLine1, primaryAddress.City, primaryAddress.PostalCode)
	w.Write([]byte("<p>" + addressString + "</p>"))
	w.Write([]byte("<p><a href=\"/invoices?tenantId=" + tenantID + "\">View Invoices</a></p>"))
	w.Write([]byte(returnToHomepageLink))
}
