package combat

import (
	"log"
	"fmt"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/pathing"
	. "bugnuts/util"
)

const PDist = 7  // MAGIC -- enemy max distance for partition
const FPDist = 3 // MAGIC -- friendly distance to pull into partition

type AntPartition struct {
	Ants []Location
	PS   *PartitionState
}

type Partitions map[Location]*AntPartition
// PartitionMap maps an ant to the partitions it belongs to
type PartitionMap map[Location]map[Location]struct{}

func NewAntPartition() *AntPartition {
	p := &AntPartition{
		Ants: make([]Location, 0, 8),
	}
	return p
}

type PartitionState struct {
	PLive int
	ALive int
	P     []PlayerState
}

type PlayerState struct {
	Player int
	Live   int
	Moves  []AntMove
	First  [][]AntMove
	// move scoring
	Score []int
	Best  int
}

func (pmap *PartitionMap) Add(from, to Location) {
	pm, ok := (*pmap)[from]
	if !ok {
		pm = make(map[Location]struct{}, 8)
		(*pmap)[from] = pm
	}
	// map a location to a partition
	pm[to] = struct{}{}
}
func (pmap *PartitionMap) Get(from Location) []Location {
	pm, ok := (*pmap)[from]
	if !ok || len(pm) == 0 {
		return []Location{}
	}
	out := make([]Location, 0, len(pm))
	for loc := range pm {
		out = append(out, loc)
	}
	return out
}

type pStat struct {
	tot    int
	menemy int
	pn     [MaxPlayers]int
}

func (c *Combat) Partition(Ants []map[Location]int) (Partitions, PartitionMap) {
	// how many ants are there
	nant := 0
	for _, ants := range Ants {
		nant += len(ants)
	}

	origin := make(map[Location]int, nant)
	for _, ants := range Ants {
		for loc := range ants {
			origin[loc] = 1
		}
	}
	f := NewFill(c.Map)
	// will only find neighbors withing 2x8 steps.
	_, near := f.MapFillSeedNN(origin, 1, 8)

	// the actual partitions
	parts := make(Partitions, 5)
	// maps an ant to the partitions it belongs to.
	pmap := make(PartitionMap, nant)

	for ploc := range Ants[0] {
		// If the ant is not mapped to a partition
		if _, ok := pmap[ploc]; !ok {
			// Look at the nearby ants
			for eloc, nn := range near[ploc] {
				// a close enemy ant
				if nn.Steps < PDist && c.PlayerMap[eloc] != 0 {
					// if mapped to a partition already then merge
					// to that one.
					if _, mapped := pmap[eloc]; mapped {
						ploc = pmap.Get(eloc)[0]
					} else {
						ap, ok := parts[ploc] // ap = a partition
						if !ok {
							ap = NewAntPartition()
						}
						parts[ploc] = ap
					}

					// add the enemy and any of it's neighbors
					pmap.Add(eloc, ploc)
					for nloc, nn := range near[eloc] {
						if nn.Steps < PDist {
							pmap.Add(nloc, ploc)
						}
					}
				}
			}
		}
	}

	// invert pmap to generate ap members
	for aloc := range pmap {
		for ploc := range pmap[aloc] {
			parts[ploc].Ants = append(parts[ploc].Ants, aloc)
		}
	}

	// Compute the counts per player for each partition
	pstats := make(map[Location]*pStat, len(parts))
	for ploc, ap := range parts {
		pstats[ploc] = &pStat{}
		pstat := pstats[ploc]
		for _, loc := range ap.Ants {
			if c.PlayerMap[loc] > -1 {
				pstat.pn[c.PlayerMap[loc]] += 1
				pstat.tot++
			}
		}
		pstat.menemy = Max(pstat.pn[1:])
	}

	// For each partition, potentially grow ours adding our ants which are
	// neighbors of ants already in the partition.  Don't grow if N > maxE + 3
	for ploc, ap := range parts {
		pstat := pstats[ploc]
		if pstat.pn[0] < pstat.menemy+2 { // MAGIC
		NEXT:
			for _, loc := range ap.Ants {
				if c.PlayerMap[loc] == 0 {
					// one of our ants
					for nloc, nn := range near[loc] {
						// Add our nearby ants which are not already in the partition
						if nn.Steps < FPDist && c.PlayerMap[nloc] == 0 {
							if _, in := pmap[nloc][ploc]; !in {
								pstat.pn[0]++
								pstat.tot++
								if pstat.pn[0] < pstat.menemy+2 {
									break NEXT
								}
								pmap.Add(nloc, ploc)
								ap.Ants = append(ap.Ants, nloc)
							}
						}
					}
				} else {
					continue
				}
			}
		}
	}

	// Finally disolve partitions which are 1-1
	// just rely on normal risk behavior on those.
	for ploc := range parts {
		if pstats[ploc].tot < 3 {
			parts[ploc] = &AntPartition{}, false
		}
	}

	for i := 0; i < len(Ants); i++ {
		if len(Ants[i]) > 0 {
			c.PFill[i] = ThreatFill(c.Map, c.Threat1, c.PThreat1[i], 10, 0)
		}
	}

	return parts, pmap
}

// NewPartitionState creates the move list and player states
func NewPartitionState(c *Combat, ap *AntPartition) *PartitionState {
	ps := &PartitionState{}

	players := make([]int, MaxPlayers)
	playermap := make([]int, MaxPlayers)

	for i, loc := range ap.Ants {
		if c.PlayerMap[loc] > -1 {
			players[c.PlayerMap[loc]]++
			ps.ALive++
		} else {
			log.Print("Invalid ap player loc %v removing it", c.ToPoint(loc))
			log.Print(ap, c.PlayerMap[loc])
			copy(ap.Ants[i:len(ap.Ants)-1], ap.Ants[i+1:])
		}
	}

	for i, n := range players {
		if n > 0 {
			playermap[i] = ps.PLive
			ps.PLive++
		} else {
			playermap[i] = -1
		}
	}

	// sanity
	if ps.PLive < 2 {
		log.Panic("Partition with less than 2 players")
		return ps
	}

	// Populate the actual player states
	ps.P = make([]PlayerState, ps.PLive)
	for _, loc := range ap.Ants {
		np := c.PlayerMap[loc]
		if np > -1 {
			if players[np] > 0 {
				ps.P[playermap[np]] = PlayerState{
					Player: np,
					Moves:  make([]AntMove, 0, players[np]),
					Live:   players[np],
				}
				players[np] = 0

			}
			ps.P[playermap[np]].Moves = append(ps.P[playermap[np]].Moves,
				AntMove{From: loc, To: loc, D: NoMovement, Player: np})
		}
	}

	return ps
}

func (ps *PlayerState) bestScore() (best []int) {
	ms := Max(ps.Score[:])
	for i, s := range ps.Score {
		if s == ms {
			best = append(best, i)
		}
	}
	return
}
func DumpPartitionState(ps *PartitionState) string {
	s := fmt.Sprintf("Players %d Ants %d: ", ps.PLive, ps.ALive)
	for i := range ps.P {
		s += "\n  " + DumpPlayerState(&ps.P[i])
	}
	return s
}

func DumpMoves(moves []AntMove) string {
	s := ""
	for n, am := range moves {
		if n == 0 {
		} else if n%6 == 0 {
			s += "\n    "
		} else {
			s += ";"
		}
		s += fmt.Sprintf("%v %v %v", am.From, am.D, am.To)
	}
	return s
}
func DumpPlayerState(p *PlayerState) string {
	s := fmt.Sprintf("Player %d Ants %d: ", p.Player, p.Live)
	s += DumpMoves(p.Moves)
	return s
}
