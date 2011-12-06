package state

import (
	"log"
	"math"
	"fmt"
	"sort"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/debug"
	. "bugnuts/pathing"
	. "bugnuts/util"
)

type Target struct {
	Item     Item
	Loc      Location // Target Location
	Count    int      // how many do we want at this location
	Pri      int      // target priority.
	Terminal bool     // Is it a terminating target 

	arrivals []int      // Inbound Arrival time
	player   []int      // Inbound player
	ant      []Location // Inbound Source
}

type TargetSet map[Location]*Target

func (tset *TargetSet) String() string {
	str := ""
	for loc, target := range *tset {
		str += fmt.Sprintf("%d: %#v\n", loc, target)
	}
	return str
}

func (tset *TargetSet) Add(item Item, loc Location, count, pri int) {
	if Debug[DBG_Targets] {
		log.Printf("Adding target: item %v loc %v count %d pri %d", item, loc, count, pri)
	}
	t, found := (*tset)[loc]
	if pri < 1 {
		log.Panicf("Target pri must be > 1")
	}
	if !found || t.Pri < pri {
		// We already have this point in the target set, replace if pri is higher
		(*tset)[loc] = &Target{
			Item:     item,
			Loc:      loc,
			Count:    count,
			Pri:      pri,
			Terminal: TerminalItem[item],
		}
	}

}

func (tset *TargetSet) Remove(loc Location) {
	t, ok := (*tset)[loc]
	if Debug[DBG_Targets] {
		if ok {
			log.Printf("Removing target: item %v loc %v count %d pri %d", t.Item, t.Loc, t.Count, t.Pri)
		} else {
			log.Printf("Removing target: not found Loc %d", loc)
		}
	}
	if ok {
		// We already have this point in the target set, replace if pri is higher
		(*tset)[loc] = t, false
	}
}

func (tset TargetSet) Merge(src *TargetSet) {
	if src == nil {
		return
	}
	// Run through explore targets
	for loc, tgt := range *src {
		nt, found := tset[loc]
		if !found || nt.Pri > tgt.Pri {
			tset[loc] = tgt
		}
	}
}

func (tset *TargetSet) Pending() int {
	n := 0
	for _, t := range *tset {
		if t.Count > 0 {
			n++
		}
	}

	return n
}

func (tset *TargetSet) Active() map[Location]int {
	tp := make(map[Location]int, tset.Pending())
	for _, t := range *tset {
		if t.Count > 0 {
			tp[t.Loc] = t.Pri
		}
	}

	return tp
}

func MakeExplorers(s *State, scale float64, count, pri int) *TargetSet {

	// Set an initial group of targetpoints
	if scale <= 0 {
		scale = 1.0
	}

	rstride := int(math.Sqrt(float64(s.ViewRadius2)) * 3 / 2 * scale)
	cstride := int(math.Sqrt(float64(s.ViewRadius2)*3) * scale)

	tset := make(TargetSet, (s.Rows * s.Cols / (rstride * cstride)))
	for r := 0; r < s.Rows; r += rstride {
		for c := cstride / 2 * (r / rstride % 2); c < s.Cols; c += cstride {
			loc := s.ToLocation(Point{R: r, C: c})
			tset.Add(EXPLORE, loc, count, pri)
		}
	}

	return &tset
}

func (s *State) AddBalanceTragets(N int, tset *TargetSet, explore *TargetSet, pri int) {
	f, _, _ := MapFillSeed(s.Map, s.Ants[0], 1)
	basins := make(map[Location]int, len(s.Ants[0])+10)
	for i, loc := range f.Seed {
		if f.Depth[i] > 0 {
			basins[loc]++
		}
	}
	sc := make([]DefScore, 0, len(basins))
	for loc, score := range basins {
		sc = append(sc, DefScore{loc: loc, score: score})
	}
	sort.Sort(DefScoreSlice(sc))
	for i := 0; i < N/2; i++ {
		(*tset).Add(EXPLORE, sc[i].loc, 2, pri)
	}
}

