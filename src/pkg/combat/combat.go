package combat

import (
	"log"
	"time"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/state"
	. "bugnuts/pathing"
)

type AntState struct {
	Start    Location
	End      Location
	NStep    int
	Steps    [8]Direction
	Prefered Direction
}

type AntPartition struct {
	Ants    map[Location]struct{}
	Players []int
}

type Partitions map[Location]*AntPartition
type PartitionMap map[Location][]Location

func NewAntPartition() *AntPartition {
	p := &AntPartition{
		Ants: make(map[Location]struct{}, 8),
	}
	return p
}

type AntMove struct {
	From   Location
	To     Location
	D      Direction
	Player int
}

type Combat struct {
	*Map
	AttackMask *Mask
	PlayerMap  []int
	AntCount   []int
	Threat     []int
	PThreat    [][]int
}

var _minusone [MAXMAPSIZE]int
var _zero [MAXMAPSIZE]int

func init() {
	for i := range _minusone {
		_minusone[i] = -1
	}
}

func NewCombat(m *Map, am *Mask, np int) *Combat {
	c := &Combat{
		Map:        m,
		AttackMask: am,
		PlayerMap:  make([]int, m.Size()),
		AntCount:   make([]int, m.Size()),
		Threat:     make([]int, m.Size()),
		PThreat:    make([][]int, np),
	}
	copy(c.PlayerMap, _minusone[:len(c.PlayerMap)])

	return c
}

func (c *Combat) Reset() {
	copy(c.PlayerMap, _minusone[:len(c.PlayerMap)])
	copy(c.AntCount, _zero[:len(c.AntCount)])
	copy(c.Threat, _zero[:len(c.AntCount)])
	for i := range c.PThreat {
		copy(c.PThreat[i], _zero[:len(c.PThreat[i])])
	}
}

// Compute initial ant threat. returns a count of dead found.
// Should be 0 unless something has gone horribly wrong.
func (c *Combat) Setup(ants []map[Location]int) (dead []PlayerLoc) {
	dead = make([]PlayerLoc, 0)

	for np := range ants {
		if len(c.PThreat[np]) == 0 {
			c.PThreat[np] = make([]int, c.Map.Size())
		}
		for loc := range ants[np] {
			c.AddAnt(dead, np, loc)
		}
	}

	return
}
// Compute initial ant threat. returns a count of dead found.
// Should be 0 unless something has gone horribly wrong.
func (c *Combat) SetupReplay(ants [][]AntMove) (moves, spawn []AntMove) {
	dead := make([]PlayerLoc, 0)
	n := 0
	for np := range ants {
		n += len(ants[np])
	}
	moves = make([]AntMove, 0, n)
	spawn = make([]AntMove, 0, len(ants)*3)

	for np := range ants {
		if len(c.PThreat[np]) == 0 {
			c.PThreat[np] = make([]int, c.Map.Size())
		}
		for i := range ants[np] {
			// ants from from > -1 and to == -1 are ants that died in the 
			// previous turn. Ignored for now.
			if ants[np][i].From > -1 && ants[np][i].To > -1 {
				moves = append(moves, ants[np][i])
				c.AddAnt(dead, np, ants[np][i].From)
			} else if ants[np][i].To > -1 {
				spawn = append(spawn, ants[np][i])
			}
		}
	}

	if len(dead) > 0 {
		log.Panic("Dead ants found in replay:", len(dead))
	}

	return
}

func (c *Combat) AddAnt(dead []PlayerLoc, np int, loc Location) {

	c.AntCount[loc]++

	// If we encounter suicides on the first pair we remove
	// the original ant's threat, mark both the original and
	// current as dead, and reset the PlayerMap.  On
	// subsequent suicides at the same location we simply mark
	// the next ant as dead.
	inc := 1 // inc threat unless suicide then decr
	tp := np // threat for player
	if c.AntCount[loc] == 1 {
		c.PlayerMap[loc] = np
	} else if c.AntCount[loc] == 2 {
		inc = -1
		tp = c.PlayerMap[loc] // need to remove the original players threat
		c.PlayerMap[loc] = -1
		dead = append(dead, PlayerLoc{Loc: loc, Player: tp})
		dead = append(dead, PlayerLoc{Loc: loc, Player: np})
	} else {
		dead = append(dead, PlayerLoc{Loc: loc, Player: np})
	}

	if c.AntCount[loc] < 3 {
		c.ApplyOffsets(loc, &c.AttackMask.Offsets, func(nloc Location) {
			c.Threat[nloc] += inc
			c.PThreat[tp][nloc] += inc
		})
	}

	return
}

