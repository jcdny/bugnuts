package state

import (
	"log"
	"strconv"
	. "bugnuts/maps"
	. "bugnuts/util"
)

func (s *State) FoodLocations() (l []Location) {
	for loc, _ := range s.Food {
		l = append(l, Location(loc))
	}

	return l
}

func (s *State) HillLocations(player int) (l []Location) {
	for loc, hill := range s.Hills {
		if hill.Player == player && hill.Killed == 0 {
			l = append(l, Location(loc))
		}
	}

	return l
}

func (s *State) EnemyHillLocations(player int) (l []Location) {
	for loc, hill := range s.Hills {
		if hill.Player != player && hill.Killed == 0 {
			l = append(l, Location(loc))
		}
	}

	return l
}

func (s *State) Threat(turn int, l Location) int8 {
	i := len(s.Met.Threat) - turn + s.Turn - 1
	if i < 0 {
		log.Printf("Threat for turn %d on turn %d we only keep %d turns", turn, s.Turn, len(s.Met.Threat))
		return 0
	}
	return s.Met.Threat[i][l]
}

func (s *State) PThreat(turn int, l Location) uint16 {
	i := len(s.Met.PThreat) - turn + s.Turn - 1
	if i < 0 {
		log.Printf("Threat for turn %d on turn %d we only keep %d turns", turn, s.Turn, len(s.Met.Threat))
		return 0
	}
	return s.Met.PThreat[i][l]
}

func (s *State) ThreatMap(turn int) []int8 {
	i := len(s.Met.Threat) - turn + s.Turn - 1
	if i < 0 {
		log.Printf("Threat for turn %d on turn %d we only keep %d turns", turn, s.Turn, len(s.Met.Threat))
		return nil
	}
	return s.Met.Threat[i]
}

func (s *State) ValidStep(loc Location) bool {
	i := s.Map.Grid[loc]

	return i != WATER && i != BLOCK && i != OCCUPIED && i != FOOD && i != MY_ANT && i != MY_HILLANT
}

func (s *State) Stepable(loc Location) bool {
	i := s.Map.Grid[loc]

	return i != WATER && i != BLOCK && i != FOOD
}

func (m *Metrics) DumpSeen() string {
	max := Max(m.Seen)
	str := ""

	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			str += strconv.Itoa(m.Seen[r*m.Cols+c] * 10 / (max + 1))
		}
		str += "\n"
	}

	return str
}
