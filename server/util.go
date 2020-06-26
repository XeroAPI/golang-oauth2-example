package server

import (
	"log"
	"net/http"
	"time"

	"github.com/XeroAPI/golang-oauth2-example/config"
	"golang.org/x/oauth2"
)

func (s *Server) handleOAuthCallback(w http.ResponseWriter, req *http.Request) {
	s.oAuthAuthorisationCode = req.URL.Query().Get("code")
	if config.DebugMode {
		log.Println("Received authorisation code:", s.oAuthAuthorisationCode)
	}
	s.httpClient = *s.config.OAuth2Config.Client(s.context, s.getAuthorisationToken())
	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// refreshAccessToken - Refresh logic inspired by https://github.com/golang/oauth2/issues/84#issuecomment-520099526 and
// https://stackoverflow.com/a/46487481
func (s *Server) refreshAccessToken() error {
	// We create a new token source that only has the refresh token, to force the OAuth2 client to retrieve a new access
	// token.
	src := s.config.OAuth2Config.TokenSource(s.context, &oauth2.Token{RefreshToken: s.oAuthToken.RefreshToken})
	newToken, err := src.Token()
	if err != nil {
		return err
	}
	// Also update the Server struct properties
	s.oAuthToken = newToken
	return nil
}

func (s *Server) getAuthorisationToken() *oauth2.Token {
	tok, err := s.config.OAuth2Config.Exchange(
		s.context,
		s.oAuthAuthorisationCode,
		oauth2.SetAuthURLParam(s.getAuthorisationHeader()),
	)
	if err != nil {
		log.Println("An error occurred while trying to exchange the authorisation code with the Xero API.")
		log.Fatalln(err)
	}
	// Also update the server struct
	s.oAuthToken = tok
	if config.DebugMode {
		log.Println("Got OAuth2 Token from API.")
		log.Println("Token expiry:", tok.Expiry.String())
	}
	return tok
}

// preFlightCheck is run before any call that requires access to the API to ensure that we still have tokens that are
// up to date.
// Returns true if upstreaming processing should be stopped (e.g. if we're triggering a redirect here).
func (s *Server) preFlightCheck(w http.ResponseWriter) bool {
	if !s.oAuthToken.Valid() {
		w.Header().Add("Location", loginPath)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return true
	}
	s.refreshAccessTokenIfNeeded()
	return false
}

// Refreshes the OAuth2 access token if it's within `renewalWindow` minutes of expiry.
func (s *Server) refreshAccessTokenIfNeeded() {
	renewalWindow := 15 * time.Minute
	now := time.Now()
	if !s.oAuthToken.Valid() {
		s.refreshAccessToken()
		return
	}
	if now.Add(renewalWindow).After(s.oAuthToken.Expiry) {
		s.refreshAccessToken()
	}
}
