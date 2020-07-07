package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/XeroAPI/golang-oauth2-example/config"
	"github.com/XeroAPI/golang-oauth2-example/server/ui"
	"github.com/XeroAPI/golang-oauth2-example/xero"
)

func (s *Server) handleInvoicePage(w http.ResponseWriter, req *http.Request) {
	if s.preFlightCheck() {
		s.redirectToLogin(w)
		return
	}
	ui.WriteGlobalStylesTag(w)
	tenantID := req.URL.Query().Get("tenantId")
	if tenantID == "" {
		w.Write([]byte("<p>Missing query string parameter 'tenantId'.</p>"))
		w.Write([]byte(returnToHomepageLink))
		return
	}
	if config.DebugMode {
		log.Println("Attempting to retrieve invoices for organisation ID:", tenantID)
	}
	// Note the filter here - We're retrieving invoices for accounts receivable.
	// If you want to retrieve invoices for accounts payable, change "ACCREC" to "ACCPAY".
	invoicesRequest, err := http.NewRequest("GET", "https://api.xero.com/api.xro/2.0/invoices?where=Type=\"ACCREC\"", nil)
	if err != nil {
		errMsg := "An error occurred while trying to create a request to send to the Xero API."
		w.Write([]byte("<p>" + errMsg + "</p>"))
		w.Write([]byte("<p>Please check the log for more details.</p>"))
		w.Write([]byte(returnToHomepageLink))
		log.Println(errMsg)
		log.Println(err)
		return
	}
	// We use httpClient.Do() instead of s.httpClient.Get(), because we need to inject the xero-tenant-id header value for
	// these calls that retrieve information specific to one organisation.
	invoicesRequest.Header.Add("xero-tenant-id", tenantID)
	invoicesRequest.Header.Add(s.getAuthorisationHeader())
	invoicesResponse, err := s.httpClient.Do(invoicesRequest)

	if err != nil && invoicesResponse.StatusCode != 200 {
		errMsg := "An error occurred while trying to retrieve the invoices connected to this organisation."
		w.Write([]byte(errMsg + " Please check the log for details"))
		log.Println(errMsg)
		log.Println(err)
	}
	respBody, err := ioutil.ReadAll(invoicesResponse.Body)
	if err != nil {
		errMsg := "An error occurred while trying to read the response from the Xero API"
		w.Write([]byte(errMsg + " Please check the log for details"))
		log.Println(errMsg)
		log.Fatalln(err)
		return
	}
	var invoices xero.InvoiceResponse
	err = json.Unmarshal(respBody, &invoices)
	if err != nil {
		errMsg := "There was an error attempting to unmarshal the data from the organisation endpoint."
		w.Write([]byte(errMsg + " Please check the log for details."))
		log.Println(errMsg)
		log.Println(err)
		if config.DebugMode {
			log.Println("Response Status:", invoices.Status)
			log.Println("Response Body:")
			fmt.Println(string(respBody))
		}
		return
	}

	w.Write([]byte("<h1>Recent Invoices (Accounts Receivable)</h1>"))
	for _, invoice := range invoices.Invoices {
		var invoiceStatusString string
		if invoice.Status == "PAID" {
			invoiceStatusString = " (PAID)"
		}
		line := fmt.Sprintf(
			"<p>%s - %0.2f%s - %s</p>",
			invoice.InvoiceNumber,
			invoice.Total,
			invoiceStatusString,
			invoice.Contact.Name,
		)
		w.Write([]byte(line))
	}
	w.Write([]byte(returnToHomepageLink))
}
