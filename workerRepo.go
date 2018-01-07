package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var currentWorkerID int
var workers Workers
var workerMutex *sync.Mutex

// Initialize mutex and fake DB
func init() {
	workerMutex = &sync.Mutex{}
	workerMutex.Lock()
	defer workerMutex.Unlock()
	currentWorkerID = 0
}

// RepoFindWorker searches for a worker with id inside mock DB
func RepoFindWorker(id int) (Worker, error) {
	workerMutex.Lock()
	defer workerMutex.Unlock()
	if id >= 0 && len(workers) != 0 && id <= currentWorkerID { // currentID? or len(workers), this is jank
		return workers[id], nil
	}
	w := NewWorker()
	w.Valid = false
	return *w, fmt.Errorf("can find worker: %d", id)
}

// RepoCreateWorker takes a worker and assigns it the next ID, then adds to workers slice
func RepoCreateWorker(w Worker) Worker {
	w.ID = currentWorkerID
	log.Print("ReporCreateWorker", w)
	workerMutex.Lock()
	defer workerMutex.Unlock()
	workers = append(workers, w)
	//queueWorker(w.ID) // TODO should we queue on creating worker? Or after creating, then worker says ready?
	currentWorkerID++
	return w
}

// RepoUpdateWorker updates a worker that matches input worker.ID, only updating updateable fields
func RepoUpdateWorker(worker Worker) error {
	// check sanity first
	if worker.ID < 0 || worker.Valid == false {
		return fmt.Errorf("worker is not valid or has illegal ID")
	}
	workerMutex.Lock()
	defer workerMutex.Unlock()
	for i, w := range workers {
		if w.ID == worker.ID {
			workers[i].Destroyed = worker.Destroyed
			workers[i].LastUpdate = time.Now()
			return nil
		}
	}
	return fmt.Errorf("worker ID not found")
}
