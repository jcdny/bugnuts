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

func (f *Fill) ToLocation(p Point) Location {
	p = f.Donut(p)
	return Location(p.r*f.Cols + p.c)
}

func (f *Fill) ToPoint(l Location) (p Point) {
	p = Point{r: int(l) / f.Cols, c: int(l) % f.Cols}

	return
}
func (f *Fill) PointAdd(p1, p2 Point) Point {
	return f.Donut(Point{r: p1.r + p2.r, c: p1.c + p2.c})
}

func (f *Fill) Donut(p Point) Point {
	if p.r < 0 {
		p.r += f.Rows
	}
	if p.r >= f.Rows {
		p.r -= f.Rows
	}
	if p.c < 0 {
		p.c += f.Cols
	}
	if p.c >= f.Cols {
		p.c -= f.Cols
	}

	return p
}

func (f *Fill) PathIn(loc Location) (Location, int) {
	steps := 0
	done := false
	for !done {
		depth := f.Depth[loc]
		p := f.ToPoint(loc)
		
		done = true
		for _, d := range Steps {
			np := f.PointAdd(p, d)
			nl := f.ToLocation(np)
			if (Debug > 4) {
				log.Printf("step from %v to %v depth %d to %d\n", p, np, depth, f.Depth[nl])
			}

			if f.Depth[nl] < depth && f.Depth[nl] > 0 {

				loc = nl
				steps++
				done = false
				break
			}
		}
	}

	return loc, steps
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

func MapFill(m *Map, origin map[Location]int) (*Fill, int, int) {
	Directions := []Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}} // w n e s

	safe := 0

	f := m.NewFill()

	q := QNew(100)

	for loc, pri := range origin {
		// log.Printf("Q loc %v pri %d", f.ToPoint(loc), pri)
		q.Q(f.ToPoint(loc))
		f.Depth[loc] = uint16(pri)
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

			if m.Grid[floc] != WATER && m.Grid[floc] != BLOCK &&
				(f.Depth[floc] == 0 || f.Depth[floc] > newDepth) {
				q.Q(fillp)
				f.Depth[floc] = newDepth
			}
		}
	}

	return f, 0, 0
}


type Target struct {
	item  Item
	loc   Location // Target Location
	count int      // how many do we want at this location
	pri   int      // target priority.

	arrivals []int      // Inbound Arrival time
	player   []int      // Inbound player
	ant      []Location // Inbound Source
}

type TargetSet map[Location]*Target

func (tset *TargetSet) Add(item Item, loc Location, count, pri int) {
	t, ok := (*tset)[loc]
	if !ok || t.pri < pri {
		// We already have this point in the target set, replace if pri is higher
		(*tset)[loc] = &Target{
			item:  item,
		        loc:   loc,
			count: count,
		pri: pri,
		}
	}

}

func (tset *TargetSet) Pending() int {
	n := 0
	for _, t := range *tset {
		if t.count > 0 {
			n++
		}
	}

	return n
}

func (tset *TargetSet) Active() map[Location]int {
	tp := make(map[Location]int, tset.Pending())
	for _, t := range *tset {
		if t.count > 0 {
			tp[t.loc] = t.pri
		}
	}

	return tp
}

