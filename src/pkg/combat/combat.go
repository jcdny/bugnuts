package combat

import (
	"log"
	"time"
	"reflect"
	"strconv"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/state"
)

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

func (c *Combat) Copy() *Combat {
	cc := NewCombat(c.Map, c.AttackMask, len(c.PThreat))
	copy(cc.PlayerMap, c.PlayerMap)
	copy(cc.AntCount, c.AntCount)
	copy(cc.Threat, c.Threat)
	for i := range c.PThreat {
		cc.PThreat[i] = make([]int, len(c.PThreat[i]))
		copy(cc.PThreat[i], c.PThreat[i])
	}

	return cc
}

func CombatCheck(c, c2 *Combat) (equal bool, diffs []string) {
	equal = reflect.DeepEqual(c, c2)
	if !equal {
		diffs = make([]string, 16)
		if !reflect.DeepEqual(c2.PlayerMap, c.PlayerMap) {
			diffs = append(diffs, "PlayerMap")
		}
		if !reflect.DeepEqual(c2.AntCount, c.AntCount) {
			diffs = append(diffs, "AntCount")
		}
		if !reflect.DeepEqual(c2.Threat, c.Threat) {
			diffs = append(diffs, "Threat")
		}
		for i := range c.PThreat {
			if !reflect.DeepEqual(c2.PThreat[i], c.PThreat[i]) {
				key := "PThreat[" + strconv.Itoa(i) + "]"
				diffs = append(diffs, key)
			}
		}
	}

	return
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
func (c *Combat) Setup(ants []map[Location]int) {
	dead := make([]PlayerLoc, 0)
	for np := range ants {
		if len(ants[np]) > 0 && len(c.PThreat[np]) == 0 {
			c.PThreat[np] = make([]int, c.Map.Size())
		}
		for loc := range ants[np] {
			c.AddAnt(np, loc, dead)
		}
	}

	if len(dead) > 0 {
		log.Panic("Dead ants found in setup", len(dead))
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
				c.AddAnt(np, ants[np][i].From, dead)
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

func (c *Combat) AddAnt(np int, loc Location, dead []PlayerLoc) {

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

// Takes a collection of moves and undoes them
func (c *Combat) Unresolve(moves []AntMove) {
	// Revert Threat 
	for _, m := range moves {
		if c.AntCount[m.To] == 1 {
			if m.From != m.To {
				// good move update threat
				c.ApplyOffsets(m.From, &c.AttackMask.Add[m.D], func(nloc Location) {
					c.Threat[nloc]--
					c.PThreat[m.Player][nloc]--
				})
				c.ApplyOffsets(m.From, &c.AttackMask.Remove[m.D], func(nloc Location) {
					c.Threat[nloc]++
					c.PThreat[m.Player][nloc]++
				})
			}
		} else {
			// suicide replace threat
			c.ApplyOffsets(m.From, &c.AttackMask.Offsets, func(nloc Location) {
				c.Threat[nloc]++
				c.PThreat[m.Player][nloc]++
			})
		}
	}

	// Now revert counts and playermap
	for _, m := range moves {
		c.AntCount[m.To] = 0
		c.PlayerMap[m.To] = -1
	}
	for _, m := range moves {
		c.AntCount[m.From] = 1
		c.PlayerMap[m.From] = m.Player
	}
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

func CombatRun(s *State, ants []*AntStep, part Partitions, pmap PartitionMap) {
	if len(part) == 0 {
		return
	}
	// Setup combat engine
	c := NewCombat(s.Map, s.AttackMask, MaxPlayers) // TODO player counts?
	c.Setup(s.Ants)

	budget := (s.Cutoff - time.Nanoseconds()) / int64(len(part)) / 4

	// sim to compute best moves
	for {
		for ploc, ap := range part {
			t := time.Nanoseconds() + budget
			if t > s.Cutoff {
				break
			}
			c.Sim(s, ploc, ap, t)
		}
		// TODO REMOVE ME WHEN WORKING
		break
	}

	// Move combat moves back to antstep

	// vis
}

func (c *Combat) Sim(s *State, ploc Location, ap *AntPartition, cutoff int64) {
	log.Printf("Simulate for ap: %v %d ants, cutoff %.2fms",
		s.ToPoint(ploc),
		len(ap.Ants),
		float64(cutoff-time.Nanoseconds())/1e6)

}
