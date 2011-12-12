package combat

import (
	"log"
	"fmt"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/pathing"
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
	First  [4][]AntMove
	// move scoring
	Score [4]int
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

	c.PFill[0] = ThreatFill(c.Map, c.Threat1, c.PThreat1[0], 10, 0)

	// the actual partitions
	parts := make(Partitions, 5)
	// maps an ant to the partitions it belongs to.
	pmap := make(PartitionMap, nant)

	for ploc := range Ants[0] {
		// If the ant is not mapped to a partition
		if _, ok := pmap[ploc]; !ok {
			// Look at the nearby ants 
			for eloc, nn := range near[ploc] {
				if nn.Steps < PDist && c.PlayerMap[eloc] != 0 {
					// a close enemy an, create a partition if needed
					ap, ok := parts[ploc] // ap = a partition
					if !ok {
						ap = NewAntPartition()
						pmap.Add(ploc, ploc)
					}
					parts[ploc] = ap

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

	// For each partition, grow ours adding friendly neighbors of ants
	// already in the partition
	for ploc, ap := range parts {
		for _, loc := range ap.Ants {
			if c.PlayerMap[loc] == 0 {
				// one of our ants
				for nloc, nn := range near[loc] {
					// Add our nearby ants which are not already in the partition
					if nn.Steps < FPDist && c.PlayerMap[nloc] == 0 {
						if _, in := pmap[nloc][ploc]; !in {
							pmap.Add(ploc, nloc)
							ap.Ants = append(ap.Ants, nloc)
						}
					}
				}
			}
		}
	}

	// Now disolve partitions which are 1-1
	for ploc, ap := range parts {
		if len(ap.Ants) == 2 {
			parts[ploc] = &AntPartition{}, false
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
