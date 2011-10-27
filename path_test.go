package main

import (
	"testing"
	"os"
	"log"
	"bufio"
)

func MapLoadFile(file string) (*Map, os.Error) {
	var m *Map = nil

	f, err := os.Open(file)

	if err != nil {
		return nil, err
	} else {
		defer f.Close()

		in := bufio.NewReader(f)
		m, err = MapLoad(in)
	}

	return m, err
}

func TestMapLoad(t *testing.T) {

	file := "testdata/maps/fill.1"
	m, err := MapLoadFile(file)

	if err != os.EOF {
		t.Errorf("Invalid load of map error os.Error == %v", err)
	}
	if m == nil {
		t.Errorf("Invalid load of map m == nil")
	}
}

func TestMapFill(t *testing.T) {
	file := "testdata/maps/maze_04p_01.map" // fill.2 Point{r:4, c:5}
	//file := "testdata/maps/fill.2"
	m, err := MapLoadFile(file)

	if err != os.EOF {
		t.Errorf("Read failed for %s: %v", file, err)
	} else if m == nil {
		t.Errorf("Invalid load of map m == nil")
	}

	// log.Printf("%v", m) // TODO test String() func round trip.

	fs, mQ, mD := SlowMapFill(m, []Point{{r: 0, c: 0}})
	log.Printf("SlowFill: mQ: %v mD: %v f::\n%v\n", mQ, mD, fs)

	// find a hill for start
	p := []Point{}
	for i, item := range m.Grid {
		if item.IsHill() {
			p = append(p, m.ToPoint(Location(i)))
		}
	}

	fs, mQ, mD = SlowMapFill(m, p)

	log.Printf("SlowFill: mQ: %v mD: %v f::\n%v\n", mQ, mD, fs)

}

func BenchmarkSlowMapFill(b *testing.B) {

	file := "testdata/maps/maze_04p_01.map"
	m, _ := MapLoadFile(file)

	//f, _, _ := SlowMapFill(m, []Point{{r: 3, c: 3}})
	//log.Printf("%v", f)

	// TODO find a hill for start
	for i := 0; i < b.N; i++ {
		SlowMapFill(m, []Point{{r: 3, c: 3}})
	}

}
