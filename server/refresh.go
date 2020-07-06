package server

import (
	"log"
	"net/http"
)

func (s *Server) handleTokenRefreshRequest(w http.ResponseWriter, req *http.Request) {
	if s.preFlightCheck() {
		s.redirectToLogin(w)
		return
	}

	err := s.refreshAccessToken()

	if err != nil {
		errMsg := "An error occurred while trying to refresh the oAuth 2.0 token."
		w.Write([]byte(errMsg))
		log.Println(errMsg)
		log.Println(err)
		return
	}
	if !s.oAuthToken.Valid() {
		errMsg := "An error occurred while trying to refresh the tokens - The token that was returned is not valid."
		w.Write([]byte(errMsg))
		log.Println(errMsg)
		w.Write([]byte(returnToHomepageLink))
		return
	}
	w.Write([]byte("<p>Tokens refreshed! Token now expires at " + s.oAuthToken.Expiry.String() + "</p>"))
	log.Println("Token refreshed - Now expires at", s.oAuthToken.Expiry.String())
	w.Write([]byte(returnToHomepageLink))
}
