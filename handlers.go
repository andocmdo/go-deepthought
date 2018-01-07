package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// max size of upload
const uploadLimit = 1048576

/*
// Index returns a Status OK header and plain text string to
// verify server is working
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "It's working...")
}
*/

// Index returns files in the working directory
func Index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["file"]
	log.Printf(file)
	if file == "" {
		file = "index.html"
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}
