package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/XeroAPI/golang-oauth2-example/config"
	"golang.org/x/oauth2"
)

// Server - Represents the running app, and all the bits and pieces we need to make API calls.
type Server struct {
	config                 *config.Config
	context                context.Context
	httpServer             *http.Server
	httpClient             http.Client
	oAuthAuthorisationCode string
	oAuthToken             *oauth2.Token
}

// Some package level variables
const callbackURI = "/auth/xero"
const loginPath = "/login"
const returnToHomepageLink = "<a href=\"/\">Return to the homepage</a>"

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
	http.HandleFunc("/organisation", s.handleOrganisationPage)
	http.HandleFunc(loginPath, s.redirectToAuthorisationEndpoint)
	http.HandleFunc(callbackURI, s.handleOAuthCallback)

	return s
}

// Start - Calls ListenAndServe() on the http server.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// getAuthorisationHeader returns the header used to authorise API requests.
func (s *Server) getAuthorisationHeader() (string, string) {
	return "authorisation", base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("Basic %s:%s", s.config.ClientID, s.config.ClientSecret),
	))
}

func (s *Server) redirectToAuthorisationEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Location", s.config.OAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
