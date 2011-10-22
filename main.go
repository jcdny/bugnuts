package main

import (
	"log"
	"bufio"
	"os"
)

func main() {
	var s State

	in := bufio.NewReader(os.Stdin)

	err := s.Start(in)
	if err != nil {
		log.Panicf("Start(in) failed (%s)", err)
	}

	log.Printf("State:\n%v\n", &s)

	me := NewBot(&s)

	for {
		// TURN PARSE
		line, err := s.ParseTurn()
		if err == os.EOF || line == "end" {
			break
		}

		//generate orders

		//emit orders

		//think while hanging out.
	}

	log.Printf("Bot Result %v", me)

}
