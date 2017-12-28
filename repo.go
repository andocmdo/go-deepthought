package main

import (
	"fmt"
	"sync"
	"time"
)

var currentID int
var jobs Jobs
var mutex *sync.Mutex

// Give us some seed data
func init() {
	mutex = &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	currentID = 1
	// TODO REMOVE THESE FAKE JOBS!!!
	//RepoCreateJob(Job{Running: false, Start: time.Now(), Symbol: "TECD"})
	//RepoCreateJob(Job{Running: false, Symbol: "AAPL"})

}

// RepoFindJob searches for a job with id inside mock DB
func RepoFindJob(id int) (Job, error) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, j := range jobs {
		if j.ID == id {
			return j, nil
		}
	}
	return Job{}, fmt.Errorf("can find job: %d", id)
}

// RepoCreateJob takes a job and assigns it the next ID, then adds to jobs slice
func RepoCreateJob(j Job) Job {
	j.ID = currentID
	mutex.Lock()
	defer mutex.Unlock()
	jobs = append(jobs, j)
	queueJob(j.ID)
	currentID++
	return j
}

// RepoUpdateJob updates a job that matches input job.ID, only updating updateable fields
func RepoUpdateJob(job Job) error {
	// check sanity first
	if job.ID == 0 || job.Valid == false {
		return fmt.Errorf("job is not valid or has ID of zero")
	}
	mutex.Lock()
	defer mutex.Unlock()
	for i, j := range jobs {
		if j.ID == job.ID {
			jobs[i].Running = job.Running
			jobs[i].Started = job.Started
			jobs[i].Ended = job.Ended
			jobs[i].Completed = job.Completed
			jobs[i].LastUpdate = time.Now()
			return nil
		}
	}
	return fmt.Errorf("job ID not found")
}

// RepoDestroyJob searches for job with id to delete. If found, it is removed
// if not found it returns an error
// DON"T ACTUALLY USE!
func RepoDestroyJob(id int) error {
	mutex.Lock()
	defer mutex.Unlock()
	for i, j := range jobs {
		if j.ID == id {
			jobs = append(jobs[:i], jobs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find job with ID of %d to delete", id)
}
