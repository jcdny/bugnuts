package main

import (
	"os"
	"io/ioutil"
	"log"
	"fmt"
)

type JobData struct {
	Name   string
	NBytes int64
}

type Result struct {
	Key    string
	N      int
	NBytes int64
}

type Job struct {
	File  string
	Finfo *os.FileInfo
	Data  *JobData
}

var (
	MaxJobs   = 1000 // Maximum # of files queued to process
	NStage    = 10   // # of simultaneous staging goroutines
	MaxMapper = 10   // Max # of jobs pending for the mapper
	MaxReduce = 1    // Max # of jobs pending per reducer
)

// Take the list of files, set up the channels to connect the
// stages of the processing, and then block until reading
// final results.
func Process(files []string) []*Result {
	stage := make(chan *Job, MaxJobs)
	mapper := make(chan *Job, MaxMapper)
	results := make(chan []*Result)

	go Mapper(mapper, results)
	go Stage(stage, mapper, NStage)
	Walk(files, stage)

	return <-results
}

// Given a slice of filenames, walk the tree and generate the Job
// entries which are passed off on the jobs channel
func Walk(files []string, out chan *Job) {
	for _, file := range files {
		walk(file, out)
	}
	close(out)
}

func walk(file string, out chan<- *Job) {
	finfo, _ := os.Lstat(file)
	if finfo.IsRegular() {
		out <- &Job{File: file, Finfo: finfo}
	} else if finfo.IsDirectory() {
		flist, _ := ioutil.ReadDir(file)
		for _, finfo := range flist {
			walk(file+"/"+finfo.Name, out)
		}
	} else {
		log.Printf("Skipped \"%s\" not a regular file or directory", file)
	}
}

// Take a Job from the jobs channel, Stage it, and hand it to the data
// reduction step.
func Stage(in <-chan *Job, out chan<- *Job, N int) {
	// ring buffer for stagers and done
	done := make(chan int)
	ring := make(chan chan *Job, N)

	// Create the stager goroutines
	for i := 0; i < NStage; i++ {
		in := make(chan *Job)
		go stager(in, out, ring, done)
		ring <- in
	}

	// Loop, handle jobs til closed.  stage() puts itself back in the
	// ring buffer when done.
	for job := range in {
		stager := <-ring
		stager <- job
	}

	// Close all the stagers
	for i := 0; i < N; i++ {
		close(<-ring)
	}

	// Wait for them to drain
	for i := 0; i < N; i++ {
		<-done
	}

	// no more jobs pending...
	close(out)
}

func stager(in chan *Job, out chan<- *Job, ring chan chan *Job, done chan<- int) {
	for job := range in {
		// Populate the job data
		job.Data = &JobData{NBytes: job.Finfo.Size, Name: job.Finfo.Name}
		// Pass job off to next step
		out <- job
		// Release this proc back into the ring buffer
		ring <- in
	}
	done <- 1
}

func Mapper(in <-chan *Job, out chan<- []*Result) {
	reducers := make(map[string]chan *Job)
	// The channel reducers use to pass back their results at
	// completion
	rout := make(chan *Result)

	for job := range in {
		// Just us the first letter of the filename as the key for the
		// reducers
		key := job.Finfo.Name[:1]

		rin, ok := reducers[key]
		if !ok {
			rin = make(chan *Job, MaxReduce)
			go reducer(key, rin, rout)
			reducers[key] = rin
		}
		rin <- job
	}

	for _, r := range reducers {
		close(r)
	}

	// We closed all the reducers, loop here until we have read all
	// their result data then return it via the out channel which was
	// passed in.
	results := make([]*Result, 0, len(reducers))
	for _, _ = range reducers {
		results = append(results, <-rout)
	}

	out <- results
	close(out)
}

// Take staged Jobs in and accumulate results until the channel is
// closed
func reducer(key string, in <-chan *Job, out chan<- *Result) {
	res := Result{Key: key}

	for r := range in {
		res.N++
		res.NBytes += r.Data.NBytes
	}

	out <- &res
}

func (r *Result) String() string {
	return fmt.Sprintf("Files starting with %s: %d files %d bytes",
		r.Key, r.N, r.NBytes)
}
