package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"
)

type Application struct {
	tasks         map[string]time.Time
	inputChannel  chan TaskRequest
	outputChannel chan TaskResult
	wg            sync.WaitGroup
}

func main() {
	var numWorkers int
	var maxNumOfTasks int
	flag.IntVar(&numWorkers, "n", 1, "Number of workers in the pool")
	flag.IntVar(&maxNumOfTasks, "m", 1, "Maximum number of tasks queued")
	flag.Parse()

	if maxNumOfTasks < numWorkers {
		maxNumOfTasks = numWorkers
	}

	app := &Application{
		tasks:         make(map[string]time.Time),
		inputChannel:  make(chan TaskRequest, maxNumOfTasks),
		outputChannel: make(chan TaskResult, maxNumOfTasks),
	}

	app.startWorkers(numWorkers)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", app.router())
	close(app.inputChannel)
	app.wg.Wait()
	close(app.outputChannel)
	log.Fatal(err)
}

func (app *Application) startWorkers(numWorkers int) {
  log.Printf("Starting %v workers", numWorkers)
	for i := 0; i < numWorkers; i++ {
		app.wg.Add(1)
		go app.taskWorker(app.inputChannel, app.outputChannel, &app.wg)
	}

	go app.taskResponder(app.outputChannel)
}
