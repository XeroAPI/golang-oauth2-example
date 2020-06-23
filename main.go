package main

import (
	"github.com/XeroAPI/golang-oauth2-example/server"
)

func main() {

	server := server.New()
	server.ListenAndServe()
}
