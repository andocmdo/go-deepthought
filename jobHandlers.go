package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	gostock "github.com/andocmdo/gostockd/common"
	"github.com/gorilla/mux"
)

// JobIndex gets all jobs as JSON
// TODO maybe here is where we implement the filtering for running, cancelled, failed, etc jobs
// search through jobs slice, and create a temp slice/array to encode into JSON to return
// use r.FormValue to get the filter.
func JobIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}

// JobShow attemps to get a specific job based on ID
func JobShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := strconv.Atoi(vars["jobID"])
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	job, err := RepoFindJob(jobID)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(job); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}

// JobCreateJSON creates a job from JSON POST data to /jobs endpoint
func JobCreateJSON(w http.ResponseWriter, r *http.Request) {
	var job gostock.Job
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, uploadLimit))
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if err := r.Body.Close(); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if err := json.Unmarshal(body, &job); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	//TODO Meat and potatoes here until I refactor
	job.Created = time.Now()
	job.Valid = true
	j := RepoCreateJob(job)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}

// JobCreateURLEnc creates a job from JSON POST data to /jobs endpoint
func JobCreateURLEnc(w http.ResponseWriter, r *http.Request) {
	//job := Job{}
	job := gostock.NewJob()

	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("no symbol in urlencoded form")
		fmt.Fprintln(w, "error parsing form values")
		return
	}
	for key, values := range r.PostForm {
		job.Args[key] = values[0] // only using the first occurence of the parameter
	}

	if (job.Args["symbol"] == "") || (job.Args["startDate"] == "") || (job.Args["endDate"] == "") {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("no symbol in urlencoded form")
		fmt.Fprintln(w, "no symbol in urlencoded form")
		return
	}

	// Meat and Potatoes here, until I refactor this mess....
	job.Created = time.Now()
	job.Valid = true
	j := RepoCreateJob(*job)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}

// JobUpdateJSON updates a job from JSON POST data to /jobs endpoint
func JobUpdateJSON(w http.ResponseWriter, r *http.Request) {
	var job gostock.Job
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, uploadLimit))
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if err := r.Body.Close(); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if err := json.Unmarshal(body, &job); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	// We make sure the JSON and request URL match
	vars := mux.Vars(r)
	jobID, err := strconv.Atoi(vars["jobID"])
	if err != nil || jobID != job.ID {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	// TODO Meat and Potatoes here, until I refactor this mess....
	j, _ := RepoUpdateJob(job) // check this error
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}
