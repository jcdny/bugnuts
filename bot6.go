package main
// The v6 Bot -- Now Officially not terrible
//
// Lesons from v5:
// The "Explore" concept was a failure.
//
// Need to be smarter about target priority
//
// Need to track chicken bots vs aggressive bots.
//
// Need to guess hills

import (
	"os"
	"fmt"
	"log"
)

type BotV6 struct {
	P        *Parameters
	Primap   []int // array mapping Item to priority
	Explore  *TargetSet
	IdleAnts []int
}

func (bot *BotV6) Priority(i Item) int {
	return bot.Primap[i]
}

//NewBot creates a new instance of your bot
func NewBotV6(s *State) Bot {
	if paramKey == "" {
		paramKey = "V6"
	}
	if _, ok := ParameterSets[paramKey]; !ok {
		log.Panicf("Unknown parameter key %s", paramKey)
	}

	mb := &BotV6{
		P:        ParameterSets[paramKey],
		IdleAnts: make([]int, 0, s.Turns+2),
	}

	mb.Primap = mb.P.MakePriMap()

	if true {
		mb.Explore = MakeExplorers(s, .8, 1, mb.Priority(EXPLORE))
	} else {
		ts := make(TargetSet, 0)
		mb.Explore = &ts
	}
	return mb
}

func (bot *BotV6) GenerateTargets(s *State) *TargetSet {
	tset := &TargetSet{}

	// TODO figure out how to handle multihill defense.
	if s.NHills[0] < 3 {
		s.AddEnemyPathinTargets(tset, bot.Priority(DEFEND), bot.P.DefendDistance)
		s.AddMCBlock(tset, bot.Priority(DEFEND), bot.P.DefendDistance)
	}

	// Generate list of food and enemy hill points.
	// Food locations should be set after ant list is done since we
	// remove adjacent food at that step.
	for _, loc := range s.FoodLocations() {
		tset.Add(FOOD, loc, 1, bot.Priority(FOOD))
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
		tset.Add(HILL1, loc, ndefend+2, bot.Priority(HILL1))
	}

	if s.NHills[0] < 3 {
		for _, loc := range s.Map.HBorder {
			depth := s.Map.FHill.Depth[loc]
			if depth > 2 && depth < uint16(bot.P.MinHorizon) {
				// Just add these as transients.
				tset.Add(WAYPOINT, loc, 1, bot.Priority(WAYPOINT))
			}
		}
	}

	return tset
}

