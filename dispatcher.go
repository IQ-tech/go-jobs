package jobs

import (
	"runtime"
)

// Dispatcher represents struct for the Dispatcher
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	maximumWorkers int
	workers        []*Worker
	workerPool     chan chan Job
	jobQueue       chan Job
	syncMode       bool
}

// NewDispatcher instanciates a new dispatcher
func NewDispatcher(maxWorkers int, maxQueue int) *Dispatcher {
	queue := make(chan Job, maxQueue)
	pool := make(chan chan Job, maxWorkers)

	dispatcher := &Dispatcher{
		workerPool:     pool,
		jobQueue:       queue,
		maximumWorkers: maxWorkers,
		syncMode:       false,
	}

	dispatcher.start()

	runtime.SetFinalizer(dispatcher, stopWorkers)

	return dispatcher
}

func stopWorkers(d *Dispatcher) {
	for _, worker := range d.workers {
		if worker != nil {
			worker.Stop()
		}
	}

	close(d.workerPool)
	close(d.jobQueue)
}

// Dispatcher will pass jobs to workers when in async mode(it is the default mode).
func (d *Dispatcher) Async() {
	d.syncMode = false
}

// Dispatcher will execute jobs synchronously without putting them into the queue
// when in sync mode.
//
// It is useful for tests where we care about jobs being executed before
// making assertions.
func (d *Dispatcher) Sync() {
	d.syncMode = true
}

// start workers
func (d *Dispatcher) start() {
	// starting n number of workers
	for i := 0; i < d.maximumWorkers; i++ {
		worker := NewWorker(d.workerPool)
		worker.Start()
		d.workers = append(d.workers, &worker)
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for job := range d.jobQueue {
		// a job request has been received
		go func(job Job) {
			// try to obtain a worker job channel that is available.
			// this will block until a worker is idle
			jobChannel := <-d.workerPool

			// dispatch the job to the worker job channel
			jobChannel <- job
		}(job)
	}
}

// Run starts a new job
func (d *Dispatcher) Run(execFunc func()) {
	if d.syncMode {
		execFunc()
		return
	}

	work := Job{
		Execute: execFunc,
	}

	d.jobQueue <- work
}
