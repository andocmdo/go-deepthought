package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	gostock "github.com/andocmdo/gostockd/common"
	"github.com/gorilla/mux"
)

// WorkerIndex gets all workers as JSON
func WorkerIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(workers); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}

// WorkerShow attemps to get a specific worker based on ID
func WorkerShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workerID, err := strconv.Atoi(vars["workerID"])
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	worker, err := RepoFindWorker(workerID)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(worker); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}

// WorkerCreateJSON creates a worker from JSON POST data to /workers endpoint
func WorkerCreateJSON(w http.ResponseWriter, r *http.Request) {
	var worker gostock.Worker
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, uploadLimit))
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if err = r.Body.Close(); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if err = json.Unmarshal(body, &worker); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if worker.Port == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		log.Printf("missing port for worker creation request, it is required.")
		fmt.Fprintln(w, err.Error())
		return
	}

	worker.Created = time.Now()
	worker.Valid = true
	ipAndPort := r.RemoteAddr
	worker.IPAddr = strings.Split(ipAndPort, ":")[0]
	wrkr := RepoCreateWorker(worker)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(wrkr); err != nil {
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}

// WorkerCreateURLEnc creates a worker from JSON POST data to /workers endpoint
func WorkerCreateURLEnc(w http.ResponseWriter, r *http.Request) {
	//worker := Worker{}
	worker := gostock.NewWorker()

	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("error parsing worker form in urlencoded form")
		fmt.Fprintln(w, "error parsing worker form values")
		return
	}
	worker.IPAddr = r.FormValue("ipaddr")
	worker.Port = r.FormValue("port")
	if worker.Port == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		log.Printf("missing port for worker creation request, it is required.")
		fmt.Fprintln(w, "missing port for worker creation request, it is required.")
		return
	}

	worker.Created = time.Now()
	worker.Valid = true
	wrkr := RepoCreateWorker(*worker)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(wrkr); err != nil {
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}

// WorkerUpdateJSON creates a worker from JSON POST data to /workers endpoint
func WorkerUpdateJSON(w http.ResponseWriter, r *http.Request) {
	var worker gostock.Worker
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, uploadLimit))
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if errr := r.Body.Close(); errr != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	if errrr := json.Unmarshal(body, &worker); errrr != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	// Worker ID and URL should match
	vars := mux.Vars(r)
	workerID, err := strconv.Atoi(vars["workerID"])
	if err != nil || workerID != worker.ID {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	// TODO remove this debug stuff after fixing bug for index out of bounds for worker

	//TODO meat and potatoes here until I refactor
	wrkr, err := RepoUpdateWorker(worker) // check this error
	if err != nil || workerID != worker.ID {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(wrkr); err != nil {
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}
