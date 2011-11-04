package main

import (
	"testing"
	"os"
	"log"
	"time"
)

func TestMapFill(t *testing.T) {
	file := "testdata/fill2.map" // fill.2 Point{r:4, c:5}
	m, err := MapLoadFile(file)

	if err != os.EOF {
		t.Errorf("Read failed for %s: %v", file, err)
	} else if m == nil {
		t.Errorf("Invalid load of map m == nil")
	}

	// log.Printf("%v", m) // TODO test String() func round trip.
	l := make(map[Location]int, 0)
	for _, hill := range m.HillLocations(-1) {
		l[hill] = 1
	}

	fs, mQ, mD := MapFill(m, l, 1)
	log.Printf("SlowFill: mQ: %v mD: %v f::\n%v\n", mQ, mD, fs)

}

func BenchmarkMapFill(b *testing.B) {
	file := "testdata/maps/mmaze_05p_01.map"
	m, err := MapLoadFile(file)
	if m == nil || err != os.EOF {
		log.Panicf("Error reading %s: err %v map: %v", file, err, m)
	}


	l := make(map[Location]int, 40)

	for _, hill := range m.HillLocations(-1) {
		l[hill] = 1
	}

	for i := 0; i < b.N; i++ {
		MapFill(m, l, 1)
	}
}

var maps = []string{
	"maze_02p_01",
	"maze_02p_02",
	"maze_03p_01",
	"maze_04p_01",
	"maze_04p_02",
	"maze_05p_01",
	"maze_06p_01",
	"maze_07p_01",
	"maze_08p_01",
	"mmaze_02p_01",
	"mmaze_02p_02",
	"mmaze_03p_01",
	"mmaze_04p_01",
	"mmaze_04p_02",
	"mmaze_05p_01",
	"mmaze_07p_01",
	"mmaze_08p_01",
	"random_walk_02p_01",
	"random_walk_02p_02",
	"random_walk_03p_01",
	"random_walk_03p_02",
	"random_walk_04p_01",
	"random_walk_04p_02",
	"random_walk_05p_01",
	"random_walk_05p_02",
	"random_walk_06p_01",
	"random_walk_06p_02",
	"random_walk_07p_01",
	"random_walk_07p_02",
	"random_walk_08p_01",
	"random_walk_08p_02",
	"random_walk_09p_01",
	"random_walk_09p_02",
	"random_walk_10p_01",
	"random_walk_10p_02",
}

func TestMapFillDist(t *testing.T) {
	
	for _, name := range maps {
		filename := "testdata/maps/" + name + ".map"
		m, err := MapLoadFile(filename)
		if m == nil || err != os.EOF {
			log.Panicf("Error: failed to read %s: %v", filename, err)
		}
		l := make(map[Location]int)
		for _, hill := range m.HillLocations(0) {
			l[hill] = 1
			break
		}
		pre := time.Nanoseconds()
		f, mQ, mD := MapFill(m, l, 1)
		post := time.Nanoseconds()
		log.Printf("Fill: mQ: %2d mD: %3d %4.1f ms %s\n", mQ, mD, float64(post-pre)/1000000, name)
		hist := make([]int,mD+1)
		for _, d := range f.Depth {
			hist[d]++
		}
		for i, k := range hist {
			if i == 10 {
				log.Printf("\"%s\",%d,%d\n", name, i, k)
			}
		}
	}
}
