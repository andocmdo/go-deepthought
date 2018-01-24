package gostock

import "time"

// Job contains state data for jobs
type Job struct {
	ID         int               `json:"id"`
	WorkerID   int               `json:"workerID"` // updateable by server
	Valid      bool              `json:"valid"`
	Dispatched bool              `json:"dispatched"` // updateable by server
	Running    bool              `json:"running"`    // updateable by server
	Completed  bool              `json:"completed"`  // updateable by server
	Created    time.Time         `json:"created"`
	Started    time.Time         `json:"started"` // updateable by server
	Ended      time.Time         `json:"ended"`   // updateable by server
	Args       map[string]string `json:"args"`
	Result     string            `json:"result"`     // updateable by client
	LastUpdate time.Time         `json:"lastUpdate"` // updateable by server
}

// Jobs is a slice of Job
type Jobs []Job

// NewJob is a constructor for Job structs (init Args map)
func NewJob() *Job {
	var j Job
	j.Args = make(map[string]string)
	return &j
}
