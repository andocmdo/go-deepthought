package main

import (
	"log"
	"runtime"
	"strconv"
	"time"
)

var toRun chan int

func init() {
	// number of processor cores to keep free, the rest will be used to run jobs
	const keepFreeCores = 1
	cores := 1
	coresAvailable := runtime.NumCPU()
	log.Println("Number of processor cores available: " + strconv.FormatInt(int64(coresAvailable), 10))
	if coresAvailable > keepFreeCores {
		cores = coresAvailable - keepFreeCores
	}
	log.Println("Number of processor cores to use: ", cores)

	toRun = make(chan int, 100)

	for i := 0; i < cores; i++ {
		go worker(i, toRun)
	}
}

func queueJob(id int) {
	toRun <- id
}

func worker(w int, jobChan <-chan int) {
	log.Printf("started worker %d", w)
	for id := range jobChan {
		log.Printf("worker %d started job %d", w, id)
		job, err := RepoFindJob(id)
		if err != nil {
			log.Printf("error on job %d", id)
			log.Printf(err.Error())
			return
		}
		job.Running = true
		job.Started = time.Now()
		err = RepoUpdateJob(job)
		if err != nil {
			log.Printf("error on job %d", id)
			log.Printf(err.Error())
			return
		}

		// This is where we would process our job
		time.Sleep(time.Second * 20)

		// And when finished, note the time, check for errors, etc
		job.Ended = time.Now()
		job.Running = false
		job.Completed = true

		err = RepoUpdateJob(job)
		if err != nil {
			log.Printf("error on job %d", id)
			log.Printf(err.Error())
			return
		}
		log.Printf("worker %d finished job %d", w, id)
	}
}
