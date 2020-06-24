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
	if s.preFlightCheck(w) {
		return
	}
	orgID := req.URL.Query().Get("id")
	if orgID == "" {
		w.Write([]byte("<p>Missing query string parameter 'id'.</p>"))
		w.Write([]byte(returnToHomepageLink))
		return
	}
	orgsRequest, err := http.NewRequest("GET", "https://api.xero.com/api.xero/2.0/organisation", nil)
	if err != nil {
		errMsg := "An error occurred while trying to create a request to send to the Xero API."
		w.Write([]byte("<p>" + errMsg + "</p>"))
		w.Write([]byte("<p>Please check the log for more details.</p>"))
		w.Write([]byte(returnToHomepageLink))
		log.Println(errMsg)
		log.Println(err)
		return
	}
	orgsRequest.Header.Add("xero-tenant-id", orgID)
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
		log.Println("Response Status:", orgsResponse.Status)
		w.Write([]byte(errMsg + " Please check the log for details."))
		log.Println(errMsg)
		log.Println(err)
		fmt.Println(respBody)
		return
	}

	org, err := orgs.GetOrgByID(orgID)
	if err != nil {
		w.Write([]byte("<p>Sorry, we couldn't find the specified organisation!</p>"))
		w.Write([]byte(returnToHomepageLink))
		log.Println("Unable to find organisation with ID:", orgID)
		if config.DebugMode {
			log.Println("Available orgs:")
			for _, org := range orgs.Organisations {
				log.Println("-", org.OrganisationID, "-", org.LegalName)
			}
			fmt.Println(string(respBody))
		}
		return
	}
	w.Write([]byte("<h1>" + org.LegalName + "</h1>"))
	primaryAddress := org.Addresses[0]
	addressString := fmt.Sprintf("%s, %s %s", primaryAddress.AddressLine1, primaryAddress.City, primaryAddress.PostalCode)
	w.Write([]byte("<p>" + addressString + "</p>"))
}
