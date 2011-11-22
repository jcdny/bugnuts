package state

import (
	"log"
	. "bugnuts/maps"
	. "bugnuts/util"
	. "bugnuts/pathing"
)

const (
	MAGICAGE = 10
)

func (s *State) ProcessFood(food []Location, turn int) {
	// Update food, possibly creating unseen food from symmetry stuff.
	for _, loc := range food {
		if _, ok := s.Food[loc]; !ok {
			// new food - also add via symmetry
		} else {
			s.Food[loc] = turn
		}
	}

	// Better would be to compute expected pickup time give neighbors
	// in the pathing step and only revert to this when there were no
	// visible neighbors.
	//
	// Should track that anyway since does not make sense to run for 
	// food another bot will certainly get unless its to enter combat.
	for loc, seen := range s.Food {
		if s.Met.Seen[loc] > seen || seen < turn-MAGICAGE {
			s.Food[loc] = 0, false
			if s.Map.Grid[loc] == FOOD {
				s.Map.Grid[loc] = LAND
			}
		}
	}
}

func (s *State) ProcessTurn(t *Turn) {
	s.Turn++

	s.TSet(t.W, WATER)
	s.ProcessVisible(t.A, 0, s.Turn)

	s.UpdateSymmetryData()

	s.ProcessFood(t.F, s.Turn)
	s.ProcessAnts(t.A, 0, s.Turn)
	s.ProcessHills(t.H, 0, s.Turn)
	s.ProcessDeadAnts(t.D, 0, s.Turn)

	if s.Turn == 1 {
		s.Turn1()
	}

	for player, ants := range s.Ants {
		for loc, seen := range ants {
			if seen < s.Met.Seen[loc] || (player == 0 && seen < s.Turn) {
				ants[loc] = 0, false
			} else {
				if seen < s.Turn && player != 0 {
					s.Met.Horizon[loc] = false
				}
				if s.Map.Grid[loc].IsHill() {
					s.Map.Grid[loc] = MY_HILLANT + Item(player)
				} else {
					s.Map.Grid[loc] = MY_ANT + Item(player)
				}
			}
			// TODO Bug here since if an ant steps out of seen we don't assume it still exists
			// unless it was out move that remove it from vision

			// TODO if the ant was visble last turn, not now and there is an ant
			// we can see 1 step away from where it was assume the new ant is
			// the same ant.

			// TODO Think about this -- assuming appearing ants match missing ones,
			// tracking max ants in a region.
		}
	}

	for loc, _ := range s.Ants[0] {
		// Update the one step land count and unseen count for my ants
		s.Met.SumVisCount(loc, s.ViewMask)
		for _, nloc := range s.Map.LocStep[loc] {
			s.Met.SumVisCount(nloc, s.ViewMask)
			if s.Met.Unknown[nloc] > 0 {
				s.Met.UpdateCounts(nloc, s.ViewMask)
			}
		}
	}

	s.Met.HBorder = s.StepHorizon(s.Met.HBorder)

	s.UpdateHillMaps()

	s.MonteCarloDensity()

	s.ComputeThreat(1, 0, s.AttackMask.MM, s.Met.Threat[len(s.Met.Threat)-1], s.Met.PThreat[len(s.Met.PThreat)-1])
}

// Given list of player/location update Land visible
// Also updates: Met.Unknown, Met.Land, Met.Seen, and Met.VisCount.
func (s *State) ProcessVisible(antloc []PlayerLoc, player, turn int) {
	for _, pl := range antloc {
		if pl.Player != player {
			continue
		}

		loc := pl.Loc
		unk := s.Met.Unknown[loc] != 0
		nland := 0

		if s.BorderDist[loc] > s.ViewMask.R {
			// In interior of map so use loc offsets
			for _, offset := range s.ViewMask.Loc {
				s.Met.Seen[loc+offset] = turn
				s.Met.VisCount[loc+offset]++
				if unk {
					if s.TGrid[loc+offset] == UNKNOWN {
						s.Map.TGrid[loc+offset] = LAND
						nland++
					} else if s.Map.TGrid[loc+offset] != WATER {
						nland++
					}
				}
			}
		} else {
			// non interior point lets go slow
			p := s.ToPoint(loc)
			for _, op := range s.ViewMask.P {
				l := s.ToLocation(s.PointAdd(p, op))
				s.Met.VisCount[l]++
				s.Met.Seen[l] = turn
				if unk {
					if s.TGrid[l] == UNKNOWN {
						s.Map.TGrid[l] = LAND
						nland++
					} else if s.TGrid[l] != WATER {
						nland++
					}
				}
			}
		}

		if unk {
			s.Met.Unknown[loc] = 0
			s.Met.Land[loc] = nland
		}
	}
}

