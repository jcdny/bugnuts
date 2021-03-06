// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package pathing

import (
	"testing"
	"os"
	"log"
	"time"
	"fmt"
	"strconv"
	"reflect"
	"rand"
	"bugnuts/replay"
	. "bugnuts/maps"
	. "bugnuts/torus"
)

func TestMapFill(t *testing.T) {
	file := "testdata/fill2.map" // fill.2 Point{r:4, c:5}
	m, err := MapLoadFile(file)

	if err != nil {
		t.Errorf("Read failed for %s: %v", file, err)
	} else if m == nil {
		t.Errorf("Invalid load of map m == nil")
	}

	// log.Printf("%v", m) // TODO test String() func round trip.
	l := make(map[Location]int, 0)
	for _, hill := range m.Hills(-1) {
		l[hill] = 1
	}

	sfs, mQ, mD := MapFillSeed(m, l, 1)
	log.Printf("SeedFill: mQ: %v mD: %v f::\n%v\n", mQ, mD, sfs)
	fs, mQ, mD := MapFill(m, l, 1)
	log.Printf("Fill: mQ: %v mD: %v f::\n%v\n", mQ, mD, fs)
	ff, mQ, mD := MapFillSlow(m, l, 1)
	log.Printf("SlowFill: mQ: %v mD: %v f::\n%v\n", mQ, mD, ff)

}

type resultNN struct {
	s1, s2 Location
	steps  int
	L      [4]Location
}

func TestMapFillSeedNN(t *testing.T) {
	// file := "../maps/testdata/maps/cell_maze_p06_01.map" 
	// file := "testdata/fill2.map" // fill.2 Point{r:4, c:5}
	file := "testdata/fillNN.map" // fill.2 Point{r:4, c:5}

	expect := []resultNN{
		{52, 136, 6, [4]Location{52, 100, 88, 136}},
		{37, 52, 5, [4]Location{37, 62, 62, 52}},
		{37, 136, 12, [4]Location{37, 142, 143, 136}},
		{28, 52, 1, [4]Location{28, 40, 40, 52}},
		{28, 136, 2, [4]Location{28, 4, 16, 136}},
		{28, 29, 0, [4]Location{28, 29, 28, 29}},
		{29, 52, 2, [4]Location{29, 40, 41, 52}},
		{29, 136, 3, [4]Location{29, 5, 5, 136}},
		{29, 37, 8, [4]Location{29, 34, 33, 37}},
	}
	m, err := MapLoadFile(file)

	if m == nil || err != nil {
		t.Errorf("Read failed for %s: %v", file, err)
	}

	l := make(map[Location]int, 0)
	for _, hill := range m.Hills(-1) {
		l[hill] = 1
	}
	f := NewFill(m)
	sfs, nn := f.MapFillSeedNN(l, 1, 0)
	if false {
		log.Printf("MapFillSeedNN:\n%v\nNN:\n", sfs)
	}

	nnexp := make(Neighbors, 0)
	for _, e := range expect {
		if _, ok := nnexp[e.s1]; !ok {
			nnexp[e.s1] = make(map[Location]Neighbor, 0)
		}
		if _, ok := nnexp[e.s2]; !ok {
			nnexp[e.s2] = make(map[Location]Neighbor, 0)
		}
		nnexp[e.s1][e.s2] = Neighbor{L: e.L, Steps: e.steps}
		nnexp[e.s2][e.s1] = Neighbor{L: [4]Location{e.L[3], e.L[2], e.L[1], e.L[0]}, Steps: e.steps}
	}
	if !reflect.DeepEqual(nn, nnexp) {
		t.Errorf("Mismatch\ngot: %v, expected %v", nn, nnexp)
	}

	/* 
		i := 0
		for s1, m1 := range nn {
				for s2, N := range m1 {
					if len(expect) < i ||
						expect[i].s1 != s1 ||
						expect[i].s2 != s2 ||
						expect[i].steps != N.Steps ||
						!reflect.DeepEqual(N.L, expect[i].L) {
						//log.Printf("%v %v (%d): %v", m.ToPoint(s1), m.ToPoint(s2), N.Steps, m.ToPoints(N.L[:]))
						//log.Printf("%d,%v,%v,%d,%#v", i, s1, s2, N.Steps, N.L[:], expect[i])
						t.Errorf("Mismatch\ngot:%v,%v,%d,%#v, expected %#v", s1, s2, N.Steps, N.L[:], expect[i])
					}
					i++
				}
			}
	*/
}

// Benchmark the version which does not maintain a seed array
// but allocates per fill
func getBenchMap() (*Map, map[Location]int) {
	return getBenchReplay()
}
func getBenchReplay() (*Map, map[Location]int) {
	// ai challenge game 148978 turn 100 236 ants, 300 506 ants
	// see 161429 for mongo # of ants...
	file := "testdata/replay.big.json.gz"
	match, err := replay.Load(file)
	if err != nil {
		log.Panicf("Load of %s failed: %v", file, err)
	}
	m := match.GetMap()
	al, _ := match.AntLocations(m, 300, 300)
	l := make(map[Location]int, len(al[0][0])*len(al[0]))
	for p := range al[0] {
		for _, loc := range al[0][p] {
			l[loc] = 1
		}
	}

	return m, l
}
func getBenchFile() (*Map, map[Location]int) {
	file := MapFile("cell_maze_p06_01")
	m, err := MapLoadFile(file)
	if m == nil || err != nil {
		log.Panicf("Error reading %s: err %v map: %v", file, err, m)
	}
	l := make(map[Location]int, 40)

	for _, hill := range m.Hills(-1) {
		l[hill] = 1
	}
	return m, l
}

