package gostock

import "time"

// Job contains state data for jobs
type Job struct {
	ID         int               `json:"id"`
	WorkerID   int               `json:"workerID"` // updateable
	Valid      bool              `json:"valid"`
	Dispatched bool              `json:"dispatched"` // updateable
	Running    bool              `json:"running"`    // updateable
	Completed  bool              `json:"completed"`  // updateable
	Created    time.Time         `json:"created"`
	Cancelled  time.Time         `json:"cancelled"` // updateable
	Started    time.Time         `json:"started"`   // updateable
	Ended      time.Time         `json:"ended"`     // updateable
	Args       map[string]string `json:"args"`
	Result     string            `json:"result"`     // updateable
	Success    bool              `json:"success"`    // updateable
	Cancel     bool              `json:"cancel"`     // updateable
	LastUpdate time.Time         `json:"lastUpdate"` // updateable
}

// Jobs is a slice of Job
type Jobs []Job

// NewJob is a constructor for Job structs (init Args map)
func NewJob() *Job {
	var j Job
	j.Args = make(map[string]string)
	return &j
}
