package ui

import (
	"net/http"
	"strings"
)

// GlobalStyles is an array of styles that get concatenated and added to the global styles tag
var GlobalStyles = []string{
	"body {font-family: sans-serif}",
}

// WriteGlobalStylesTag will write the style tag to the response.
func WriteGlobalStylesTag(w http.ResponseWriter) {
	w.Write([]byte("<style>" + strings.Join(GlobalStyles, "") + "</style>"))
}
