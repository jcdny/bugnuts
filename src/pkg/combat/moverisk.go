package combat

import (
	"log"
	"sort"
	"fmt"
	. "bugnuts/game"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/util"
	. "bugnuts/watcher"
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
	free   int
	freed  int // for single degree of freedom this is the direction...
	nrisk  int // count of positions with risk.
	met    [5]rMet
}

func (ps *PartitionState) FirstStepRisk(c *Combat, full bool) int {
	TPush("@firststeprisk")
	defer TPop()

	tdepth := []int{0, 0, 65535, 65535}
	trisk := []int{RiskNeutral, RiskAverse, RiskNeutral, RiskAverse}
	antperm := 1

	for np := range ps.P {
		if ps.P[np].rs == nil {
			ps.P[np].rs, ps.P[np].davg, ps.P[np].dmin = riskmet(c, ps.P[np].Moves)
		}

		if (full && len(ps.P[np].Moves) <= 4) || len(ps.P[np].Moves) == 1 {
			risk := RiskNeutral
			if np != 0 || len(ps.P[np].Moves) == 1 {
				risk = Suicidal
			}
			ps.P[np].First = allMoves(ps.P[np].rs, risk, c)
			if Debug[DBG_Combat] {
				log.Print("generated ", len(ps.P[np].First), " moves for ", len(ps.P[np].rs), " ants")
			}
		}

		if len(ps.P[np].First) > 16 || len(ps.P[np].First) == 0 {
			tdepth[2] = ps.P[np].davg
			ps.P[np].First = make([][]AntMove, len(tdepth))
			for d := 0; d < len(tdepth); d++ {
				ps.P[np].First[d] = moveEmRisk(ps.P[np].rs, tdepth[d], c, trisk[d])
			}
		}
		if len(ps.P[np].Moves) > 0 && len(ps.P[np].First) > 0 {
			antperm *= len(ps.P[np].Moves) * len(ps.P[np].First)
		}
	}

	return antperm
}

func riskmet(c *Combat, ants []AntMove) (rs []rScore, davg, dmin int) {
	dtot := 0
	dmin = 65535
	rs = make([]rScore, len(ants))

	for i := range rs {
		r := &rs[i]
		am := &ants[i]
		r.am = am
		r.free = 4
		r.freed = 4
		for d := 0; d < 5; d++ {
			r.met[d].d = Direction(d)
			loc := c.Map.LocStep[am.From][d]
			r.met[d].to = loc
			if risk, ok := c.Risk[am.Player][loc]; ok {
				r.met[d].risk = risk
				r.nrisk++
			}
			r.met[d].depth = int(c.PFill[am.Player].Depth[loc])
			if c.Ants1[loc]&PlayerFlag[am.Player] != 0 {
				r.met[d].step = true
			}
			if d < 5 {
				if r.met[d].risk == Suicidal || !r.met[d].step {
					r.free--
				} else {
					r.freed = d
				}
			}
			// for orig loc we save depth for prioritizing moves
			if d == 4 {
				dtot += r.met[d].depth
				dmin = MinV(r.met[d].depth, dmin)
				r.depth = r.met[d].depth
			}
		}
	}
	if Debug[DBG_Combat] {
		log.Print("Partition Depth: dtot: ", dtot, " dmin: ", dmin, " Ants:", len(ants))
	}

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
	if p[i].nrisk != 0 && p[j].nrisk == 0 {
		return true
	}
	if p[i].target != p[j].target {
		return p[i].target < p[j].target
	}
	return p[i].target < p[j].target
}

