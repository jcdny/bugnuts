package state

import (
	. "bugnuts/torus"
	. "bugnuts/pathing"
	. "bugnuts/util"
)

// Update expected locations and flows for enemy ants
func (s *State) MonteCarloDensity() {
	tgt := make(map[Location]int, 32)
	for _, loc := range s.HillLocations(0) {
		tgt[loc] = 1
	}

	if false {
		for _, loc := range s.FoodLocations() {
			tgt[loc] = 20
		}
	}
	if len(tgt) > 0 {
		ants := make([]Location, 0, 200)
		f, _, _ := MapFillSeed(s.Map, tgt, 0)

		for player := 1; player < len(s.Ants); player++ {
			for loc, _ := range s.Ants[player] {
				endloc := f.Seed[loc]
				steps := f.Depth[loc] - f.Depth[f.Seed[loc]]
				hill, ok := s.Hills[endloc]
				if !ok || hill.Player != player {
					if (ok && steps < 50) || steps < 16 {
						ants = append(ants, Location(loc))
					}
				}
			}
		}

		if len(ants) > 0 {
			paths := 64
			for paths*len(ants) > 2048 {
				paths = paths >> 1
			}
			s.Met.MCDist, s.Met.MCFlow = f.MontePathIn(s.Map, ants, paths, 1)
			s.Met.MCDistMax = Max(s.Met.MCDist)
			s.Met.MCPaths = paths
		} else {
			s.Met.MCPaths = 0
		}
	}
}
