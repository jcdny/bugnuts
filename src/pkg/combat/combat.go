package combat

import (
	"log"
	"reflect"
	"strconv"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/pathing"
)

type Combat struct {
	*Map
	AttackMask *Mask
	PlayerMap  []int
	AntCount   []int
	Threat     []int      // Global threat
	PThreat    [][]int    // Player's Own Threat
	Ants1      []int      // bitmask of possible ants
	Threat1    []int      // One Step Threat
	PThreat1   [][]int    // Player's one step threat
	PThreat1Pr []int      // one step threat prob (only computed for player 0)
	PFill      []*Fill    // Distance to Threat1 surface
	TBPathin   []int      // The pathin count to the threat border player 0 enemies
	TBPathout  []int      // The pathin count to the threat border player 0
	Border     []Location // The threat border player 0
	Risk       []map[Location]int
	ECutoff    int
	FCutoff    int
	Score      func(dead []AntMove, np int) (score float64)
}

var _minusone [MAXMAPSIZE]int
var _zero [MAXMAPSIZE]int

func init() {
	for i := range _minusone {
		_minusone[i] = -1
	}
}

func NewCombat(m *Map, am *Mask, np int, score func(dead []AntMove, np int) (score float64)) *Combat {
	c := &Combat{
		Map:        m,
		AttackMask: am,
		PlayerMap:  make([]int, m.Size()),
		AntCount:   make([]int, m.Size()),
		Threat:     make([]int, m.Size()),
		PThreat:    make([][]int, np),
		Ants1:      make([]int, m.Size()), //
		Threat1:    make([]int, m.Size()), //
		PThreat1:   make([][]int, np),
		PThreat1Pr: make([]int, m.Size()),
		PFill:      make([]*Fill, np),
		Score:      score,
	}
	copy(c.PlayerMap, _minusone[:len(c.PlayerMap)])

	return c
}

// Copy a Combat struct.  Just copies the things the combat section uses 
// metrics.
func (c *Combat) Copy() *Combat {
	cc := NewCombat(c.Map, c.AttackMask, len(c.PThreat), c.Score)
	copy(cc.PlayerMap, c.PlayerMap)
	copy(cc.AntCount, c.AntCount)
	copy(cc.Threat, c.Threat)
	for i := range c.PThreat {
		cc.PThreat[i] = make([]int, len(c.PThreat[i]))
		copy(cc.PThreat[i], c.PThreat[i])
	}

	return cc
}

func CombatCheck(c, c2 *Combat) (equal bool, diffs map[string]struct{}) {

	diffs = make(map[string]struct{}, 8)
	if !reflect.DeepEqual(c2.PlayerMap, c.PlayerMap) {
		diffs["PlayerMap"] = struct{}{}
	}
	if !reflect.DeepEqual(c2.AntCount, c.AntCount) {
		diffs["AntCount"] = struct{}{}
	}
	if !reflect.DeepEqual(c2.Threat, c.Threat) {
		diffs["Threat"] = struct{}{}
	}
	for i := range c.PThreat {
		if !reflect.DeepEqual(c2.PThreat[i], c.PThreat[i]) {
			key := "PThreat[" + strconv.Itoa(i) + "]"
			diffs[key] = struct{}{}
		}
	}
	equal = len(diffs) == 0

	return
}

func (c *Combat) Reset() {
	copy(c.PlayerMap, _minusone[:len(c.PlayerMap)])
	copy(c.AntCount, _zero[:len(c.AntCount)])
	copy(c.Threat, _zero[:len(c.Threat)])
	copy(c.Ants1, _zero[:len(c.Ants1)])
	copy(c.Threat1, _zero[:len(c.Threat1)])
	for i := range c.PThreat {
		copy(c.PThreat[i], _zero[:len(c.PThreat[i])])
		copy(c.PThreat1[i], _zero[:len(c.PThreat[i])])
	}
}

func (c *Combat) ResetAlloc() {
	copy(c.PlayerMap, _minusone[:len(c.PlayerMap)])
	c.AntCount = make([]int, len(c.AntCount))
	c.Threat = make([]int, len(c.Threat))
	c.Ants1 = make([]int, len(c.Threat1))
	c.Threat1 = make([]int, len(c.Threat1))
	for i := range c.PThreat {
		c.PThreat[i] = make([]int, len(c.PThreat[i]))
		c.PThreat1[i] = make([]int, len(c.PThreat1[i]))
	}
}

// Compute initial ant threat and threat fill.
func (c *Combat) Setup(ants []map[Location]int) {
	dead := make([]PlayerLoc, 0)
	for np := range ants {
		if len(ants[np]) > 0 && len(c.PThreat[np]) == 0 {
			c.PThreat[np] = make([]int, c.Map.Size())
			c.PThreat1[np] = make([]int, c.Map.Size())
		}
		for loc := range ants[np] {
			c.AddAnt(np, loc, dead)
		}
	}

	for i := 0; i < len(ants); i++ {
		maxdepth := uint16(10)
		var border, interior []Location
		if len(ants[i]) > 0 {
			if i == 0 {
				maxdepth = 0
			}
			c.PFill[i], border, interior = ThreatFill(c.Map, c.Threat1, c.PThreat1[i], maxdepth, 0)

			if false && i == 0 {
				c.TBPathin, c.ECutoff = ThreatPathin(c.PFill[0], ants[1:])
				c.TBPathout, c.FCutoff = ThreatPathin(c.PFill[0], ants[:1])
				c.Border = border
			}

			for _, loc := range interior {
				c.PFill[i].Depth[loc] = 1
			}
		}
	}

	c.Risk = c.Riskly(ants)

	if len(dead) > 0 {
		log.Panic("Dead ants found in setup", len(dead))
	}
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
			c.PThreat1[np] = make([]int, c.Map.Size())
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
	// If we encounter suicides on the first pair we remove the
	// original ant's threat, mark both the original and current as
	// dead, and reset the PlayerMap.  On subsequent suicides at the
	// same location we simply mark the next ant as dead.

	c.AntCount[loc]++
	inc := 1 // inc threat unless suicide then decr
	tp := np // threat for player
	if c.AntCount[loc] == 1 {
		c.PlayerMap[loc] = np
		c.Ants1[loc] |= PlayerFlag[np]
	} else if c.AntCount[loc] == 2 {
		inc = -1
		tp = c.PlayerMap[loc] // need to remove the original players threat
		c.Ants1[loc] &= PlayerMask[tp]
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
			c.Threat1[nloc] += inc
			c.PThreat1[tp][nloc] += inc
		})

		// Get the freedomkey and mask the player in to the permissible steps at the same time.
		key := 0
		for i := uint(0); i < 4; i++ {
			l := c.Map.LocStep[loc][i]
			if StepableItem[c.Map.Grid[l]] {
				key += 1 << i
				if inc > 0 {
					c.Ants1[l] |= PlayerFlag[tp]
				} else {
					c.Ants1[l] &= PlayerMask[tp]
				}
			}
		}

		c.ApplyOffsets(loc, &c.AttackMask.MM[key].Add, func(nloc Location) {
			c.Threat1[nloc] += inc
			c.PThreat1[tp][nloc] += inc
		})
	}

	return
}

// Unresolve takes a collection of ant moves and undoes them.
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
			log.Print("Illegal step for move set ", moves)
			log.Panicf("Source ant location %v, PlayerMap %d, AntCount %d",
				c.ToPoint(m.From), c.PlayerMap[m.From], c.AntCount[m.From])

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
