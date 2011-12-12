package state

import (
	"strconv"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/util"
)

func (s *State) FoodLocations() (l []Location) {
	for loc := range s.Food {
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
