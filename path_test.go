package main

import (
	"testing"
	"os"
	"log"
	"bufio"
	"strconv"
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

	Directions := []Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}} // w n e s
	Diagonals := []Point{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}}

	f := m.NewFill()

	q := QNew(100) // TODO think more about q cap

	q.Q(origin)
	f.Depth[m.ToLocation(origin)] = newDepth

	for !q.Empty() {
		p := q.DQ()

		Depth := f.Depth[m.ToLocation(p)]
		newDepth := Depth + 1

		log.Printf("DQ'd %v step %d to %d Map:\n%v", p, Depth, newDepth, f)

		validlast := false
		d, ed := 0, 4

		// loop over directions
		for d, ed := 0, 4; d < ed; d++ {
			if safe++; safe > 1000 { log.Panicf("Oh No Crazytime") }
			dir = Directions[d%4]

			fillp := m.PointAdd(p, dir)
			floc := m.ToLocation(fillp)

			log.Printf("%s", PrettyFill(m, f, p, fillp, q, Depth))

			if m.Grid[nloc] != WATER && f.Depth[nloc] == 0 {
				if !validlast {
					q.Q(fillp)
				}

				validlast == true
				diagdone := false

				// We have at least 1 valid square.  Try and move diagonally
			D:
				for _, diag := range Diag {
					dp := p
					dfillp := fp

					for {
						dp = m.PointAdd(dp, diag)
						dloc = m.ToLocation(dp)
						dfillp = m.PointAdd(dfillp, diag)
						dfloc = m.ToLocation(dfillp)

						if m.Grid[dfloc] != WATER && f.Depth[dloc] == Depth {
							// Set the previous point and step the floc forward
							f.Depth[floc] == newDepth
							floc = dfloc
							diagdone = true
						} else if diagdone {
							// TODO set d, ed, p, dir, validlast to continue
							// need to be careful to make sure properly
							// oriented.
							d--
							ed = d + 4
							p = dp

							// since we did a diagonal quit looking for others.
							break D
						}
					}
				}

				// Fallen out of diagonal code.  We always emerge here
				// with p set to the new location and d and ed set to
				// apply normal rules to the new loc pointer and set the
				// current position
				f.Depth[floc] == newDepth
			} else {
				validlast == false
			}

	}

	return f, 0, int(newDepth - 1)

}

func 

// Program to dump the fill and q state in a pretty format.
// @ or # is current pos, . is unvisited, % is water
// A is a point in the queue
func PrettyFill(m *Map, f *Fill, p, fillp Point, q *Queue, Depth uint16) string {
	s := ""
	for i, d := range f.Depth {
		curp := Point{r: i / f.Cols, c: i % f.Cols}

		if curp.c == 0 {
			switch curp.r {
			case 1:
				s += "  Depth: " + strconv.Itoa(int(Depth))
			case 2:
				s += "  QSize: " + strconv.Itoa(q.Size())
			}
			s += "\n"
		}

		qpos := q.Position(curp)

		if m.PointEqual(p, curp) {
			if qpos < 0 {
				s += "@" // point
			} else {
				s += "#" // point with point already in q
			}
		} else if m.PointEqual(fillp, curp) {
				s += "*"
		} else if qpos < 0 {
			if d == 0 {
				if m.Grid[i] == WATER {
					s += "%"
				} else {
					s += "."
				}
			} else {
				s += string('0' + byte((d-1)%10))
			}
		} else {
			s += string('A' + qpos%26)
		}
	}

	return s
}
