// Pathing computes distance maps, possibly with seed locations and nearest neighbors.
package pathing

import (
	"log"
	"sort"
	"rand"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/debug"
	. "bugnuts/util"
)

// Compute a path in to a point and return location and steps to minima.
func (f *Fill) PathIn(loc Location) (Location, int) {
	return f.NPathIn(loc, -1)
}

// NPathIn computes an N step path to a minima.  If steps == -1 then
// go to minima and return steps taken; if steps == 0 it's a noop, more for
// clean logic elsewhere.
func (f *Fill) NPathIn(loc Location, steps int) (Location, int) {
	if steps == 0 {
		return loc, steps
	} else if steps < -1 {
		steps = -1
	}

	origloc := loc

OUT:
	for {
		depth := f.Depth[loc]
		for _, d := range Permute4() {
			nl := f.LocStep[loc][d]
			if f.Depth[nl] < depth && f.Depth[nl] > 0 {
				loc = nl
				steps--
				if steps == 0 {
					break OUT
				} else {
					break
				}
			}
		}
	}

	if Debug[DBG_PathIn] && WS.Watched(loc, -1, 0) {
		log.Printf("step from %v to %v depth %d to %d, steps %d\n", f.ToPoint(origloc), f.ToPoint(loc), f.Depth[origloc], f.Depth[loc], steps)
	}

	return loc, -(steps + 1)
}

// MontePathIn computes montecarlo distribution and flow for pathing
// in to the set minimum depth, N samples per start location.
func (f *Fill) MontePathIn(m *Map, start []Location, N int, MinDepth uint16) (dist []int, flow [][4]int) {
	dist = make([]int, len(f.Depth))
	flow = make([][4]int, len(f.Depth))

	for _, origloc := range start {
		for n := 0; n < N; n++ {
			loc := origloc
			d := 0

			for d < 4 {
				depth := f.Depth[loc]
				nperm := rand.Intn(24)
				for d = 0; d < 4; d++ {
					nloc := m.LocStep[loc][Perm4[nperm][d]]
					if f.Depth[nloc] < depth && f.Depth[nloc] > MinDepth {
						flow[loc][Perm4[nperm][d]]++
						loc = nloc
						dist[loc]++
						break
					}
				}
			}
		}
	}
	return
}

// Closest builds a list of locations ordered by depth from closest to furthest
// TODO see if perm on the per depth list helps.
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