// given a list of ants generates all permissible moves
// discard permutations that are not valid and remove duplicates
// also collapse non risk states
func allMoves(rs []rScore, risktype int, c *Combat) [][]AntMove {
	t := Torus{44, 84}
	// because of how I dedup this cant be > 5
	if len(rs) > 4 {
		return [][]AntMove{}
	}

	pl := Permutations(len(rs), 5)
	plv := make([][]int, 0, len(pl))

	dedup := make(map[int64]struct{}, len(pl))
	tolist := make([]Location, len(rs))

	if false {
		s := ""
		for _, r := range rs {
			s += fmt.Sprintf("%d ", r.am.From)
		}
	}
	for _, p := range pl {
		var m int
		for m = 0; m < len(p); m++ {
			d := p[m]
			if !rs[m].met[d].step || rs[m].met[d].risk > risktype {
				if Debug[DBG_Combat] {
					log.Print("** Drop ", m, d, rs[m].am.From, !rs[m].met[d].step, rs[m].met[d].risk > risktype)
				}
				break
			}
			c.AntCount[rs[m].met[d].to]++
			c.AntCount[rs[m].am.From]--
			if rs[m].met[d].risk != RiskNone {
				tolist[m] = rs[m].met[d].to
			} else {
				if Debug[DBG_Combat] {
					log.Print("Safe ", t.ToPoint(rs[m].met[d].to))
				}
				tolist[m] = 65535
			}
		}
		if m == len(p) {
			// check if to counts are all 1
			var mc int
			for mc = 0; mc < len(p); mc++ {
				d := p[mc]
				if c.AntCount[rs[mc].met[d].to] != 1 {
					if Debug[DBG_Combat] {
						log.Print("** Drop ", mc, p, p[mc], " collision ", rs[mc].met[d].to)
					}
					break
				}
			}
			if mc == len(p) {
				if len(p) > 1 {
					sort.Sort(LocationSlice(tolist))
					sig := int64(0)
					for _, l := range tolist {
						sig = (sig << 16) + int64(l)
					}
					if _, ok := dedup[sig]; !ok {
						plv = append(plv, p)
						dedup[sig] = struct{}{}
					} else {
						if Debug[DBG_Combat] {
							log.Print("Drop ", p, " dup ")
						}
					}
				} else {
					plv = append(plv, p)
				}
			}
		}
		// reset any ant we provisionally moved.
		for mr := 0; mr < m; mr++ {
			d := p[mr]
			c.AntCount[rs[mr].met[d].to]--
			c.AntCount[rs[mr].am.From]++
		}
	}
	// plv now contains the set of permissible permutations, gen 1st move set.
	// TODO remove swaps....

	moves := make([][]AntMove, len(plv))
	mbuf := make([]AntMove, len(rs)*len(plv))
	n := 0
	for pn, p := range plv {
		for i, d := range p {
			mbuf[n] = AntMove{From: rs[i].am.From, D: rs[i].met[d].d, To: rs[i].met[d].to, Player: rs[i].am.Player}
			n++
		}
		moves[pn] = mbuf[n-len(p) : n]
	}

	return moves
}

// MoveEm given a list of AntMove update D and To for the move in the given direction
func moveEmRisk(rs []rScore, tdepth int, c *Combat, risktype int) []AntMove {
	nrmove := 0
	for i := range rs {
		if rs[i].nrisk != 0 {
			nrmove++
		}
		rs[i].am.D = InvalidMove
		rs[i].am.To = -1
		rs[i].target = Abs(rs[i].depth - tdepth)
		for d := 0; d < 5; d++ {
			rs[i].met[d].target = Abs(rs[i].met[d].depth - tdepth)
			rs[i].met[d].netrisk = MaxV(rs[i].met[d].risk-risktype, 0)
		}

		if Debug[DBG_Combat] || WS.Watched(rs[i].am.From, rs[i].am.Player) {
			log.Print("Ant #", i, "/", len(rs), " ", rs[i].am.From, ": depth, tdepth, target ", rs[i].depth, rs[i].target)
			sort.Sort(rMetSlice(rs[i].met[:]))
			for d := 0; d < 5; d++ {
				log.Print("\t", rs[i].met[d])
			}
		}
	}

	sort.Sort(rScoreSlice(rs))

	var moved, nm int
	for nm, moved = 1, 0; nrmove > 0 && moved < len(rs) && nm != 0; nm = 0 {
		// TODO resort/shuffle per time through.  really need to do constraint propigagtion
		// but the constraint of the deadline is binding.
		for i := moved; nrmove > 0 && i < len(rs); i++ {
			sort.Sort(rMetSlice(rs[i].met[:]))
			am := rs[i].am
			if WS.Watched(am.From, am.Player) {
				log.Print(am.From, " best ", rs[i].met[0], c.AntCount[rs[i].met[0].to])
			}

			if rs[i].met[0].step == true && c.AntCount[rs[i].met[0].to] == 0 {
				rs[i].am.D = rs[i].met[0].d
				rs[i].am.To = rs[i].met[0].to
			}

			if am.D < NoMovement {
				if WS.Watched(am.From, am.Player) {
					log.Print(am.From, " moved ", am.To)
				}
				c.AntCount[am.To]++
				c.AntCount[am.From]--

				if rs[i].nrisk > 0 {
					nrmove--
				}

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

	moves := make([]AntMove, moved)
	for i := 0; i < moved; i++ {
		// log.Print(*rs[i].am)
		moves[i] = *rs[i].am
	}

	return moves
}
