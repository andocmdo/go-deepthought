package main

import "time"

// Job contains state data for jobs
type Job struct {
	ID         int         `json:"id"`
	Valid      bool        `json:"valid"`
	Running    bool        `json:"running"`   // updateable
	Completed  bool        `json:"completed"` // updateable
	Recieved   time.Time   `json:"recieved"`
	Started    time.Time   `json:"started"`    // updateable
	Ended      time.Time   `json:"ended"`      // updateable
	LastUpdate time.Time   `json:"lastUpdate"` // updateable
	Command    Commandline `json:"command"`
	Result     string      `json:"result"`
}

type Commandline struct {
	Program string            `json:"program"`
	Args    map[string]string `json:"args"`
}

// Jobs is a slice of Job
type Jobs []Job
