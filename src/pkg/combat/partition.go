package combat

import (
	"log"
	"fmt"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/pathing"
)

type AntPartition struct {
	Ants map[Location]struct{}
	PS   *PartitionState
}

type Partitions map[Location]*AntPartition
type PartitionMap map[Location][]Location

func NewAntPartition() *AntPartition {
	p := &AntPartition{
		Ants: make(map[Location]struct{}, 8),
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

func CombatPartition(Ants []map[Location]int, m *Map) (Partitions, PartitionMap) {
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
	f := NewFill(m)
	// will only find neighbors withing 2x8 steps.
	_, near := f.MapFillSeedNN(origin, 1, 8)

	parts := make(Partitions, 5)
	// maps an ant to the partitions it belongs to.
	pmap := make(PartitionMap, nant)

	for ploc := range Ants[0] {
		if _, ok := pmap[ploc]; !ok {
			// for any of my ants not already in a partition
			for eloc, nn := range near[ploc] {
				if nn.Steps < 7 {
					if _, ok := Ants[0][eloc]; !ok {
						// a close enemy ant, add it and it's nearest neighbors to the partition

						ap, ok := parts[ploc] // ap = a partition
						if !ok {
							ap = NewAntPartition()
						}
						parts[ploc] = ap

						for nloc, nn := range near[eloc] {
							if nn.Steps < 7 {
								ap.Ants[nloc] = struct{}{}

								pm, ok := pmap[nloc]
								if !ok {
									pm = make([]Location, 0, 8)
								}
								pmap[nloc] = append(pm, ploc)
							}
						}
						pm, ok := pmap[eloc]
						if !ok {
							pm = make([]Location, 0, 8)
						}

						pmap[eloc] = append(pm, ploc)
						ap.Ants[eloc] = struct{}{}
					}
				}
			}
		}

		if ap, ok := parts[ploc]; ok {
			// If we created a partition centered on this ant add any
			// close neighbors of the friendly ants already in the
			// partition
			for loc := range ap.Ants {
				if _, ok := Ants[0][loc]; ok {
					// one of our friendly ants, add any close neigbors of our friendly guy
					for nloc, nn := range near[loc] {
						if nn.Steps < 2 {
							_, me := Ants[0][nloc]
							_, in := ap.Ants[nloc]
							if me && !in {
								ap.Ants[nloc] = struct{}{}
								pm, ok := pmap[nloc]
								if !ok {
									pm = make([]Location, 0, 8)
								}
								pmap[nloc] = append(pm, ploc)
							}
						}
					}
				}
			}
		}
	}

	/*
		for loc := range enemy {
			for floc, nn := range near[loc] {
				if nn.Steps < 6 {
					if _, ok := Ants[0][floc]; !ok {
						// a close not me ant
						enemy[floc] = 0
					}
				}
			}
		}
	*/

	return parts, pmap
}

// NewPartitionState creates the move list and player states
func NewPartitionState(c *Combat, ap *AntPartition) *PartitionState {
	ps := &PartitionState{}

	players := make([]int, MaxPlayers)
	playermap := make([]int, MaxPlayers)

	for loc := range ap.Ants {
		if c.PlayerMap[loc] > -1 {
			players[c.PlayerMap[loc]]++
			ps.ALive++
		} else {
			log.Print("Invalid ap player loc %v removing it", c.ToPoint(loc))
			log.Print(ap, c.PlayerMap[loc])
			ap.Ants[loc] = struct{}{}, false
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
	for loc := range ap.Ants {
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
