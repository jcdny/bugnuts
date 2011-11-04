package main

import (
	"testing"
	"os"
	"log"
)

func TestMapFill(t *testing.T) {
	file := "testdata/maps/fill.2" // fill.2 Point{r:4, c:5}
	// file := "testdata/maps/fill.2"
	m, err := MapLoadFile(file)

	if err != os.EOF {
		t.Errorf("Read failed for %s: %v", file, err)
	} else if m == nil {
		t.Errorf("Invalid load of map m == nil")
	}

	// log.Printf("%v", m) // TODO test String() func round trip.
	l := make(map[Location]int, 0)
	for _, hill := range m.HillLocations() {
		l[hill] = 1
	}

	fs, mQ, mD := MapFill(m, l, 1)
	log.Printf("SlowFill: mQ: %v mD: %v f::\n%v\n", mQ, mD, fs)

}

func BenchmarkSlowMapFill(b *testing.B) {

	file := "testdata/maps/maze_04p_01.map"
	m, _ := MapLoadFile(file)
	l := make(map[Location]int, 0)
	for _, hill := range m.HillLocations() {
		l[hill] = 1
	}

	//f, _, _ := SlowMapFill(m, []Point{{r: 3, c: 3}})
	//log.Printf("%v", f)

	// TODO find a hill for start
	for i := 0; i < b.N; i++ {
		MapFill(m, l, 1)
	}
}
