// Pathing computes distance maps, possibly with seed locations and nearest neighbors.
package pathing

import (
	"log"
	"rand"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/debug"
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
