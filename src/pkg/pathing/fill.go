// Pathing computes distance maps, possibly with seed locations and nearest neighbors.
package pathing

import (
	"log"
	"strconv"
	. "bugnuts/maps"
	. "bugnuts/torus"
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
		s += "\nSeed:\n"
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
	return f.MapFill(origin, pri)
}

func MapFillSeed(m *Map, origin map[Location]int, pri uint16) (*Fill, int, uint16) {
	f := NewFill(m)
	return f.MapFillSeed(origin, pri)
}

func (f *Fill) slowReset() {
	// this takes 2.5x as long as the copy version below
	for i := range f.Depth {
		f.Depth[i] = 0
	}
}

var _zero [MAXMAPSIZE]uint16

func (f *Fill) Reset() {
	copy(f.Depth, _zero[:len(f.Depth)])
}

// Generate a BFS Fill.  if pri is > 0 then use it for the point pri otherwise
// use origin map value
func (f *Fill) MapFill(origin map[Location]int, pri uint16) (*Fill, int, uint16) {
	q := make([]Location, 0, 200+len(origin)*2)
	m := f.Map

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
// use origin map value
func (f *Fill) MapFillSeed(origin map[Location]int, pri uint16) (*Fill, int, uint16) {
	return f.MapFillSeedMD(origin, pri, 0)
}

// Generate a BFS Fill.  if pri is > 0 then use it for the point pri otherwise
// use origin map value
func (f *Fill) MapFillSeedMD(origin map[Location]int, pri, maxDepth uint16) (*Fill, int, uint16) {
	f.Seed = make([]Location, len(f.Depth))
	m := f.Map
	q := make([]Location, 0, 200+len(origin)*2)
	if maxDepth == 0 {
		maxDepth = 65535
	}

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
		if newDepth > maxDepth {
			break
		}

		for i := 0; i < 4; i++ {
			floc := m.LocStep[loc][i]
			// TODO: block is local.  Ugly but making it traversible in some # of turns might help 
			if (f.Depth[floc] == 0 || f.Depth[floc] > newDepth) &&
				m.Grid[floc] != WATER && m.Grid[floc] != BLOCK {
				q = append(q, floc)
				f.Depth[floc] = newDepth
				f.Seed[floc] = Seed
			}
		}
	}

	return f, cap(q), newDepth
}

type Neighbor struct {
	L     [4]Location
	Steps int
}

type Neighbors map[Location]map[Location]Neighbor

func NewNeighbors(origin []Location) Neighbors {
	nn := make(Neighbors, len(origin))
	for _, loc := range origin {
		nn[loc] = make(map[Location]Neighbor, 0)
	}
	return nn
}

func (nn Neighbors) Add(s1, s2, l1, l2 Location, steps int) {
	nn[s1][s2] = Neighbor{L: [4]Location{s1, l1, l2, s2}, Steps: steps}
}

func (nn Neighbors) doubleNN() {
	for s1 := range nn {
		for s2, n := range nn[s1] {
			nn[s2][s1] = Neighbor{L: [4]Location{n.L[3], n.L[2], n.L[1], n.L[0]}, Steps: n.Steps}
		}
	}
}

// MapFillSeedNN computes a flood fill BFS.  If pri is > 0 then use it for the point pri otherwise
// use map value.  Returns the fill together with the q size and max depth.
func (f *Fill) MapFillSeedNN(origin map[Location]int, pri, maxDepth uint16) (*Fill, Neighbors) {
	f.Seed = make([]Location, len(f.Depth))
	m := f.Map

	q := make([]Location, 0, 100+len(origin)*4)
	if maxDepth == 0 {
		maxDepth = 65535
	}
	for loc, opri := range origin {
		q = append(q, loc)
		f.Seed[loc] = loc
		if pri > 0 {
			f.Depth[loc] = pri
		} else {
			f.Depth[loc] = uint16(opri)
		}
	}

	nn := NewNeighbors(q)

	newDepth := uint16(0)
	for len(q) > 0 {
		loc := q[0]
		q = q[1:len(q)]

		Depth := f.Depth[loc]
		Seed := f.Seed[loc]
		newDepth = Depth + 1
		if newDepth > maxDepth {
			break
		}

		for i := 0; i < 4; i++ {
			floc := m.LocStep[loc][i]
			// Order of tests here matters for performance.
			if f.Depth[floc] > 0 &&
				f.Seed[loc] != f.Seed[floc] &&
				(f.Depth[floc] == newDepth || f.Depth[floc] == Depth) {
				// We are at a point where we are half way between two Seeds.
				// Check if this is a new minima and if so update the NN map.
				var Seed1, Seed2, L1, L2 Location
				Seed1 = f.Seed[floc]
				if Seed1 < Seed {
					Seed2 = Seed
					L1 = loc
					L2 = floc
				} else {
					Seed1, Seed2 = Seed, Seed1
					L1 = floc
					L2 = loc
				}

				nSteps := 2 * int(Depth-1)
				if f.Depth[floc] == newDepth {
					nSteps++
					L1, L2 = floc, floc
				}

				if N, ok := nn[Seed1][Seed2]; !ok || N.Steps > nSteps {
					// either we have not seen this before or we have a shorter distance.
					// TODO I should check if it's even possible to have a shorter distance...
					nn.Add(Seed1, Seed2, L1, L2, nSteps)
				}
			} else if (f.Depth[floc] == 0 || f.Depth[floc] > newDepth) &&
				m.Grid[floc] != WATER && m.Grid[floc] != BLOCK {
				q = append(q, floc)
				f.Depth[floc] = newDepth
				f.Seed[floc] = Seed
			}
		}
	}

	// add max/min pair.
	nn.doubleNN()

	return f, nn
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
