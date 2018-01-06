package main

import "time"

// Job contains state data for jobs
type Job struct {
	ID         int       `json:"id"`
	Valid      bool      `json:"valid"`
	Running    bool      `json:"running"`   // updateable
	Completed  bool      `json:"completed"` // updateable
	Recieved   time.Time `json:"recieved"`
	Started    time.Time `json:"started"`    // updateable
	Ended      time.Time `json:"ended"`      // updateable
	LastUpdate time.Time `json:"lastUpdate"` // updateable
	Symbol     string    `json:"symbol"`
	Result     string    `json:"result"` //updateable
	StartDate  string    `json:"startDate"`
	EndDate    string    `json:"endDate"`
	PopSize    int       `json:"popSize"`
	MutRate    float64   `json:"mutRate"`
	MaxGen     int       `json:"maxGen"`
}

// Jobs is a slice of Job
type Jobs []Job
