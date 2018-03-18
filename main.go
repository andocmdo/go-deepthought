package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func main() {
	// first get command args
	workerNode := flag.Bool("worker", false, "run as worker node")
	numWorkers := flag.Int("workers", 1, "max number of workers to spawn")
	startPort := flag.Int("startPort", 12345, "starting port to accept jobs")
	ipPort := flag.String("master", "127.0.0.1:8080",
		"IP address and port of the API server (master node)")
	api := flag.String("api", "/api/v1/", "api root")
	flag.Parse()

	if *workerNode {
		// number of processor cores on system
		coresAvailable := runtime.NumCPU()
		log.Println("Number of processor cores available: " +
			strconv.FormatInt(int64(coresAvailable), 10))
		log.Println("Number of workers: " +
			strconv.FormatInt(int64(*numWorkers), 10))
		log.Println("Starting port (for listening TCP port to accept jobs): " +
			strconv.FormatInt(int64(*startPort), 10))

		// init the server struct to hold master server info
		master := Server{URLroot: "http://" + *ipPort, URLjobs: "http://" + *ipPort +
			*api + "jobs", URLworkers: "http://" + *ipPort + *api + "workers"}

		// TODO also when the threads have started, we will wait as well if we lose connection?
		for {
			resp, err := http.Get(master.URLroot)
			//defer resp.Body.Close()
			if err == nil {
				resp.Body.Close()
				break
			}
			log.Print("Error connecting to master server. Is it running?")
			log.Print("Error: ", err.Error())
			log.Print("Retry connection to master in 30 secs")
			time.Sleep(time.Second * 30)

		}

		for i := 0; i < *numWorkers; i++ {
			worker := &Worker{Port: strconv.Itoa(*startPort + i)}
			go worker.run(i, master)
			time.Sleep(time.Millisecond * 100) // TODO remove this after testing
		}

	} else {
		// if we are not a worker node we are the master
		// TODO add a thread to check the status of workers (ping them essentially) an
		// set their "ready" status accordingly.

		router := NewRouter()
		log.Fatal(http.ListenAndServe(":8080", router))

	}

	select {} // run forever

}
