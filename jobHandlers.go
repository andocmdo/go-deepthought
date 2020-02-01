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

// JobSummary provides running, queued, completed, and failed summary as JSON
func JobSummary(w http.ResponseWriter, r *http.Request) {
	answer := make(map[string]int)
	answer["running"] = 0
	answer["queued"] = 0
	answer["completed"] = 0
	answer["failed"] = 0
	answer["total"] = 0

	for i := 0; i < len(jobs); i++ {
		if jobs[i].Running {
			answer["running"] += 1
		}
		if !jobs[i].Dispatched {
			answer["queued"] += 1
		}
		if jobs[i].Completed {
			answer["completed"] += 1
		}
		if jobs[i].Completed && !jobs[i].Success {
			answer["failed"] += 1
		}
	}
	answer["total"] = len(jobs)
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf(err.Error())
		fmt.Fprintln(w, err.Error())
	}
}

// JobShow attemps to get a specific job based on ID
func JobShow(w http.ResponseWriter, r *http.Request) {
	// this is the variables passed in URL
	vars := mux.Vars(r)

	// try to read the URL as a jobID
	jobID, err := strconv.Atoi(vars["jobID"])
	if err != nil {
		// must not be a number...

		// this will be the encoded list of jobs, based on filter,
		// (or no filter / all jobs) defaults to all jobs
		var answer Jobs

		// check for filters TODO use cases here?
		if vars["jobID"] == "running" {
			for i := 0; i < len(jobs); i++ {
				if jobs[i].Running {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "dispatched" {
			for i := 0; i < len(jobs); i++ {
				if jobs[i].Dispatched {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "completed" {
			for i := 0; i < len(jobs); i++ {
				if jobs[i].Completed {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "cancelled" {
			for i := 0; i < len(jobs); i++ {
				if jobs[i].Cancel {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "successful" {
			for i := 0; i < len(jobs); i++ {
				if jobs[i].Success {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "notRunning" {
			for i := 0; i < len(jobs); i++ {
				if !jobs[i].Running {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "notDispatched" {
			for i := 0; i < len(jobs); i++ {
				if !jobs[i].Dispatched {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "notCompleted" {
			for i := 0; i < len(jobs); i++ {
				if !jobs[i].Completed {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "notCancelled" {
			for i := 0; i < len(jobs); i++ {
				if !jobs[i].Cancel {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "notSuccessful" {
			for i := 0; i < len(jobs); i++ {
				if !jobs[i].Success {
					answer = append(answer, jobs[i])
				}
			}
		} else if vars["jobID"] == "failed" {
			for i := 0; i < len(jobs); i++ {
				if jobs[i].Completed && !jobs[i].Success {
					answer = append(answer, jobs[i])
				}
			}
		} else {
			// final fall through case.
			// not number, not any of above
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			log.Printf(err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(answer); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf(err.Error())
			fmt.Fprintln(w, err.Error())
		}
	} else {
		// jobID was a number, so go get the job
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
	job := NewJob()

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

	if (job.Args["command"] == "")  {
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
