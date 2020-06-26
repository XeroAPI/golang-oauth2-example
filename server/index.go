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

func (s *Server) handleIndexPage(w http.ResponseWriter, req *http.Request) {
	if s.preFlightCheck(w) {
		return
	}

	connectionsResponse, err := s.httpClient.Get("https://api.xero.com/connections")

	if err != nil && connectionsResponse.StatusCode != 200 {
		errMsg := "An error occurred while trying to retrieve the organisations connected to this account."
		w.Write([]byte(errMsg + " Please check the log for details"))
		log.Println(errMsg)
		log.Println(err)
	}
	respBody, err := ioutil.ReadAll(connectionsResponse.Body)
	if err != nil {
		errMsg := "An error occurred while trying to read the response from the Xero API"
		w.Write([]byte(errMsg + " Please check the log for details"))
		log.Println(errMsg)
		log.Fatalln(err)
		if config.DebugMode {
			fmt.Println(string(respBody))
		}
		return
	}
	var connections []xero.Connection
	err = json.Unmarshal(respBody, &connections)
	if err != nil {
		errMsg := "There was an error attempting to unmarshal the data from the connections endpoint."
		w.Write([]byte(errMsg + " Please check the log for details."))
		log.Println(errMsg)
		log.Println(err)
		if config.DebugMode {
			fmt.Println(string(respBody))
		}
		return
	}
	w.Write([]byte("<h1>Connected Entities™</h1>"))
	w.Write([]byte("<ul>"))
	for _, org := range connections {
		w.Write([]byte("<li><a href=\"/organisation?tenantId=" + org.TenantID + "\">" + org.TenantName + "</a></li>"))
	}
	w.Write([]byte("</ul>"))

	w.Write([]byte("<h2>Other Fun Things</h2>"))
	w.Write([]byte("<ul>"))
	w.Write([]byte("<li><a href=\"/refresh\">Force Refresh Access Token</a></li>"))
	w.Write([]byte("<li><a href=\"/segfault\">Cause a Segfault™</a></li>"))
	w.Write([]byte("</ul>"))
}
