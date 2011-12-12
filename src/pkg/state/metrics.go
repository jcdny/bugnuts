package state

import (
	. "bugnuts/pathing"
	. "bugnuts/torus"
	. "bugnuts/maps"
	. "bugnuts/util"
)

const NTHREAT = 10

type Metrics struct {
	*Map
	Seen     []int      // Turn on which cell was last visible.
	VisCount []int      // How many ants see this cell.
	Horizon  []bool     // Inside the event horizon.  false means there could be an ant there we have not seen
	HBorder  []Location // List of border points
	Land     []int      // Count of land tiles visible from a given tile
	PrFood   []int      // Count of Turn * Land adjusted for # that see it.
	Unknown  []int      // Count of Unknown visible from a given location
	VisSum   []int      // sum of count of visibles for overlap.
	Runs     [][4]uint8 // What is the run distance in a given direction for a location
	// Fills
	FDownhill *Fill // Downhill from my own hills
	FHill     *Fill // Distance to my hills
	// FViewHill *Fill // Distance to visibility boundary for my hills
	// FAll      *Fill // All hill fill
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
		Runs:     make([][4]uint8, size),
		HBorder:  make([]Location, 0, 1000),
	}

	for i := range met.Unknown {
		met.Unknown[i] = 999
	}

	return &met
}

func (m *Metrics) UpdateRuns() {
	for r := 0; r < m.Rows; r++ {
		loc := r * m.Cols
		cs, ce := 0, 0
		for cs < m.Cols {
			for cs < m.Cols && !StepableItem[m.Grid[loc+cs]] {
				cs++
			}
			for ce = cs; ce < m.Cols && StepableItem[m.Grid[loc+ce]]; ce++ {
				m.Runs[loc+ce][West] = uint8(ce - cs)
			}
			if cs == 0 && ce == m.Cols {
				// Special case empty row
				for c := 0; c < m.Cols; c++ {
					m.Runs[loc+c][East] = uint8(m.Cols)
					m.Runs[loc+c][West] = uint8(m.Cols)
				}
			} else {
				for cb := ce - 1; cb >= cs; cb-- {
					m.Runs[loc+cb][East] = uint8(ce - 1 - cb)
				}
			}
			cs = ce
		}
		// Fall out here at end, if both 0 and Cols are
		// steppable bridge the two unless its an empty row
		if StepableItem[m.Grid[loc]] && StepableItem[m.Grid[loc+m.Cols-1]] {
			ce := int(m.Runs[loc][East])
			cw := int(m.Runs[loc+m.Cols-1][West])
			if ce != m.Cols {
				for c := 0; c <= ce; c++ {
					m.Runs[loc+c][West] += uint8(cw + 1)
				}
				for c := 0; c <= cw; c++ {
					m.Runs[loc+m.Cols-1-c][East] += uint8(ce + 1)
				}
			}
		}
	}
	for c := 0; c < m.Cols; c++ {
		rs, re := 0, 0
		for rs < m.Rows {
			for rs < m.Rows && !StepableItem[m.Grid[rs*m.Cols+c]] {
				rs++
			}
			for re = rs; re < m.Rows && StepableItem[m.Grid[re*m.Cols+c]]; re++ {
				m.Runs[re*m.Cols+c][North] = uint8(re - rs)
			}
			if rs == 0 && re == m.Rows {
				// Special case empty col
				for r := 0; r < m.Rows; r++ {
					m.Runs[r*m.Cols+c][North] = uint8(m.Rows)
					m.Runs[r*m.Cols+c][South] = uint8(m.Rows)
				}
			} else {
				for rb := re - 1; rb >= rs; rb-- {
					m.Runs[rb*m.Cols+c][South] = uint8(re - 1 - rb)
				}
			}
			rs = re
		}
		// Fall out here at end, if both 0 and Row end are steppable bridge the two.
		if StepableItem[m.Grid[c]] && StepableItem[m.Grid[(m.Rows-1)*m.Cols+c]] {
			rs := int(m.Runs[c][South])
			rn := int(m.Runs[(m.Rows-1)*m.Cols+c][North])
			if rn != m.Rows {
				for r := 0; r <= rs; r++ {
					m.Runs[r*m.Cols+c][North] += uint8(rn + 1)
				}
				for r := 0; r <= rn; r++ {
					m.Runs[(m.Rows-r-1)*m.Cols+c][South] += uint8(rs + 1)
				}
			}
		}
	}
}

func (m *Metrics) UpdateCounts(loc Location, o *Offsets) {
	var nunknown, nland int
	m.Map.ApplyOffsets(loc, o, func(nloc Location) {
		if m.TGrid[nloc] == UNKNOWN {
			nunknown++
		}
		if m.Grid[nloc] != WATER {
			nland++
		}
	})
	m.Unknown[loc] = nunknown
	m.Land[loc] = nland
}

func (m *Metrics) SumVisCount(loc Location, o *Offsets) {
	nvis := 0
	m.Map.ApplyOffsets(loc, o, func(nloc Location) { nvis += m.VisCount[nloc] })
	m.VisSum[loc] = nvis
}

func (m *Metrics) ComputePrFood(loc, sloc Location, turn int, o *Offsets, f *Fill) int {
	prfood := 0
	turn++

	m.Map.ApplyOffsets(loc, o, func(nloc Location) {
		var horizonwt int

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
	})
	m.PrFood[loc] = prfood

	return prfood
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
