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
	"rand"
	"fmt"
	"log"
	"sort"
)

type BotV6 struct {
	P        *Parameters
	Primap   []int // array mapping Item to priority
	Explore  *TargetSet
	IdleAnts []int
}

type Neighborhood struct {
	//TODO add hill distance step
	valid  bool
	threat int
	goal   int
	prfood int
	//vis     int
	//unknown int
	//land    int
	perm   int // permuter
	d      Direction
	safest bool
}

type AntStep struct {
	source  Location   // our original location
	move    Direction  // the next step
	dest    []Location // track routing
	steps   []int      // and distance
	steptot int        // and sum total distance
	N       []*Neighborhood
	foodp   bool
	goalp   bool
	perm    int // to randomize ants when sorting
	nfree   int
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
	if bot.Explore == nil {
		return
	}
	for loc, _ := range *bot.Explore {
		if s.Map.Seen[loc] == s.Turn {
			bot.Explore.Remove(loc)
		} else {
			(*bot.Explore)[loc].Count = 1
		}
	}
}

// Stores the neighborhood of the ant.
func (s *State) Neighborhood(loc Location, nh *Neighborhood, d Direction) {
	nh.threat = int(s.Threat(s.Turn, loc))
	//nh.vis = s.Map.VisSum[loc]
	//nh.unknown = s.Map.Unknown[loc]
	//nh.land = s.Map.Land[loc]
	nh.prfood = s.Map.PrFood[loc]
	nh.d = d
}

