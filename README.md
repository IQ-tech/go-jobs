# Bounded job execution with in memory queue

Run jobs with controlled concurrency.

# Installation

To install, use `go get`

```terminal
go get github.com/IQ-tech/go-jobs
```

# Usage

## Creating a new dispatcher

Creating a dispatcher that will run 2 concurrent jobs and queue 5 jobs.

```go
dispatcher := jobs.NewDispatcher(2, 5)
```

## Running jobs

`Run` enqueues a job to that will be picked up by one of the workers.  
`Run` will block if the queue is full.

```go
dispatcher := jobs.NewDispatcher(2, 1)

// Does not block because a worker picks the job from
// the queue immediately after we add the job to it.
//
// We have 2 workers and 1 is busy now.
dispatcher.Run(func() {
  time.Sleep(10 * time.Second)
})


// Does not block because a worker picks the job from
// the queue immediately after we add the job to it.
//
// We have 2 workers and both are busy now.
dispatcher.Run(func() {
  time.Sleep(10 * time.Second)
})

// Does not block because even though both workers are busy,
// we have a queue that can have at most one job in it.
dispatcher.Run(func() {
  time.Sleep(10 * time.Second)
})

// Blocks because both workers are busy and the queue is full.
dispatcher.Run(func() {
  time.Sleep(10 * time.Second)
})
```

## Testing

Asynchronicity makes testing more annoying, because of that we have the `Sync`  
method to make testing easier.

Calling `Sync` will change the dispatcher to sync mode(default is async)  
which makes the dispatcher execute jobs synchronously instead of enqueuing them.

```go
package package_test

func Test_Sync(t*testing.T){
  t.Parallel()
	// We dont need workers or a queue if we are running in sync mode
	dispatcher := NewDispatcher(0, 0)

	dispatcher.Sync()

	ran := false
	dispatcher.Run(func() {
		ran = true
	})

  assert.True(t, ran)
}
```
