package main

import (
	"log"
	"net/http"
	"net/url"
)

func main() {

	// TODO add a thread to check the status of workers (ping them essentially) an
	// set their "ready" status accordingly.

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

	select {} // run forever
}
