package main

import "time"

// Worker contains state data for workers
type Worker struct {
	ID         int       `json:"id"`
	Valid      bool      `json:"valid"`
	Created    time.Time `json:"created"`
	IPAddr     string    `json:"ipaddr"`
	Port       string    `json:"port"`
	Ready      bool      `json:"ready"`      // updateable
	Working    bool      `json:"working"`    // updateable
	LastUpdate time.Time `json:"lastUpdate"` // updateable
}

// Workers is a slice of worker
type Workers []Worker

// NewWorker is a constructor for Worker structs (init Args map)
func NewWorker() *Worker {
	var w Worker
	//j.Args = make(map[string]string)
	return &w
}
