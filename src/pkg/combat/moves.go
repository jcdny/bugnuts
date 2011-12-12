package combat

import (
	. "bugnuts/game"
	. "bugnuts/maps"
)

// MoveEm given a list of AntMove update D and To for the move in the given direction
func moveEm(moves []AntMove, d Direction, c *Combat) {
	m := c.Map
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

func (ps *PartitionState) FirstStep(c *Combat) {
	for np := range ps.P {
		am := ps.P[np].Moves
		pm := make([]AntMove, ps.P[np].Live*4)
		for d := 0; d < 4; d++ {
			copy(pm[d*len(am):(d+1)*len(am)], am)
			ps.P[np].First[d] = pm[d*len(am) : (d+1)*len(am)]
			moveEm(ps.P[np].First[d], Direction(d), c)
		}
	}
}
