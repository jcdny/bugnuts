package combat

import (
	"log"
	. "bugnuts/game"
	. "bugnuts/util"
	. "bugnuts/watcher"
)

func ScoreRiskAverse(dead []AntMove, np int) (score float64) {
	for _, da := range dead {
		if da.Player == np {
			score -= 2
		} else {
			score += 1
		}
	}

	return
}

func ScoreRiskNeutral(dead []AntMove, np int) (score float64) {
	for _, da := range dead {
		if da.Player == np {
			score -= 1
		} else {
			score += 1
		}
	}
	return
}

func (ps *PartitionState) ComputeScore(c *Combat) {
	if ps == nil || len(ps.P) == 0 {
		return
	}
	TPush("+scoring")
	N := make([]int, len(ps.P))
	for i := range ps.P {
		N[i] = len(ps.P[i].First)
		ps.P[i].Score = make([]float64, N[i])
	}
	perm := PermuteList(N)
	defer TPopn(len(perm) * ps.ALive)

	move := make([][]AntMove, len(perm))
	for ip, p := range perm {
		move[ip] = make([]AntMove, ps.ALive)
		ib, ie := 0, 0
		for np := 0; np < ps.PLive; np++ {
			if len(p) == 0 {
				log.Print("Bad P:", len(ps.P), ps.ALive, N)
			}
			ib, ie = ie, ie+len(ps.P[np].First[p[np]])
			//log.Print("Perm,Player,Scenario,ib,ie,total", ip, np, p[np], ib, ie, ps.ALive)
			//log.Print(ps.P[np].First[p[np]])
			copy(move[ip][ib:ie], ps.P[np].First[p[np]])
		}
		_, dead := c.Resolve(move[ip][:ie])
		// just score player 0 for now....
		ps.P[0].Score[p[0]] = c.Score(dead, 0)
		c.Unresolve(move[ip])

		if Debug[DBG_Combat] && len(dead) > 0 {
			log.Print("perm ", ip, ":", p, ": ", DumpMoves(move[ip]))
			log.Print("DEAD: ", len(dead), ": ", DumpMoves(dead))
		}
	}
	if Debug[DBG_Combat] {
		log.Print("Score: ", ps.P[0].Score)
	}
}
