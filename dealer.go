package main

import (
	"log"
	"runtime"
	"strconv"
)

func init() {
	jobsToRun = make(chan int, 10000)
	readyWorkers = make(chan int, 500)

	// number of processor cores to keep free, the rest will be used to run jobs
	const keepFreeCores = 1
	cores := 1
	coresAvailable := runtime.NumCPU()
	log.Println("Number of processor cores available: " + strconv.FormatInt(int64(coresAvailable), 10))
	if coresAvailable > keepFreeCores {
		cores = coresAvailable - keepFreeCores
	}
	log.Println("Number of processor cores to use for job dealers: ", cores)

	for i := 0; i < cores; i++ {
		go dealer(i, jobsToRun, readyWorkers)
	}
}

// Dealer depyloys jobs to waiting workers
func dealer(d int, jobChan <-chan int, workerChan <-chan int) {
	log.Printf("started dealer %d", d)
	for {
		workerID := <-workerChan
		jobID := <-jobChan
		log.Printf("dealer %d is sending job %d to worker %d", d, jobID, workerID)
		job, err := RepoFindJob(jobID)
		if err != nil {
			log.Printf("dealer %d encountered an error finding job %d to send to worker %d", d, jobID, workerID)
			log.Printf(err.Error())
		}

		/* TODO these are update by the worker itself through api
		job.Running = true
		job.Started = time.Now()
		_, err = RepoUpdateJob(job)
		if err != nil {
			log.Printf("error on job %d", jobID)
			log.Printf(err.Error())
		}
		*/

		// This is where we send out job
		// connect to tcp port and send job data
		log.Printf("Dealer %d sent job %d to worker %d ", d, jobID, workerID)

		// And when finished, note the time, check for errors, etc
		job.Dispatched = true
		/* TODO These will be updated by the API
		job.Running = false
		job.Completed = true
		job.Result = string(out)
		*/
		_, err = RepoUpdateJob(job)
		if err != nil {
			log.Printf("error dispatching job %d", jobID)
			log.Printf(err.Error())
			return
		}
		log.Printf("dealer %d successfully dispatched job %d to worker %d", d, jobID, workerID)

	}
}
