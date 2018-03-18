package main

import (
	"fmt"
	"sync"
	"time"
)

var nextJobID int
var jobs Jobs
var jobMutex *sync.Mutex

// Give us some seed data
func init() {
	jobMutex = &sync.Mutex{}
	jobMutex.Lock()
	defer jobMutex.Unlock()
	nextJobID = 1
}

// RepoFindJob searches for a job with id inside mock DB
func RepoFindJob(id int) (Job, error) {
	jobMutex.Lock()
	defer jobMutex.Unlock()
	if validJobID(id) {
		return jobs[id-1], nil
	}
	return Job{}, fmt.Errorf("can find job: %d", id)
}

// RepoCreateJob takes a job and assigns it the next ID, then adds to jobs slice
func RepoCreateJob(j Job) Job {
	j.ID = nextJobID
	jobMutex.Lock()
	defer jobMutex.Unlock()
	jobs = append(jobs, j)
	jobsToRun <- j.ID
	nextJobID++
	return j
}

// RepoUpdateJob updates a job that matches input job.ID, only updating updateable fields
func RepoUpdateJob(job Job) (Job, error) {
	// check sanity first
	if validJobID(job.ID) {
		jobMutex.Lock()
		defer jobMutex.Unlock()
		index := job.ID - 1

		jobs[index].WorkerID = job.WorkerID
		jobs[index].Dispatched = job.Dispatched
		jobs[index].Running = job.Running
		jobs[index].Completed = job.Completed
		jobs[index].Cancelled = job.Cancelled
		jobs[index].Started = job.Started
		jobs[index].Ended = job.Ended
		jobs[index].Result = job.Result
		jobs[index].Success = job.Success
		jobs[index].Cancel = job.Cancel
		jobs[index].LastUpdate = time.Now()
		return jobs[index], nil
	}
	job.Valid = false
	return job, fmt.Errorf("job ID not found")
}

func validJobID(id int) bool {
	if id > 0 && len(jobs) != 0 && id < nextJobID { // nextJobID? or len(jobs), this is jank
		return true
	}
	return false
}
