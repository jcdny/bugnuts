package bot8

import (
	"time"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/state"
	. "bugnuts/pathing"
	. "bugnuts/viz"
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

func CombatPartition(s *State) {
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

	ap := make(Partitions, 5)
	// maps an ant to the partitions it belongs to.
	pmap := make(map[Location][]Location, nant)

	for loc := range s.Ants[0] {
		if _, ok := pmap[loc]; !ok {
			// for any of my ants not already in a partition
			for eloc, nn := range near[loc] {
				if nn.Steps < 8 {
					if _, ok := s.Ants[0][eloc]; !ok {
						// a close enemy ant, add it and it's nearest neighbors to the partition
						if p, ok := ap[loc]; !ok {
							p = NewAntPartition()
							ap[loc] = p
						}
						for nloc, nn := range near[eloc] {
							if nn.Steps < 10 {
								p.Ants[nloc] = struct{}
								pm, ok := pmap[nloc]
								if !ok {
									pm = make([]Location, 8)
								}
								pmap[nloc] = append(pm, loc)
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

	// Now visualize the frenemies.
	if Viz["combat"] {
		VizFrenemies(s, ap, pmap)
	}
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
