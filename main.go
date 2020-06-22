package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/XeroAPI/golang-oauth2-example/config"

	"golang.org/x/oauth2"
)

func main() {
	log.Println("Works on my machine!")

	oAuthScopes := []string{
		"openid",
		"profile",
		"email",
		"accounting.transactions",
		"accounting.settings",
		"offline_access",
	}

	appConfig := config.New("")
	if config.DebugMode {
		log.Println("Loaded config:")
		appConfig.Print()
	}

	os.Exit(0)

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		Scopes:       oAuthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.xero.com/identity/connect/authorize",
			TokenURL: "https://identity.xero.com/connect/token",
		},
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	client.Get("...")
}
