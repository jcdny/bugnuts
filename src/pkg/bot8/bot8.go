// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

// The v8 Bot -- Now with combat (no really - not good combat but combat)
package bot8

import (
	"os"
	"fmt"
	"log"
	"time"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/state"
	. "bugnuts/MyBot"
	. "bugnuts/parameters"
	. "bugnuts/watcher"
	. "bugnuts/pathing"
	. "bugnuts/viz"
	. "bugnuts/game"
	. "bugnuts/util"
)

type BotV8 struct {
	P          *Parameters
	PriMap     *[256]int // array mapping Item to priority
	Explore    *TargetSet
	IdleAnts   []int
	StaticAnts []int
	RiskOff    bool
}

func init() {
	RegisterABot(ABot{Key: "v8", Desc: "V8 - combat bot", PKey: "v8", NewBot: NewBotV8})
}

//NewBot creates a new instance of your bot
func NewBotV8(s *State, pset *Parameters) Bot {
	if pset == nil {
		log.Panic("Nil parameter set")
	}

	mb := &BotV8{
		P:          pset,
		IdleAnts:   make([]int, 0, s.Turns+2),
		StaticAnts: make([]int, s.Turns+2),
	}

	mb.PriMap = mb.P.MakePriMap()

	if mb.P.Explore {
		mb.Explore = MakeExplorers(s, .8, 1, mb.PriMap[EXPLORE])
	} else {
		ts := make(TargetSet, 0)
		mb.Explore = &ts
	}

	return mb
}

func BlockEm(s *State) []Location {
	// bash in blocks for any location where enemies only have 1 degree of freedom...
	ablock := make([]Location, 0, 500)
	for loc, t := range s.C.Threat1 {
		if t > 7 && s.C.PThreat[0][loc] < 2 &&
			s.Map.Grid[loc] == LAND {
			s.Map.Grid[loc] = BLOCK
			ablock = append(ablock, Location(loc))
		}
	}

	return ablock
}

func ChopEm(s *State) {
	origin := make(map[Location]int, len(s.Ants[0]))
	for loc := range s.Ants[0] {
		origin[loc] = 1
	}
	f := NewFill(s.Map)
	// will only find neighbors withing 2x8 steps.
	f.MapFillSeed(origin, 1)

	for _, hill := range s.EnemyHillLocations(0) {
		if f.Depth[hill] != 0 {
			s.HChop[hill] = []Location{}, false
			s.HCRisk[hill] = 0, false
		} else {
			// a hill we cant reach
			brisk := 1000000
			var bpath []Location
			if nb, ok := s.HCRisk[hill]; ok {
				brisk = nb
				bpath = s.HChop[hill]
			}
			locs, _ := f.Sample(s.Rand, 20, 0, 100)
			for _, loc := range locs {
				if s.Met.EHill.Seed[loc] == hill {
					nrisk := 0
					path := s.Met.EHill.NPathInPath(s.Rand, loc, -1)

					for _, ploc := range path {
						nrisk += s.C.Threat[ploc] - s.C.PThreat[0][ploc]
					}
					if nrisk < brisk {
						bpath = path
						brisk = nrisk
					}
				}
			}
			if len(bpath) > 0 {
				// log.Print(s.ToPoints(bpath))
				s.HChop[hill] = bpath
				s.HCRisk[hill] = brisk
				for _, loc := range bpath {
					if s.Map.Grid[loc] == BLOCK {
						s.Map.Grid[loc] = LAND
					}
				}
			}
		}
	}
}

