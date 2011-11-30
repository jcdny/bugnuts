package main

import (
	"os"
	"io/ioutil"
	"log"
)

type Job struct {
	File  string
	Finfo *os.FileInfo
	Key   string
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

func Mapper(in <-chan *Job, out chan<- []*Result) {
	reducers := make(map[string]chan *Job)
	// The channel reducers use to pass back their results at
	// completion
	rout := make(chan *Result)

	for job := range in {
		// Just us the first letter of the filename as the key for the
		// reducers

		rin, ok := reducers[job.Key]
		if !ok {
			rin = make(chan *Job, MaxReduce)
			go reducer(job.Key, rin, rout)
			reducers[job.Key] = rin
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
