package combat

import (
	"log"
	"fmt"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/pathing"
	. "bugnuts/util"
	. "bugnuts/watcher"
)

const PDist = 6  // MAGIC -- enemy max distance for partition
const FPDist = 4 // MAGIC -- friendly distance to pull into partition

type AntPartition struct {
	PLoc  Location
	Ants  []Location
	Pants []Location
	PS    *PartitionState
}

type Partitions map[Location]*AntPartition
// PartitionMap maps an ant to the partitions it belongs to
type PartitionMap map[Location]map[Location]struct{}

func NewAntPartition(ploc Location) *AntPartition {
	p := &AntPartition{
		PLoc: ploc,
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

// Returns the list of partition keys for a given ant
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
	pn     [MaxPlayers]int // count of ants per player
	pp     [MaxPlayers]int // count of ants per player post prune
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
							ap = NewAntPartition(ploc)
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
				pstat.pn[c.PlayerMap[loc]]++
				pstat.tot++
			}
		}
		pstat.menemy = Max(pstat.pn[1:])
	}

	// For each partition, potentially grow ours adding our ants which are
	// neighbors of ants already in the partition.  Don't grow if N > maxE + 3
	for ploc, ap := range parts {
		if WS.Watched(ploc, 0) {
			log.Print("growing ", ploc, pstats[ploc])
		}
		pstat := pstats[ploc]
		if pstat.pn[0] < pstat.menemy+3 { // MAGIC
		NEXT:
			for _, loc := range ap.Ants {
				if c.PlayerMap[loc] == 0 {
					// one of our ants
					for nloc, nn := range near[loc] {
						// Add our nearby ants which are not already in the partition
						if nn.Steps < FPDist && c.PlayerMap[nloc] == 0 {
							if WS.Watched(ploc, 0) {
								log.Print("checking: ", nloc, nn)
							}
							if _, in := pmap[nloc][ploc]; !in {
								if WS.Watched(ploc, 0) {
									log.Print("Adding: ", nloc, nn)
								}

								pstat.pn[0]++
								pstat.tot++
								if pstat.pn[0] > pstat.menemy+3 {
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

	// Now prune any ant which need not be involved in combat...
	for ploc, ap := range parts {
		log.Print("Pruning ", ploc)
		found := false
		for np, n := range pstats[ploc].pn {
			if n > 0 {
				found = found || ap.prune(pstats[ploc], np, c)
			}
		}
		// Nuke any partition with no ants in combat
		if !found {
			parts[ploc] = &AntPartition{}, false
		}
	}

	// Finally disolve partitions which are 1-1
	// just rely on normal risk behavior on those.
	for ploc := range parts {
		if pstats[ploc].tot < 3 {
			parts[ploc] = &AntPartition{}, false
		}
	}

	return parts, pmap
}

func (ap *AntPartition) prune(stat *pStat, np int, c *Combat) bool {
	// skip the ants not part of np
	// also a bit of a hack... set playermap to -1 for now.
	// this is so pruning is easier.
	ants := make([]Location, 0, len(ap.Ants))
	n := 0
	for _, loc := range ap.Ants {
		if c.PlayerMap[loc] == np {
			ants = append(ants, loc)
			c.PlayerMap[loc] = -1
		} else {
			ap.Ants[n] = loc
			n++
		}
	}
	ap.Ants = ap.Ants[:n]
	// now have 2 lists, ants not me (ap.Ants) and ants to check (ants)
	log.Print("P ", np, " Others ", n, " me ", len(ants))

	// now go through and for any ant in combat add it to cant and put it in map
	// iterate until we have a round with no adds
	cants := make([]Location, 0, len(ants))
	for {
		n = 0
		for _, loc := range ants {
			var d int
			for d = 0; d < 5; d++ {
				nl := c.LocStep[loc][d]
				if c.PlayerMap[nl] > -1 {
					break
				}
				if _, ok := c.Risk[np][nl]; ok {
					break
				}
			}
			if d < 5 {
				cants = append(cants, loc)
				c.PlayerMap[loc] = np
			} else {
				ants[n] = loc
				n++
			}
		}
		if len(ants) == n {
			break
		} else {
			ants = ants[:n]
		}
	}
	// Trunc and bash player id back in for ants not in combat
	for _, loc := range ants {
		c.PlayerMap[loc] = np
	}

	// update stats
	stat.tot -= stat.pn[np] - len(ants)
	stat.pp[np] = stat.pn[np] - len(ants)

	// set pants to ants
	ap.Pants = append(ap.Pants, ants...)
	ap.Ants = append(ap.Ants, cants...)

	// at this point ants are ants that were not in combat or connected to 
	// an ant in combat are in ants, cants are combat ants
	log.Print("player ", np, " ants nc ", len(ants), " combat ", len(cants))

	return len(cants) > 0
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
