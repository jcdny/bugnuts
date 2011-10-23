package main

import (
	"log"
	"bufio"
	"os"
	"fmt"
)

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

	//log.Printf("State:\n%v\n", &s)

	me := NewBot(&s)
	fmt.Fprintf(os.Stdout, "go\n")

	for {
		// Reset for Next Parse

		line, err := s.ParseTurn() // TURN PARSE
		if err == os.EOF || line == "end" {
			break
		}

		//log.Printf("Generating orders turn %d", s.Turn)
		s.DoTurn() // generate orders

		// additional thinking til timeout

		// emit orders
	}

	//s.DumpSeen()
	//s.DumpMap()

	log.Printf("Bot Result %v", me)

}
