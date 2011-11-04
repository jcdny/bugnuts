package main

import (
	"log"
	"strconv"
	"sort"
	"rand"
)

type Fill struct {
	// add offset and wrap flag for subfill work
	Rows    int
	Cols    int
	Depth   []uint16
	SeedLoc []Location
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
	origloc := loc
	done := false
	for !done {
		depth := f.Depth[loc]
		p := f.ToPoint(loc)
		done = true
		for _, d := range Steps {
			np := f.PointAdd(p, d)
			nl := f.ToLocation(np)

			if f.Depth[nl] < depth && f.Depth[nl] > 0 {
				loc = nl
				steps++
				done = false
				break
			}
		}
	}

	if Debug > 4 {
		log.Printf("step from %v to %v depth %d to %d, steps %d\n", f.ToPoint(origloc), f.ToPoint(loc), f.Depth[origloc], f.Depth[loc], steps)
	}

	return loc, steps
}

func (f *Fill) MontePathIn(loc Location, N int, MinDepth uint16) (dist map[Location]int) {
	steps := 0
	origloc := loc
	done := false
	for n := 0; n < N; n++ {
		loc = origloc
		for !done {
			depth := f.Depth[loc]
			p := f.ToPoint(loc)
			done = true
			for _, d := range rand.Perm(4) {
				np := f.PointAdd(p, Steps[d])
				nl := f.ToLocation(np)

				if f.Depth[nl] < depth && f.Depth[nl] > MinDepth {
					loc = nl
					steps++
					done = false
					break
				}
			}
		}
		if Debug > 4 {
			log.Printf("step from %v to %v depth %d to %d, steps %d\n", f.ToPoint(origloc), f.ToPoint(loc), f.Depth[origloc], f.Depth[loc], steps)
		}

		dist[loc]++
	}

	log.Printf("Pathin map: %#v", dist)

	return dist
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

// Generate a BFS Fill.  if pri is > 0 then use it for the point pri otherwise
// use map value
func MapFill(m *Map, origin map[Location]int, pri uint16) (*Fill, int, int) {
	Directions := []Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}} // w n e s

	safe := 0

	f := m.NewFill()

	q := QNew(100)

	for loc, opri := range origin {
		// log.Printf("Q loc %v pri %d", f.ToPoint(loc), pri)
		q.Q(f.ToPoint(loc))
		if pri > 0 {
			f.Depth[loc] = uint16(pri)
		} else {
			f.Depth[loc] = uint16(opri)
		}
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

// Build list of locations ordered by depth from closest to furthest
// TODO see if perm on the per depth list helps
func (f *Fill) Closest(slice []Location) []Location {
	llist := make(map[int][]Location) // List of locations keyed by depth
	dlist := make([]int, 0, 128)      // List of depths encountered

	if len(slice) < 1 {
		return slice
	}
	log.Printf("Closest slice %v", slice)

	for _, loc := range slice {
		depth := int(f.Depth[loc])
		if _, ok := llist[depth]; !ok {
			llist[depth] = make([]Location, 0)
			dlist = append(dlist, depth)
		}
		llist[depth] = append(llist[depth], loc)
	}

	sort.Sort(IntSlice(dlist))

	n := 0
	for _, depth := range dlist {
		copy(slice[n:n+len(llist[depth])], llist[depth])
		n += len(llist[depth])
	}

	if n != len(slice) {
		log.Panicf("Output length does not match input length (%d, %d)", n, len(slice))
	}

	return slice
}

// Return N random points sampled from a fill with steps between low and hi inclusive.
// it will return a count > 1 if the sample size is smaller than N
func (f *Fill) Sample(n, low, high int) ([]Location, []int) {
	pool := make([]Location, 0, 200)
	lo, hi := uint16(low), uint16(high)
	for i, depth := range f.Depth {
		if depth >= lo && depth <= hi {
			pool = append(pool, Location(i))
		}
	}
	if len(pool) == 0 {
		return nil, nil
	}

	over := n / len(pool)
	perm := rand.Perm(len(pool))[0 : n%len(pool)]
	log.Printf("Looking for %d explore points %d-%d, have %d possible", n, low, hi, len(pool))

	var count []int
	if over > 0 {
		count = make([]int, len(pool))
		for i, _ := range count {
			count[i] = over
		}
	} else {
		count = make([]int, len(perm))
	}

	for i, _ := range perm {
		count[i]++
	}

	if over > 0 {
		return pool, count
	} else {
		pout := make([]Location, len(perm))
		for i, pi := range perm {
			log.Printf("adding location %d to output pool", pool[pi])
			pout[i] = pool[pi]
		}
		return pout, count
	}

	return nil, nil
}

// Build list of locations ordered by depth from closest to furthest
// TODO see if perm on the per depth list helps
type Segment struct {
	src   Location
	end   Location
	steps int
}
type SegSlice []Segment

func (p SegSlice) Len() int { return len(p) }
// Metric is Manhattan distance from origin.
func (p SegSlice) Less(i, j int) bool { return p[i].steps < p[j].steps }
func (p SegSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (f *Fill) ClosestStep(seg []Segment) {
	if len(seg) < 1 {
		return
	}

	for i, _ := range seg {
		end, nsteps := f.PathIn(seg[i].src)
		seg[i].steps += nsteps
		seg[i].end = end
	}
	sort.Sort(SegSlice(seg))
}