func (bot *BotV8) GenerateTargets(s *State) *TargetSet {
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

	// we have a ton of food targets when symmetry is high....
	if false && len(s.Map.SMap) > 0 && len(s.Map.SMap[0]) > 5 {
		ts := make(TargetSet, 0)
		bot.Explore = &ts
	} else if s.Turn > 150 && s.Turn%30 == 0 {
		ts := make(TargetSet, 0)
		bot.Explore = &ts
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

func (bot *BotV8) DoTurn(s *State) os.Error {
	bot.Explore.RemoveSeen(s, 1)

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

	// Lets combat a bit.
	ap, pmap, ctargets := s.C.Partition(s.Ants)
	for _, loc := range ctargets {
		tset.Add(RALLY, loc, 1, bot.PriMap[RALLY])
	}
	// Now visualize the frenemies.
	if Viz["combat"] {
		VizFrenemies(s, ap, pmap)
	}
	// s.C.Risk = s.C.Riskly(s.Ants) // done in setup now.
	s.C.Run(ants, ap, pmap, s.Cutoff, s.Rand)

	ablock := BlockEm(s)
	ChopEm(s)

	if Debug[DBG_Sample] {
		VizTiles(s, ablock, 0, 0, 255)
		log.Print("Set blocks on ", len(ablock), " cells ", s.ToPoints(ablock))
		for loc, chop := range s.HChop {
			VizTiles(s, chop, 0, 255, 0)
			log.Print("Set cuts to ", s.ToPoint(loc))
		}
	}

	iter := -1
	maxiter := 50
	nMove := 1
	for len(ants) > 0 && tset.Pending() > 0 && nMove != 0 {
		if time.Nanoseconds() > s.Cutoff {
			if Debug[DBG_Timeouts] {
				log.Print("Turn cutoff on iteration ", iter, " in DoTurn")
			}
			break
		}
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
			ant, aok := ants[seg.Src]
			tgt, ok := (*tset)[seg.End]
			if !aok || (!ok && seg.End != 0) {
				if Debug[DBG_MoveErrors] {
					log.Printf("Move from %v(%d) to %v(%d) no target ant: %#v",
						s.ToPoint(seg.Src), seg.Src, s.ToPoint(seg.End), seg.End, ant)
					log.Printf("Source item \"%v\", pending=%d", s.Map.Grid[seg.Src], tset.Pending())
				}
				if Viz["error"] {
					p := s.ToPoint(seg.Src)
					VizLine(s.Torus, p, s.ToPoint(seg.End), false)
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
							}
						}
					}
				}

				if WS.Watched(ant.Source, 0) {
					for i := 0; i < 5; i++ {
						log.Printf("TURN %d: %v -> %v (%s): \"%s\" steps %d: \"%s\":%#v",
							s.Turn, s.ToPoint(ant.Source), s.ToPoint(s.Map.LocStep[seg.Src][ant.N[i].D]),
							ant.N[i].D, tgt.Item, seg.Steps, ant.N[i].Item, ant.N[i])
					}
				}
				if good {
					// A good move exists so assume we step to the target
					if Viz["path"] {
						VizPath(s.ToPoint(seg.Src), f.NPathInString(nil, seg.Src, -1, 0), 1)
						VizPath(s.ToPoint(seg.Src), f.NPathInString(nil, seg.Src, -1, 1), 2)
					}
					if Viz["goals"] {
						VizLine(s.Torus, s.ToPoint(seg.Src), s.ToPoint(seg.End), false)
					}
					tgt.Count--
					nMove++
					ant.Goalp = true
					ant.Steps = append(ant.Steps, seg.Steps-ant.Steptot)
					ant.Dest = append(ant.Dest, seg.End)
					ant.Steptot = seg.Steps

					if tgt.Terminal || len(ant.Dest) > 15 || ant.Steptot > 50 {
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
				if idle > 4*len(eh) {
					newview := MinV(idle-2*len(eh), MaxV(s.Stats.LTS.Horizon/s.ViewRadius2-len(*bot.Explore), 2))
					if Debug[DBG_Targets] {
						log.Print("Adding explore points currently ", len(*bot.Explore), " looking to add ", newview, s.Size(), s.Stats.LTS.Horizon)
					}
					nadded = s.AddBorderTargets(newview, tset, bot.Explore, bot.PriMap[EXPLORE])
				}
				for _, loc := range eh {
					(*tset)[loc].Count += (idle - nadded) / len(eh)
					tc := MinV((s.Turns-s.Turn)/60, 3)
					(*tset)[loc].Count = MinV((*tset)[loc].Count, len(s.Ants[0])/(5-tc))
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
			tp[ant.Source] = 5 // MAGIC
		} else {
			tp[ant.Source] = 1
		}
	}

	// TODO walk away if N in > N enemy in.

	// Walk away from Hills - only if hill # < 3
	// Bloodbath maps best defense is fast gathering.
	if s.NHills[0] < 3 {
		fa, _, _ := MapFillSeed(s.Map, tp, 0)
		for _, ant := range endants {
			ant.N[4].PrFood = s.Met.ComputePrFood(ant.Source, ant.Source, s.Turn, &s.ViewMask.Offsets, fa)
			for d := 0; d < 4; d++ {
				ant.N[d].PrFood = s.Met.ComputePrFood(s.Map.LocStep[ant.Source][d], ant.Source, s.Turn, &s.ViewMask.Offsets, fa)
			}
			if !ant.Goalp && s.Met.FDownhill.Depth[ant.Source] > 1 {
				dh := int(s.Met.FHill.Depth[ant.Source])
				ant.Goalp = true
				ant.N[4].Goal = 0
				// TODO May need to set a dest as well
				ant.Steps = append(ant.Steps, dh)
				for d := Direction(0); d < 5; d++ {
					ant.N[d].Goal = s.Met.FDownhill.DistanceStep(ant.Source, ant.N[d].D)
					if WS.Watched(ant.Source, 0) {
						log.Printf("DOWNHILL: %v %s %#v", s.ToPoint(ant.Source), ant.N[d].D, ant.N[d])
					}
				}
			}

		}
	}

	s.GenerateMoves(endants)

	for _, ant := range endants {
		if WS.Watched(ant.Source, 0) {
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

	s.TurnDone()

	return nil
}
