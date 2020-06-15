# streamsearcher
Go lib that searches a stream for a phrase and outputs its findings. Uses a
Job/Worker queue implementation to utilize concurrency and optimize performance.

## Build
Run `sh ./build.sh` to build the executable.

Alternatively, run `go build -o streamsearcher *.go` and then run with

## Run
To run the StreamSearcher, build the executable then run

`./streamsearcher`

If no flags are provided, the StreamSearcher will generate a file with
pseudo-random data as well as needles to search for. It will then attempt to
find the needles by spinning off workers and jobs, which report their status
as they complete.

## Help
For help with flags, run `./streamsearcher --help` or `./streamsearcher -h`.
Flags include:
- `filename`: The name of the file to be used as input for searching. (`./_input.txt`)
- `term`: The term to search for. If input data is generated, this term is
          interpolated throughout the generated input. (`Leapfn`)
- `workers`: Number of workers to spawn (`10`)
- `timeout`: Number of seconds before jobs should timeout (`60`)
- `chunksize`: Number of bytes to be searched by each job.

## Assumptions
In order to complete this assignment, I made certain assumptions.
- The input is a file: since I do not know the format of the input and it must
  be able to accept an infinite input, the test input is on a generated file.
  However, this does not change the implementation; streams (Readers) in
  Go are modular, meaning the only change necessary from using input files to
  input streams is how the input is opened. All other logic works the same.

- The number of Jobs is equal to the number of Workers: The instructions say
  that only 11 lines should be printed, 1 for each worker. Because of this,
  each worker is only given 1 job. That is not a requirement, however; if the
  `chunksize` flag is provided, then the number of jobs will equal:
  `fileSize/chunkSize`. This breaks the requirement of the instructions,
  but it more practical and realistic in terms of general solutions.

- Obvious error handling: Go's error handling paradigm means
  most well-written functions are of the type `func(...) (..., error)`. Error
  handling is baked into how good Go code is written. Here, however, there
  was no explanation of what to do with errors other than to say the job failed.
  Under normal circumstances, errors would be Queued out the same way the job's
  output their data.

## Initial Instructions
Please complete a code exercise using your own workspace, IDE, references, etc. This should take approximately a few hours to complete, though feel free to take as much time as you need, as long as you're happy with the result. Your code will be evaluated based on (in this order): simplicity, readability, code style, use of best practices, efficiency, and thoughtful error handling and logging.

Coding Exercise:

Write a program in a language of your choice that spawns 10 workers (threads, processes, actors, whatever), where each worker simultaneously searches a stream of random (or pseudo-random) data for the string 'Lpfn', then informs the parent of the following data fields via some form of inter-process communication or shared data structure:
* elapsed time
* count of bytes read
* status

The parent collects the results of each worker (confined to a timeout, explained below) and writes a report to stdout for each worker sorted in descending order by [elapsed]:
[elapsed] [byte_cnt] [status]

Where [elapsed] is the elapsed time for that worker in ms, [byte_cnt] is the number of random bytes read to find the target string and [status] should be one of {SUCCESS, TIMEOUT, FAILURE}. FAILURE should be reported for any error/exception of the worker and the specific error messages should go to stderr. TIMEOUT is reported if that worker exceeds a given time limit, where the program should support a command-line option for the timeout value that defaults to 60s. If the status is not SUCCESS, the [elapsed] and [byte_cnt] fields will be empty.

The parent should always write a record for each worker and the total elapsed time of the program should not exceed the timeout limit. If a timeout occurs for at least one worker, only those workers that could not complete in time should report TIMEOUT, other workers may have completed in time and should report SUCCESS. Note that the source of random bytes must be a stream such that a worker will continue indefinitely until the target string is found, or a timeout or exception occurs. A final line of output will show the average bytes read per time unit in a time unit of your choice where failed/timeout workers will not report stats. 11 lines of output total to stdout.

Please package your submission with tar or zip. The package must include a README with these instructions, a UNIX shell executable file (or instructions on how to build one) that runs your program and responds appropriately to -h, which gives instructions on running the program (including the timeout option).
