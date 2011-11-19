package state

import (
	"strconv"
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
	FDownhill *Fill
	FHill     *Fill
	FAll      *Fill
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
	return &met
}

func (m *Metrics) UpdateCounts(loc Location, mask *Mask) {
	nunknown := 0
	nland := 0
	p := m.ToPoint(loc)
	for _, op := range mask.P {
		nloc := m.ToLocation(m.PointAdd(p, op))
		if m.Grid[nloc] == UNKNOWN {
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
func (m *Metrics) DumpSeen() string {
	max := Max(m.Seen)
	str := ""

	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			str += strconv.Itoa(m.Seen[r*m.Cols+c] * 10 / (max + 1))
		}
		str += "\n"
	}

	return str
}
