package main

import (
	"testing"
	"os"
	"log"
	"bufio"
)

func TestMapLoad(t *testing.T) {
	var m *Map = nil

	f, err := os.Open("testdata/maps/fill.1")

	if err != nil {
		t.Errorf("Open failed: %v", err)
	} else {
		defer f.Close()

		in := bufio.NewReader(f)

		m, err = MapLoad(in)

		if err != os.EOF {
			t.Errorf("Invalid load of map error == %v", err)
		}

		if m == nil {
			t.Errorf("Invalid load of map m == nil")
		}

		// log.Printf("%v", m)

		// TODO test String() func round trip.
		// TODO test error handling make err return
	}
}

func TestMapFill(t *testing.T) {
	var m *Map = nil

	// fill.2 Point{r:4, c:5}
	f, err := os.Open("testdata/maps/fill.2")

	if err != nil {
		t.Errorf("Open failed: %v", err)
	} else {
		defer f.Close()

		in := bufio.NewReader(f)

		m, err = MapLoad(in)

		if err != os.EOF {
			t.Errorf("Invalid load of map error == %v", err)
		}

		if m == nil {
			t.Errorf("Invalid load of map m == nil")
		}

		log.Printf("%v", m) // TODO test String() func round trip.

		// find a hill for start
		fs, mQ, mD := SlowMapFill(m, Point{r: 3, c: 3})

		log.Printf("SlowFill: mQ: %v mD: %v f::\n%v\n", mQ, mD, fs)
		ff, mQ, mD := MapFill(m, Point{r: 3, c: 3})
		log.Printf("FastFill: mQ: %v mD: %v f::\n%v\n", mQ, mD, ff)
	}
}

func BenchmarkSlowMapFill(b *testing.B) {
	var m *Map = nil

	// fill.2 Point{r:4, c:5}
	f, _ := os.Open("testdata/maps/fill.2")

	defer f.Close()
	in := bufio.NewReader(f)
	m, _ = MapLoad(in)

	// find a hill for start
	for i := 0; i < b.N; i++ {
		SlowMapFill(m, Point{r: 3, c: 3})
	}

}

func BenchmarkFastMapFill(b *testing.B) {
	var m *Map = nil

	// fill.2 Point{r:4, c:5}
	f, _ := os.Open("testdata/maps/fill.2")

	defer f.Close()
	in := bufio.NewReader(f)
	m, _ = MapLoad(in)

	// find a hill for start
	for i := 0; i < b.N; i++ {
		MapFill(m, Point{r: 3, c: 3})
	}

}
