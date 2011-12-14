package combat

import (
	"log"
	"time"
	"rand"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
)

func (c *Combat) Run(ants map[Location]*AntStep, part Partitions, pmap PartitionMap, cutoff int64, rng *rand.Rand) {
	if len(part) == 0 {
		return
	}
	budget := (cutoff - time.Nanoseconds()) / int64(len(part)) / 2
	// TODO prioritize partition resolution...
	for {
		for ploc, ap := range part {
			if ap.PS == nil {
				ap.PS = NewPartitionState(c, ap)
				ap.PS.FirstStepRisk(c)
			}

			t := time.Nanoseconds() + budget
			if t > cutoff {
				log.Print("Out of time in Run combat")
				break
			}

			// c.Sim(ap, ploc, t, rng)
			log.Print("Scoring ", ploc)
			ap.PS.ComputeScore(c)

		}
		break
	}

	setMoves(ants, part, rng)
}

func setMoves(ants map[Location]*AntStep, part Partitions, rng *rand.Rand) {
	mm := make(map[Location]AntMove, 100)
	for ploc, ap := range part {
		// this can happen if we run out of time...
		if ap.PS != nil {
			ps := &ap.PS.P[0]

			best := ps.bestScore()
			if len(best) == len(ps.Score) && ps.Score[0] < 0 {
				ps.Best = -1
			} else {
				ps.Best = best[0]
			}
			log.Print(ploc, " best state is ", ps.Best)

			if ps.Best != -1 {
				for _, am := range ps.First[ps.Best] {
					log.Print(am.From, " move is ", am)
					mm[am.From] = am
				}
			}
		}
	}

	togoo := make(map[Location]struct{}, len(mm))
	for loc, move := range mm {
		if am, ok := ants[loc]; !ok {
			log.Print("Attempt to move an unfound ant", loc)
		} else {
			if _, found := togoo[move.To]; !found {
				am.Move = move.D
				am.Dest = append(am.Dest, move.To)
				am.Steps = append(am.Steps, 1)
				am.Steptot += 5 // MAGIC - ants in combat tend not to path to anything 
				am.Goalp = true
				am.Combatp = true
			} else {
				log.Print("COLLISION ", move.To)
			}
		}
	}
}

func (c *Combat) Sim(ap *AntPartition, ploc Location, cutoff int64, rng *rand.Rand) {
	log.Printf("Simulate for ap: %v %d ants, cutoff %.2fms",
		c.ToPoint(ploc),
		len(ap.Ants),
		float64(cutoff-time.Nanoseconds())/1e6)
	MonteSim(c, ap.PS, rng)
}

func MonteSim(c *Combat, ps *PartitionState, rng *rand.Rand) {
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
