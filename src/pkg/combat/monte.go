package combat

import (
	"log"
	"time"
	"rand"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/watcher"
)

func (c *Combat) Run(ants map[Location]*AntStep, part Partitions, pmap PartitionMap, cutoff int64, rng *rand.Rand) {
	TPush("@combat")
	defer TPop()

	if len(part) == 0 {
		return
	}
	// count active partitions
	N := 0
	for _, ap := range part {
		if len(ap.Ants) > 0 {
			N++
		}
	}
	if N == 0 {
		return
	}

	for ploc, ap := range part {
		if Debug[DBG_Combat] {
			log.Print("Starting processing for partition ", ploc, " len(ap.Ants) ", len(ap.Ants))
		}
		if len(ap.Ants) == 0 {
			continue
		}

		ap.PS = NewPartitionState(c, ap)
		ap.PS.FirstStepRisk(c, 0)
	}

	// compute the per partition time budget
	budget := (cutoff - time.Nanoseconds() - 30*MS) / int64(N*3) / 2
	if budget < 0 {
		return
	}

	// TODO prioritize partition resolution...
	for {
		for ploc, ap := range part {
			if Debug[DBG_Combat] {
				log.Print("Starting processing for partition ", ploc, " len(ap.Ants) ", len(ap.Ants))
			}
			t := time.Nanoseconds() + budget
			if t > cutoff {
				now := time.Nanoseconds()
				if Debug[DBG_Timeouts] {
					log.Print("Out of time in Run combat parts, cutoff, budget (ms):",
						len(part), cutoff-now/1000000, budget/1000000)
				}
				break
			}

			if Debug[DBG_Combat] {
				log.Print("****************************** Scoring ", ploc)
			}
			ap.PS.ComputeScore(c)
		}
		break
	}

	setMoves(ants, part, rng)
}

func setMoves(ants map[Location]*AntStep, part Partitions, rng *rand.Rand) {
	mm := make(map[Location]AntMove, 100)
	mp := make(map[Location]Location, 100)
	for ploc, ap := range part {
		if ap.PS == nil || len(ap.PS.P) == 0 || len(ap.PS.P[0].Score) == 0 {
			// Skip ones with no results...
			// this happens if we had no combat ants in the partition,
			// or possibly ran out of time simulating.
			continue
		}
		ps := &ap.PS.P[0]

		best := ps.bestScore()
		if len(best) == len(ps.Score) && ps.Score[0] < 0 {
			ps.Best = -1
		} else {
			ps.Best = best[0]
		}
		if Debug[DBG_Combat] {
			log.Print(ploc, " best state is ", ps.Best)
		}

		if ps.Best != -1 {
			for _, am := range ps.First[ps.Best] {
				if Debug[DBG_Combat] {
					log.Print(am.From, " move is ", am)
				}
				mm[am.From] = am
				mp[am.From] = ploc
			}
		}
	}

	togoo := make(map[Location]struct{}, len(mm))
	for loc, move := range mm {
		if am, ok := ants[loc]; !ok {
			log.Print("Attempt to move an unfound ant", loc)
		} else {
			if _, found := togoo[move.To]; !found {
				am.N[move.D].Combat = 1
				am.Dest = append(am.Dest, move.To)
				am.Steps = append(am.Steps, 1)
				am.Steptot += 5 // MAGIC - ants in combat tend not to path to anything
				am.Goalp = true
				am.Combat = mp[loc]
				togoo[move.To] = struct{}{}
			} else {
				log.Print("COLLISION ", move.To)
			}
		}
	}
}

func (c *Combat) xSim(ap *AntPartition, ploc Location, cutoff int64, rng *rand.Rand) {
	log.Printf("Simulate for ap: %v %d ants, cutoff %.2fms",
		c.ToPoint(ploc),
		len(ap.Ants),
		float64(cutoff-time.Nanoseconds())/1e6)
	xMonteSim(c, ap.PS, rng)
}

func xMonteSim(c *Combat, ps *PartitionState, rng *rand.Rand) {
	ps.ComputeScore(c)
	best := ps.P[0].bestScore()
	if len(best) == len(ps.P[0].Score) {
		ps.P[0].Best = 1
	} else {
		ps.P[0].Best = best[rng.Intn(len(best))]
	}
}

// generate the list of permuted directions for n players
func genPerm4(n uint) [][]Direction {
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
