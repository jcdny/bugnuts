// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package state

import (
	"log"
	"math"
	. "bugnuts/game"
	. "bugnuts/maps"
	. "bugnuts/combat"
	. "bugnuts/torus"
	. "bugnuts/util"
	. "bugnuts/watcher"
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
	Suicide     [MaxPlayers]int
	StaticCount [MaxPlayers]int // Count of unmoved ants
	PRisk       [MaxPlayers][MaxRiskStat]int
	Score       [MaxPlayers]float64 // a scaling factor for combat
}

type Statistics struct {
	NP      int
	DiedMap []int // Death count per location
	//StaticMap      []int            // Turns given cell has been occupied by the same player.
	//PlayerMap      []int            // used to drive static map
	RiskTot        [MaxPlayers][MaxRiskStat]int
	SuicideTot     [MaxPlayers]int
	DiedTot        [MaxPlayers]int  // Chronicle of a Death Foretold, Running total of deaths seen
	DiedTotAll     int              // Total number of deaths seen this game
	SeenMax        [MaxPlayers]int  // Maximum number seen
	SeenMaxTurn    [MaxPlayers]int  // Turn on which max number was seen
	HorizonMax     int              // Maximum known state extent
	HorizonMaxTurn int              // turn on which we had max knowledge
	TStats         []TurnStatistics // per turn statistics
	LTS            *TurnStatistics  // last computed turn stats.
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
	s.Stats.ProcessSeen(ts, turn.A, turn.Turn, s.Cprev)
	s.Stats.ProcessDeadAnts(ts, turn.D, s.Cprev)
	if Debug[DBG_Statistics] {
		log.Print("Rtot  ", s.Stats.RiskTot[0:s.Stats.NP])
		log.Print("PRisk ", ts.PRisk[0:s.Stats.NP])
	}

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

	s.Stats.LTS = ts
	s.Stats.ComputeScore(s.Cprev)
}

func (s *Statistics) ComputeScore(c *Combat) {
	// we are more aggresive towards players we outnumber esp if 
	// we have the most ant however if we have not seen much of the
	// map we are risk averse 
	if s.LTS == nil {
		return
	}
	// fraction of map seen
	vis := .1
	if c != nil {
		vis = float64(c.Size()-s.LTS.Horizon) / float64(c.Size())
	}

	frac := float64(s.LTS.Seen[0]) / float64(s.LTS.SeenAll)

	for i := range s.LTS.Seen {
		if TurnGet() < 60 { // MAGIC
			s.LTS.Score[i] = .4
		} else {
			frace := float64(s.LTS.Seen[i]) / float64(s.LTS.SeenAll)
			s.LTS.Score[i] = math.Fmax(math.Fmin(frac*vis*vis/frace, 2.0), .3)
		}
	}
	s.LTS.Score[0] = -1.5

	if Debug[DBG_Scoring] {
		log.Print("Score per player is ", s.LTS.Score[0:s.NP])
	}
}

func (s *Statistics) ProcessDeadAnts(ts *TurnStatistics, deadants []PlayerLoc, c *Combat) {
	if len(deadants) == 0 {
		return
	}

	suicide := make(map[Location]int, 20)
	for _, pl := range deadants {
		ts.Died[pl.Player]++
		ts.DiedAll++
		s.DiedTot[pl.Player]++
		s.DiedMap[pl.Loc]++
		if _, ok := suicide[pl.Loc]; ok {
			ts.Suicide[pl.Player]++
			s.SuicideTot[pl.Player]++
			suicide[pl.Loc]++
			if Debug[DBG_Statistics] {
				log.Print("suicide", pl)
			}
		} else {
			suicide[pl.Loc] = 1
		}
	}
	s.DiedTotAll += ts.DiedAll

	if c != nil {
		for _, pl := range deadants {
			if n, ok := suicide[pl.Loc]; ok && n == 1 {
				r, ok := c.Risk[pl.Player][pl.Loc]
				if !ok {
					r = RiskNone
				} else if Debug[DBG_Statistics] {
					log.Print("rdeath ", r, pl)
				}
				s.RiskTot[pl.Player][r]++
			}
		}
	}

}

func (s *Statistics) ProcessSeen(ts *TurnStatistics, ants []PlayerLoc, turn int, c *Combat) {
	for _, pl := range ants {
		ts.Seen[pl.Player]++
		ts.SeenAll++
	}

	for i, n := range ts.Seen {
		if n >= s.SeenMax[i] {
			if n > 0 {
				s.NP = MaxV(s.NP, i+1)
			}
			s.SeenMax[i] = n
			s.SeenMaxTurn[i] = turn
		}
	}

	if c != nil {
		for _, pl := range ants {
			r, ok := c.Risk[pl.Player][pl.Loc]
			if !ok {
				r = RiskNone
			} else {
				// log.Print("rlife ",r, pl)
			}
			ts.PRisk[pl.Player][r]++
		}
	}
}
