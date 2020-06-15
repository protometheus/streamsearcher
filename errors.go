package main

// A TimeoutError ...
type TimeoutError struct {
	Err error
}

// Error
func (te TimeoutError) Error() string {
	return te.Err.Error()
}

// An ExecutionError is
type ExecutionError struct {
	Err error
}

// Error
func (ee ExecutionError) Error() string {
	return ee.Err.Error()
}
