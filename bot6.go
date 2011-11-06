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
	"fmt"
	"os"
	"rand"
	"log"
)

type BotV6 struct {
	P *Parameters

	Explore *TargetSet

	IdleAnts []int
}

type AntStep struct {
	source  Location
	done    bool
	steptot int
	dest    []Location
	steps   []int
	nloc    [5]Location
	threat  [5]int8
	foodp   bool
	safest  [5]bool
}

//NewBot creates a new instance of your bot
func NewBotV6(s *State) Bot {
	mb := &BotV6{
		P:        ParameterSets["V6"],
		IdleAnts: make([]int, 0, s.Turns),
	}

	mb.Explore = MakeExplorers(s, .8, 1, mb.P.Priority[EXPLORE])
	return mb
}

func (bot *BotV6) ExploreUpdate(s *State) {
	// Any explore point which is visible should be nuked
	for loc, _ := range *bot.Explore {
		if s.Map.Seen[loc] == s.Turn {
			bot.Explore.Remove(loc)
		} else {
			(*bot.Explore)[loc].Count = 1
		}
	}
}

// Stores the neighborhood of the ant.
func (s *State) AntStep(loc Location) *AntStep {
	as := &AntStep{
		source:  loc,
		steptot: 0,
		dest:    make([]Location, 0, 4),
		steps:   make([]int, 0, 4),
	}

	p := s.Map.ToPoint(loc)
	for i, dir := range DirectionOffset {
		np := s.Map.PointAdd(p, dir)
		nloc := s.Map.ToLocation(np)
		as.nloc[i] = nloc
		as.threat[i] = s.Threat(s.Turn, nloc)

		if s.Item(nloc) == FOOD {
			as.foodp = true
		}
	}

	// Compute set of valid moves
	minthreat := MinInt8(as.threat[:])
	for i, _ := range DirectionOffset {
		as.safest[i] = as.threat[i] == minthreat
	}

	return as
}

func (s *State) EnemyPathinTargets(tset *TargetSet, priority int) {
	hills := make(map[Location]int, 6)
	for _, loc := range s.MyHillLocations() {
		hills[loc] = 1
	}

	f, _, _ := MapFill(s.Map, hills, 0)

	for i := 1; i < len(s.Ants); i++ {
		for loc, _ := range s.Ants[i] {
			_, steps := f.PathIn(Location(loc))
			if steps < 15 {
				(*tset).Add(DEFEND, Location(loc), 2, priority)
			}
		}
	}
}

func (tset *TargetSet) String() string {
	str := ""
	for loc, target := range *tset {
		str += fmt.Sprintf("%d: %#v\n", loc, target)
	}
	return str
}

