package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/XeroAPI/golang-oauth2-example/config"
	"github.com/XeroAPI/golang-oauth2-example/xero"
	"golang.org/x/oauth2"
)

// Handles HTTP requests from the browser.

var appConfig *config.Config
var oauth2Config *oauth2.Config
var authorisationCode string
var appContext context.Context
var apiClient *http.Client

func init() {
	appContext = context.Background()
}

// New - Returns an instance of the HTTP server.
func New(c *config.Config, oc *oauth2.Config, callbackURI string) *http.Server {
	appConfig = c
	oauth2Config = oc
	http.HandleFunc("/", handleIndexPage)
	http.HandleFunc("/login", redirectToAuthorisationEndpoint)
	http.HandleFunc(callbackURI, handleOAuthCallback)
	return &http.Server{Addr: fmt.Sprintf(":%d", c.AppPort)}
}

func redirectToAuthorisationEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Location", oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func handleIndexPage(w http.ResponseWriter, req *http.Request) {
	if authorisationCode == "" {
		w.Header().Add("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	// Todo - Something more efficient than re-initialising the client each time this endpoint is hit
	client := oauth2Config.Client(appContext, getAuthorisationToken())
	orgsResponse, err := client.Get("https://api.xero.com/connections")

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
	var organisations []xero.Organisation
	err = json.Unmarshal(respBody, &organisations)
	if err != nil {
		errMsg := "There was an error attempting to unmarshal the data from the connections endpoint."
		w.Write([]byte(errMsg + " Please check the log for details."))
		log.Println(errMsg)
		log.Println(err)
		return
	}
	w.Write([]byte("<h1>Connected organisations</h1>"))
	w.Write([]byte("<ul>"))
	for _, org := range organisations {
		w.Write([]byte("<li>" + org.TenantName + "</li>"))
	}
	w.Write([]byte("</ul>"))
}

func handleOAuthCallback(w http.ResponseWriter, req *http.Request) {
	authorisationCode = req.URL.Query().Get("code")
	if config.DebugMode {
		log.Println("Received authorisation code:", authorisationCode)
	}
	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func getAuthorisationToken() *oauth2.Token {
	encodedAuthHeaderValue := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("Basic %s:%s", appConfig.ClientID, appConfig.ClientSecret),
	))
	tok, err := oauth2Config.Exchange(
		appContext,
		authorisationCode,
		oauth2.SetAuthURLParam("authorization", encodedAuthHeaderValue),
	)
	if err != nil {
		log.Println("An error occurred while trying to exchange the authorisation code with the Xero API.")
		log.Fatal(err)
	}
	return tok
}
