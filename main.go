package main

import (
	"flag"
	"log"
	"bufio"
	"os"
	"fmt"
)

var Debug int = 0

func init() {
	flag.IntVar(&Debug, "d", 0, "Debug level 0 none 1 game 2 per turn 3 per ant 4 excessive")
	flag.Parse()
}

func main() {

	in := bufio.NewReader(os.Stdin)

	Run(in)
}

func Run(in *bufio.Reader) {
	var s State

	err := s.Start(in)
	if err != nil {
		log.Panicf("Start(in) failed (%s)", err)
	}

	if Debug > 1 {
		log.Printf("State:\n%v\n", &s)
	}

	me := NewBot(&s)
	fmt.Fprintf(os.Stdout, "go\n")

	for {
		// Reset for Next Parse

		line, err := s.ParseTurn() // TURN PARSE

		// State Validation

		if err == os.EOF || line == "end" {
			break
		}

		if Debug > 1 {
			log.Printf("TURN %d Generating orders turn", s.Turn)
		}
		// Generate order list
		s.DoTurn()

		// Validation of orders

		// additional thinking til timeout

		// emit orders
	}

	// Read end of game data.

	// Do end of game diagnostics

	//s.DumpSeen()
	//s.DumpMap()

	if Debug > 0 {
		log.Printf("Bot Result %v", me)
	}
}
