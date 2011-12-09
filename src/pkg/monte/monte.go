package combat

import (
	"log"
	"time"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/state"
	. "bugnuts/game"
)

type AntState struct {
	Start    Location
	End      Location
	NStep    int
	Steps    [8]Direction
	Prefered Direction
}

func (c *Combat) PartitionMoves(ap *AntPartition) []AntMove {
	// take the partition and setup the move list.
	am := make([]AntMove, len(ap.Ants))
	i := 0
	for loc := range ap.Ants {
		am[i].From = loc
		am[i].To = loc
		am[i].D = NoMovement
		if c.PlayerMap[loc] > -1 {
			am[i].Player = c.PlayerMap[loc]
			i++
		} else {
			log.Printf("Invalid ap player loc %v", c.ToPoint(am[i].From))
		}
	}

	return am
}

func (c *Combat) Sim(s *State, ploc Location, ap *AntPartition, cutoff int64) {
	log.Printf("Simulate for ap: %v %d ants, cutoff %.2fms",
		s.ToPoint(ploc),
		len(ap.Ants),
		float64(cutoff-time.Nanoseconds())/1e6)
}
