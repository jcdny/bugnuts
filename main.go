package main

import (
	"os"
	"log"
	"bufio"
)

//main initializes the state and starts the processing loop
func main() {
	var s State

	in := bufio.NewReader(os.Stdin)
	
	err := s.Start(in)
	if err != nil {
		log.Panicf("Start() failed (%s)", err)
	}

	mb := NewBot(&s)

	err = s.Loop(mb, func() {
		//if you want to do other between-turn debugging things, you can do them here
	})
	if err != nil && err != os.EOF {
		log.Panicf("Loop() failed (%s)", err)
	}
}
