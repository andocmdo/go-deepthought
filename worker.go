package main

import "time"

// Worker contains state data for workers
type Worker struct {
	ID         int       `json:"id"`
	Valid      bool      `json:"valid"`
	Ready      bool      `json:"ready"`
	Working    bool      `json:"working"`
	IPAddr     string    `json:"ipaddr"`
	Port       int       `json:"port"`
	Created    time.Time `json:"created"`
	Destroyed  time.Time `json:"destroyed"`  // updateable
	LastUpdate time.Time `json:"lastUpdate"` // updateable
}

// Workers is a slice of worker
type Workers []Worker

var readyWorkers chan int

// NewWorker is a constructor for Worker structs (init Args map)
func NewWorker() *Worker {
	var w Worker
	//j.Args = make(map[string]string)
	return &w
}

func queueWorker(id int) {
	readyWorkers <- id
}
