package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/browser"

	"github.com/XeroAPI/golang-oauth2-example/config"

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
	const callbackURI = "/auth/xero"

	appConfig := config.New("")
	if config.DebugMode {
		log.Println("Loaded config:")
		appConfig.Print()
	}

	ctx := context.Background()
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

	server := &http.Server{Addr: fmt.Sprintf(":%d", appConfig.AppPort), Handler: nil}
	var authorisationCode string

	// We want to spawn the HTTP server at this point to make sure we're ready for the redirect from the Xero API.
	redirectHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
		authorisationCode = req.URL.Query().Get("code")
		if config.DebugMode {
			log.Println("Received authorisation code:", authorisationCode)
		}
		server.Shutdown(ctx)
	}

	http.HandleFunc(callbackURI, redirectHandler)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	browser.OpenURL(url)
	log.Println("If a browser window did not open, please open the following link in a browser to continue:")
	fmt.Println(url)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln(
			fmt.Sprintf("An error occurred while trying to spawn an HTTP server on port %d:", appConfig.AppPort),
			err,
		)
	}

	encodedAuthHeaderValue := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("Basic %s:%s", appConfig.ClientID, appConfig.ClientSecret),
	))
	tok, err := conf.Exchange(
		ctx,
		authorisationCode,
		oauth2.SetAuthURLParam("authorization", encodedAuthHeaderValue),
	)
	if err != nil {
		log.Println("An error occurred while trying to exchange the authorisation code with the Xero API.")
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)

	orgsRequest, err := http.NewRequest(http.MethodGet, "https://api.xero.com/connections", nil)
	// orgsRequest, err := http.NewRequest(http.MethodGet, "https://api.xero.com/api.xro/2.0/Organisations", nil)
	// orgsRequest.Header.Add("authorization", fmt.Sprintf("Bearer %s", tok.AccessToken))
	if err != nil {
		log.Println("An error occurred while trying to build a request for Organisations data.")
		log.Fatalln(err)
	}
	resp, err := client.Do(orgsRequest)

	if err != nil {
		log.Println("An error occurred while trying to get organisation data from the Xero API.")
		log.Fatalln(err)
	}
	if resp.StatusCode != 200 {
		log.Println("Got a non-200 response code:", resp.Status)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("An error occurred while trying to read the response from the Xero API")
		log.Fatalln(err)
	}
	fmt.Println(string(respBody))
}
