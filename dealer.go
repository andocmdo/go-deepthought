package main

import (
	"encoding/json"
	"log"
	"net"
)

const jobQueueSize = 10000
const workerQueueSize = 500

var jobsToRun chan int
var readyWorkers chan int

func init() {
	jobsToRun = make(chan int, jobQueueSize)
	readyWorkers = make(chan int, workerQueueSize)
	log.Println("Job queue size: ", jobQueueSize)
	log.Println("Worker queue size: ", workerQueueSize)

	go dealer(0, jobsToRun, readyWorkers)
}

// Dealer depyloys jobs to waiting workers
func dealer(d int, jobChan <-chan int, workerChan <-chan int) {
	log.Printf("started dealer %d", d)

	for {
		workerID := <-workerChan
		jobID := <-jobChan
		log.Printf("dealer %d is sending job %d to worker %d", d, jobID, workerID)
		job, err := RepoFindJob(jobID)
		log.Printf("dealer %d found job %d : %+v", d, jobID, job)
		if err != nil {
			log.Printf("dealer %d encountered an error finding job %d to send to worker %d", d, jobID, workerID)
			log.Printf(err.Error())
		}
		wrkr, err := RepoFindWorker(workerID)
		if err != nil {
			log.Printf("dealer %d encountered an error finding worker %d", d, workerID)
			log.Printf(err.Error())
		}

		// This is where we send out job
		// connect to tcp port and send job data
		conn, err := net.Dial("tcp", wrkr.IPAddr+":"+wrkr.Port)

		if err != nil {
			log.Printf("dealer %d encountered an error connecting to worker %d", d, workerID)
			log.Printf(err.Error())
		}

		// using JSON now instead of gob, keeping this for notes. Remove before committing Andy
		enc := json.NewEncoder(conn) // Will write to network.
		dec := json.NewDecoder(conn) // Will read from network.
		err = enc.Encode(job)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		log.Printf("Dealer %d sent job %d to worker %d ", d, jobID, workerID)

		err = dec.Decode(&job)
		if err != nil {
			// TODO check if job is valid, if not there was some other error, and we
			// should reschedule the job here
			log.Fatal("decode error 1: ", err)
		}

		log.Printf("Dealer %d got confirmation from worker %d ", d, workerID)
		log.Printf("%+v", job)

		// And when finished, note the time, check for errors, etc
		job.Dispatched = true
		job.WorkerID = job.ID
		//worker.JobID //TODO update the worker with the jobID it is running.
		_, err = RepoUpdateJob(job)
		if err != nil {
			log.Printf("error dispatching job %d", jobID)
			log.Printf(err.Error())
			return
		}
		log.Printf("dealer %d successfully dispatched job %d to worker %d", d, jobID, workerID)
		conn.Close()

	}
}
