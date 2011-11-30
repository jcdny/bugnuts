// Pathing computes distance maps, possibly with seed locations and nearest neighbors.
package pathing

import (
	"log"
	"strconv"
	"sort"
	"rand"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/debug"
	. "bugnuts/util"
)

// A Fill is an distance metric on a given Map.
type Fill struct {
	Depth []uint16   // Steps from a given Seed point
	Seed  []Location // The Seed point associated with a given point.
	*Map             // The map referenced for the construction of this fill 
}

// NewFill returns a pointer to a new Fill
func NewFill(m *Map) *Fill {
	f := &Fill{
		Depth: make([]uint16, m.Size(), m.Size()),
		Map:   m,
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

	if f.Seed != nil {
		s += "Seed:\n"
		for i, d := range f.Seed {
			if i%f.Cols == 0 {
				s += "\n"
			}
			if d == 0 {
				s += "."
			} else {
				s += string('a' + byte((d-1)%26))
			}
		}
	}

	return s
}

// Program to dump the fill and q state in a pretty format.  @ or # is
// current pos, . is unvisited, % is water A is a point in the queue
func PrettyFill(m *Map, f *Fill, p, fillp Point, q *Queue, Depth uint16) string {
	s := ""
	for i, d := range f.Depth {
		curp := Point{R: i / f.Cols, C: i % f.Cols}

		if curp.C == 0 {
			switch curp.R {
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

// mapFillSlow generates a flood fill of a map given a collection of origin points.
// If pri is > 0 then use it for the initial points depth otherwise use origin map value.
// Mostly still here as a sanity checker.
func MapFillSlow(m *Map, origin map[Location]int, pri uint16) (*Fill, int, uint16) {
	Directions := []Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}} // w n e s

	safe := 0

	f := NewFill(m)

	q := NewQueue(100)

	for loc, opri := range origin {
		// log.Printf("Q loc %v pri %d", f.ToPoint(loc), pri)
		q.Q(f.ToPoint(loc))
		if pri > 0 {
			f.Depth[loc] = uint16(pri)
		} else {
			f.Depth[loc] = uint16(opri)
		}
	}

	newDepth := uint16(0)
	for !q.Empty() {
		// just for sanity...
		if safe++; safe > 100*len(f.Depth) {
			log.Panicf("Oh No Crazytime %d %d", len(f.Depth), safe)
		}

		p := q.DQ()

		Depth := f.Depth[m.ToLocation(p)]
		newDepth = Depth + 1

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

	return f, q.Cap(), newDepth
}

// Generate a BFS Fill.  if pri is > 0 then use it for the point pri otherwise
// use map value
func MapFill(m *Map, origin map[Location]int, pri uint16) (*Fill, int, uint16) {
	f := NewFill(m)
	return f.MapFill(m, origin, pri)
}

func MapFillSeed(m *Map, origin map[Location]int, pri uint16) (*Fill, int, uint16) {
	f := NewFill(m)
	return f.MapFillSeed(m, origin, pri)
}

func (f *Fill) Reset() {
	for i := 0; i < len(f.Depth); i++ {
		f.Depth[i] = 0
	}
}

// Generate a BFS Fill.  if pri is > 0 then use it for the point pri otherwise
// use map value
func (f *Fill) MapFill(m *Map, origin map[Location]int, pri uint16) (*Fill, int, uint16) {
	if f.Rows != m.Rows || f.Cols != m.Cols {
		log.Panicf("Map and fill mismatch")
	}

	q := make([]Location, 0, 200+len(origin)*2)

	for loc, opri := range origin {
		// log.Printf("Q loc %v pri %d", f.ToPoint(loc), pri)
		q = append(q, loc)
		if pri > 0 {
			f.Depth[loc] = pri
		} else {
			f.Depth[loc] = uint16(opri)
		}
	}

	newDepth := uint16(0)
	for len(q) > 0 {
		loc := q[0]
		q = q[1:len(q)]

		Depth := f.Depth[loc]
		newDepth = Depth + 1

		for i := 0; i < 4; i++ {
			floc := m.LocStep[loc][i]
			if m.Grid[floc] != WATER && m.Grid[floc] != BLOCK &&
				(f.Depth[floc] == 0 || f.Depth[floc] > newDepth) {
				q = append(q, floc)
				f.Depth[floc] = newDepth
			}
		}
	}

	return f, cap(q), newDepth
}

// Generate a BFS Fill.  if pri is > 0 then use it for the point pri otherwise
// use map value
func (f *Fill) MapFillSeed(m *Map, origin map[Location]int, pri uint16) (*Fill, int, uint16) {
	if f.Rows != m.Rows || f.Cols != m.Cols {
		log.Panicf("Map and fill mismatch")
	}

	f.Seed = make([]Location, len(f.Depth))

	q := make([]Location, 0, 200+len(origin)*2)

	for loc, opri := range origin {
		// log.Printf("Q loc %v pri %d", f.ToPoint(loc), pri)
		q = append(q, loc)
		f.Seed[loc] = loc
		if pri > 0 {
			f.Depth[loc] = pri
		} else {
			f.Depth[loc] = uint16(opri)
		}
	}

	newDepth := uint16(0)
	for len(q) > 0 {
		loc := q[0]
		q = q[1:len(q)]

		Depth := f.Depth[loc]
		Seed := f.Seed[loc]
		newDepth = Depth + 1

		for i := 0; i < 4; i++ {
			floc := m.LocStep[loc][i]
			// TODO: block is local.  Ugly but making it traversible in some # of turns might help 
			if m.Grid[floc] != WATER && m.Grid[floc] != BLOCK &&
				(f.Depth[floc] == 0 || f.Depth[floc] > newDepth) {
				q = append(q, floc)
				f.Depth[floc] = newDepth
				f.Seed[floc] = Seed
			}
		}
	}

	return f, cap(q), newDepth
}

func (f *Fill) Distance(from, to Location) int {
	return int(f.Depth[from]) - int(f.Depth[to])
}

func (f *Fill) DistanceStep(loc Location, d Direction) int {
	if d == NoMovement {
		return 0
	}
	return int(f.Depth[loc]) - int(f.Depth[f.LocStep[loc][d]])
}

// Closest builds a list of locations ordered by depth from closest to furthest
// TODO see if perm on the per depth list helps.
// 
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

// Sample returns N random points sampled from a fill with step
// distance between low and hi inclusive.  it will return a count > 1
// if the sample size is smaller than N.  If n < 1 then return all
// points.
func (f *Fill) Sample(n, low, high int) ([]Location, []int) {
	pool := make([]Location, 0, 200)
	lo, hi := uint16(low), uint16(high)
	for i, depth := range f.Depth {
		if depth >= lo && depth <= hi {
			pool = append(pool, Location(i))
		}
	}
	if n < 1 {
		return pool, nil
	}

	if len(pool) == 0 {
		return nil, nil
	}

	over := n / len(pool)
	perm := rand.Perm(len(pool))[0 : n%len(pool)]
	if Debug[DBG_Sample] {
		log.Printf("Sample: Looking for %d explore points %d-%d, have %d possible", n, low, hi, len(pool))
	}

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
			if Debug[DBG_Sample] {
				log.Printf("Sample: adding location %d to output pool", pool[pi])
			}
			pout[i] = pool[pi]
		}
		return pout, count
	}

	return nil, nil
}

// Segment is a pathing segment which has a Src location, End location
// and distance in Steps.
type Segment struct {
	Src   Location
	End   Location
	Steps int
}

type SegSlice []Segment

func (p SegSlice) Len() int           { return len(p) }
func (p SegSlice) Less(i, j int) bool { return p[i].Steps < p[j].Steps }
func (p SegSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ClosestStep takes a slice of Segment with Src populated and 
// computes the End and Steps returning the slice order by steps from closest to furthest.
// Return true if any segments were found.
func (f *Fill) ClosestStep(seg []Segment) bool {
	if len(seg) < 1 {
		return false
	}

	any := false
	for i, _ := range seg {
		seg[i].End = f.Seed[seg[i].Src]
		seg[i].Steps += Abs(int(f.Depth[seg[i].Src]) - int(f.Depth[seg[i].End]))
		if seg[i].End != 0 || f.Depth[seg[i].End] != 0 {
			any = true
		}
	}

	sort.Sort(SegSlice(seg))

	return any
}
