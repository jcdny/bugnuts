// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

// The v5 Bot -- Marginally less Terrible!!!!
// Leson from v4: if you run out of goals don't have your ants just go
// to sleep.
//
// v5 adds priority for goals, a collection of exploration
// goals, and adding more explore goals if we have idle ants.
//
// Also adds enemy ant avoidance though still no combat besides the
// willingness to sacrfice.
//
// Does not match v5 on aichallenge.org since I uploaded before I wanted to
// make a BotV6.go but the git tag matches what was uploaded.
package bot5

import (
	"fmt"
	"os"
	"log"
	. "bugnuts/maps"
	. "bugnuts/pathing"
	. "bugnuts/torus"
	. "bugnuts/state"
	. "bugnuts/parameters"
	. "bugnuts/MyBot"
	. "bugnuts/watcher"
	. "bugnuts/viz"
)

type BotV5 struct {
	P *Parameters

	Explore *TargetSet

	IdleAnts []int
	PriMap   *[256]int
}

func init() {
	RegisterABot(ABot{Key: "v5", Desc: "V5 - goal seeker", PKey: "v5", NewBot: NewBotV5})
}

func (bot *BotV5) Priority(i Item) int {
	return bot.PriMap[i]
}

//NewBot creates a new instance of your bot
func NewBotV5(s *State, p *Parameters) Bot {
	mb := &BotV5{
		P:        p,
		IdleAnts: make([]int, 0, s.Turns),
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

func (bot *BotV5) ExploreUpdate(s *State) {
	// Any explore point which is visible should be nuked
	for loc := range *bot.Explore {
		if s.Met.Seen[loc] == s.Turn {
			bot.Explore.Remove(loc)
		} else {
			(*bot.Explore)[loc].Count = 1
		}
	}
}

func (bot *BotV5) DoTurn(s *State) os.Error {
	bot.ExploreUpdate(s)

	tset := TargetSet{}

	tset.Merge(bot.Explore)

	// Generate list of food and enemy hill points.
	for _, loc := range s.FoodLocations() {
		if Debug[DBG_Targets] {
			log.Printf("adding target %v(%d) food pri %d", s.Map.ToPoint(loc), loc, bot.Priority(FOOD))
		}
		tset.Add(FOOD, loc, 1, bot.Priority(FOOD))
	}

	// TODO handle different priorities for different enemies.
	eh := s.EnemyHillLocations(0)
	idle := 0
	if len(bot.IdleAnts) > s.Turn-1 && len(eh) > 0 {
		//log.Printf("IDLE: %d : %d : %d : %v", s.Turn, len(bot.IdleAnts), len(eh), bot.IdleAnts)
		idle = bot.IdleAnts[s.Turn-1] / len(eh)
		if idle > 15 {
			idle = 15
		}
	}

	for _, loc := range eh {
		tset.Add(HILL1, loc, 1+idle, bot.Priority(HILL1))
	}

	if Debug[DBG_Targets] {
		log.Printf("Target set %v", tset)
	}

	// List of available ants
	ants := make(map[Location]int, len(s.Ants[0]))
	for k, v := range s.Ants[0] {
		ants[k] = v
	}

	// TODO remove dedicated ants eg sentinel, capture, defense guys

	moves := make(map[Location]Direction, len(ants))

	for iter := 0; iter < 15 && len(ants) > 0 && tset.Pending() > 0; iter++ {
		if Debug[DBG_Iterations] {
			log.Printf("Location iteration %d, ants: %d, tset.Pending %d", iter, len(ants), tset.Pending())
		}

		// TODO: Here should update map for fixed ants.
		f, _, _ := MapFill(s.Map, tset.Active(), 0)

		// Build list of locations sorted by depth
		ccl := make([]Location, 0, len(ants))
		for loc := range ants {
			ccl = append(ccl, loc)
		}
		for _, loc := range f.Closest(ccl) {
			depth := f.Depth[loc]
			threat := s.C.Threat1[loc] - s.C.PThreat1[0][loc]
			p := s.Map.ToPoint(loc)
		STEP:
			// Perm here so our bots are not biased to move in particular directions
			for _, d := range Permute4(s.Rand) {
				if Debug[DBG_Iterations] {
					log.Printf("Allocating ant %d dir %d", loc, d)
				}
				// find a direction we can step in thats stepable.
				np := s.Map.PointAdd(p, Steps[d])
				nl := s.Map.ToLocation(np)
				item := s.Map.Grid[nl]
				nthreat := s.C.Threat1[nl] - s.C.PThreat1[0][nl]
				// TODO clean up risk aversion...
				// TODO Parameterize willingness to sacrifice
				if f.Depth[nl] < uint16(depth) &&
					((s.Turn > 70 && nthreat < 2) || nthreat == 0 || nthreat < threat) &&
					(item == LAND || item == FOOD || item.IsEnemyHill(0)) {
					// We have a valid next step, path in to dest and see if
					// We should remove ant and possibly target.

					// Need to handle the special case where food spawns next to us
					// in that case we don't move, mark the food as a block, and
					// remove ourself and the food from the list

					if item == FOOD {
						// We are next to food, remove this ant from the
						// available list, mark this food as taken, and mark the
						// food as a block.
						s.SetOccupied(nl)
						moves[loc] = Direction(5) // explicitly say we will not move
						if Debug[DBG_Targets] {
							log.Printf("Removing %v food adjacent", s.Map.ToPoint(loc))
						}
						ants[loc] = 0, false
						tset[nl].Count = 0
					} else {
						endloc, steps := f.PathIn(s.Rand, nl)
						ep := s.Map.ToPoint(endloc)

						tgt, ok := tset[endloc]
						if Debug[DBG_Targets] {
							log.Printf("Found target %v -> %v, end %v %d steps, tgt: %v",
								p, np, ep, steps, tgt)
						}
						if !ok {
							if Debug[DBG_Targets] {
								log.Printf("pathin from %v(%d)->%v to %v(%d) no matching target.",
									p, loc, np, ep, endloc)
							}
						} else if steps >= 0 && tgt.Count > 0 {
							if Debug[DBG_Threat] {
								if threat != 0 || nthreat != 0 {
									log.Printf("#%d: %v -> %v threat %d -> %d\n", s.Turn, p, np, threat, nthreat)
								}
							}
							moves[loc] = Direction(d)
							tgt.Count -= 1
							s.SetOccupied(nl)
							ants[loc] = 0, false

							if Viz["path"] {
								fmt.Fprintf(os.Stdout, "v line %d %d %d %d\n", p.R, p.C, ep.R, ep.C)
							}
						}
					}
					break STEP
				}
			}
		}
		if Debug[DBG_Iterations] {
			log.Printf("Done allocating ants, Pending %d, ants %d", tset.Pending(), len(ants))
		}

		// TODO If we have more ants than targets we have bored ants, try to expand viewable area, slice
		if tset.Pending() < 1 && len(ants) > 0 {
			if len(bot.IdleAnts) < s.Turn {
				bot.IdleAnts = bot.IdleAnts[0 : s.Turn+1]
				bot.IdleAnts[s.Turn] = len(ants)
			}
			//log.Printf("%d ants with nothing to do", len(ants))

			fexp, _, _ := MapFill(s.Map, s.Ants[0], 1)
			loc, N := fexp.Sample(s.Rand, len(ants), 14, 14)
			for i := range loc {
				bot.Explore.Add(EXPLORE, loc[i], N[i], bot.Priority(EXPLORE))
				tset.Add(EXPLORE, loc[i], N[i], bot.Priority(EXPLORE))
			}
		} else {
			if len(bot.IdleAnts) < s.Turn {
				bot.IdleAnts = bot.IdleAnts[0 : s.Turn+1]
				bot.IdleAnts[s.Turn] = len(ants)
			}
		}
	}

	for loc, d := range moves {
		if d < 5 {
			p := s.Map.ToPoint(loc)
			fmt.Fprintf(os.Stdout, "o %d %d %s\n", p.R, p.C, DirectionChar[d])
		}
	}

	fmt.Fprintf(os.Stdout, "go\n")

	// TODO tiebreak on global goals.

	return nil
}
