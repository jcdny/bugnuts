package combat

// See benchmark/combat_test.go for functional testing.

import (
	"testing"
	"reflect"
	"rand"
	"bugnuts/game"
	"bugnuts/state"
	"bugnuts/torus"
)

const (
	NP   = 4
	NANT = 100
)

func TestCopy(t *testing.T) {
	g := game.NewGameInfo(100, 100)
	s := state.NewState(g)
	ants := make(map[torus.Location]struct{}, NP*NANT)
	for len(ants) < NP*NANT {
		ants[torus.Location(rand.Intn(s.Map.Size()))] = struct{}{}
	}
	for i := 0; i < 4; i++ {
		s.Ants[i] = make(map[torus.Location]int, 100)
	}
	j := 0
	for loc := range ants {
		s.Ants[j%4][loc] = 1
		j++
	}

	c := NewCombat(s.Map, s.AttackMask, NP)
	c.Setup(s.Ants)
	cc := c.Copy()

	if !reflect.DeepEqual(c, cc) {
		t.Error("Copy not deep equal")
	}

	cc.Threat[0] = 9
	eq, diff := CombatCheck(c, cc)
	if eq {
		t.Error("Altered copy deep equal")
	} else {
		if _, ok := diff["Threat"]; !ok || len(diff) > 1 {
			t.Error("Expected Threat as diff got ", diff)
		}
	}
	cc.Threat[0] = c.Threat[0]

	c.PThreat[1][0] = 9
	eq, diff = CombatCheck(c, cc)
	if eq {
		t.Error("Altered copy deep equal")
	} else {
		if _, ok := diff["PThreat[1]"]; !ok || len(diff) > 1 {
			t.Error("Expected PThreat[1] as diff got ", diff)
		}
	}
}
