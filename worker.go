package main

import (
	"bytes"
	"context"
	"io"
	"time"
)

// A Worker performs the actual searching for the StreamSearcher.
// A single Worker will iterate over the jobs Queue, pull one off,
// and attempt to read that job's chunk of the file.
func Worker(ss *StreamSearcher) {
	for job := range ss.jobs {
		func() {
			// mark the job as complete at the end of its execution
			defer ss.CompleteJob(job)

			// get a context with a timeout
			ctx, cancelFunc := context.WithTimeout(
				context.Background(),
				time.Duration(ss.timeoutMillis)*time.Millisecond,
			)

			// make the ouput a channel so we can wait to see what arrives first:
			// - the job outputs a response
			// - the worker times out
			// Must be done in a goroutine to avoid deadlock.
			output := make(chan JobOutput, 1)
			go func(output chan JobOutput) {
				// Call the context's cancel function regardless of success in execution
				// Same goes for the closing of channels.
				defer cancelFunc()
				defer close(output)

				// start the timer
				start := time.Now()

				// Read into the buffer
				buffer := make([]byte, ss.chunkSize)
				_, err := ss.file.ReadAt(buffer, job.startByte)
				if err != nil && err != io.EOF {
					output <- JobOutput{
						id:        job.id,
						err:       err,
						bytesRead: 0,
						elapsed:   0,
					}
					return
				}

				// Get the index of the search term if it exists.
				// Note: The number of bytes read is not exactly equal to
				// the number of bytes until the index of the search term.
				// This is because the Index() may read the same bytes twice.
				// This is known; the Index() function is known to outperform
				// the naive iteration approach, especially for large strings.
				b := int64(bytes.Index(buffer, []byte(ss.searchTerm)))
				elapsed := time.Since(start).Seconds()
				output <- JobOutput{
					id:        job.id,
					err:       nil,
					bytesRead: b,
					elapsed:   elapsed,
				}

			}(output)

			// Based on whatever channel responds first,
			// set the job's status, elapsed, and bytesRead
			select {
			case <-ctx.Done():
				job.status = TIMEOUT
				job.elapsed = 0
				job.bytesRead = 0
				return
			case jobOutput := <-output:
				// This is where error handling would happen. The instructions said to
				// only print output, so nothing is done with the job's error
				if jobOutput.bytesRead < 0 || jobOutput.err != nil {
					job.status = FAILURE
					job.elapsed = 0
					job.bytesRead = 0
					return
				} else {
					// Successful Execution
					job.status = SUCCESS
					job.elapsed = jobOutput.elapsed
					job.bytesRead = jobOutput.bytesRead
					return
				}
			}
		}()
	}
}
