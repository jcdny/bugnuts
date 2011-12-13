package combat

import (
	"log"
	. "bugnuts/game"
	. "bugnuts/util"
)

func (ps *PartitionState) ComputeScore(c *Combat) {
	N := make([]int, len(ps.P))
	for i := range ps.P {
		N[i] = len(ps.P[i].First)
		ps.P[i].Score = make([]int, N[i])
	}
	perm := PermuteList(N)

	move := make([][]AntMove, len(perm))
	for ip, p := range perm {
		move[ip] = make([]AntMove, ps.ALive)
		ib, ie := 0, 0
		for np := 0; np < ps.PLive; np++ {
			ib, ie = ie, ie+len(ps.P[np].First[p[np]])
			//log.Print("ip,np,P, ib,ie,ps.Alive:", ip, np, p[np], ib, ie, ps.ALive)
			copy(move[ip][ib:ie], ps.P[np].First[p[np]])
		}
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

		if true && len(dead) > 0 {
			log.Print("perm ", ip, ":", p, ": ", DumpMoves(move[ip]))
			log.Print("DEAD: ", len(dead), ": ", DumpMoves(dead))
		}
	}
	log.Print("Score: ", ps.P[0].Score)
}
