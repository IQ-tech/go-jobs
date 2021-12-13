package jobs

// Dispatcher represents struct for the Dispatcher
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	maximumWorkers int
	workerPool     chan chan Job
	jobQueue       chan Job
}

// NewDispatcher instanciates a new dispatcher
func NewDispatcher(maxWorkers int, maxQueue int) *Dispatcher {
	queue := make(chan Job, maxQueue)
	pool := make(chan chan Job, maxWorkers)

	dispatcher := &Dispatcher{workerPool: pool, maximumWorkers: maxWorkers, jobQueue: queue}

	dispatcher.start()

	return dispatcher
}

// start workers
func (d *Dispatcher) start() {
	// starting n number of workers
	for i := 0; i < d.maximumWorkers; i++ {
		worker := NewWorker(d.workerPool)
		worker.Start()
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
	work := Job{
		Execute: execFunc,
	}
	d.jobQueue <- work
}
