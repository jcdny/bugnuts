package main

import (
	"log"
	"strconv"
)

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
				if q != nil {
					s += "  QSize: " + strconv.Itoa(q.Size())
				}
			}
			s += "\n"
		}

		qpos := -1
		if q != nil {
			qpos = q.Position(curp)
		}

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
				s += string('0' + byte(d%10))
			}
		} else {
			s += string('A' + qpos%26)
		}
	}

	return s
}

func MapFill(m *Map, origin []Point) (*Fill, int, int) {
	Directions := []Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}} // w n e s
	newDepth := uint16(1)                                   // dont start with 0 since 0 means not visited.

	safe := 0

	f := m.NewFill()

	q := QNew(100)

	for _, p := range origin {
		q.Q(p)
		f.Depth[m.ToLocation(p)] = newDepth
	}

	for !q.Empty() {
		// just for sanity...
		if safe++; safe > 100*len(f.Depth) {
			log.Panicf("Oh No Crazytime %d %d", len(f.Depth), safe)
		}

		p := q.DQ()

		Depth := f.Depth[m.ToLocation(p)]
		newDepth := Depth + 1

		for _, d := range Directions {
			fillp := m.PointAdd(p, d)
			floc := m.ToLocation(fillp)

			if m.Grid[floc] != WATER && (f.Depth[floc] == 0 || f.Depth[floc] > newDepth) {
				q.Q(fillp)
				f.Depth[floc] = newDepth
			}
		}
	}

	return f, 0, 0
}

// Generate a fill from Map m return fill slice, max Q size, max depth
func ExpMapFill(m *Map, origin Point) (*Fill, int, int) {

	newDepth := uint16(1) // dont start with 0 since 0 means not visited.
	safe := 0

	Directions := []Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	Diagonals := []Point{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}, {-1, 1}}

	f := m.NewFill()

	q := QNew(100) // TODO think more about q cap

	q.Q(origin)
	f.Depth[m.ToLocation(origin)] = newDepth

	for !q.Empty() {
		p := q.DQ()

		Depth := f.Depth[m.ToLocation(p)]
		newDepth := Depth + 1

		validlast := false

		for i, s := range Diagonals {
			// on a given diagonal just go until we
			// stop finding same depth
			for {
				// Debug lets not infinite loop
				if safe++; safe > 1000 {
					log.Panicf("Oh No Crazytime")
				}

				fillp := m.PointAdd(p, Directions[i])
				nloc := m.ToLocation(fillp)

				// log.Printf("p: %v np: %v i: %d item: %c d: %d", p, fillp, i, m.Grid[nloc].ToSymbol(), f.Depth[nloc])
				//log.Printf("%s", PrettyFill(m, f, p, fillp, q, Depth))

				if m.Grid[nloc] != WATER && f.Depth[nloc] == 0 {
					f.Depth[nloc] = newDepth
					// Queue a new start point
					if !validlast {
						q.Q(fillp)
						//log.Printf("Q %v", fillp)
					}
					validlast = true
				} else {
					validlast = false
				}

				np := m.PointAdd(p, s)
				nloc = m.ToLocation(np)
				//log.Printf("p %v np %v fillp %v %v", p, np, fillp, f)

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
