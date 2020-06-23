package main

import (
	"fmt"
	"log"

	"github.com/XeroAPI/golang-oauth2-example/config"
	"github.com/XeroAPI/golang-oauth2-example/server"

	"golang.org/x/oauth2"
)

func main() {
	var oAuthScopes = []string{
		"openid",
		"profile",
		"email",
		"accounting.transactions",
		"accounting.settings",
		"offline_access",
	}

	appConfig := config.New("")
	const callbackURI = "/auth/xero"

	conf := &oauth2.Config{
		ClientID:     appConfig.ClientID,
		ClientSecret: appConfig.ClientSecret,
		Scopes:       oAuthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.xero.com/identity/connect/authorize",
			TokenURL: "https://identity.xero.com/connect/token",
		},
		RedirectURL: fmt.Sprintf("http://localhost:%d%s", appConfig.AppPort, callbackURI),
	}
	if config.DebugMode {
		log.Println("RedirectURL:", conf.RedirectURL)
	}

	server := server.New(appConfig, conf, callbackURI)
	server.ListenAndServe()
}
