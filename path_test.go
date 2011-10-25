package main

import (
	"testing"
	"os"
	"log"
	"bufio"
)

func TestMapLoad(t *testing.T) {
	var m *Map = nil

	f, err := os.Open("testdata/test.map")

	if err != nil {
		t.Errorf("Open failed: %v", err)
	} else {
		in := bufio.NewReader(f)

		m, err = MapLoad(in)

		if err != os.EOF {
			t.Errorf("Invalid load of map error == %v", err)
		}

		if m == nil {
			t.Errorf("Invalid load of map m == nil")
		}

		// log.Printf("%v", m) // TODO test String() func round trip.
	}
}
