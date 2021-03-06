package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
	defaultAppPort := 8080

	// Set the port the webserver will listen on
	if c.AppPort == 0 {
		if envAppPort := os.Getenv("APP_PORT"); envAppPort != "" {
			var err error
			c.AppPort, err = strconv.Atoi(envAppPort)
			if err != nil {
				log.Fatalln("An error occurred while trying to read the APP_PORT environment variable:", err)
			}
		} else {
			c.AppPort = defaultAppPort
		}
	}

	// Set the redirect URL
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
	http.HandleFunc("/invoices", s.handleInvoicePage)
	http.HandleFunc("/refresh", s.handleTokenRefreshRequest)
	http.HandleFunc("/segfault", s.handleSegfaultRequest)
	http.HandleFunc(loginPath, s.redirectToAuthorisationEndpoint)
	http.HandleFunc(callbackURI, s.handleOAuthCallback)

	return s
}

// Start - Calls ListenAndServe() on the http server.
func (s *Server) Start() error {
	log.Printf("Hey there! I'm up and running, and can be accessed at: http://localhost:%d\n", s.config.AppPort)
	return s.httpServer.ListenAndServe()
}

// getAuthorisationHeader returns the header used to authorise API requests.
func (s *Server) getAuthorisationHeader() (string, string) {
	return "authorization", base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("Basic %s:%s", s.config.ClientID, s.config.ClientSecret),
	))
}

func (s *Server) redirectToAuthorisationEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Location", s.config.OAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
