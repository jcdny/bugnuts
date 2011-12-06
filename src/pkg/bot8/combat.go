package bot8

import (
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
	Ants [][]AntMove
}

func CombatPartition(s *State) {
	nant := 0
	for player, ants := range s.Ants {
		for loc := range ants {
			nant++
		}
	}
	origin := make(map[Location]int, nant)
	for player, ants := range s.Ants {
		for loc := range ants {
			origin[loc] = 1
		}
	}
	f := NewFill(s.Map)
	_, near := f.MapFillSeedNN(origin, 1)
	enemy := make(map[Location]int, 30)
	friend := make(map[Location]int, 30)
	for loc := range s.Ants[0] {
		for eloc, nn := range near[loc] {
			if nn.Steps < 8 {
				if _, ok := s.Ants[0][eloc]; !ok {
					// a close enemy ant
					enemy[eloc] = 0
					friend[loc] = 0
				}
			}
		}
	}
	for loc := range friend {
		for floc, nn := range near[loc] {
			if nn.Steps < 3 {
				if _, ok := s.Ants[0][eloc]; ok {
					// a close friend
					friend[floc] = 0
				}
			}
		}
	}
	for loc := range enemy {
		for floc, nn := range near[loc] {
			if nn.Steps < 3 {
				if _, ok := s.Ants[0][floc]; !ok {
					// a close not me
					enemy[floc] = 0
				}
			}
		}
	}

	// Now visualize the frenemies.
	VizFrenemies(s, friend, enemy)
}

func Combat(s *State, ants []*AntStep) {
	// partition by connectedness
	CombatPartition(s)

	// Compute available time

	// Compute per partition time budget

	// sim to compute best moves
	if false {
		for {
			for _, ap := range pants {
				if time.NanoSeconds()+ap.estSimTime > s.TurnEnd {
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
