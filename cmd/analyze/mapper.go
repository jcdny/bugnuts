package main

import (
	"os"
	"io/ioutil"
	"sync"
	"log"
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
	fin := make(chan *os.FileInfo, MaxFiles)

	go walk(fin, fprocess)

	log.Printf("%#v", files)
	for _, file := range files {
		log.Printf("Sending %s to walk", file)
		finfo, _ := os.Lstat(file)
		fin <- finfo
	}
}

func walk(fin chan *os.FileInfo, fprocess chan *os.FileInfo) {
	for {
		finfo := <-fin
		if finfo.IsRegular() {
			log.Printf("Walk: File %s", finfo.Name)
			fprocess <- finfo
		} else if finfo.IsDirectory() {
			log.Printf("Walk: Directory %s - descending", finfo.Name)
			flist, _ := ioutil.ReadDir(finfo.Name)
			for _, finfo := range flist {
				fin <- finfo
			}
		}
		log.Printf("Huh? %#v", finfo)
	}
}

func Stage(fprocess chan *os.FileInfo, replays chan *Replay) {
	for {
		finfo := <-fprocess
		replays <- &Replay{NBytes: finfo.Size, Name: finfo.Name}
	}
}

func mapper(replays chan *Replay) {
	mapper := &Mapper{reduce: make(map[string]chan *Replay)}

	for {
		r := <-replays
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
}

func Reduce(key string, in chan *Replay) {
	out := Results{}
	for {
		r := <-in
		log.Printf("Reducing %s on key %s", r.Name, key)
		out.N++
		out.NBytes += r.NBytes
	}
}