func (c *Combat) Resolve(moves []AntMove) (live, dead []AntMove) {
	// walk through the moves update counts
	for i := range moves {
		c.AntCount[moves[i].From]--
		c.AntCount[moves[i].To]++
	}

	// go through the moves and collect suicides and update threat for valid moves.
	ndead := 0
	for i := range moves {
		m := &moves[i]
		// TODO This should never happen consider removing in prod
		if c.AntCount[m.From] < 0 || c.PlayerMap[m.From] < 0 {
			log.Panic("Illegal step, source ant location %v, np %d", c.ToPoint(m.From), c.PlayerMap[m.From])
			if ndead < i {
				moves[ndead], moves[i] = moves[i], moves[ndead]
			}
			ndead++
			continue
		}

		if c.AntCount[m.To] == 1 {
			if m.From != m.To {
				// good move update threat
				c.ApplyOffsets(m.From, &c.AttackMask.Add[m.D], func(nloc Location) {
					c.Threat[nloc]++
					c.PThreat[m.Player][nloc]++
				})
				c.ApplyOffsets(m.From, &c.AttackMask.Remove[m.D], func(nloc Location) {
					c.Threat[nloc]--
					c.PThreat[m.Player][nloc]--
				})
			}
		} else {
			// suicide remove threat add as dead
			c.ApplyOffsets(m.From, &c.AttackMask.Offsets, func(nloc Location) {
				c.Threat[nloc]--
				c.PThreat[m.Player][nloc]--
			})
			if ndead < i {
				moves[ndead], moves[i] = moves[i], moves[ndead]
			}
			ndead++
		}
	}
	// log.Printf("len %d lmove %d\nsuicide is %v\nmoves is %v", len(moves), ndead, moves[:ndead], moves[ndead:])

	// now update player map for suicides and moves
	for i := range moves {
		c.PlayerMap[moves[i].From] = -1
	}
	for i := range moves[:ndead] {
		c.PlayerMap[moves[i].To] = -1
	}
	for i := ndead; i < len(moves); i++ {
		c.PlayerMap[moves[i].To] = moves[i].Player
	}

	// now do actual combat resolution

	for i := ndead; i < len(moves); i++ {
		// log.Printf("Combat for %v", c.ToPoint(moves[i].To))
		loc := moves[i].To
		np := moves[i].Player
		t := c.Threat[loc] - c.PThreat[np][loc]
		if t > 0 {
			c.ApplyOffsetsBreak(loc, &c.AttackMask.Offsets, func(nloc Location) bool {
				ntp := c.PlayerMap[nloc]
				if ntp >= 0 && ntp != np && t >= c.Threat[nloc]-c.PThreat[ntp][nloc] {
					if ndead < i {
						moves[ndead], moves[i] = moves[i], moves[ndead]
					}
					ndead++
					return false
				}
				return true
			})
		}
	}

	dead = moves[:ndead]
	live = moves[ndead:]

	// finally update player map for all dead
	for _, m := range dead {
		c.PlayerMap[m.To] = -1
	}

	return
}

func CombatPartition(s *State) (Partitions, PartitionMap) {
	// how many ants are there
	nant := 0
	for _, ants := range s.Ants {
		nant += len(ants)
	}

	origin := make(map[Location]int, nant)
	for _, ants := range s.Ants {
		for loc := range ants {
			origin[loc] = 1
		}
	}
	f := NewFill(s.Map)
	// will only find neighbors withing 2x8 steps.
	_, near := f.MapFillSeedNN(origin, 1, 8)

	parts := make(Partitions, 5)
	// maps an ant to the partitions it belongs to.
	pmap := make(PartitionMap, nant)

	for ploc := range s.Ants[0] {
		if _, ok := pmap[ploc]; !ok {
			// for any of my ants not already in a partition
			for eloc, nn := range near[ploc] {
				if nn.Steps < 7 {
					if _, ok := s.Ants[0][eloc]; !ok {
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
				if _, ok := s.Ants[0][loc]; ok {
					// one of our friendly ants, add any close neigbors of our friendly guy
					for nloc, nn := range near[loc] {
						if nn.Steps < 2 {
							_, me := s.Ants[0][nloc]
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
					if _, ok := s.Ants[0][floc]; !ok {
						// a close not me ant
						enemy[floc] = 0
					}
				}
			}
		}
	*/

	return parts, pmap
}

func CombatRun(s *State, ants []*AntStep, part Partitions, pmap PartitionMap) {
	if len(part) == 0 {
		return
	}
	budget := (s.Cutoff - time.Nanoseconds()) / int64(len(part)) / 4

	// sim to compute best moves
	for {
		for ploc, ap := range part {
			t := time.Nanoseconds() + budget
			if t > s.Cutoff {
				break
			}
			Sim(s, ploc, ap, t)
		}
		// TODO REMOVE ME WHEN WORKING
		break
	}

	// Move combat moves back to antstep

	// vis
}

func Sim(s *State, ploc Location, ap *AntPartition, cutoff int64) {
	log.Printf("Simulate for ap: %v %d ants, cutoff %.2fms",
		s.ToPoint(ploc),
		len(ap.Ants),
		float64(cutoff-time.Nanoseconds())/1e6)

}

// AntMove sorted by To then Player
type AntMoveSlice []AntMove

func (p AntMoveSlice) Len() int { return len(p) }
func (p AntMoveSlice) Less(i, j int) bool {
	return p[i].To < p[j].To || (p[i].To == p[j].To && p[i].Player < p[j].Player)
}
func (p AntMoveSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// AntMove sorted by Player then To
type AntMovePlayerSlice []AntMove

func (p AntMovePlayerSlice) Len() int { return len(p) }
func (p AntMovePlayerSlice) Less(i, j int) bool {
	return p[i].Player < p[j].Player || (p[i].Player == p[j].Player && p[i].To < p[j].To)
}
func (p AntMovePlayerSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func DumpAntMove(m *Map, am []AntMove, p int, turn int) {
	for _, a := range am {
		if p == a.Player || p < 0 {
			log.Printf("Move t=%d p=%d %#v %v %#v", turn, a.Player, m.ToPoint(a.From), a.D, m.ToPoint(a.To))
		}
	}
}
