package combat

import (
	"log"
	"time"
	"rand"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/util"
)

func (c *Combat) Run(ants []*AntStep, part Partitions, pmap PartitionMap, cutoff int64, rng *rand.Rand) {
	if len(part) == 0 {
		return
	}

	budget := (cutoff - time.Nanoseconds()) / int64(len(part)) / 4
	// sim to compute best moves
	// c2 := c.Copy() // Debug state not changing...
	for {
		for ploc, ap := range part {
			t := time.Nanoseconds() + budget
			if t > cutoff {
				break
			}
			c.Sim(ap, ploc, t, rng)
			/* 
				 eq, diffs := CombatCheck(c, c2)
				if !eq {
					log.Print("Not equal: ", diffs)
				}*/
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
