package main

import (
	"context"
	"log"
	"os"
	"sync"
)

// Status represents the status of a job
// This is how enums are declared in go
// (or with iota keyword)
type Status string

var (
	// Successful execution
	SUCCESS Status = "SUCCESS"

	// Timeout occurred during execution
	TIMEOUT Status = "TIMEOUT"

	// Failure occurred during execution
	FAILURE Status = "FAILURE"
)

// A Job is
type Job struct {
	id        int64
	startByte int64
	endByte   int64
	bytesRead int64
	elapsed   int64
	status    Status
	ctx       *context.Context
}

// A Queue is
type Queue (chan *Job)

// SpawnWorkers takes a Reader stream and creates workers
func (ss *StreamSearcher) SpawnWorkers() error {
	for i := 0; i < ss.jobWorkers; i++ {
		go Worker(ss)
	}
	return nil
}

// CompleteJob tells the StreamSearcher that a job has completed or timed out.
func (ss *StreamSearcher) CompleteJob() {
	ss.wg.Done()
}

// StreamSearcher
type StreamSearcher struct {
	jobWorkers int
	jobs       Queue
	chunkSize  int
	file       *os.File
	wg         *sync.WaitGroup
	ctx        context.Context
	outputs    chan WorkerOutput
}

// NewStreamSearcher
func NewStreamSearcher() (*StreamSearcher, error) {
	file, err := os.Open("./_input.txt")
	if err != nil {
		return nil, err
	}

	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}

	jobWorkers := 10
	chunkSize := int(fs.Size()) / jobWorkers

	var wg sync.WaitGroup

	return &StreamSearcher{
		jobWorkers: jobWorkers,
		jobs:       make(Queue, chunkSize),
		chunkSize:  chunkSize,
		file:       file,
		wg:         &wg,
		outputs:    make(chan WorkerOutput, jobWorkers),
	}, nil
}

// Run
func (ss *StreamSearcher) Run() {
	ss.ctx = context.Background()

	// For each job, we enqueue it.
	// If we had an infinite stream, this would be an infinite loop
	// which would listen for valid input
	ss.wg.Add(10)
	for i := 0; i < 10; i++ {
		// create the Job
		j := &Job{
			id:        int64(i),
			startByte: int64(i * ss.chunkSize),
			endByte:   int64((i+1)*ss.chunkSize - 1),
			status:    SUCCESS,
		}

		// enqueue the job onto our Queue
		go func() {
			ss.jobs <- j
		}()
	}

	if err := ss.SpawnWorkers(); err != nil {
		log.Fatal(err)
	}

	ss.wg.Wait()
	// close(ss.outputs)
	//
	// var totalBytesRead float64
	// var totalElapsed float64
	// for {
	// 	select {
	// 	case output := <-ss.outputs:
	// 		if output.bytesRead < 0 {
	// 			fmt.Printf("Job %v did not find the needle", output.id)
	// 		} else {
	// 			totalBytesRead += float64(output.bytesRead)
	// 			totalElapsed += output.elapsed
	// 			fmt.Printf("Found needle in %v bytes read by job %v in %v secs \n",
	// 				output.bytesRead,
	// 				output.id,
	// 				output.elapsed,
	// 			)
	// 		}
	// 	default:
	// 		break
	// 	}
	// }

	// fmt.Printf("Total bytes read: %v\nBytes Read/sec: %v\n", totalBytesRead, (totalBytesRead / totalElapsed))

}
