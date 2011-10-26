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

		log.Printf("%v", m) // TODO test String() func round trip.
	}
}

func TestMapFill(t *testing.T) {
	var m *Map = nil

	// fill.2 Point{r:4, c:5}
	f, err := os.Open("testdata/maps/fill.3")

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
		f, mQ, mD := MapFill(m, Point{r: 3, c: 3})

		log.Printf("mQ: %v mD: %v f::\n%v\n", mQ, mD, f)
	}
}

type Fill struct {
	// add offset and wrap flag for subfill work
	Rows  int
	Cols  int
	Depth []uint16
}

func (m *Map) NewFill() *Fill {
	f := &Fill{
		Depth: make([]uint16, m.Size(), m.Size()),
		Rows:  m.Rows,
		Cols:  m.Cols,
	}

	return f
}

func (f *Fill) String() string {
	s := ""
	for i, d := range f.Depth {
		if i%f.Cols == 0 {
			s += "\n"
		}
		if d == 0 {
			s += "."
		} else {
			s += string('a' + byte((d-1)%26))
		}
	}

	return s
}

// Generate a fill from Map m return fill slice, max Q size, max depth
func MapFill(m *Map, origin Point) (*Fill, int, int) {

	newDepth := uint16(1) // dont start with 0 since 0 means not visited.
	safe := 0

	// CW search for next step
	// need an extra rotate to handle gap at end...
	//

	cw := []Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	diag := []Point{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}, {-1, 1}}

	f := m.NewFill()

	q := QNew(100) // TODO think more about q cap

	q.Q(origin)
	f.Depth[m.ToLocation(origin)] = newDepth

	for !q.Empty() {
		p := q.DQ()
		log.Printf("DQ %v", p)

		Depth := f.Depth[m.ToLocation(p)]
		newDepth := Depth + 1

		log.Printf("Start from %v step %d to %d Map:\n%v", p, Depth, newDepth, f)

		validlast := false

		for i, s := range diag {
			// on a given diagonal just go until we 
			// stop finding same depth
			for {
				// Debug lets not infinite loop
				if safe++; safe > 1000 {
					log.Panicf("Oh No Crazytime")
				}

				fillp := m.PointAdd(p, cw[i])
				nloc := m.ToLocation(fillp)

				// log.Printf("p: %v np: %v i: %d item: %c d: %d", p, fillp, i, m.Grid[nloc].ToSymbol(), f.Depth[nloc])


				if m.Grid[nloc] != WATER && f.Depth[nloc] == 0 {
					f.Depth[nloc] = newDepth
					// Queue a new start point
					if !validlast {
						q.Q(fillp)
						log.Printf("Q %v", fillp)
					}
					validlast = true
				} else {
					validlast = false
				}

				np := m.PointAdd(p, s)
				nloc = m.ToLocation(np)
				log.Printf("p %v np %v fillp %v %v", p, np, fillp, f)
				if f.Depth[nloc] == Depth {
					p = np
				} else {
					break
				}
			}
		}
	}

	return f, 0, int(newDepth - 1)

}

// ...#.....
// ...#.3...
// ...#323..
// ...#2123.
// ...#323..
// .....3...
// #########

// ...#.....
// ...#.....
// ...#.2...
// ...#212..
// ...#.2...
// .........
// #########
