package state

import (
	. "bugnuts/game"
	. "bugnuts/maps"
)

type TurnStatistics struct {
	N       int             // max ant # seen
	Food    int             // Food count
	Unknown int             // Count of unknown tiles (from TGrid, before
	Horizon int             // Count of tiles where we know state
	Seen    [MaxPlayers]int // How many ants we have seen
	SeenAll int             // sum total ants seen
	Died    [MaxPlayers]int // Chronicle of a Death Foretold
	DiedAll int             // Total deaths seen
	//Gathered    [MaxPlayers]int // How much food did we see gathered
	StaticCount [MaxPlayers]int // Count of unmoved ants
}

type Statistics struct {
	DiedMap []int // Death count per location
	//StaticMap      []int            // Turns given cell has been occupied by the same player.
	//PlayerMap      []int            // used to drive static map
	DiedTot        [MaxPlayers]int  // Chronicle of a Death Foretold, Running total of deaths seen
	DiedTotAll     int              // Total number of deaths seen this game
	SeenMax        [MaxPlayers]int  // Maximum number seen
	SeenMaxTurn    [MaxPlayers]int  // Turn on which max number was seen
	HorizonMax     int              // Maximum known state extent
	HorizonMaxTurn int              // turn on which we had max knowledge
	TStats         []TurnStatistics // per turn statistics
}

func NewStatistics(g *GameInfo) *Statistics {
	stats := &Statistics{
		DiedMap: make([]int, g.Rows*g.Cols),
		//StaticMap: make([]int, g.Rows*g.Cols),
		//PlayerMap: make([]int, g.Rows*g.Cols),
		TStats:     make([]TurnStatistics, g.Turns+2),
		HorizonMax: MAXMAPSIZE,
	}

	return stats
}

func (s *State) UpdateStatistics(turn *Turn) {
	ts := &s.Stats.TStats[turn.Turn]
	s.Stats.ProcessDeadAnts(ts, turn.D)
	s.Stats.ProcessSeen(ts, turn.A, turn.Turn)
	ts.Food = len(turn.F)

	// Horizon count
	nh := 0
	for _, h := range s.Met.Horizon {
		if !h {
			nh++
		}
	}
	ts.Horizon = nh
	if nh <= s.Stats.HorizonMax {
		s.Stats.HorizonMax = nh
		s.Stats.HorizonMaxTurn = turn.Turn
	}

	// Unknown, quit computing if we knew the whole map.
	if turn.Turn <= 10 || s.Stats.TStats[turn.Turn-1].Unknown > 0 {
		nunk := 0
		for _, i := range s.Map.TGrid {
			if i == UNKNOWN {
				nunk++
			}
		}
		ts.Unknown = nunk
	}
}

func (s *Statistics) ProcessDeadAnts(ts *TurnStatistics, deadants []PlayerLoc) {
	for _, pl := range deadants {
		ts.Died[pl.Player]++
		ts.DiedAll++
		s.DiedTot[pl.Player]++
		s.DiedMap[pl.Loc]++
	}
	s.DiedTotAll += ts.DiedAll

}

func (s *Statistics) ProcessSeen(ts *TurnStatistics, ants []PlayerLoc, turn int) {
	for _, pl := range ants {
		ts.Seen[pl.Player]++
		ts.SeenAll++
	}
	for i, n := range ts.Seen {
		if n >= s.SeenMax[i] {
			s.SeenMax[i] = n
			s.SeenMaxTurn[i] = turn
		}
	}
}
