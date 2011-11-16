package main

import (
	"flag"
	"log"
	"fmt"
	"os"
)

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

var (
	verbose = flag.Bool("v", false, "Verbose")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s file [file ...]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		usage()
	}

	files, replays := NewDispatch()

	go Stage(files, replays)
	go mapper(replays)
	Walk(args, files)
}
