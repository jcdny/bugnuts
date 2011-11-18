package state

import (
	"log"
	"rand"
	"fmt"
	"sort"
	"os"
	. "bugnuts/maps"
	. "bugnuts/debug"
)

type Neighborhood struct {
	//TODO add hill distance step
	valid   bool
	threat  int
	pthreat int
	goal    int
	prfood  int
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

func (s *State) GenerateAnts(tset *TargetSet, risk int) (ants map[Location]*AntStep) {
	ants = make(map[Location]*AntStep, len(s.Ants[0]))

	for loc, _ := range s.Ants[0] {
		ants[loc] = s.AntStep(loc, risk)

		fixed := false

		// If I am on my hill and there is an adjacent enemy don't move
		hill, ok := s.Hills[loc]
		if ok && hill.Player == 0 {
			for _, nloc := range s.Map.LocStep[loc] {
				if s.Map.Item(nloc).IsEnemyAnt(0) {
					fixed = true
					break
				}
			}
		}

		// Handle the special case of adjacent food, pause a step unless
		// someone already paused for this food.
		if ants[loc].foodp && ants[loc].steptot == 0 {
			for _, nloc := range s.Map.LocStep[loc] {
				if s.Map.Item(nloc) == FOOD && (*tset)[nloc].Count > 0 {
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

// Stores the neighborhood of the ant.
func (s *State) Neighborhood(loc Location, nh *Neighborhood, d Direction) {
	nh.threat = int(s.Threat(s.Turn, loc))
	nh.pthreat = int(s.PThreat(s.Turn, loc))
	//nh.vis = s.Map.VisSum[loc]
	//nh.unknown = s.Met.Unknown[loc]
	//nh.land = s.Met.Land[loc]
	nh.prfood = s.Met.PrFood[loc]
	nh.d = d
}

func (s *State) AntStep(loc Location, risk int) *AntStep {
	as := &AntStep{
		source:  loc,
		steptot: 0,
		move:    6,
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
		as.N[d].perm = int(permute[d])

		if s.Map.Item(nloc) == FOOD {
			as.foodp = true
		}
		if s.ValidStep(nloc) {
			as.N[d].valid = true
			as.nfree++
		}
	}
	s.Neighborhood(loc, as.N[4], Direction(4))
	as.N[4].perm = int(permute[4])
	as.N[4].valid = true

	// Compute the min threat moves.
	if risk > 0 {
		for i := 0; i < 5; i++ {
			as.N[i].threat -= 4
			if as.N[i].threat <= 0 {
				as.N[i].threat = 0
				as.N[i].pthreat = 0
			}
		}
	}
	minthreat := as.N[4].threat*100 + as.N[4].pthreat
	for i := 0; i < 4; i++ {
		nt := as.N[i].threat*100 + as.N[i].pthreat
		if nt < minthreat {
			minthreat = nt
		}
	}
	for i := 0; i < 5; i++ {
		as.N[i].safest = (as.N[i].threat == minthreat)
	}

	return as
}

func (s *State) EmitMoves(ants []*AntStep) {
	for _, ant := range ants {
		if ant.move >= 0 && ant.move < NoMovement {
			p := s.ToPoint(ant.source)
			fmt.Fprintf(os.Stdout, "o %d %d %s\n", p.R, p.C, DirectionChar[ant.move])
		} else if ant.move != NoMovement {
			p := s.ToPoint(ant.source)
			log.Printf("Invalid move %d %d\n", p.R, p.C)
		}
	}
}

func (s *State) GenerateMoves(antsIn []*AntStep) {
	// make a copy of the ant slice
	ants := make([]*AntStep, len(antsIn))
	copy(ants, antsIn)
	lastants := len(ants)

	// loop until we move all the ants.
	for {
		if Debug[DBG_Movement] {
			log.Printf("ants: %d: %v", len(ants), ants)
			for i, ant := range ants {
				log.Printf("ants #%d: %v", i, ant)
			}
			sort.Sort(AntSlice(ants))
		}
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
				N.perm = int(perm[i])
				if N.d == NoMovement {
					N.valid = true
				} else {
					N.valid = s.ValidStep(s.Map.LocStep[ant.source][N.d])
				}
			}
		}
	}
}

func (s *State) Step(ant *AntStep) bool {
	if ant.move < 0 {
		sort.Sort(ENSlice(ant.N))
		if Debug[DBG_Movement] {
			for i, N := range ant.N {
				log.Printf("STEP %d %#v", i, N)
			}
		}
		ant.move = ant.N[0].d
		if ant.move == NoMovement || s.MoveAnt(ant.source, s.Map.LocStep[ant.source][ant.N[0].d]) {
			return true
		}
		ant.move = InvalidMove
	}
	return false
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
	if p[i].pthreat != p[j].pthreat {
		return p[i].pthreat < p[j].pthreat
	}
	if p[i].goal != p[j].goal {
		return p[i].goal > p[j].goal
	}
	if p[i].prfood != p[j].prfood {
		return p[i].prfood > p[j].prfood
	}
	return p[i].perm < p[j].perm
}
