package main

import (
	"context"
	"fmt"
)

// Status represents the status of a job
// This is how enums are declared in go
// (or with iota keyword)
type Status string

const (
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
	elapsed   float64
	status    Status
	ctx       *context.Context
}

func (j Job) String() string {
	// if the status was not SUCCESS,
	// hide the elapsed and bytes Read
	if j.status != SUCCESS {
		return fmt.Sprintf("    %v ", j.status)
	}
	return fmt.Sprintf("%v %v %v", j.elapsed, j.bytesRead, j.status)
}

// JobOutput represents the output of a job
type JobOutput struct {
	id        int64
	err       error
	elapsed   float64
	bytesRead int64
}

// A Queue represents a channel (concurrent array) of Jobs.
type Queue (chan *Job)
