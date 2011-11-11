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
	"log"
	"rand"
	"math"
)

type BotV6 struct {
	P        *Parameters
	Primap   []int // array mapping Item to priority
	Explore  *TargetSet
	IdleAnts []int
}

type Neighborhood struct {
	//TODO add hill distance step
	threat  int
	safest  bool
	vis     int
	unknown int
	land    int
	goal    int
}

type AntStep struct {
	source  Location   // our original location
	move    Direction  // the next step
	dest    []Location // track routing
	steps   []int      // and distance
	steptot int        // and sum total distance
	N       [5]Neighborhood
	foodp   bool
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
		IdleAnts: make([]int, 0, s.Turns),
	}

	mb.Primap = mb.P.MakePriMap()

	mb.Explore = MakeExplorers(s, .8, 1, mb.Priority(EXPLORE))
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
func (s *State) Neighborhood(loc Location, n *Neighborhood) {
	n.threat = int(s.Threat(s.Turn, loc))
	n.vis = s.Map.VisSum[loc]
	n.unknown = s.Map.Unknown[loc]
	n.land = s.Map.Land[loc]
}

func (s *State) AntStep(loc Location) *AntStep {
	as := &AntStep{
		source:  loc,
		steptot: 0,
		dest:    make([]Location, 0, 4),
		steps:   make([]int, 0, 4),
	}

	nloc := loc
	for i := 4; i >= 0; i-- {
		s.Neighborhood(nloc, &as.N[i])
		if s.Item(nloc) == FOOD {
			as.foodp = true
		}
		// do last since we are tacking on loc at [4]
		nloc = s.Map.LocStep[loc][i]
	}

	// Compute the min threat moves.
	minthreat := as.N[0].threat
	for i := 1; i < 5; i++ {
		if as.N[i].threat < minthreat {
			minthreat = as.N[0].threat
		}
	}
	for i := 1; i < 5; i++ {
		as.N[i].safest = (as.N[i].threat == minthreat)
	}

	return as
}

