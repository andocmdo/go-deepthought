package main

import (
	"log"
	"os/exec"
	"time"
)

func init() {

	jobsToRun = make(chan int, 1000)

}

func worker(w int, jobChan <-chan int) {
	log.Printf("started worker %d", w)

	for id := range jobChan {
		log.Printf("worker %d started job %d", w, id)
		job, err := RepoFindJob(id)
		if err != nil {
			log.Printf("error on job %d", id)
			log.Printf(err.Error())
		}
		job.Running = true
		job.Started = time.Now()
		_, err = RepoUpdateJob(job)
		if err != nil {
			log.Printf("error on job %d", id)
			log.Printf(err.Error())
		}

		// This is where we would process our job
		cmd := exec.Command("bash", "-c", "sleep 10; date")
		out, err := cmd.Output()
		if err != nil {
			log.Printf("error on job %d", id)
			log.Printf(err.Error())
		}
		log.Printf("Job %d output: %s", id, out)

		// And when finished, note the time, check for errors, etc
		job.Ended = time.Now()
		job.Running = false
		job.Completed = true
		job.Result = string(out)

		_, err = RepoUpdateJob(job)
		if err != nil {
			log.Printf("error on job %d", id)
			log.Printf(err.Error())
			return
		}
		log.Printf("worker %d finished job %d", w, id)
	}
}
