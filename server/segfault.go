package server

import (
	"net/http"
)

func (s *Server) handleSegfaultRequest(w http.ResponseWriter, req *http.Request) {
	if s.preFlightCheck() {
		s.redirectToLogin(w)
		return
	}

	w.Write([]byte("<h1>Why would you do that!?</h1>"))
	w.Write([]byte(returnToHomepageLink))
}