func (s *State) EnemyPathinTargets(tset *TargetSet, priority int, DefendDist int) {
	hills := make(map[Location]int, 6)
	for _, loc := range s.HillLocations(0) {
		hills[loc] = 1
	}

	f, _, _ := MapFill(s.Map, hills, 0)

	for i := 1; i < len(s.Ants); i++ {
		for loc, _ := range s.Ants[i] {
			// TODO: use seed rather than PathIn
			_, steps := f.PathIn(Location(loc))
			if steps < DefendDist {
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

func (bot *BotV6) GenerateTargets(s *State) *TargetSet {
	tset := &TargetSet{}

	s.EnemyPathinTargets(tset, bot.Priority(DEFEND), bot.P.DefendDistance)

	// Generate list of food and enemy hill points.
	// Food locations should be set after ant list is done since we
	// remove adjacent food at that step.
	for _, loc := range s.FoodLocations() {
		if Debug > 4 {
			log.Printf("adding target %v(%d) food pri %d", s.ToPoint(loc), loc, bot.Priority(FOOD))
		}
		tset.Add(FOOD, loc, 1, bot.Priority(FOOD))
	}

	tset.Merge(bot.Explore)

	// TODO handle different priorities/attack counts
	// TODO compute defender count
	eh := s.EnemyHillLocations(0)
	for _, loc := range eh {
		// ndefend := s.PathinCount(loc, 10)
		tset.Add(HILL1, loc, 8, bot.Priority(HILL1))
	}

	return tset
}

func (s *State) GenerateAnts(tset *TargetSet) (ants map[Location]*AntStep) {
	ants = make(map[Location]*AntStep, len(s.Ants[0]))

	for loc, _ := range s.Ants[0] {
		ants[loc] = s.AntStep(loc)

		fixed := false

		// If I am on my hill and there is an adjacent enemy don't move
		hill, ok := s.Hills[loc]
		if ok && hill.Player == 0 {
			for _, nloc := range s.Map.LocStep[loc] {
				if s.Item(nloc).IsEnemyAnt(0) {
					fixed = true
					break
				}
			}
		}

		// Handle the special case of adjacent food, pause a step unless
		// someone already paused for this food.
		if ants[loc].foodp && ants[loc].steptot == 0 {
			for _, nloc := range s.Map.LocStep[loc] {
				if s.Item(nloc) == FOOD && (*tset)[nloc].Count > 0 {
					(*tset)[nloc].Count = 0
					s.SetOccupied(nloc) // food cant move but it will be gone.
					fixed = true
				}
			}
		}

		if fixed {
			ants[loc].steptot = 1
			ants[loc].dest = append(ants[loc].dest, loc) // staying for now.
			ants[loc].steps = append(ants[loc].steps, 1)
			ants[loc].move = Direction(5)
		}
	}

	return ants
}

func (bot *BotV6) DoTurn(s *State) os.Error {
	// TODO this still seems clunky.  need to figure where this belongs.
	s.FoodUpdate(bot.P.ExpireFood)
	bot.ExploreUpdate(s)

	tset := bot.GenerateTargets(s)
	ants := s.GenerateAnts(tset)

	endants := make([]*AntStep, 0, len(ants))

	// List of available ants, with local neighborhood
	segs := make([]Segment, 0, len(ants))

	for _, i := range rand.Perm(len(s.Map.HBorder)) {
		loc := s.Map.HBorder[i]
		depth := s.Map.FHill.Depth[loc]
		if int(depth) < bot.P.MinHorizon {
			// Just add these as transients.
			tset.Add(WAYPOINT, loc, 1, bot.Priority(WAYPOINT))
		}
	}

	if Viz["targets"] {
		s.VizTargets(tset)
	}

	var iter, maxiter int = 0, 50
	for iter = 0; iter < maxiter && len(ants) > 0 && tset.Pending() > 0; iter++ {
		if Debug > 4 {
			log.Printf("TURN %d ITER %d TGT PENDING %d", s.Turn, iter, tset.Pending())
			// log.Printf("ACTIVE SET: %v", tset.Active())
		}

		// TODO: Here should update map for fixed ants.
		f, _, _ := MapFillSeed(s.Map, tset.Active(), 0)

		segs = segs[0:0]
		for loc, _ := range ants {
			segs = append(segs, Segment{src: loc, steps: ants[loc].steptot})
		}

		f.ClosestStep(segs)
		for _, seg := range segs {
			tgt, ok := (*tset)[seg.end]
			if !ok {
				log.Printf("Move from %v(%d) to %v(%d) no target ant: %#v",
					s.ToPoint(seg.src), seg.src, s.ToPoint(seg.end), seg.end, ants[seg.src])
				log.Printf("Source item \"%v\", pending=%d", s.Map.Grid[seg.src], tset.Pending())
				if Viz["error"] {
					p := s.ToPoint(seg.src)
					VizLine(s.Map, p, s.ToPoint(seg.end), false)
					fmt.Fprintf(os.Stdout, "v tileBorder %d %d MM\n", p.r, p.c)
				}
			} else if ok && tgt.Count > 0 {
				ants[seg.src].N[5].goal = 0
				good := false
				for i := 0; i < 4; i++ {
					nloc := s.Map.LocStep[seg.src][i]
					// Don't mark target as taken unless its a valid step and risk = 0
					goal := int(f.Depth[seg.src] - f.Depth[nloc])
					ants[seg.src].N[i].goal = goal
					if s.ValidStep(nloc) &&
						((ants[seg.src].N[i].safest && goal > 0) ||
							((tgt.Item == DEFEND || tgt.Item.IsHill()) &&
								ants[seg.src].N[i].threat < 2 && seg.steps < 10)) {
						good = true
					}
				}

				if good {
					tgt.Count--
					ants[seg.src].steps = append(ants[seg.src].steps, seg.steps-ants[seg.src].steptot)
					ants[seg.src].steptot = seg.steps
					ants[seg.src].dest = append(ants[seg.src].dest, seg.end)
					if tgt.Terminal {
						endants = append(endants, ants[seg.src])
					} else {
						ants[seg.end] = ants[seg.src]
						ants[seg.src] = &AntStep{}, false
					}
				}
			}
		}

		// TODO If we have more ants than targets we have bored ants, try to expand viewable area, etc
		if tset.Pending() < 1 && len(ants) > 0 {
			if len(bot.IdleAnts) < s.Turn {
				bot.IdleAnts = bot.IdleAnts[0 : s.Turn+1]
				bot.IdleAnts[s.Turn] = len(ants)
				if Debug > 3 {
					log.Printf("BotV6: %d ants with nothing to do", len(ants))
				}
			}

			if false {
				// Generate a target list for unseen areas and exploration
				// tset.Add(RALLY, s.Map.ToLocation(Point{58, 58}), len(ants), bot.Priority(RALLY))
				fexp, _, _ := MapFill(s.Map, s.Ants[0], 1)
				loc, N := fexp.Sample(len(ants), 18, 18)
				for i, _ := range loc {
					exp := s.ToPoint(loc[i])
					fmt.Fprintf(os.Stdout, "v star %d %d .5 1.5 5 true\n", exp.r, exp.c)

					bot.Explore.Add(EXPLORE, loc[i], N[i], bot.Priority(EXPLORE))
					tset.Add(EXPLORE, loc[i], N[i], bot.Priority(EXPLORE))
				}
			}
		} else {
			if len(bot.IdleAnts) < s.Turn {
				bot.IdleAnts = bot.IdleAnts[0 : s.Turn+1]
				bot.IdleAnts[s.Turn] = len(ants)
			}
		}
	}

	// for any ant coming out with no move use the metrics to find a best next step.
	dbest := make([]int, 0, 5)
	for loc, ant := range ants {
		best := math.MaxInt32
		d := 0
		nloc := loc
		best = X[0]
		for _, run := range [2]bool{false, true} {
			for d = 0; d < 4; d++ {
				// Compute a metric for the best move given:
				// * unknown cells
				// * visibility overlap
				// * land visible * turns unseen (proxy for food prob)
				nloc = s.Map.LocStep[loc][d]
			}
		}

		ants[loc].move = Direction(d)
		s.MoveAnt(loc, nloc)
	}

	if Debug > 0 {
		log.Printf("TURN %d %d iterations", s.Turn, iter)
	}

	for _, ant := range ants {
		endants = append(endants, ant)
	}
	for _, ant := range endants {
		if ant.move < 5 {
			p := s.ToPoint(ant.source)
			fmt.Fprintf(os.Stdout, "o %d %d %s\n", p.r, p.c, DirectionChar[ant.move])
		}
	}

	s.Viz()

	fmt.Fprintf(os.Stdout, "go\n")
	// TODO Flush ??

	return nil
}
