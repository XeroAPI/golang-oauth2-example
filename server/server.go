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

// Server - Represents the running app, and all the bits and pieces we need to make API calls.
type Server struct {
	config                 *config.Config
	context                context.Context
	httpServer             *http.Server
	httpClient             *http.Client
	oAuthAuthorisationCode string
	oAuthToken             *oauth2.Token
}

// Some package level variables
var callbackURI = "/auth/xero"
var oAuthScopes = []string{
	"openid",
	"profile",
	"email",
	"accounting.transactions",
	"accounting.settings",
	"offline_access",
}

// New - Returns an instance of the HTTP server.
func New(c *config.Config) *Server {
	c.OAuth2Config.RedirectURL = fmt.Sprintf("http://localhost:%d%s", c.AppPort, callbackURI)
	if config.DebugMode {
		log.Println("RedirectURL:", c.OAuth2Config.RedirectURL)
	}

	s := &Server{
		config:     c,
		context:    context.Background(),
		httpServer: &http.Server{Addr: fmt.Sprintf(":%d", c.AppPort)},
	}

	http.HandleFunc("/", s.handleIndexPage)
	http.HandleFunc("/login", s.redirectToAuthorisationEndpoint)
	http.HandleFunc(callbackURI, s.handleOAuthCallback)

	return s
}

// Start - Calls ListenAndServe() on the http server.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) redirectToAuthorisationEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Location", s.config.OAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) handleIndexPage(w http.ResponseWriter, req *http.Request) {
	if s.oAuthAuthorisationCode == "" {
		w.Header().Add("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	orgsResponse, err := s.httpClient.Get("https://api.xero.com/connections")

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

func (s *Server) handleOAuthCallback(w http.ResponseWriter, req *http.Request) {
	s.oAuthAuthorisationCode = req.URL.Query().Get("code")
	if config.DebugMode {
		log.Println("Received authorisation code:", s.oAuthAuthorisationCode)
	}
	s.httpClient = s.config.OAuth2Config.Client(s.context, s.getAuthorisationToken())
	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) getAuthorisationToken() *oauth2.Token {
	encodedAuthHeaderValue := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("Basic %s:%s", s.config.ClientID, s.config.ClientSecret),
	))
	tok, err := s.config.OAuth2Config.Exchange(
		s.context,
		s.oAuthAuthorisationCode,
		oauth2.SetAuthURLParam("authorization", encodedAuthHeaderValue),
	)
	if err != nil {
		log.Println("An error occurred while trying to exchange the authorisation code with the Xero API.")
		log.Fatal(err)
	}
	// Also update the server object
	s.oAuthToken = tok
	if config.DebugMode {
		log.Println("Got OAuth2 Token from API.")
		log.Println("Token expiry:", tok.Expiry.String())
	}
	return tok
}
