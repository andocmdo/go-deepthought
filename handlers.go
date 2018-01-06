package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

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

// FrontEnd returns files in the working directory
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

// JobIndex gets all jobs as JSON
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
	var job Job
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
	job.Recieved = time.Now()
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
	job := Job{}
	symbol := r.FormValue("symbol")
	startDate := r.FormValue("startDate")
	endDate := r.FormValue("endDate")
	popSizeS := r.FormValue("popSize")
	popSize, err := strconv.Atoi(popSizeS)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("popSize unable to be parsed in urlencoded form")
		fmt.Fprintln(w, "popSize unable to be parsed in urlencoded form")
		return
	}
	mutRateS := r.FormValue("mutRate")
	mutRate, err := strconv.ParseFloat(mutRateS, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("mutRate unable to be parsed in urlencoded form")
		fmt.Fprintln(w, "mutRate unable to be parsed in urlencoded form")
		return
	}
	maxGenS := r.FormValue("maxGen")
	maxGen, err := strconv.Atoi(maxGenS)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("maxGen unable to be parsed in urlencoded form")
		fmt.Fprintln(w, "maxGen unable to be parsed in urlencoded form")
		return
	}
	if (symbol == "") || (startDate == "") || (endDate == "") {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("no symbol in urlencoded form")
		fmt.Fprintln(w, "no symbol in urlencoded form")
		return
	}
	job.Symbol = symbol
	job.StartDate = startDate
	job.EndDate = endDate
	job.PopSize = popSize
	job.MutRate = mutRate
	job.MaxGen = maxGen
	job.Recieved = time.Now()
	job.Valid = true
	j := RepoCreateJob(job)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}