func (s *State) AddBorderTargets(N int, tset *TargetSet, explore *TargetSet, pri int) int {
	// Generate a target list for unseen areas and exploration
	// tset.Add(RALLY, s.Map.ToLocation(Point{58, 58}), len(ants), bot.Priority(RALLY))
	fexp, _, _ := MapFill(s.Map, s.Ants[0], 1)
	loc, n := fexp.Sample(s.Rand, N, 14, 20)
	added := 0
	for i := range loc {
		if s.Met.Seen[loc[i]] < s.Turn-1 {
			if Debug[DBG_BorderTargets] {
				log.Printf("Adding %d", i)
				log.Printf("Adding %d %v %v", i, s.ToPoint(loc[i]), n[i])
			}

			/* 
				     exp := s.ToPoint(loc[i])
					 if Viz["targets"] {
							fmt.Fprintf(os.Stdout, "v star %d %d .5 1.5 9 true\n", exp.R, exp.C)
					}
			*/
			if explore != nil {
				(*explore).Add(EXPLORE, loc[i], 1, pri)
			}
			(*tset).Add(EXPLORE, loc[i], 1, pri)
			added++
		}
	}
	return added
}

type DefScore struct {
	loc   Location
	score int
}
type DefScoreSlice []DefScore

func (p DefScoreSlice) Len() int           { return len(p) }
func (p DefScoreSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p DefScoreSlice) Less(i, j int) bool { return p[i].score > p[j].score }

func (s *State) AddMCBlock(tset *TargetSet, priority int, DefendDist int) {
	if true || len(s.Met.MCDist) == 0 || s.NHills[0] > 2 || s.Turn < 30 {
		return
	}

	hills := make(map[Location]int, 1)
	for _, loc := range s.HillLocations(0) {
		hills[loc] = 1

		f, _, _ := MapFill(s.Map, hills, 0)

		loclist, _ := f.Sample(s.Rand, 0, 2, 8)
		Def := make([]DefScore, len(loclist))
		for i, loc := range loclist {
			(*tset).Remove(loc)
			//log.Printf("DIST: %d %d", loc, len(s.Met.MCDist))
			Def[i] = DefScore{loc: loc, score: s.Met.MCDist[loc]}
		}
		sort.Sort(DefScoreSlice(Def))
		for i := 0; i < MinV(4, len(Def)); i++ {
			if Def[i].score/s.Met.MCPaths > 2 {
				(*tset).Add(DEFEND, Def[i].loc, 1, priority)
			}
		}
		hills[loc] = 0, false
	}
}

func (s *State) AddEnemyPathinTargets(tset *TargetSet, priority int, DefendDist int) {
	hills := make(map[Location]int, 6)
	for _, loc := range s.HillLocations(0) {
		hills[loc] = 1
	}

	f, _, _ := MapFill(s.Map, hills, 0)

	for i := 1; i < len(s.Ants); i++ {
		for loc := range s.Ants[i] {
			// TODO: use seed rather than PathIn
			steps := int(f.Depth[Location(loc)] - 1)
			if steps < DefendDist {
				tloc, _ := f.NPathIn(s.Rand, loc, MaxV(steps-8, steps/2))
				if Debug[DBG_PathInDefense] {
					log.Printf("Enemy Pathin: defense: %v @ %v", s.ToPoint(loc), s.ToPoint(tloc))
				}
				(*tset).Add(DEFEND, tloc, 1, priority)
				if len(s.Met.MCDist) > 0 {
					maxf := s.Met.MCFlow[tloc][0]
					d := 0
					for i := 1; i < 4; i++ {
						if s.Met.MCFlow[tloc][i] > maxf {
							d = i
						}
					}
					dirs := [2][2]Direction{{1, 3}, {0, 2}}
					for _, da := range dirs[d%2] {
						nl := s.Map.LocStep[tloc][da]
						if s.Map.Grid[nl] != WATER {
							(*tset).Add(DEFEND, nl, 1, priority)
						}
						if Debug[DBG_PathInDefense] {
							log.Printf("Maxflow %v: %s, adding dirs %v %v", s.Met.MCFlow[tloc], Direction(d), dirs[d%2], s.ToPoint(nl))
						}
					}
				}
			}
		}
	}
}

func (tset *TargetSet) RemoveSeen(s *State, count int) {
	for loc := range *tset {
		if s.Met.Seen[loc] == s.Turn {
			(*tset).Remove(loc)
		} else {
			(*tset)[loc].Count = count
		}
	}
}