func (bot *BotV6) DoTurn(s *State) os.Error {
	// TODO this still seems clunky.  need to figure out a better way
	s.FoodUpdate(bot.P.ExpireFood)

	bot.ExploreUpdate(s)

	tset := TargetSet{}

	s.EnemyPathinTargets(&tset, bot.P.Priority[DEFEND])

	// Generate list of food and enemy hill points.
	for _, loc := range s.FoodLocations() {
		if Debug > 4 {
			log.Printf("adding target %v(%d) food pri %d", s.Map.ToPoint(loc), loc, bot.P.Priority[FOOD])
		}
		tset.Add(FOOD, loc, 1, bot.P.Priority[FOOD])
	}
	tset.Merge(bot.Explore)

	// TODO handle different priorities/attack counts
	eh := s.EnemyHillLocations()
	for _, loc := range eh {
		// ndefend := s.PathinCount(loc, 10)
		tset.Add(HILL1, loc, 8, bot.P.Priority[HILL1])
	}

	// List of available ants, with local neighborhood
	ants := make(map[Location]*AntStep, len(s.Ants[0]))
	endants := make([]*AntStep, 0, len(s.Ants[0]))
	moves := make(map[Location]Direction, len(ants))
	for loc, _ := range s.Ants[0] {
		ants[loc] = s.AntStep(loc)

		// Handle the special case of adjacent food, pause a step unless
		// someone already paused for this food.
		if ants[loc].foodp && ants[loc].steptot == 0 {
			found := false
			for _, nloc := range ants[loc].nloc[0:4] {
				if s.Item(nloc) == FOOD && tset[nloc].Count > 0 {
					tset[nloc].Count = 0
					s.SetBlock(nloc)
					found = true
				}
			}
			if found {
				ants[loc].steptot = 1
				ants[loc].dest = append(ants[loc].dest, loc) // staying for now.
				ants[loc].steps = append(ants[loc].steps, 1)
				moves[loc] = Direction(5)
			}
		}
	}

	segs := make([]Segment, 0, len(ants))

	for iter := 0; iter < 15 && len(ants) > 0 && tset.Pending() > 0; iter++ {
		if Debug > 4 {
			log.Printf("Location iteration %d, ants: %d, tset.Pending %d", iter, len(ants), tset.Pending())
		}

		// TODO: Here should update map for fixed ants.
		// log.Printf("Tset.Active: %v", tset.Active())

		f, _, _ := MapFill(s.Map, tset.Active(), 0)
		// log.Printf("Depth: %d %d", f.Depth[2381], f.Depth[2382])

		segs = segs[0:0]
		for loc, _ := range ants {
			segs = append(segs, Segment{src: loc, steps: ants[loc].steptot})
		}

		f.ClosestStep(segs)
		// log.Printf("Segments: %v", segs)

		for _, seg := range segs {
			loc := seg.src
			p := s.Map.ToPoint(loc)
			ep := s.Map.ToPoint(seg.end)

			tgt, ok := tset[seg.end]
			if !ok {
				log.Printf("Move from %v(%d) to %v(%d) no target ant: %#v", s.Map.ToPoint(seg.src), seg.src, s.Map.ToPoint(seg.end), seg.end, ants[loc])
				if Viz {
					fmt.Fprintf(os.Stdout, "v line %d %d %d %d\n", p.r, p.c, ep.r, ep.c)
					fmt.Fprintf(os.Stdout, "v tileBorder %d %d MM\n", p.r, p.c)
				}
			}

			moved := false
			if ok && tgt.Count > 0 {
				moved = true
				if ants[loc].steptot == 0 {
					// Perm here so our bots are not biased to move in particular directions
					d := 0
					var nloc Location
				WAYOUT:
					for _, run := range [2]bool{false, true} {
						for _, d = range rand.Perm(4) {
							moved = false
							nloc = ants[loc].nloc[d]
							if (tgt.Item == DEFEND || tgt.Item.IsHill() || ants[loc].safest[d]) &&
								(run || f.Depth[nloc] < f.Depth[loc]) &&
								s.ValidStep(nloc) {
								moved = true
								break WAYOUT
							}
						}
					}
					if moved {
						s.MoveAnt(loc, nloc)
						moves[loc] = Direction(d)
					}
				}
			}

			if moved {
				fmt.Fprintf(os.Stdout, "v line %d %d %d %d\n", p.r, p.c, ep.r, ep.c)
				tgt.Count--
				ants[loc].steps = append(ants[loc].steps, seg.steps-ants[loc].steptot)
				ants[loc].steptot = seg.steps
				ants[loc].dest = append(ants[loc].dest, seg.end)

				if tgt.Terminal {
					endants = append(endants, ants[loc])
					ants[loc] = &AntStep{}, false
				} else {
					ants[seg.end] = ants[loc]
					ants[loc] = &AntStep{}, false
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
			loc, N := fexp.Sample(len(ants), 18, 18)
			for i, _ := range loc {
				exp := s.Map.ToPoint(loc[i])
				fmt.Fprintf(os.Stdout, "v star %d %d .5 1.5 5 true\n", exp.r, exp.c)

				bot.Explore.Add(EXPLORE, loc[i], N[i], bot.P.Priority[EXPLORE])
				tset.Add(EXPLORE, loc[i], N[i], bot.P.Priority[EXPLORE])
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

	if Viz {
		lthreat := int8(10)
		for i, threat := range s.ThreatMap(s.Turn) {
			if threat > 0 {
				if lthreat != threat {
					fmt.Fprintf(os.Stdout, "v setFillColor 255 0 0 %.1f\n", float64(threat)*.2)
					lthreat = threat
				}
				p := s.Map.ToPoint(Location(i))
				fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.r, p.c)
			}
		}
		fmt.Fprintf(os.Stdout, "v setFillColor 0 0 0 1.0\n")
	}
	fmt.Fprintf(os.Stdout, "go\n")

	// TODO tiebreak on global goals.

	return nil
}
