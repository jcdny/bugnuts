package bot7

// The v7 Bot -- Now with combat (eventually)

import (
	"os"
	"fmt"
	"log"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/state"
	. "bugnuts/MyBot"
	. "bugnuts/parameters"
	. "bugnuts/debug"
	. "bugnuts/pathing"
	. "bugnuts/viz"
)

type BotV7 struct {
	P          *Parameters
	PriMap     *[256]int // array mapping Item to priority
	Explore    *TargetSet
	IdleAnts   []int
	StaticAnts []int
	RiskOff    bool
}

func init() {
	RegisterABot(ABot{Key: "v7", Desc: "V7 - combat bot", PKey: "v7", NewBot: NewBotV7})
}

//NewBot creates a new instance of your bot
func NewBotV7(s *State, pset *Parameters) Bot {
	if pset == nil {
		log.Panic("Nil parameter set")
	}

	mb := &BotV7{
		P:          pset,
		IdleAnts:   make([]int, 0, s.Turns+2),
		StaticAnts: make([]int, s.Turns+2),
	}

	mb.PriMap = mb.P.MakePriMap()

	if true {
		mb.Explore = MakeExplorers(s, .8, 1, mb.PriMap[EXPLORE])
	} else {
		ts := make(TargetSet, 0)
		mb.Explore = &ts
	}

	return mb
}

func (bot *BotV7) GenerateTargets(s *State) *TargetSet {
	tset := &TargetSet{}

	// TODO figure out how to handle multihill defense.
	if s.NHills[0] < 3 {
		s.AddEnemyPathinTargets(tset, bot.PriMap[DEFEND], bot.P.DefendDistance)
		s.AddMCBlock(tset, bot.PriMap[DEFEND], bot.P.DefendDistance)
	}

	// Generate list of food and enemy hill points.
	// Food locations should be set after ant list is done since we
	// remove adjacent food at that step.
	for _, loc := range s.FoodLocations() {
		tset.Add(FOOD, loc, 1, bot.PriMap[FOOD])
	}

	tset.Merge(bot.Explore)

	// TODO handle different priorities/attack counts
	// TODO compute defender count
	// TODO scale as function of my #s vs defnder.
	// TODO impute food counts for hill taking.
	eh := s.EnemyHillLocations(0)
	for _, loc := range eh {
		ndefend := 2 // s.PathinCount(loc, 10)
		if s.Turn > 1000 {
			ndefend = 200
		}
		tset.Add(HILL1, loc, ndefend+2, bot.PriMap[HILL1])
	}

	if s.NHills[0] < 3 {
		for _, loc := range s.Met.HBorder {
			depth := s.Met.FHill.Depth[loc]
			if depth > 2 && depth < uint16(bot.P.MinHorizon) {
				// Just add these as transients.
				tset.Add(WAYPOINT, loc, 1, bot.PriMap[WAYPOINT])
			}
		}
	}

	return tset
}

