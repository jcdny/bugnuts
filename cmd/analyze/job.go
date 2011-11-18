package main

import (
	"os"
	"io"
	"io/ioutil"
	"json"
	"bugnuts/replay"
	"log"
	"bytes"
	"compress/gzip"
	"fmt"
)

type JobData struct {
	Name  string
	Match *replay.Match
}

type Result struct {
	Key     string
	N       int
	Games   []*replay.GameResult
	Players []*replay.PlayerResult
}

// Take a job and load the required data
func Load(job *Job) (*JobData, os.Error) {
	m := replay.Match{}

	buf, err := ioutil.ReadFile(job.File)
	if err != nil {
		return nil, err
	}

	if buf[0] == '\x1f' && buf[1] == '\x8b' {
		// decode gzip file
		unzip, err := gzip.NewReader(bytes.NewBuffer(buf[:]))
		if err != nil {
			log.Panic("Unzip: %v", err)
		}
		bout := bytes.NewBuffer(make([]byte, 0, 8*len(buf)))
		_, err = io.Copy(bout, unzip)
		if err != nil {
			log.Panic("Unzip: %v", err)
		}
		buf = bout.Bytes()
	}

	err = json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}
	return &JobData{Match: &m}, nil
}

func stager(in chan *Job, out chan<- *Job, ring chan chan *Job, done chan<- int) {
	var err os.Error
	for job := range in {
		// Populate the job data
		job.Data, err = Load(job)
		job.Key = fmt.Sprintf("%d", job.Data.Match.GameId%10)

		// Pass job off to next step
		if err == nil {
			out <- job
		} else {
			// TODO feed job to a failed channel to record somewhere.
			log.Printf("Fail: %s %v", job.File, err)
		}

		// Release this proc back into the ring buffer
		ring <- in
	}
	done <- 1
}

// Take staged Jobs in and accumulate results until the channel is
// closed
func reducer(key string, in <-chan *Job, out chan<- *Result) {
	res := Result{
		Key:     key,
		Games:   make([]*replay.GameResult, 0, 1000),
		Players: make([]*replay.PlayerResult, 0, 4000),
	}

	for job := range in {
		res.N++
		g, p := job.Data.Match.ExtractMetaData()
		res.Games = append(res.Games, g)
		res.Players = append(res.Players, p...)
	}

	out <- &res
}

func (r *Result) String() string {
	buf := make([]byte, 0, 100*len(r.Games))
	b := bytes.NewBuffer(buf)
	for _, g := range r.Games {
		fmt.Fprintf(b, "\"game\",%d,%s,%d,\"%s\",%s,%s,%s,%s,%s\n",
			g.GameId, g.Date, g.GameLength, g.Challenge, g.MatchupId, g.PostId, g.WorkerId, g.Location, g.MapId)
	}
	for _, p := range r.Players {
		fmt.Fprintf(b, "\"player\",\"%s\",%d\n", p.PlayerName, p.GameId)
	}

	return string(b.Bytes())
}
