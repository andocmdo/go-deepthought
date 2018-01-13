package main

import (
	"fmt"
	"sync"
	"time"

	gostock "github.com/andocmdo/gostockd/common"
)

var currentWorkerID int
var workers gostock.Workers
var workerMutex *sync.Mutex

// Initialize mutex and fake DB
func init() {
	workerMutex = &sync.Mutex{}
	workerMutex.Lock()
	defer workerMutex.Unlock()
	currentWorkerID = 0
}

// RepoFindWorker searches for a worker with id inside mock DB
func RepoFindWorker(id int) (gostock.Worker, error) {
	workerMutex.Lock()
	defer workerMutex.Unlock()
	if validWorkerID(id) { // currentID? or len(workers), this is jank
		return workers[id], nil
	}
	w := gostock.NewWorker()
	w.Valid = false
	return *w, fmt.Errorf("can find worker: %d", id)
}

// RepoCreateWorker takes a worker and assigns it the next ID, then adds to workers slice
func RepoCreateWorker(w gostock.Worker) gostock.Worker {
	w.ID = currentWorkerID
	//log.Print("RepoCreateWorker", w)
	workerMutex.Lock()
	defer workerMutex.Unlock()
	workers = append(workers, w)
	currentWorkerID++
	return w
}

// RepoUpdateWorker updates a worker that matches input worker.ID, only updating updateable fields
func RepoUpdateWorker(worker gostock.Worker) (gostock.Worker, error) {
	// check sanity first
	if validWorkerID(worker.ID) {
		workerMutex.Lock()
		defer workerMutex.Unlock()

		workers[worker.ID].Ready = worker.Ready
		workers[worker.ID].Working = worker.Working
		//workers[i].IPAddr = worker.IPAddr
		//workers[i].Port = worker.Port
		workers[worker.ID].LastUpdate = time.Now()

		// if this update was to notify a worker was ready, then add to the queue
		if workers[worker.ID].Ready == true {
			readyWorkers <- worker.ID
		}

		// now return the updated info
		return workers[worker.ID], nil
	}
	worker.Valid = false
	return worker, fmt.Errorf("worker ID not found")
}

func validWorkerID(id int) bool {
	if id >= 0 && len(workers) != 0 && id <= currentWorkerID { // currentID? or len(workers), this is jank
		return true
	}
	return false
}