func (bot *BotV6) DoTurn(s *State) os.Error {
	// TODO this still seems clunky.  need to figure where this belongs.
	s.FoodUpdate(bot.P.ExpireFood)
	bot.Explore.UpdateSeen(s, 1)

	tset := bot.GenerateTargets(s)
	ants := s.GenerateAnts(tset)

	endants := make([]*AntStep, 0, len(ants))
	segs := make([]Segment, 0, len(ants))

	if Viz["targets"] {
		s.VizTargets(tset)
	}

	var iter, maxiter int = 0, 50
	for iter = 0; iter < maxiter && len(ants) > 0 && tset.Pending() > 0; iter++ {
		// TODO: Here should update map for fixed ants.
		f, _, _ := MapFillSeed(s.Map, tset.Active(), 0)
		if Debug[DBG_Iterations] {
			log.Printf("TURN %d ITER %d TGT PENDING %d ANTS %d, ENDANTS %d", s.Turn, iter, tset.Pending(), len(ants), len(endants))
			log.Printf("ACTIVE SET: %v", tset.Active())
		}

		segs = segs[0:0]
		for loc, _ := range ants {
			segs = append(segs, Segment{src: loc, steps: ants[loc].steptot})
		}

		if !f.ClosestStep(segs) {
			// corner case: we added a guess or explore point which subsequently turned out to
			// be in a wall but the point has not become visible yet.
			segs = segs[0:0]
			for loc, tgt := range *tset {
				if tgt.Count > 0 {
					tgt.Count = 0
					bot.Explore.Remove(loc)
				}
			}

		}

		for _, seg := range segs {
			ant := ants[seg.src]
			tgt, ok := (*tset)[seg.end]
			if !ok && seg.end != 0 {
				if Debug[DBG_MoveErrors] {
					log.Printf("Move from %v(%d) to %v(%d) no target ant: %#v",
						s.ToPoint(seg.src), seg.src, s.ToPoint(seg.end), seg.end, ant)
					log.Printf("Source item \"%v\", pending=%d", s.Map.Grid[seg.src], tset.Pending())
				}
				if Viz["error"] {
					p := s.ToPoint(seg.src)
					VizLine(s.Map, p, s.ToPoint(seg.end), false)
					fmt.Fprintf(os.Stdout, "v tileBorder %d %d MM\n", p.r, p.c)
				}
			} else if ok && tgt.Count > 0 {
				// We have a target - make sure we can step in the direction of the target.
				good := true
				if ant.steptot == 0 {
					// if it's a real step make sure there is something we would do
					good = false
					ant.N[4].goal = 0
					dh := int(s.Map.FHill.Depth[seg.src])
					for i := 0; i < 4; i++ {
						nloc := s.Map.LocStep[seg.src][i]
						// Don't mark target as taken unless its a valid step and risk = 0
						// TODO not sure this is how I should be doing this.
						goal := f.Distance(seg.src, nloc) * 10

						// Prefer steps which are downhill from the hill if we are close.
						if dh > 0 && dh < 40 {
							goal += s.Map.FDownhill.Distance(seg.src, nloc)
						}

						ant.N[i].goal = goal
						// Check for a valid move towards the goal
						if WS.Watched(ant.source, s.Turn, 0) {
							log.Printf("TURN %d: %v->%v : %v goal:%d DHill:%d \"%s\" %d: %#v",
								s.Turn, s.ToPoint(ant.source), s.ToPoint(nloc), s.Stepable(nloc),
								goal, dh, tgt.Item, seg.steps, ant.N[i])
						}
						if s.Stepable(nloc) && goal > 0 {
							// and it needs to be a step we can take
							if ant.N[i].safest {
								good = true
							} else if ant.N[i].threat < 2 && seg.steps < 20 &&
								(tgt.Item == DEFEND || tgt.Item.IsHill()) {
								good = true
								ant.N[i].threat = 0 // TODO HACK!
							}
						}
					}
				}

				if good {
					// A good move exists so assume we step to the target
					if Viz["path"] {
						VizLine(s.Map, s.ToPoint(seg.src), s.ToPoint(seg.end), false)
					}
					tgt.Count--
					ant.goalp = true
					ant.steps = append(ant.steps, seg.steps-ant.steptot)
					ant.dest = append(ant.dest, seg.end)
					ant.steptot = seg.steps

					if tgt.Terminal {
						endants = append(endants, ant)
					} else {
						ants[seg.end] = ant
					}
					ants[seg.src] = &AntStep{}, false
				}
			}
		}

		if len(bot.IdleAnts) < s.Turn+1 && (tset.Pending() < 1 || iter+1 == maxiter) {
			bot.IdleAnts = bot.IdleAnts[0 : s.Turn+1]
			// See if we have more ants than targets
			idle := 0
			for _, ant := range ants {
				if !ant.goalp && ant.steptot < 20 {
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
				maxiter = 2 * maxiter
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
			for loc, _ := range tset.Active() {
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
		if ant.goalp {
			tp[ant.source] = 5 // TODO more magic
		} else {
			tp[ant.source] = 1
		}
	}

	// Walk away from Hills - only if hill # < 3
	// Bloodbath maps best defense is fast gathering.
	if s.NHills[0] < 3 {
		fa, _, _ := MapFillSeed(s.Map, tp, 0)
		for _, ant := range endants {
			ant.N[4].prfood = s.Map.ComputePrFood(ant.source, ant.source, s.Turn, s.viewMask, fa)
			for d := 0; d < 4; d++ {
				ant.N[d].prfood = s.Map.ComputePrFood(s.Map.LocStep[ant.source][d], ant.source, s.Turn, s.viewMask, fa)
			}
			if !ant.goalp && s.Map.FDownhill.Depth[ant.source] > 1 {
				dh := int(s.Map.FHill.Depth[ant.source])
				ant.goalp = true
				ant.N[4].goal = 0
				// TODO May need to set a dest as well
				ant.steps = append(ant.steps, dh)
				for d := Direction(0); d < 4; d++ {
					ant.N[d].goal = s.Map.FDownhill.DistanceStep(ant.source, d)
					//log.Printf("DOWNHILL: %v %s %#v", s.ToPoint(ant.source), d, ant.N[d])
				}
			}

		}
	}

	s.GenerateMoves(endants)
	s.EmitMoves(endants)
	s.Viz()
	fmt.Fprintf(os.Stdout, "go\n") // TODO Flush ??

	return nil
}
