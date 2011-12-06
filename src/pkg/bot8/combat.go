package bot8

import (
	"os"
	"fmt"
	"log"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/state"
	. "bugnuts/MyBot"
	. "bugnuts/parameters"
	. "bugnuts/debug"
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

func Combat(s *State, ants []*AntStep) {
	// compute ants which are in combat range of original or new locations.
	cants := s.InCombat()

	// partition by connectedness
	sants, pants := s.Partition(ants, cants)

	// Compute available time

	// Compute per partition time budget

	// sim to compute best moves
	for {
		for _, ap := range pants {
			if time.NanoSeconds()+ap.estSimTime > s.TurnEnd {
				break
			}
			s.Sim(ap)
		}
	}

	// Move combat moves back to antstep

	// vis

	// call for help
}
