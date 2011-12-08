package combat

import (
	"time"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/state"
	. "bugnuts/pathing"
)

type AntMove struct {
	Start    Location
	End      Location
	NStep    int
	Steps    [8]Direction
	Prefered Direction
}

type AntPartition struct {
	Ants       map[Location]struct{}
	estSimTime int64
}

type Partitions map[Location]*AntPartition

func NewAntPartition() *AntPartition {
	p := &AntPartition{
		Ants: make(map[Location]struct{}, 8),
	}
	return p
}

func CombatPartition(s *State) (Partitions, map[Location][]Location) {
	// how many ants are there
	nant := 0
	for _, ants := range s.Ants {
		nant += len(ants)
	}

	origin := make(map[Location]int, nant)
	for _, ants := range s.Ants {
		for loc := range ants {
			origin[loc] = 1
		}
	}
	f := NewFill(s.Map)
	// will only find neighbors withing 2x8 steps.
	_, near := f.MapFillSeedNN(origin, 1, 8)

	parts := make(Partitions, 5)
	// maps an ant to the partitions it belongs to.
	pmap := make(map[Location][]Location, nant)

	for ploc := range s.Ants[0] {
		if _, ok := pmap[ploc]; !ok {
			// for any of my ants not already in a partition
			for eloc, nn := range near[ploc] {
				if nn.Steps < 7 {
					if _, ok := s.Ants[0][eloc]; !ok {
						// a close enemy ant, add it and it's nearest neighbors to the partition

						ap, ok := parts[ploc] // ap = a partition
						if !ok {
							ap = NewAntPartition()
						}
						parts[ploc] = ap

						for nloc, nn := range near[eloc] {
							if nn.Steps < 7 {
								ap.Ants[nloc] = struct{}{}

								pm, ok := pmap[nloc]
								if !ok {
									pm = make([]Location, 0, 8)
								}
								pmap[nloc] = append(pm, ploc)
							}
						}
						pm, ok := pmap[eloc]
						if !ok {
							pm = make([]Location, 0, 8)
						}

						pmap[eloc] = append(pm, ploc)
						ap.Ants[eloc] = struct{}{}
					}
				}
			}
		}

		if ap, ok := parts[ploc]; ok {
			// If we created a partition centered on this ant add any
			// close neighbors of the friendly ants already in the
			// partition
			for loc := range ap.Ants {
				if _, ok := s.Ants[0][loc]; ok {
					// one of our friendly ants, add any close neigbors of our friendly guy
					for nloc, nn := range near[loc] {
						if nn.Steps < 2 {
							_, me := s.Ants[0][nloc]
							_, in := ap.Ants[nloc]
							if me && !in {
								ap.Ants[nloc] = struct{}{}
								pm, ok := pmap[nloc]
								if !ok {
									pm = make([]Location, 0, 8)
								}
								pmap[nloc] = append(pm, ploc)
							}
						}
					}
				}
			}
		}
	}

	/*
		for loc := range enemy {
			for floc, nn := range near[loc] {
				if nn.Steps < 6 {
					if _, ok := s.Ants[0][floc]; !ok {
						// a close not me ant
						enemy[floc] = 0
					}
				}
			}
		}
	*/

	return parts, pmap
}

func Combat(s *State, ants []*AntStep) {
	var pants []*AntPartition
	// partition by connectedness
	CombatPartition(s)

	// Compute available time

	// Compute per partition time budget

	// sim to compute best moves
	if false {
		for {
			for _, ap := range pants {
				if time.Nanoseconds()+ap.estSimTime > 500 { // s.TurnEnd {
					break
				}
				//s.Sim(ap)
			}
		}
	}

	// Move combat moves back to antstep

	// vis

	// call for help
}
