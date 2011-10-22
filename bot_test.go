package main

import (
	"testing"
	"os"
	"bufio"
//	"bytes"
)


func TestParse(t *testing.T) {
	var s State

	f, err := os.Open("testdata/stream1.dat")
	if err != nil {
                t.Errorf("Open failed: %v", err)
	} else {
		in := bufio.NewReader(f)
		err := s.Start(in)
		if err != nil {
			t.Errorf("Start() failed (%v)", err)
		}
		
		mb := NewBot(&s)
		
		err = s.Loop(mb, func() {
			//if you want to do other between-turn debugging things, you can do them here
		})
		if err != nil && err != os.EOF {
			t.Errorf("Loop() failed (%s)", err)
		}
	}
}
