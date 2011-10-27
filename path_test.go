package main

import (
	"testing"
	"os"
	"log"
	"image"
)

func TestMapLoad(t *testing.T) {

	file := "testdata/maps/big"
	m, err := MapLoadFile(file)

	if err != os.EOF {
		t.Errorf("Invalid load of map error os.Error == %v", err)
	}
	if m == nil {
		t.Errorf("Invalid load of map m == nil")
	}
	m.WriteDebugImage("_test", 0, func (c, r int) image.NRGBAColor { return m.At(r, c) })
}

func TestMapFill(t *testing.T) {
	file := "testdata/maps/maze_04p_01.map" // fill.2 Point{r:4, c:5}
	// file := "testdata/maps/fill.2"
	m, err := MapLoadFile(file)

	if err != os.EOF {
		t.Errorf("Read failed for %s: %v", file, err)
	} else if m == nil {
		t.Errorf("Invalid load of map m == nil")
	}

	// log.Printf("%v", m) // TODO test String() func round trip.
	p := []Point{}
	for _, hill := range m.HillLocations() {
		p = append(p, m.ToPoint(hill))
	}

	fs, mQ, mD := SlowMapFill(m, p[0:1])
	log.Printf("SlowFill: mQ: %v mD: %v f::\n%v\n", mQ, mD, fs)

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
