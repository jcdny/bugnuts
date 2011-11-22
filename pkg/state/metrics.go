package state

import (
	"log"
	. "bugnuts/pathing"
	. "bugnuts/maps"
	. "bugnuts/util"
)

const NTHREAT = 10

type Metrics struct {
	*Map
	Seen     []int      // Turn on which cell was last visible.
	VisCount []int      // How many ants see this cell.
	Threat   [][]int8   // how much threat is there on a given cell
	PThreat  [][]uint16 // Prob of n threat
	Horizon  []bool     // Inside the event horizon.  false means there could be an ant there we have not seen
	HBorder  []Location // List of border points
	Land     []int      // Count of land tiles visible from a given tile
	PrFood   []int      // Count of Turn * Land adjusted for # that see it.
	Unknown  []int      // Count of Unknown
	VisSum   []int      // sum of count of visibles for overlap.
	// Fills
	FDownhill *Fill // Downhill from my own hills
	FHill     *Fill // Distance to my hills
	FViewHill *Fill // Distance to visibility boundary for my hills
	FAll      *Fill // All hill fill
	// MC distributions
	MCDist    []int
	MCFlow    [][4]int
	MCDistMax int
	MCPaths   int
}

func NewMetrics(m *Map) *Metrics {
	size := m.Rows * m.Cols
	met := Metrics{
		Map:      m,
		Seen:     make([]int, size),
		VisCount: make([]int, size),
		Land:     make([]int, size),
		PrFood:   make([]int, size),
		Unknown:  make([]int, size),
		VisSum:   make([]int, size),
		Horizon:  make([]bool, size),
		HBorder:  make([]Location, 0, 1000),
	}
	for i := 0; i < NTHREAT; i++ {
		met.Threat = append(met.Threat, make([]int8, size))
		met.PThreat = append(met.PThreat, make([]uint16, size))
	}
	for i, _ := range met.Unknown {
		met.Unknown[i] = 999
	}

	return &met
}

func (m *Metrics) UpdateCounts(loc Location, mask *Mask) {
	nunknown := 0
	nland := 0
	p := m.ToPoint(loc)
	for _, op := range mask.P {
		nloc := m.ToLocation(m.PointAdd(p, op))
		if m.TGrid[nloc] == UNKNOWN {
			nunknown++
		}
		if m.Grid[nloc] != WATER {
			nland++
		}
	}
	m.Unknown[loc] = nunknown
	m.Land[loc] = nland
}

func (m *Metrics) SumVisCount(loc Location, mask *Mask) {
	nvis := 0
	p := m.ToPoint(loc)
	for _, op := range mask.P {
		nloc := m.ToLocation(m.PointAdd(p, op))
		nvis += m.VisCount[nloc]
	}
	m.VisSum[loc] = nvis
}

func (m *Metrics) ComputePrFood(loc, sloc Location, turn int, mask *Mask, f *Fill) int {
	prfood := 0
	turn++
	horizonwt := 0
	p := m.ToPoint(loc)

	for _, op := range mask.P {
		nloc := m.ToLocation(m.PointAdd(p, op))
		if sloc == f.Seed[nloc] &&
			f.Distance(sloc, nloc) < 15 {
			viewwt := MaxV(4-m.VisCount[nloc], 1)

			// food we compete for is more of a priority
			if m.Horizon[nloc] {
				horizonwt = 4
			} else {
				horizonwt = 5
			}

			foodp := 0
			if m.Grid[nloc] != WATER {
				// TODO Max turn magic here should maybe decline over time.
				// also should test values
				foodp = MinV(turn-m.Seen[nloc], 12)
			}

			prfood += foodp * viewwt * horizonwt
		}
	}
	m.PrFood[loc] = prfood

	return prfood
}

// Compute the threat for N turns out (currently only n = 0 or 1)
// if player > -1 then sum players not including player
func (s *State) ComputeThreat(turn, player int, mask []*MoveMask, threat []int8, pthreat []uint16) {
	if turn > 1 || turn < 0 {
		log.Panicf("Illegal turns out = %d", turn)
	}

	if len(threat) != s.Rows*s.Cols || len(threat) != len(pthreat) {
		log.Panic("ComputeThreat slice size mismatch")
	}

	m := mask[0]
	for i, _ := range s.Ants {
		if i != player {
			for loc, _ := range s.Ants[i] {
				p := s.Map.ToPoint(loc)
				if turn > 0 {
					m = mask[s.Map.FreedomKey(loc)]
				}
				for i, op := range m.Point {
					threat[s.ToLocation(s.PointAdd(p, op))]++
					pthreat[s.ToLocation(s.PointAdd(p, op))] += m.MaxPr[i]
				}
			}
		}
	}

	return
}

func (s *State) StepHorizon(hlist []Location) []Location {
	m := s.Met
	hlist = hlist[0:0]

	// Remove now visible cells; dont bother with water here.
	for loc, Seen := range m.Seen {
		if Seen >= s.Turn {
			m.Horizon[loc] = true
		}
	}

	// generate list of cells on border; exclude water cells here
	for loc, h := range m.Horizon {
		if !h && m.Grid[loc] != WATER {
			for _, nloc := range m.LocStep[loc] {
				if m.Horizon[nloc] && m.Grid[nloc] != WATER {
					// if the point has an adjacent non horizon point which is not water 
					// then add it to the border list.
					hlist = append(hlist, Location(loc))
					break
				}
			}
		}
	}

	// step one from all cells on border
	for _, loc := range hlist {
		for d := 0; d < 4; d++ {
			m.Horizon[m.LocStep[loc][d]] = false
		}
	}

	return hlist
}
