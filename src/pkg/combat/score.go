// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package combat

import (
	"log"
	"time"
	. "bugnuts/game"
	. "bugnuts/util"
	. "bugnuts/watcher"
)

func ScoreRiskAverse(dead []AntMove, np int) (score float64) {
	log.Print("Risk averse scoring")
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
	log.Print("Risk Neutral scoring")

	for _, da := range dead {
		if da.Player == np {
			score -= 1
		} else {
			score += 1
		}
	}
	return
}

func (ps *PartitionState) ComputeScore(c *Combat, cutoff int64) {
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
		if time.Nanoseconds() > cutoff {
			log.Print("Cutoff in ComputeScore")
			ps.P[0].Score = ps.P[0].Score[:0]
			return
		}
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
		ps.P[0].Score[p[0]] += c.Score(dead, 0)
		c.Unresolve(move[ip][:ie])

		if Debug[DBG_Combat] {
			log.Print("perm ", ip, ":", p, ": ", DumpMoves(move[ip][:ie]))
			log.Print("DEAD: ", len(dead), ": ", DumpMoves(dead))
		}
	}
	if Debug[DBG_Combat] {
		log.Print("Score: ", ps.P[0].Score)
	}
}
