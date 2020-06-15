package main

import (
	"flag"
	"log"
	"os"
)

const (
	defaultFilename          = "./_input.txt"
	defaultWorkers           = 10
	defaultSearchTerm        = "Lpfn"
	defaultTimeoutMillis     = 1000
	defaultGeneratedFileSize = 1000000000
)

func main() {
	// Set up flags for parsing along with their default values
	var filename string
	flag.StringVar(
		&filename,
		"filename",
		defaultFilename,
		`Filename of input stream to be read.
    If no filename is provided, the default input file is generated`,
	)

	var numWorkers int
	flag.IntVar(
		&numWorkers,
		"workers",
		defaultWorkers,
		"Number of workers to spawn for iterating over the stream.",
	)

	var searchTerm string
	flag.StringVar(
		&searchTerm,
		"term",
		defaultSearchTerm,
		`Term to be searched for in the input file.
    If streamsearcher is generating a file as input, this value will be included
    at pseudo-random intervals in the generated file.`,
	)

	var timeoutMillis int64
	flag.Int64Var(
		&timeoutMillis,
		"timeout",
		defaultTimeoutMillis,
		`Number of milliseconds to wait before timing out execution of a job's search.`,
	)

	var chunkSize int64
	flag.Int64Var(
		&chunkSize,
		"chunksize",
		-1,
		`Size of chunk (in bytes) for each job to search from the input.
    If not provided, defaults to size(input)/number_of_workers to equalize work
    across all workers.`,
	)

	var generatedFileSize int64
	flag.Int64Var(
		&generatedFileSize,
		"genfilesize",
		defaultGeneratedFileSize,
		`Size of the generated random input file (in bytes) for each job to search
    from the input.`,
	)

	flag.Parse()

	// generate the test input as needed
	if filename == "./_input.txt" {
		os.Remove(filename)
		GenerateInput(filename, searchTerm, generatedFileSize)
	}

	// build the StreamSearcher
	streamSearcher, err := NewStreamSearcher(
		filename,
		searchTerm,
		numWorkers,
		timeoutMillis,
		chunkSize)
	if err != nil {
		log.Fatal(err)
	}

	// Execute the search
	streamSearcher.Search()
}
