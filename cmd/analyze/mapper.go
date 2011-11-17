package main

import (
	"os"
	"io/ioutil"
	"sync"
	"log"
	"time"
)

type Mapper struct {
	mu     sync.Mutex
	reduce map[string]chan *Replay
}

type Replay struct {
	Name   string
	NBytes int64
}

type Results struct {
	N      int
	NBytes int64
}

type Reducer struct {
	out Results
}

var (
	ST        = int64(1000000)
	MaxFiles  = 100
	MaxReplay = 10
	MaxReduce = 10
)

func NewDispatch() (chan *os.FileInfo, chan *Replay) {
	files := make(chan *os.FileInfo, MaxFiles)
	replays := make(chan *Replay, MaxReplay)
	return files, replays
}

func Walk(files []string, fprocess chan *os.FileInfo) {
	fin := make(chan string, MaxFiles)

	go walk(fin, fprocess)

	log.Printf("%#v", files)
	for _, file := range files {
		log.Printf("Sending %s to walk", file)
		fin <- file
	}
	time.Sleep(2 * ST)
	log.Printf("Done Walk")
}

func walk(fin chan string, fprocess chan *os.FileInfo) {
	for file := range fin {
		finfo, _ := os.Lstat(file)
		if finfo.IsRegular() {
			log.Printf("walk: File %s", finfo.Name)
			fprocess <- finfo
		} else if finfo.IsDirectory() {
			flist, _ := ioutil.ReadDir(file)
			log.Printf("walk: Directory %s - descending, N %d", finfo.Name, len(flist))
			for _, finfo := range flist {
				fin <- file + "/" + finfo.Name
			}
		} else {
			log.Printf("Huh? %#v", finfo)
		}
	}
	time.Sleep(ST)
	log.Printf("Done walk")
}

func Stage(fprocess chan *os.FileInfo, replays chan *Replay) {
	for {
		finfo := <-fprocess
		replays <- &Replay{NBytes: finfo.Size, Name: finfo.Name}
	}
	time.Sleep(ST)
	log.Printf("Done Stage")
}

func mapper(replays <-chan *Replay) {
	mapper := &Mapper{reduce: make(map[string]chan *Replay)}

	for r := range replays {
		log.Printf("Looking up key for %s", r.Name)
		key := r.Name[0:1]
		m, ok := mapper.reduce[key]
		if !ok {
			m = make(chan *Replay, MaxReduce)
			go Reduce(key, m)
			mapper.reduce[key] = m
		}
		m <- r
		log.Printf("Sending %s to %s", r.Name, key)
		mapper.reduce[key] <- r
	}
	time.Sleep(ST)
	log.Printf("Done mapper")
}

func Reduce(key string, in <-chan *Replay) {
	out := Results{}
	for r := range in {
		log.Printf("Reducing %s on key %s", r.Name, key)
		out.N++
		out.NBytes += r.NBytes
	}
	time.Sleep(ST)
	log.Printf("Done reducer %s", key)
}
