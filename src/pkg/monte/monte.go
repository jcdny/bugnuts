package monte

import (
	"log"
	. "bugnuts/state"
	. "bugnuts/game"
	. "bugnuts/maps"
	. "bugnuts/combat"
)

// MoveEm given a list of AntMove update D and To for the move in the given direction
func MoveEm(moves []AntMove, d Direction, m *Map, c *Combat) {
	for i := range moves {
		moves[i].D = d
		moves[i].To = m.LocStep[moves[i].From][d]
	}
	step := PermStepD5[d][0][:]

	var moved, nm int
	for nm, moved = 1, 0; moved < len(moves) && nm != 0; nm = 0 {
		for i := moved; i < len(moves); i++ {
			if !StepableItem[m.Grid[moves[i].To]] {
				// log.Print("Not stepable ", moves[i].From, moves[i].D, moves[i].To)
				// Water or food -- make space for guy if there is one.
				behind := m.LocStep[moves[i].From][step[4]]
				if c.AntCount[behind] != 0 {
					sr := m.LocStep[moves[i].From][step[1]]
					sl := m.LocStep[moves[i].From][step[2]]
					br := m.LocStep[sr][step[4]]
					bl := m.LocStep[sl][step[4]]

					if c.AntCount[br] == 0 && StepableItem[m.Grid[sr]] {
						moves[i].D = step[1]
						moves[i].To = sr
					} else if c.AntCount[bl] == 0 && StepableItem[m.Grid[sl]] {
						moves[i].D = step[2]
						moves[i].To = sl
					} else {
						continue
					}
				} else {
					continue
				}
			} else if c.AntCount[moves[i].To] != 0 {
				continue
			}

			// If we got here we have a valid move
			c.AntCount[moves[i].To]++
			c.AntCount[moves[i].From]--

			if moved < i {
				moves[moved], moves[i] = moves[i], moves[moved]
			}
			moved++
			nm++
		}
	}

	// Reset the unmovable
	for i := moved; i < len(moves); i++ {
		moves[i].D = NoMovement
		moves[i].To = moves[i].From
	}

	// Reset the player counts for the moved
	for i := 0; i < moved; i++ {
		c.AntCount[moves[i].From]++
		c.AntCount[moves[i].To]--
	}
}

func FirstStep(m *Map, c *Combat, ps *PartitionState) [][4][]AntMove {
	moves := make([][4][]AntMove, ps.PLive)
	for np := range ps.P {
		am := ps.P[np].Moves
		pm := make([]AntMove, ps.P[np].Live*4)
		for d := 0; d < 4; d++ {
			copy(pm[d*len(am):(d+1)*len(am)], am)
			moves[np][d] = pm[d*len(am) : (d+1)*len(am)]
			MoveEm(moves[np][d], Direction(d), m, c)
		}
	}

	return moves
}

func Sim(s *State, c *Combat, ps *PartitionState) {
	// the 1 step prefered direction permute for all ants
	fs := FirstStep(s.Map, c, ps)

	perm := genperm(uint(ps.PLive))
	move := make([][]AntMove, len(perm))

	for ip, p := range perm {
		move[ip] = make([]AntMove, ps.ALive)
		ib, ie := 0, 0
		for np := 0; np < ps.PLive; np++ {
			ib, ie = ie, ie+len(fs[np][p[np]])
			//log.Print("ip,np,P, ib,ie,ps.Alive:", ip, np, p[np], ib, ie, ps.ALive)
			copy(move[ip][ib:ie], fs[np][p[np]])
		}
		c2 := c.Copy()
		_, dead := c.Resolve(move[ip])
		c.Unresolve(move[ip])
		if len(dead) > 0 {
			log.Print("perm ", ip, ":", p, ": ", DumpMoves(move[ip]))
			log.Print("DEAD: ", len(dead), ": ", DumpMoves(dead))
		}
		if same, diffs := CombatCheck(c, c2); !same {
			log.Print("****************************** Unresolve unequal: ", diffs)
		}
	}

}

// generate the list of permuted directions for n players
func genperm(n uint) [][]Direction {
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