func (s *State) AntStep(loc Location) *AntStep {
	as := &AntStep{
		source:  loc,
		steptot: 0,
		move:    -1,
		dest:    make([]Location, 0, 4),
		steps:   make([]int, 0, 4),
		N:       make([]*Neighborhood, 5),
		nfree:   1,
		perm:    rand.Int(),
	}
	nh := new([5]Neighborhood)
	for i, _ := range as.N {
		as.N[i] = &nh[i]
	}

	// Populate the neighborhood info
	permute := Permute5()
	for d := 0; d < 4; d++ {
		nloc := s.Map.LocStep[loc][d]
		s.Neighborhood(nloc, as.N[d], Direction(d))
		as.N[d].perm = permute[d]

		if s.Item(nloc) == FOOD {
			as.foodp = true
		}
		if s.ValidStep(nloc) {
			as.N[d].valid = true
			as.nfree++
		}
	}
	s.Neighborhood(loc, as.N[4], Direction(4))
	as.N[4].perm = permute[4]
	as.N[4].valid = true

	// Compute the min threat moves.
	minthreat := as.N[4].threat
	for i := 0; i < 4; i++ {
		if as.N[i].threat < minthreat {
			minthreat = as.N[i].threat
		}
	}
	for i := 0; i < 5; i++ {
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

	for _, loc := range s.Map.HBorder {
		depth := s.Map.FHill.Depth[loc]
		if depth > 2 && depth < uint16(bot.P.MinHorizon) {
			// Just add these as transients.
			tset.Add(WAYPOINT, loc, 1, bot.Priority(WAYPOINT))
		}
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
			ants[loc].move = NoMovement
			ants[loc].nfree = 0
			ants[loc].goalp = true
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

	segs := make([]Segment, 0, len(ants))

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
			ant := ants[seg.src]
			tgt, ok := (*tset)[seg.end]
			if !ok {
				log.Printf("Move from %v(%d) to %v(%d) no target ant: %#v",
					s.ToPoint(seg.src), seg.src, s.ToPoint(seg.end), seg.end, ant)
				log.Printf("Source item \"%v\", pending=%d", s.Map.Grid[seg.src], tset.Pending())
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
					for i := 0; i < 4; i++ {
						nloc := s.Map.LocStep[seg.src][i]
						// Don't mark target as taken unless its a valid step and risk = 0
						// TODO not sure this is how I should be doing this.
						goal := int(f.Depth[seg.src]) - int(f.Depth[nloc])
						ant.N[i].goal = goal
						// Check for a valid move towards the goal
						if s.ValidStep(nloc) && goal > 0 {
							// and it needs to be a step we can take
							if ant.N[i].safest ||
								((tgt.Item == DEFEND || tgt.Item.IsHill()) &&
									ant.N[i].threat < 2 && seg.steps < 10) {
								good = true
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

		// We have more ants than targets we have bored ants, keep track of their #s
		if tset.Pending() < 1 {
			if len(bot.IdleAnts) < s.Turn {
				idle := 0
				for _, ant := range ants {
					if !ant.goalp {
						idle++
					}
				}
				bot.IdleAnts = bot.IdleAnts[0 : s.Turn+1]
				bot.IdleAnts[s.Turn] = idle
				if Debug > 3 {
					log.Printf("TURN %d IDLE %d", s.Turn, len(ants))
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
		}
	}
	if Debug > 0 {
		log.Printf("TURN %d ITER %d", s.Turn, iter)
	}

	for _, ant := range ants {
		endants = append(endants, ant)
	}

	// Generate moves
	s.GenerateMoves(endants)

	for _, ant := range endants {
		if ant.move > -1 && ant.move < NoMovement {
			p := s.ToPoint(ant.source)
			fmt.Fprintf(os.Stdout, "o %d %d %s\n", p.r, p.c, DirectionChar[ant.move])
		} else {
			if ant.move != NoMovement {
				p := s.ToPoint(ant.source)
				log.Printf("Invalid move %d %d\n", p.r, p.c)
			}
		}
	}

	s.Viz()
	fmt.Fprintf(os.Stdout, "go\n") // TODO Flush ??
	return nil
}

// Order ants for trying to move.
type AntSlice []*AntStep

func (p AntSlice) Len() int      { return len(p) }
func (p AntSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p AntSlice) Less(i, j int) bool {
	if p[i].goalp != p[j].goalp {
		return p[i].goalp
	}
	if p[i].goalp && p[i].steps[0] != p[j].steps[0] {
		return p[i].steps[0] < p[j].steps[0]
	}
	if p[i].nfree != p[j].nfree {
		return p[i].nfree > p[j].nfree
	}

	return p[i].perm > p[j].perm
}

func (s *State) GenerateMoves(antsIn []*AntStep) {
	// make a copy of the ant slice
	ants := make([]*AntStep, len(antsIn))
	copy(ants, antsIn)
	lastants := len(ants)

	// loop until we move all the ants.
	for {
		sort.Sort(AntSlice(ants))
		stuck := 0
		for _, ant := range ants {
			if !s.Step(ant) {
				ants[stuck] = ant
				stuck++
			}
		}
		// if we have 0 ants remaining or we did not
		// allocate any ants this turn then terminate
		if stuck == 0 || stuck == lastants {
			break
		}
		ants = ants[0:stuck]
		lastants = stuck

		// Recompute perm and nfree
		perm := Permute5()
		for _, ant := range ants {
			for i, N := range ant.N {
				N.perm = perm[i]
				if N.d == NoMovement {
					N.valid = true
				} else {
					N.valid = s.ValidStep(s.Map.LocStep[ant.source][N.d])
				}
			}
		}
	}
}

// For ordering perspective moves...

type ENSlice []*Neighborhood

func (p ENSlice) Len() int      { return len(p) }
func (p ENSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ENSlice) Less(i, j int) bool {
	if p[i].valid != p[j].valid {
		return p[i].valid
	}
	if p[i].threat != p[j].threat {
		return p[i].threat < p[j].threat
	}
	if p[i].goal != p[j].goal {
		return p[i].goal > p[j].goal
	}
	if p[i].prfood != p[j].prfood {
		return p[i].prfood > p[j].prfood
	}
	return p[i].perm < p[j].perm
}

func (s *State) Step(ant *AntStep) bool {
	if ant.move < 0 {
		sort.Sort(ENSlice(ant.N))
		if Debug > 3 {
			for i, N := range ant.N {
				log.Printf("STEP %d %#v", i, N)
			}
		}
		if !ant.N[0].valid {
			return false
		}
		ant.move = ant.N[0].d
	}

	// sort the possible steps
	if ant.move != NoMovement {
		s.MoveAnt(ant.source, s.Map.LocStep[ant.source][ant.move])
	}

	return true
}
