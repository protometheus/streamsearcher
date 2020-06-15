package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
)

// A StreamSearcher contains all the information necessary to being a search for
// the supplied search term in the supplied input file.
type StreamSearcher struct {
	jobWorkers    int
	jobs          Queue
	chunkSize     int64
	file          *os.File
	filename      string
	fileSize      int64
	searchTerm    string
	timeoutMillis int64
	completedJobs Queue
	wg            *sync.WaitGroup
	ctx           context.Context
}

// NewStreamSearcher creates a new StreamSearcher with the given filename.
func NewStreamSearcher(filename, searchTerm string, jobWorkers int, timeoutMillis int64, chunkSize int64) (*StreamSearcher, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// chunkSize is how much of the stream that a given job will read.
	fileSize := fs.Size()
	if chunkSize <= 0 {
		chunkSize = fileSize / int64(jobWorkers)
	}

	var wg sync.WaitGroup
	return &StreamSearcher{
		jobWorkers:    jobWorkers,
		jobs:          make(Queue, chunkSize),
		chunkSize:     chunkSize,
		file:          file,
		filename:      filename,
		fileSize:      fs.Size(),
		searchTerm:    searchTerm,
		timeoutMillis: timeoutMillis,
		completedJobs: make(Queue, chunkSize),
		wg:            &wg,
	}, nil
}

// Search begins the StreamSearcher's search for the searchTerm
// in the file.
func (ss *StreamSearcher) Search() {
	ss.ctx = context.Background()

	// For each job, we enqueue it.
	// Note: If we had an infinite stream, this would be an infinite loop
	// which would listen for non-EOF of input
	numJobs := int(ss.fileSize / ss.chunkSize)
	ss.wg.Add(numJobs)
	for i := 0; i < numJobs; i++ {
		// Create the Job.
		j := &Job{
			id:        int64(i),
			startByte: int64(i) * ss.chunkSize,
			endByte:   int64(i+1) * (ss.chunkSize - 1),
			status:    SUCCESS,
		}

		// enqueue the job in a goroutine so infinite iteration is possible
		go func() {
			ss.jobs <- j
		}()
	}

	// spawns the workers
	for i := 0; i < ss.jobWorkers; i++ {
		go Worker(ss)
	}

	ss.wg.Wait()

	// gather final statistics
	// These could be stored on the StreamSearcher itself, or
	// calculated as each job finished. To avoid race conditions,
	// the simplest approach is used.
	var totalBytes int64
	var totalElapsed float64
	for idx := 0; idx < numJobs; idx++ {
		select {
		case j := <-ss.completedJobs:
			if j.status == SUCCESS {
				totalBytes += j.bytesRead
				totalElapsed += j.elapsed
			}
		}
	}

	var rate float64
	if totalElapsed != 0 {
		rate = float64(totalBytes) / totalElapsed
	}
	log.Printf("%v bytes/s read", rate)
}

// CompleteJob tells the StreamSearcher that a job has completed or timed out.
func (ss *StreamSearcher) CompleteJob(j *Job) {
	go func() { ss.completedJobs <- j }()
	ss.wg.Done()
	fmt.Println(j.String())
}
