package combat

import (
	"log"
	"time"
	"reflect"
	"strconv"
	"rand"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/util"
)

type Combat struct {
	*Map
	AttackMask *Mask
	PlayerMap  []int
	AntCount   []int
	Threat     []int   // Global threat
	PThreat    [][]int // Player's Own Threat
	Threat1    []int   // One Step Threat
	PThreat1   [][]int // Player's one step threat
	PThreat1Pr []int   // one step threat prob (only computed for player 0)
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
		Threat1:    make([]int, m.Size()), //
		PThreat1:   make([][]int, np),
		PThreat1Pr: make([]int, m.Size()),
	}
	copy(c.PlayerMap, _minusone[:len(c.PlayerMap)])

	return c
}

// Copy a Combat struct.  This does not copy 1 step threat and other single ant
// metrics.
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

func CombatCheck(c, c2 *Combat) (equal bool, diffs map[string]struct{}) {
	equal = reflect.DeepEqual(c, c2)
	if !equal {
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
	}

	return
}

func (c *Combat) Reset() {
	copy(c.PlayerMap, _minusone[:len(c.PlayerMap)])
	copy(c.AntCount, _zero[:len(c.AntCount)])
	copy(c.Threat, _zero[:len(c.AntCount)])
	for i := range c.PThreat {
		copy(c.PThreat[i], _zero[:len(c.PThreat[i])])
		copy(c.PThreat1[i], _zero[:len(c.PThreat[i])])
	}
}

// Compute initial ant threat. returns a count of dead found.
// Should be 0 unless something has gone horribly wrong.
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
			c.Threat1[nloc] += inc
			c.PThreat1[tp][nloc] += inc
		})

		c.ApplyOffsets(loc, &c.AttackMask.MM[c.Map.FreedomKey(loc)].Add, func(nloc Location) {
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

func (c *Combat) Run(ants []*AntStep, part Partitions, pmap PartitionMap, cutoff int64, rng *rand.Rand) {
	if len(part) == 0 {
		return
	}

	budget := (cutoff - time.Nanoseconds()) / int64(len(part)) / 4
	// sim to compute best moves
	for {
		for ploc, ap := range part {
			t := time.Nanoseconds() + budget
			if t > cutoff {
				break
			}
			c.Sim(ap, ploc, t, rng)
		}
		break
	}
	mm := make(map[Location]*AntMove, 100)
	for _, ap := range part {
		ps := &ap.PS.P[0]
		if ps.Best != int(InvalidMove) {
			for _, am := range ps.First[ps.Best] {
				mm[am.From] = &am
			}
		}
	}

	for _, a := range ants {
		if move, ok := mm[a.Source]; ok {
			a.Move = move.D
		}
	}

	// vis
}

func (c *Combat) Sim(ap *AntPartition, ploc Location, cutoff int64, rng *rand.Rand) {
	log.Printf("Simulate for ap: %v %d ants, cutoff %.2fms",
		c.ToPoint(ploc),
		len(ap.Ants),
		float64(cutoff-time.Nanoseconds())/1e6)

	ap.PS = NewPartitionState(c, ap)
	ap.PS.FirstStep(c)
	MonteSim(c, ap.PS, rng)

}

func MonteSim(c *Combat, ps *PartitionState, rng *rand.Rand) {
	perm := genPerm(uint(ps.PLive))
	move := make([][]AntMove, len(perm))

	for ip, p := range perm {
		move[ip] = make([]AntMove, ps.ALive)
		ib, ie := 0, 0
		for np := 0; np < ps.PLive; np++ {
			ib, ie = ie, ie+len(ps.P[np].First[p[np]])
			//log.Print("ip,np,P, ib,ie,ps.Alive:", ip, np, p[np], ib, ie, ps.ALive)
			copy(move[ip][ib:ie], ps.P[np].First[p[np]])
		}
		//c2 := c.Copy()
		_, dead := c.Resolve(move[ip])
		for _, da := range dead {
			for np := range ps.P {
				if da.Player == np {
					ps.P[np].Score[p[np]] -= 20
				} else {
					ps.P[np].Score[p[np]] += 10
				}
			}
		}
		c.Unresolve(move[ip])
		if false && len(dead) > 0 {
			log.Print("perm ", ip, ":", p, ": ", DumpMoves(move[ip]))
			log.Print("DEAD: ", len(dead), ": ", DumpMoves(dead))
		}

		/*
			 if same, diffs := CombatCheck(c, c2); !same {
				log.Print("****************************** Unresolve unequal: ", diffs)
			}
		*/
	}
	log.Print("Score: ", ps.P[0].Score)
	if Min(ps.P[0].Score[:]) != Max(ps.P[0].Score[:]) {
		ms := Max(ps.P[0].Score[:])
		for _, d := range Permute4(rng) {
			if ps.P[0].Score[d] == ms {
				ps.P[0].Best = int(d)
			}
		}
	} else {
		ps.P[0].Best = int(InvalidMove)
	}
}

// generate the list of permuted directions for n players
func genPerm(n uint) [][]Direction {
	nperm := uint(4) << (2 * (n - 1))
	dl := make([]Direction, nperm*n)
	out := make([][]Direction, nperm)
	for i := uint(0); i < nperm; i++ {
		for s := uint(0); s < n; s++ {
			dl[i*n+s] = Direction((i >> (2 * s)) & 3)
		}
		out[i] = dl[i*n : (i+1)*n]
	}

	return out
}
