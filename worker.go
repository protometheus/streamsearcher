package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"time"
)

type WorkerOutput struct {
	id        int64
	elapsed   float64
	bytesRead int
}

// Workers perform the actual searching for the StreamSearcher.
// A single Worker will iterate over the jobs Queue, pull one off,
// and attempt to read that job's chunk of the file.
func Worker(ss *StreamSearcher) {
	for job := range ss.jobs {
		func() {
			defer ss.CompleteJob()

			ctx, cancelFunc := context.WithTimeout(
				context.Background(),
				time.Duration(5000000)*time.Nanosecond,
			)

			output := make(chan WorkerOutput, 1)
			go func(output chan WorkerOutput) {
				defer cancelFunc()
				start := time.Now()

				buffer := make([]byte, ss.chunkSize)
				_, err := ss.file.ReadAt(buffer, job.startByte)
				if err != nil && err != io.EOF {
					log.Println(err)
					return
				}

				b := bytes.Index(buffer, []byte("Leapfn"))
				elapsed := time.Since(start).Seconds()
				output <- WorkerOutput{
					id:        job.id,
					bytesRead: b,
					elapsed:   elapsed,
				}

				close(output)
			}(output)

			select {
			case workerOutput := <-output:
				// go func() { ss.outputs <- workerOutput }()
				if workerOutput.bytesRead < 0 {
					log.Printf("Job %v did not find the needle", workerOutput.id)
				} else {
					log.Printf("Found needle in %v bytes read by job %v in %v secs \n",
						workerOutput.bytesRead,
						workerOutput.id,
						workerOutput.elapsed,
					)
				}
			case <-ctx.Done():
				log.Printf("Job %v timed out\n", job.id)
			}

		}()
	}
}
