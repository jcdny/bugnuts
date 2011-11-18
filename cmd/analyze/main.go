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

	results := Process(args)

	for _, r := range results {
		log.Print(r)
	}

}
