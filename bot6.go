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
	source Location
	steptot int
	dest   []Location
	steps  []int

	nloc   [5]Location
	item   [5]Item
	threat [5]int8

	foodp  bool
	validp [5]bool
	moves [5]bool
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
		source: loc,
		steptot: 0,
	dest: make([]Location, 0, 4),
	steps: make([]int, 0, 4),
	}
	

	p := s.Map.ToPoint(loc)
	for i, dir := range DirectionOffset {
		np := s.Map.PointAdd(p, dir)
		nloc := s.Map.ToLocation(np)
		as.nloc[i] = nloc
		as.threat[i] = s.Threat(s.Turn, nloc)
		
		item := s.Item(nloc)
		as.item[i] = item
		if item == FOOD {
			as.foodp = true
		}
		if item == LAND || item == FOOD || item.IsEnemyHill() {
			as.validp[i] = true
		}
	}

	// Compute set of valid moves
	minthreat := MinInt8(as.threat[:])
	for i, _ := range DirectionOffset {
		as.moves[i] = (as.validp[i] && as.threat[i] == minthreat)
	}

	return as
}


func (bot *BotV6) DoTurn(s *State) os.Error {
	// TODO this still seems clunky.  need to figure out a better way
	s.FoodUpdate(bot.P.ExpireFood)

	bot.ExploreUpdate(s)

	tset := TargetSet{}

	// List of available ants, with local neighborhood
	ants := make(map[Location]*AntStep, len(s.Ants[0]))
	for loc, _ := range s.Ants[0] {
		ants[loc] = s.AntStep(loc)
		log.Printf("Ants %#v", ants[loc])
	}


	// TODO handle different priorities/attack counts
	eh := s.EnemyHillLocations()
	idle := 0
	if len(bot.IdleAnts) > s.Turn-1 && len(eh) > 0 {
		//log.Printf("IDLE: %d : %d : %d : %v", s.Turn, len(bot.IdleAnts), len(eh), bot.IdleAnts)
		idle = bot.IdleAnts[s.Turn-1] / len(eh)
		if idle > 3 {
			idle = 3
		}
	}

	for _, loc := range eh {
		tset.Add(HILL1, loc, 1+idle, bot.P.Priority[HILL1])
	}

	// Generate list of food and enemy hill points.
	for _, loc := range s.FoodLocations() {
		if Debug > 4 {
			log.Printf("adding target %v(%d) food pri %d", s.Map.ToPoint(loc), loc, bot.P.Priority[FOOD])
		}
		tset.Add(FOOD, loc, 1, bot.P.Priority[FOOD])
	}
	tset.Merge(bot.Explore)

	// TODO remove dedicated ants eg sentinel, capture, defense guys
	moves := make(map[Location]Direction, len(ants))

	segs := make([]Segment, 0, len(ants))

	for iter := 0; iter < 15 && len(ants) > 0 && (iter == 0 || tset.Pending() > 0); iter++ {
		if Debug > 4 {
			log.Printf("Location iteration %d, ants: %d, tset.Pending %d", iter, len(ants), tset.Pending())
		}

		// TODO: Here should update map for fixed ants.
		f, _, _ := MapFill(s.Map, tset.Active(), 0)

		segs = segs[0:len(ants)]
		i := 0
		for loc, _ := range ants {		
			segs[i] = Segment{src: loc, steps: ants[loc].steptot}
			i++
		}
		
		f.ClosestStep(segs)
		log.Printf("Segments: %v", segs)

		for _, seg := range segs {
			loc := seg.src
			// p := s.Map.ToPoint(loc)
			if ants[loc].foodp && ants[loc].steptot == 0 {
				// Special case adjacent to food, pause a step
				ants[loc].moves[4] = true // sitting on our thumbs
				ants[loc].steptot = 1
				ants[loc].dest = append(ants[loc].dest, seg.src) // staying for now.
				ants[loc].steps = append(ants[loc].steps, 1)
				for i, item := range ants[loc].item {
					if item == FOOD {
						s.SetBlock(ants[loc].nloc[i])
						tset[ants[loc].nloc[i]].Count = 0
					}
				}
				moves[loc] = Direction(5)
			} else {
				moved := false
				if tset[seg.end].Count > 0 {
					moved = true
					if ants[loc].steptot == 0 {
						// Perm here so our bots are not biased to move in particular directions
						for _, d := range rand.Perm(4) {
							moved = false
							nloc = ants[loc].nloc[d]
							_, antp := ants[nloc]
							if ants[loc].moves[d] &&
								f.Depth[nloc] < f.Depth[loc] && 
								!antp {
								s.SetBlock(nloc)
								moves[loc] = Direction(d)
								moved = true
								break
							}
						}
					}
				}

				if moved {
					tset[seg.end].Count--
					ants[loc].steps = append(ants[loc].steps, seg.steps - ants[loc].steptot)
					ants[loc].steptot = seg.steps
					ants[loc].dest = append(ants[loc].dest, seg.end)
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

	fmt.Fprintf(os.Stdout, "go\n")

	// TODO tiebreak on global goals.

	return nil
}
