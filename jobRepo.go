package main

import (
	"fmt"
	"sync"
	"time"
)

var currentJobID int
var jobs Jobs
var jobMutex *sync.Mutex

// Give us some seed data
func init() {
	jobMutex = &sync.Mutex{}
	jobMutex.Lock()
	defer jobMutex.Unlock()
	currentJobID = 0
}

// RepoFindJob searches for a job with id inside mock DB
func RepoFindJob(id int) (Job, error) {
	jobMutex.Lock()
	defer jobMutex.Unlock()
	if id <= currentJobID || len(jobs) != 0 { // currentJobID? or len(jobs), this is jank
		return jobs[id], nil
	}
	return Job{}, fmt.Errorf("can find job: %d", id)
}

// RepoCreateJob takes a job and assigns it the next ID, then adds to jobs slice
func RepoCreateJob(j Job) Job {
	j.ID = currentJobID
	jobMutex.Lock()
	defer jobMutex.Unlock()
	jobs = append(jobs, j)
	queueJob(j.ID)
	currentJobID++
	return j
}

// RepoUpdateJob updates a job that matches input job.ID, only updating updateable fields
func RepoUpdateJob(job Job) error {
	// check sanity first
	if job.ID < 0 || job.Valid == false {
		return fmt.Errorf("job is not valid or has illegal ID")
	}
	jobMutex.Lock()
	defer jobMutex.Unlock()
	for i, j := range jobs {
		if j.ID == job.ID {
			jobs[i].Running = job.Running
			jobs[i].Started = job.Started
			jobs[i].Ended = job.Ended
			jobs[i].Completed = job.Completed
			jobs[i].Result = job.Result
			jobs[i].LastUpdate = time.Now()
			return nil
		}
	}
	return fmt.Errorf("job ID not found")
}
