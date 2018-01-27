package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	gostock "github.com/andocmdo/gostockd/common"
)

var nextWorkerID int
var workers gostock.Workers
var workerMutex *sync.Mutex

// Initialize mutex and fake DB
func init() {
	workerMutex = &sync.Mutex{}
	workerMutex.Lock()
	defer workerMutex.Unlock()
	nextWorkerID = 1
}

// RepoFindWorker searches for a worker with id inside mock DB
func RepoFindWorker(id int) (gostock.Worker, error) {
	workerMutex.Lock()
	defer workerMutex.Unlock()
	if validWorkerID(id) { // currentID? or len(workers), this is jank
		return workers[id-1], nil
	}
	w := gostock.NewWorker()
	w.Valid = false
	return *w, fmt.Errorf("can find worker: %d", id)
}

// RepoCreateWorker takes a worker and assigns it the next ID, then adds to workers slice
func RepoCreateWorker(w gostock.Worker) gostock.Worker {
	w.ID = nextWorkerID
	//log.Print("RepoCreateWorker", w)
	workerMutex.Lock()
	defer workerMutex.Unlock()
	workers = append(workers, w)
	nextWorkerID++
	return w
}

// RepoUpdateWorker updates a worker that matches input worker.ID, only updating updateable fields
// TODO this is horrific crap with the +1. HAVE TO FIX THIS
func RepoUpdateWorker(worker gostock.Worker) (gostock.Worker, error) {
	// check sanity first
	if validWorkerID(worker.ID) {
		workerMutex.Lock()
		defer workerMutex.Unlock()
		index := worker.ID - 1

		//TODO remove this debug loggin
		log.Printf("setting worker at index %d to have a jobID of %d", index, worker.JobID)
		workers[index].JobID = worker.JobID
		log.Printf("set worker at index %d to have a jobID of %d", index, worker.JobID)
		workers[index].Ready = worker.Ready
		workers[index].Working = worker.Working
		//workers[i].IPAddr = worker.IPAddr
		//workers[i].Port = worker.Port
		workers[index].LastUpdate = time.Now()

		// if this update was to notify a worker was ready, then add to the queue
		if workers[index].Ready == true {
			readyWorkers <- worker.ID
		}

		// now return the updated info
		return workers[index], nil
	}
	worker.Valid = false
	return worker, fmt.Errorf("worker ID not found")
}

func validWorkerID(id int) bool {
	if id > 0 && len(workers) != 0 && id < nextWorkerID { // currentID? or len(workers), this is jank
		return true
	}
	return false
}
