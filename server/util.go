package server

import (
	"log"
	"net/http"

	"github.com/XeroAPI/golang-oauth2-example/config"
	"golang.org/x/oauth2"
)

func (s *Server) handleOAuthCallback(w http.ResponseWriter, req *http.Request) {
	s.oAuthAuthorisationCode = req.URL.Query().Get("code")
	if config.DebugMode {
		log.Println("Received authorisation code:", s.oAuthAuthorisationCode)
	}
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
	s.httpClient = *s.config.OAuth2Config.Client(s.context, tok)
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

// preFlightCheck is run before any call that requires access to the API to ensure that we still have tokens that are
// up to date.
// Returns true if redirectToLogin() should be called.
func (s *Server) preFlightCheck() bool {
	// We return true if the oAuthToken is not valid, and that we should redirect to login.
	// We don't perform the redirect here, as some pages will have mixed behaviour (e.g. the index page)
	return !s.oAuthToken.Valid()
}

func (s *Server) redirectToLogin(w http.ResponseWriter) {
	w.Header().Add("Location", loginPath)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