func (s *State) ProcessHills(hl []PlayerLoc, player int, turn int) {

	for _, pl := range hl {
		if hill, found := s.Hills[pl.Loc]; found {
			hill.Player = pl.Player
			hill.Seen = turn
			hill.guess = false
		} else {
			s.Hills[pl.Loc] = &Hill{
				Location: pl.Loc,
				Player:   pl.Player,
				Found:    turn,
				Seen:     turn,
				Killed:   0,
				Killer:   -1,
				guess:    false,
			}
		}

	}

	// Update hill data in map.
	for i, _ := range s.NHills {
		s.NHills[i] = s.InitialHills
	}

	for loc, hill := range s.Hills {
		if hill.Killed == 0 && s.Met.Seen[loc] > hill.Seen {
			if hill.guess {
				// We just guessed at a location anyway, just remove it
				s.Hills[loc] = &Hill{}, false
			} else {
				hill.Killed = turn
			}
		}

		if hill.Killed > 0 {
			s.NHills[hill.Player]--
		} else if hill.Killed == 0 {
			if s.Met.Seen[loc] < turn {
				// If the hill is not visible then set Horizon false
				// since it could be a source of ant.
				s.Met.Horizon[loc] = false
			}
		}
	}

}

func (s *State) ProcessAnts(antloc []PlayerLoc, player, turn int) {
	for _, pl := range antloc {
		if s.Ants[pl.Player] == nil {
			s.Ants[pl.Player] = make(map[Location]int)
			// TODO New ant seen - start guessing hill loc
		}
		s.Ants[pl.Player][pl.Loc] = turn
	}
}

func (s *State) ResetGrid() {
	// Rotate threat maps and clear first.
	n := len(s.Met.Threat)

	if n > 1 {
		s.Met.Threat = append(s.Met.Threat[1:n], s.Met.Threat[0])
		s.Met.PThreat = append(s.Met.PThreat[1:n], s.Met.PThreat[0])
	}

	for i := range s.Met.Threat[0] {
		s.Met.Threat[0][i] = 0
		s.Met.PThreat[0][i] = 0
	}

	// Set all seen map to land
	for i, t := range s.Met.Seen {
		s.Met.VisCount[i] = 0
		if t == s.Turn && s.Map.Grid[i] > LAND {
			s.Map.Grid[i] = LAND
		}
	}
}

func (s *State) ProcessDeadAnts(deadants []PlayerLoc, player, turn int) {
	for _, pl := range deadants {
		if s.Stats.Dead[pl.Player] == nil {
			s.Stats.Dead[pl.Player] = make(map[Location]int)
		}
		s.Stats.Dead[pl.Player][pl.Loc]++
		s.Stats.Died[pl.Player]++
	}
	// TODO track suicides/sacrifices and who the killer was.
}

func (s *State) UpdateHillMaps() {
	// Generate the fill for all my hills.
	if s.NHills[0] == 0 {
		return
	}

	lend := make(map[Location]int)
	for _, hill := range s.HillLocations(0) {
		lend[hill] = 1
	}

	s.Met.FHill, _, _ = MapFillSeed(s.Map, lend, 1)

	outbound := make(map[Location]int)
	R := MaxV(MinV(s.Rows, s.Cols)/s.NHills[0]/2, 8)
	samples, _ := s.Met.FHill.Sample(0, R, R)

	for _, loc := range samples {
		outbound[loc] = 1
	}

	if len(outbound) < 1 {
		log.Panicf("UpdateHillMaps no outside border")
	} else {
		s.Met.FDownhill, _, _ = MapFillSeed(s.Map, outbound, 1)
		for i, d := range s.Met.FHill.Depth {
			if int(d) > R {
				s.Met.FDownhill.Depth[i] = 0
				s.Met.FDownhill.Seed[i] = 0
			}
		}
	}
}
