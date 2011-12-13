package combat

import (
	"log"
	"sort"
	. "bugnuts/game"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/util"
	. "bugnuts/debug"
)

type rMet struct {
	d       Direction
	to      Location
	step    bool
	risk    int
	depth   int
	netrisk int // min(0,risk - risktype)
	target  int // Abs(depth - tr)
	perm    int // permuter
}

type rScore struct {
	am     *AntMove
	depth  int
	target int
	perm   int
	met    [5]rMet
}

func (ps *PartitionState) FirstStepRisk(c *Combat) {
	tdepth := []int{0, 0, 0, 65535, 65535}
	trisk := []int{Suicidal, RiskNeutral, RiskAverse, RiskAverse, RiskAverse}
	for np := range ps.P {
		rs, davg := riskmet(c, ps.P[np].Moves)
		tdepth[2] = davg
		ps.P[np].First = make([][]AntMove, len(tdepth))
		for d := 0; d < len(tdepth); d++ {
			ps.P[np].First[d] = moveEmRisk(rs, tdepth[d], c, trisk[d])
		}
	}
}

func riskmet(c *Combat, ants []AntMove) (rs []rScore, davg int) {
	dtot := 0
	dmin := 65535
	rs = make([]rScore, len(ants))

	for i := range rs {
		r := &rs[i]
		am := &ants[i]
		r.am = am
		for d := 0; d < 5; d++ {
			r.met[d].d = Direction(d)
			loc := c.Map.LocStep[am.From][d]
			r.met[d].to = loc
			if risk, ok := c.Risk[am.Player][loc]; ok {
				r.met[d].risk = risk
			}
			r.met[d].depth = int(c.PFill[am.Player].Depth[loc])
			if c.Ants1[loc]&PlayerFlag[am.Player] != 0 {
				r.met[d].step = true
			}
			// for orig loc we save depth for prioritizing moves
			if d == 4 {
				dtot += r.met[d].depth
				dmin = MinV(r.met[d].depth, dmin)
				r.depth = r.met[d].depth
			}
		}
	}
	log.Print("dtot: ", dtot, " dmin: ", dmin, " len:", len(ants))

	davg = dtot / len(ants)

	return
}

type rMetSlice []rMet

func (p rMetSlice) Len() int      { return len(p) }
func (p rMetSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p rMetSlice) Less(i, j int) bool {
	if p[i].step != p[j].step {
		return p[i].step
	}
	if p[i].netrisk != p[j].netrisk {
		return p[i].netrisk < p[j].netrisk
	}
	if p[i].target != p[j].target {
		return p[i].target < p[j].target
	}
	if p[i].risk != p[j].risk {
		return p[i].risk < p[j].risk
	}
	return p[i].perm < p[j].perm
}

type rScoreSlice []rScore

func (p rScoreSlice) Len() int      { return len(p) }
func (p rScoreSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p rScoreSlice) Less(i, j int) bool {
	if p[i].target != p[j].target {
		return p[i].target < p[j].target
	}
	return p[i].target < p[j].target

}

// MoveEm given a list of AntMove update D and To for the move in the given direction
func moveEmRisk(rs []rScore, tr int, c *Combat, risktype int) []AntMove {
	for i := range rs {
		rs[i].am.D = NoMovement
		rs[i].am.To = rs[i].am.From
		rs[i].target = Abs(rs[i].depth - tr)
		for d := 0; d < 5; d++ {
			rs[i].met[d].target = Abs(rs[i].met[d].depth - tr)
			rs[i].met[d].netrisk = MaxV(rs[i].met[d].risk-risktype, 0)
		}
		log.Print(rs[i].am.From, ": tr ", tr, " risktype ", risktype, " depth ", rs[i].depth, rs[i].target)
		sort.Sort(rMetSlice(rs[i].met[:]))
		for d := 0; d < 5; d++ {
			log.Print("\t", rs[i].met[d])
		}
	}

	var moved, nm int
	for nm, moved = 1, 0; moved < len(rs) && nm != 0; nm = 0 {
		// TODO resort/shuffle per time through.  really need to do constraint propigagtion
		// but the constraint of the deadline is binding.
		for i := moved; i < len(rs); i++ {
			sort.Sort(rMetSlice(rs[i].met[:]))
			am := rs[i].am
			if WS.Watched(am.From, -1, am.Player) {
				log.Print(am.From, " best ", rs[i].met[0], c.AntCount[rs[i].met[0].to])
			}

			if rs[i].met[0].step == true && c.AntCount[rs[i].met[0].to] == 0 {
				rs[i].am.D = rs[i].met[0].d
				rs[i].am.To = rs[i].met[0].to
			}

			if am.D != NoMovement {
				if WS.Watched(am.From, -1, am.Player) {
					log.Print(am.From, " moved ", am.To)
				}
				c.AntCount[am.To]++
				c.AntCount[am.From]--

				if moved < i {
					rs[moved], rs[i] = rs[i], rs[moved]
				}

				moved++
				nm++
			}
		}
	}

	// Reset the player counts for the moved
	for i := 0; i < moved; i++ {
		c.AntCount[rs[i].am.From]++
		c.AntCount[rs[i].am.To]--
	}

	moves := make([]AntMove, len(rs))
	for i := range rs {
		moves[i] = *rs[i].am
	}

	return moves
}
