// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package state

import (
	"log"
	"fmt"
	"sort"
	"os"
	. "bugnuts/torus"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/watcher"
)

func (s *State) MoveAnt(from, to Location) bool {
	if from == to {
		return true
	}
	if s.ValidStep(to) {
		s.Map.Grid[from], s.Map.Grid[to] = LAND, OCCUPIED
		return true
	}
	return false
}

func (s *State) GenerateAnts(tset *TargetSet, riskOff bool) (ants map[Location]*AntStep) {
	ants = make(map[Location]*AntStep, len(s.Ants[0]))

	for loc := range s.Ants[0] {
		ants[loc] = s.AntStep(loc, riskOff)
		//log.Printf("Ant: %#v", ants[loc])

		fixed := false

		// If I am on my hill and there is an adjacent enemy don't move
		hill, ok := s.Hills[loc]
		if ok && hill.Player == 0 {
			for _, nloc := range s.Map.LocStep[loc] {
				if nloc != loc && s.Map.Grid[nloc].IsEnemyAnt(0) {
					fixed = true
					break
				}
			}
		}

		// Handle the special case of adjacent food, pause a step unless
		// someone already paused for this food.
		if ants[loc].Foodp && ants[loc].Steptot == 0 {
			for _, nloc := range s.Map.LocStep[loc] {
				if s.Map.Grid[nloc] == FOOD && (*tset)[nloc].Count > 0 {
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
	nh.Threat = s.C.Threat1[loc] - s.C.PThreat1[0][loc]
	//nh.PrThreat = s.C.PrThreat[loc]
	nh.PrFood = s.Met.PrFood[loc]
	nh.D = d
	nh.Item = s.Map.Grid[loc]
}

func (s *State) AntStep(loc Location, riskOff bool) *AntStep {
	as := &AntStep{
		Source:  loc,
		Steptot: 0,
		Move:    InvalidMove,
		Dest:    make([]Location, 0, 4),
		Steps:   make([]int, 0, 4),
		N:       make([]*Neighborhood, 5),
		NFree:   0,
		Perm:    s.Rand.Int(),
	}
	nh := new([5]Neighborhood)
	for i := range as.N {
		as.N[i] = &nh[i]
	}

	// Populate the neighborhood info
	permute := Permute5(s.Rand)
	for d := 0; d < 5; d++ {
		nloc := s.Map.LocStep[loc][d]
		s.Neighborhood(nloc, as.N[d], Direction(d))
		if riskOff {
			as.N[d].Threat -= 3 // MAGIC
			if as.N[d].Threat <= 0 {
				as.N[d].Threat = 0
				as.N[d].PrThreat = 0
			}
		}
		as.N[d].Perm = int(permute[d])
		//log.Printf("%v %v %v", s.ToPoint(nloc), Direction(d), s.Map.Grid[nloc])

		if s.Map.Grid[nloc] == FOOD {
			as.Foodp = true
		}
		if nloc == loc {
			as.N[d].Valid = true
		} else if s.ValidStep(nloc) {
			as.N[d].Valid = true
			if as.N[d].Threat == 0 {
				as.NFree++
			}
		}
	}

	// Compute the min threat moves and flag as safest.
	minthreat := as.N[4].Threat*100 + as.N[4].PrThreat
	for i := 0; i < 4; i++ {
		nt := as.N[i].Threat*100 + as.N[i].PrThreat
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
			log.Printf("Encountered Invalid move %d %d turn %d\n", p.R, p.C, s.Turn)
		}
	}
}

func (s *State) TurnDone() {
	fmt.Fprintf(os.Stdout, "go\n") // TODO Flush ??
}

func (s *State) GenerateMoves(antsIn []*AntStep) {
	// make a copy of the ant slice
	ants := make([]*AntStep, len(antsIn))
	copy(ants, antsIn)
	lastants := len(ants)

	// loop until we move all the ants.
	for {
		sort.Sort(AntSlice(ants))
		if Debug[DBG_Movement] {
			log.Printf("ants: %d: %v", len(ants), ants)
			for i, ant := range ants {
				log.Printf("ants #%d: %#v", i, ant)
			}
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
		perm := Permute5(s.Rand)
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
			log.Printf("move %#v", ant)
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
