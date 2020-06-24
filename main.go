package main

import (
	"log"

	"github.com/XeroAPI/golang-oauth2-example/config"
	"github.com/XeroAPI/golang-oauth2-example/server"
	"golang.org/x/oauth2"
)

func main() {
	const callbackURI = "/auth/xero"
	oAuthScopes := []string{
		"openid",
		"profile",
		"email",
		"accounting.transactions",
		"accounting.settings",
		"offline_access",
	}
	appConfig := config.New("")
	appConfig.OAuth2Config = &oauth2.Config{
		ClientID:     appConfig.ClientID,
		ClientSecret: appConfig.ClientSecret,
		Scopes:       oAuthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.xero.com/identity/connect/authorize",
			TokenURL: "https://identity.xero.com/connect/token",
		},
	}
	server := server.New(appConfig)
	err := server.Start()
	if err != nil {
		log.Fatalln("An error with the web server occurred:", err)
	}
}
