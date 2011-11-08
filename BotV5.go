package main
// The v5 Bot -- Marginally less Terrible!!!!
//
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

import (
	"fmt"
	"os"
	"rand"
	"log"
)

type BotV5 struct {
	P *Parameters

	Explore *TargetSet

	IdleAnts []int
	Primap   []int
}

func (bot *BotV5) Priority(i Item) int {
	return bot.Primap[i]
}

//NewBot creates a new instance of your bot
func NewBotV5(s *State) Bot {
	mb := &BotV5{
		P:        ParameterSets["V5"],
		IdleAnts: make([]int, 0, s.Turns),
	}

	mb.Explore = MakeExplorers(s, .8, 1, mb.Priority(EXPLORE))
	return mb
}

func (bot *BotV5) ExploreUpdate(s *State) {
	// Any explore point which is visible should be nuked
	for loc, _ := range *bot.Explore {
		if s.Map.Seen[loc] == s.Turn {
			bot.Explore.Remove(loc)
		} else {
			(*bot.Explore)[loc].Count = 1
		}
	}
}

func (bot *BotV5) DoTurn(s *State) os.Error {
	// TODO this still seems clunky.  need to figure out a better way
	s.FoodUpdate(bot.P.ExpireFood)
	bot.ExploreUpdate(s)

	tset := TargetSet{}

	tset.Merge(bot.Explore)

	// Generate list of food and enemy hill points.
	for _, loc := range s.FoodLocations() {
		if Debug > 4 {
			log.Printf("adding target %v(%d) food pri %d", s.Map.ToPoint(loc), loc, bot.Priority(FOOD))
		}
		tset.Add(FOOD, loc, 1, bot.Priority(FOOD))
	}

	// TODO handle different priorities for different enemies.
	eh := s.EnemyHillLocations()
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

	if Debug > 4 {
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
		if Debug > 4 {
			log.Printf("Location iteration %d, ants: %d, tset.Pending %d", iter, len(ants), tset.Pending())
		}

		// TODO: Here should update map for fixed ants.
		f, _, _ := MapFill(s.Map, tset.Active(), 0)

		// Build list of locations sorted by depth
		ccl := make([]Location, len(ants))
		for loc, _ := range ants {
			ccl = append(ccl, loc)
		}
		for _, loc := range f.Closest(ccl) {
			depth := f.Depth[loc]
			threat := s.Threat(s.Turn, loc)
			p := s.Map.ToPoint(loc)
		STEP:
			// Perm here so our bots are not biased to move in particular directions
			for _, d := range rand.Perm(4) {
				// find a direction we can step in thats stepable.
				np := s.Map.PointAdd(p, Steps[d])
				nl := s.Map.ToLocation(np)
				item := s.Item(nl)
				nthreat := s.Threat(s.Turn, nl)
				// TODO clean up risk aversion...
				// TODO Parameterize willingness to sacrifice
				if f.Depth[nl] < uint16(depth) &&
					((s.Turn > 70 && nthreat < 2) || nthreat == 0 || nthreat < threat) &&
					(item == LAND || item == FOOD || item.IsEnemyHill()) {
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
						if Debug > 4 {
							log.Printf("Removing %v food adjacent", s.Map.ToPoint(loc))
						}
						ants[loc] = 0, false
						tset[nl].Count = 0
					} else {
						endloc, steps := f.PathIn(nl)
						ep := s.Map.ToPoint(endloc)

						tgt, ok := tset[endloc]
						if Debug > 4 {
							log.Printf("Found target %v -> %v, end %v %d steps, tgt: %v",
								p, np, ep, steps, tgt)
						}
						if !ok {
							if Debug > 3 {
								log.Printf("pathin from %v(%d)->%v to %v(%d) no matching target.",
									p, loc, np, ep, endloc)
							}
						} else if steps >= 0 && tgt.Count > 0 {
							if s.Threat(s.Turn, nl) > 0 {
								log.Printf("#%d: %v -> %v threat %d -> %d\n", s.Turn, p, np, threat, nthreat)
							}
							moves[loc] = Direction(d)
							tgt.Count -= 1
							s.SetOccupied(nl)
							ants[loc] = 0, false

							if Viz["path"] {
								fmt.Fprintf(os.Stdout, "v line %d %d %d %d\n", p.r, p.c, ep.r, ep.c)
							}
						}
					}
					break STEP
				}
			}
		}

		// TODO If we have more ants than targets we have bored ants, try to expand viewable area, slice
		if tset.Pending() < 1 && len(ants) > 0 {
			if len(bot.IdleAnts) < s.Turn {
				bot.IdleAnts = bot.IdleAnts[0 : s.Turn+1]
				bot.IdleAnts[s.Turn] = len(ants)
			}
			//log.Printf("%d ants with nothing to do", len(ants))

			fexp, _, _ := MapFill(s.Map, s.Ants[0], 1)
			loc, N := fexp.Sample(len(ants), 14, 14)
			for i, _ := range loc {
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
			fmt.Fprintf(os.Stdout, "o %d %d %s\n", p.r, p.c, DirectionChar[d])
		}
	}

	fmt.Fprintf(os.Stdout, "go\n")

	// TODO tiebreak on global goals.

	return nil
}