func BenchmarkMapFillAlloc(b *testing.B) {
	m, l := getBenchMap()
	for i := 0; i < b.N; i++ {
		MapFill(m, l, 1)
	}
}

// Benchmark resetting the fill struct in a loop not with a copy
func BenchmarkResetSlow(b *testing.B) {
	m, _ := getBenchMap()

	f := NewFill(m)
	for i := 0; i < b.N; i++ {
		f.slowReset()
	}
}
// Benchmark resetting the fill struct
func BenchmarkReset(b *testing.B) {
	m, _ := getBenchMap()

	f := NewFill(m)
	for i := 0; i < b.N; i++ {
		f.Reset()
	}
}

func BenchmarkMapFill(b *testing.B) {
	m, l := getBenchMap()

	f := NewFill(m)
	for i := 0; i < b.N; i++ {
		f.Reset()
		f.MapFill(l, 1)
	}
}

// Benchmark not reusing the fill struct.
func BenchmarkMapFillSeedNN(b *testing.B) {
	m, l := getBenchMap()

	for i := 0; i < b.N; i++ {
		f := NewFill(m)
		f.MapFillSeedNN(l, 1, 0)
	}
}
// Benchmark not reusing the fill struct.
func BenchmarkMapFillSeedNNMD16(b *testing.B) {
	m, l := getBenchMap()

	for i := 0; i < b.N; i++ {
		f := NewFill(m)
		f.MapFillSeedNN(l, 1, 16)
	}
}

func BenchmarkMapFillSeedNNMD8(b *testing.B) {
	m, l := getBenchMap()

	for i := 0; i < b.N; i++ {
		f := NewFill(m)
		f.MapFillSeedNN(l, 1, 8)
	}
}

func BenchmarkMapFillSeedNNMD4(b *testing.B) {
	m, l := getBenchMap()

	for i := 0; i < b.N; i++ {
		f := NewFill(m)
		f.MapFillSeedNN(l, 1, 4)
	}
}
// Benchmark allocating fill + computing seed.
func BenchmarkMapFillSeed(b *testing.B) {
	m, l := getBenchMap()

	for i := 0; i < b.N; i++ {
		MapFillSeed(m, l, 1)
	}
}

func TestMapFillDist(t *testing.T) {
	out, _ := os.Create("tmp/dist.csv")
	defer out.Close()

	for _, name := range AllMaps {
		filename := MapFile(name)
		m, err := MapLoadFile(filename)
		if m == nil || err != nil {
			log.Panicf("Error: failed to read %s: %v", filename, err)
		}
		for _, player := range []int{-1, 0} {
			l := make(map[Location]int)
			for _, hill := range m.Hills(player) {
				l[hill] = 1
			}
			pre := time.Nanoseconds()
			f, mQ, mD := MapFillSlow(m, l, 1)
			post := time.Nanoseconds()
			ff, mQ, mD := MapFill(m, l, 1)
			postff := time.Nanoseconds()
			ffs, mQ, mD := MapFillSeed(m, l, 1)
			postffs := time.Nanoseconds()
			diff := 0
			for i, f := range f.Depth {
				if f != ff.Depth[i] || f != ffs.Depth[i] || ffs.Depth[i] != ff.Depth[i] {
					diff++
				}
			}
			log.Printf("Fill: mQ:%3d mD: %3d %4.1f/%4.1f/%4.1f ms %d diffs player %d points %d %s", mQ, mD, float64(post-pre)/1000000, float64(postff-post)/1000000, float64(postffs-postff)/1000000, diff, player, len(l), name)

			// Generate histograms.
			empty := NewMap(m.Rows, m.Cols, 1)
			fe, _, mDe := MapFill(empty, l, 1)
			if mD > mDe {
				mDe = mD
			}

			histe := make([]int, mDe+1)
			hist := make([]int, mDe+1)
			for i, d := range f.Depth {
				hist[d]++
				histe[fe.Depth[i]]++
			}
			if player == 0 {
				for i, k := range hist {
					fmt.Fprintf(out, "\"%s\",%d,%d,%d\n", name, i, k, histe[i])
				}
			}
		}
	}
}

// Take a map and generate montecarlo ant densities...
func TestMonteCarloPathing(t *testing.T) {
	for _, name := range AllMaps {
		filename := MapFile(name)
		m, err := MapLoadFile(filename)
		if m == nil || err != nil {
			log.Panicf("Error: failed to read %s: %v", filename, err)
		}

		lend := make(map[Location]int)
		for _, hill := range m.Hills(0) {
			lend[hill] = 1
		}

		f, _, _ := MapFillSeed(m, lend, 1)

		/* lsrc := make([]Location,0,len(m.Grid))
		for loc, item := range m.Grid {
			if item != WATER {
				lsrc = append(lsrc, Location(loc))
			}
		}*/

		lsrc := m.Hills(1)

		d := 10000 / len(lsrc)
		pre := time.Nanoseconds()
		rng := rand.New(rand.NewSource(1))
		paths, _ := f.MontePathIn(rng, lsrc, d, 1)
		post := time.Nanoseconds()
		log.Printf("Montecarlo %d paths in %.2f ms", d*len(lsrc), float64(post-pre)/1000000)

		str := ""
		// Write out in the annoying R image layout
		for c := 0; c < m.Cols; c++ {
			for r := m.Rows - 1; r >= 0; r-- {
				loc := Location(r*m.Cols + c)
				if r != m.Rows-1 {
					str += ","
				}
				str += strconv.Itoa(int(paths[loc]))
			}
			str += "\n"
		}
		out, _ := os.Create("tmp/" + name + ".csv")
		fmt.Fprint(out, str)
		out.Close()
	}
}
