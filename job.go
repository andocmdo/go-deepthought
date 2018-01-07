package main

import "time"

// Job contains state data for jobs
type Job struct {
	ID         int               `json:"id"`
	Valid      bool              `json:"valid"`
	Running    bool              `json:"running"`   // updateable
	Completed  bool              `json:"completed"` // updateable
	Recieved   time.Time         `json:"recieved"`
	Started    time.Time         `json:"started"`    // updateable
	Ended      time.Time         `json:"ended"`      // updateable
	LastUpdate time.Time         `json:"lastUpdate"` // updateable
	Args       map[string]string `json:"args"`
	Result     string            `json:"result"` // updateable
}

// Jobs is a slice of Job
type Jobs []Job

// NewJob is a constructor for Job structs (init Args map)
func NewJob() *Job {
	var j Job
	j.Args = make(map[string]string)
	return &j
}
