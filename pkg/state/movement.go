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
	Valid   bool
	Threat  int
	PThreat int
	Goal    int
	PrFood  int
	//Vis     int
	//Unknown int
	//Land    int
	Perm   int // permuter
	D      Direction
	Safest bool
}

type AntStep struct {
	Source  Location   // our original location
	Move    Direction  // the next step
	Dest    []Location // track routing
	Steps   []int      // and distance
	Steptot int        // and sum total distance
	N       []*Neighborhood
	Foodp   bool
	Goalp   bool
	Perm    int // to randomize ants when sorting
	NFree   int
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
		if ants[loc].Foodp && ants[loc].Steptot == 0 {
			for _, nloc := range s.Map.LocStep[loc] {
				if s.Map.Item(nloc) == FOOD && (*tset)[nloc].Count > 0 {
					(*tset)[nloc].Count = 0
					s.SetOccupied(nloc) // food cant move but it will be gone.
					fixed = true
				}
			}
		}

		if fixed {
			ants[loc].Steptot = 1
			ants[loc].Dest = append(ants[loc].Dest, loc) // staying for now.
			ants[loc].Steps = append(ants[loc].Steps, 1)
			ants[loc].Move = NoMovement
			ants[loc].NFree = 0
			ants[loc].Goalp = true
		}
	}
	return ants
}

// Stores the neighborhood of the ant.
func (s *State) Neighborhood(loc Location, nh *Neighborhood, d Direction) {
	nh.Threat = int(s.Threat(s.Turn, loc))
	nh.PThreat = int(s.PThreat(s.Turn, loc))
	//nh.Vis = s.Map.VisSum[loc]
	//nh.Unknown = s.Met.Unknown[loc]
	//nh.Land = s.Met.Land[loc]
	nh.PrFood = s.Met.PrFood[loc]
	nh.D = d
}

func (s *State) AntStep(loc Location, risk int) *AntStep {
	as := &AntStep{
		Source:  loc,
		Steptot: 0,
		Move:    InvalidMove,
		Dest:    make([]Location, 0, 4),
		Steps:   make([]int, 0, 4),
		N:       make([]*Neighborhood, 5),
		NFree:   1,
		Perm:    rand.Int(),
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
		as.N[d].Perm = int(permute[d])

		if s.Map.Item(nloc) == FOOD {
			as.Foodp = true
		}
		if s.ValidStep(nloc) {
			as.N[d].Valid = true
			as.NFree++
		}
	}
	s.Neighborhood(loc, as.N[4], Direction(4))
	as.N[4].Perm = int(permute[4])
	as.N[4].Valid = true

	// Compute the min threat moves.
	if risk > 0 {
		for i := 0; i < 5; i++ {
			as.N[i].Threat -= 4
			if as.N[i].Threat <= 0 {
				as.N[i].Threat = 0
				as.N[i].PThreat = 0
			}
		}
	}
	minthreat := as.N[4].Threat*100 + as.N[4].PThreat
	for i := 0; i < 4; i++ {
		nt := as.N[i].Threat*100 + as.N[i].PThreat
		if nt < minthreat {
			minthreat = nt
		}
	}
	for i := 0; i < 5; i++ {
		as.N[i].Safest = (as.N[i].Threat == minthreat)
	}

	return as
}

func (s *State) EmitMoves(ants []*AntStep) {
	for _, ant := range ants {
		if ant.Move >= 0 && ant.Move < NoMovement {
			p := s.ToPoint(ant.Source)
			fmt.Fprintf(os.Stdout, "o %d %d %s\n", p.R, p.C, DirectionChar[ant.Move])
		} else if ant.Move != NoMovement {
			p := s.ToPoint(ant.Source)
			log.Printf("Invalid move %d %d %d\n", p.R, p.C, int(ant.Move))
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
				N.Perm = int(perm[i])
				if N.D == NoMovement {
					N.Valid = true
				} else {
					N.Valid = s.ValidStep(s.Map.LocStep[ant.Source][N.D])
				}
			}
		}
	}
}

func (s *State) Step(ant *AntStep) bool {
	if ant.Move == InvalidMove {
		sort.Sort(ENSlice(ant.N))
		if Debug[DBG_Movement] {
			for i, N := range ant.N {
				log.Printf("STEP %d %#v", i, N)
			}
		}
		ant.Move = ant.N[0].D
		if ant.Move == NoMovement || s.MoveAnt(ant.Source, s.Map.LocStep[ant.Source][ant.N[0].D]) {
			return true
		}
		ant.Move = InvalidMove
	}
	return false
}

// Order ants for trying to move.
type AntSlice []*AntStep

func (p AntSlice) Len() int      { return len(p) }
func (p AntSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p AntSlice) Less(i, j int) bool {
	if p[i].Goalp != p[j].Goalp {
		return p[i].Goalp
	}
	if p[i].Goalp && p[i].Steps[0] != p[j].Steps[0] {
		return p[i].Steps[0] < p[j].Steps[0]
	}
	if p[i].NFree != p[j].NFree {
		return p[i].NFree > p[j].NFree
	}

	return p[i].Perm > p[j].Perm
}

// For ordering perspective moves...
type ENSlice []*Neighborhood

func (p ENSlice) Len() int      { return len(p) }
func (p ENSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ENSlice) Less(i, j int) bool {
	if p[i].Valid != p[j].Valid {
		return p[i].Valid
	}
	if p[i].Threat != p[j].Threat {
		return p[i].Threat < p[j].Threat
	}
	if p[i].PThreat != p[j].PThreat {
		return p[i].PThreat < p[j].PThreat
	}
	if p[i].Goal != p[j].Goal {
		return p[i].Goal > p[j].Goal
	}
	if p[i].PrFood != p[j].PrFood {
		return p[i].PrFood > p[j].PrFood
	}
	return p[i].Perm < p[j].Perm
}