func (bot *BotV7) DoTurn(s *State) os.Error {

	bot.Explore.UpdateSeen(s, 1)
	tset := bot.GenerateTargets(s)

	if float64(bot.StaticAnts[s.Turn-1])/float64(len(s.Ants[0])) > bot.P.RiskOffThreshold {
		bot.RiskOff = true
	} else {
		bot.RiskOff = false
	}

	ants := s.GenerateAnts(tset, bot.RiskOff)

	endants := make([]*AntStep, 0, len(ants))
	segs := make([]Segment, 0, len(ants))

	if Viz["targets"] {
		VizTargets(s, tset)
	}

	iter := -1
	maxiter := 50
	nMove := 1
	for len(ants) > 0 && tset.Pending() > 0 && nMove != 0 {
		iter++
		nMove = 0

		// TODO: Here should update map for fixed ants.
		f, _, _ := MapFillSeed(s.Map, tset.Active(), 0)
		if Debug[DBG_Iterations] {
			log.Printf("TURN %d ITER %d TGT PENDING %d ANTS %d, ENDANTS %d", s.Turn, iter, tset.Pending(), len(ants), len(endants))
		}

		segs = segs[0:0]
		for loc := range ants {
			segs = append(segs, Segment{Src: loc, Steps: ants[loc].Steptot})
		}

		if !f.ClosestStep(segs) {
			// corner case: we added a guess or explore point which subsequently turned out to
			// be in a wall but the point has not become visible yet.
			segs = segs[0:0]
			for loc, tgt := range *tset {
				if tgt.Count > 0 {
					nMove++ // try another iteration
					tgt.Count = 0
					bot.Explore.Remove(loc)
				}
			}

		}

		for _, seg := range segs {
			ant := ants[seg.Src]
			tgt, ok := (*tset)[seg.End]
			if !ok && seg.End != 0 {
				if Debug[DBG_MoveErrors] {
					log.Printf("Move from %v(%d) to %v(%d) no target ant: %#v",
						s.ToPoint(seg.Src), seg.Src, s.ToPoint(seg.End), seg.End, ant)
					log.Printf("Source item \"%v\", pending=%d", s.Map.Grid[seg.Src], tset.Pending())
				}
				if Viz["error"] {
					p := s.ToPoint(seg.Src)
					VizLine(s.Map, p, s.ToPoint(seg.End), false)
					fmt.Fprintf(os.Stdout, "v tileBorder %d %d MM\n", p.R, p.C)
				}
			} else if ok && tgt.Count > 0 {
				// We have a target - make sure we can step in the direction of the target.
				good := true
				if ant.Steptot == 0 {
					// if it's a real step make sure there is something we would do
					good = false
					ant.N[4].Goal = 0
					dh := int(s.Met.FHill.Depth[seg.Src])
					for i := 0; i < 4; i++ {
						nloc := s.Map.LocStep[seg.Src][ant.N[i].D]
						// Don't mark target as taken unless its a valid step and risk = 0
						// TODO not sure this is how I should be doing this.
						goal := f.Distance(seg.Src, nloc) * 10

						// Prefer steps which are downhill from the hill if we are close.
						if dh > 0 && dh < 40 {
							goal += s.Met.FDownhill.Distance(seg.Src, nloc)
						}

						ant.N[i].Goal = goal
						// Check for a valid move towards the goal
						if s.Stepable(nloc) && goal > 0 {
							// and it needs to be a step we can take
							if ant.N[i].Safest {
								good = true
							} else if ant.N[i].Threat < 2 && seg.Steps < 20 &&
								(tgt.Item == DEFEND || tgt.Item.IsHill()) {
								good = true
								ant.N[i].Threat = 0 // TODO HACK!
							}
						}
					}
				}

				if WS.Watched(ant.Source, s.Turn, 0) {
					for i := 0; i < 5; i++ {
						log.Printf("TURN %d: %v -> %v (%s): \"%s\" steps %d: \"%s\":%#v",
							s.Turn, s.ToPoint(ant.Source), s.ToPoint(s.Map.LocStep[seg.Src][ant.N[i].D]),
							ant.N[i].D, tgt.Item, seg.Steps, ant.N[i].Item, ant.N[i])
					}
				}
				if good {
					// A good move exists so assume we step to the target
					if Viz["path"] {
						VizLine(s.Map, s.ToPoint(seg.Src), s.ToPoint(seg.End), false)
					}
					tgt.Count--
					nMove++
					ant.Goalp = true
					ant.Steps = append(ant.Steps, seg.Steps-ant.Steptot)
					ant.Dest = append(ant.Dest, seg.End)
					ant.Steptot = seg.Steps

					if tgt.Terminal {
						endants = append(endants, ant)
					} else {
						ants[seg.End] = ant
					}
					ants[seg.Src] = &AntStep{}, false
				}
			}
		}
		if len(bot.IdleAnts) < s.Turn+1 && (tset.Pending() < 1 || nMove == 0) {
			nMove++ // to loop again
			bot.IdleAnts = bot.IdleAnts[0 : s.Turn+1]
			// See if we have more ants than targets
			idle := 0
			for _, ant := range ants {
				if !ant.Goalp && ant.Steptot < 20 {
					idle++
				}
			}
			bot.IdleAnts[s.Turn] = idle

			if Debug[DBG_Iterations] {
				log.Printf("TURN %d IDLE %d", s.Turn, len(ants))
			}

			// we have N idle ants put the to work
			// Do this by marking ants with large basins as targets.
			if idle > 0 {
				maxiter = maxiter + 50
				eh := s.EnemyHillLocations(0)
				nadded := 0
				if idle > 2*len(eh) {
					nadded = s.AddBorderTargets(idle-2*len(eh), tset, bot.Explore, 1)
				}
				for _, loc := range eh {
					(*tset)[loc].Count += (idle - nadded) / len(eh)
				}
			}
		}
	}

	if Debug[DBG_Iterations] {
		log.Printf("TURN %d ITER %d ANTS %d END %d", s.Turn, iter, len(ants), len(endants))
		if iter == maxiter {
			for loc := range tset.Active() {
				log.Printf("Active: %v %v", s.ToPoint(loc), (*tset)[loc])
			}
		}
	}

	for _, ant := range ants {
		endants = append(endants, ant)
	}

	// Generate food prob given fill basins for existing ants.
	// also set outbound as goal for goalless ants near hill
	tp := make(map[Location]int, len(endants))
	for _, ant := range endants {
		if ant.Goalp {
			tp[ant.Source] = 5 // TODO more magic
		} else {
			tp[ant.Source] = 1
		}
	}

	// Walk away from Hills - only if hill # < 3
	// Bloodbath maps best defense is fast gathering.
	if s.NHills[0] < 3 {
		fa, _, _ := MapFillSeed(s.Map, tp, 0)
		for _, ant := range endants {
			ant.N[4].PrFood = s.Met.ComputePrFood(ant.Source, ant.Source, s.Turn, s.ViewMask, fa)
			for d := 0; d < 4; d++ {
				ant.N[d].PrFood = s.Met.ComputePrFood(s.Map.LocStep[ant.Source][d], ant.Source, s.Turn, s.ViewMask, fa)
			}
			if !ant.Goalp && s.Met.FDownhill.Depth[ant.Source] > 1 {
				dh := int(s.Met.FHill.Depth[ant.Source])
				ant.Goalp = true
				ant.N[4].Goal = 0
				// TODO May need to set a dest as well
				ant.Steps = append(ant.Steps, dh)
				for d := Direction(0); d < 5; d++ {
					ant.N[d].Goal = s.Met.FDownhill.DistanceStep(ant.Source, ant.N[d].D)
					if WS.Watched(ant.Source, s.Turn, 0) {
						log.Printf("DOWNHILL: %v %s %#v", s.ToPoint(ant.Source), ant.N[d].D, ant.N[d])
					}
				}
			}

		}
	}

	s.GenerateMoves(endants)
	for _, ant := range endants {
		if WS.Watched(ant.Source, s.Turn, 0) {
			log.Printf("ANT: %#v", ant)
			for d := Direction(0); d < 5; d++ {
				log.Printf("MOVE: %v %s d:%s :: %#v", s.ToPoint(ant.Source), ant.Move, ant.N[d].D, ant.N[d])
			}
		}
		if ant.Move > 3 || ant.Move < 0 {
			bot.StaticAnts[s.Turn]++
		}
	}
	s.EmitMoves(endants)
	Visualize(s)
	fmt.Fprintf(os.Stdout, "go\n") // TODO Flush ??

	return nil
}
