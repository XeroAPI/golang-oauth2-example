package server

import (
	"net/http"

	"github.com/XeroAPI/golang-oauth2-example/server/ui"
)

func (s *Server) handleSegfaultRequest(w http.ResponseWriter, req *http.Request) {
	if s.preFlightCheck() {
		s.redirectToLogin(w)
		return
	}

	ui.WriteGlobalStylesTag(w)
	w.Write([]byte("<h1>Why would you do that!?</h1>"))
	w.Write([]byte(returnToHomepageLink))
}
