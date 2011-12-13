package combat

// See benchmark/combat_test.go for functional testing.

import (
	"log"
	"testing"
	"rand"
	"bugnuts/maps"
	"bugnuts/torus"
	"bugnuts/game"
)

const (
	chatty = false
)

const (
	NP   = 4
	NANT = 800
	ROWS = 200
	COLS = 125
)

func makeAnts(m *maps.Map, np, nant int) []map[torus.Location]int {
	ants := make(map[torus.Location]struct{}, np*nant)
	for len(ants) < np*nant {
		ants[torus.Location(rand.Intn(m.Size()))] = struct{}{}
	}
	Ants := make([]map[torus.Location]int, np)
	for i := 0; i < np; i++ {
		Ants[i] = make(map[torus.Location]int, 100)
	}
	j := 0
	for loc := range ants {
		Ants[j%np][loc] = 1
		j++
	}

	return Ants
}

func stageCombat(np, rows, cols int) (c *Combat, m *maps.Map, mask *maps.Mask) {
	g := game.NewGameInfo(rows, cols)
	m = maps.NewMap(g.Rows, g.Cols, np)
	mask = maps.MakeMask(g.AttackRadius2, g.Rows, g.Cols)
	c = NewCombat(m, mask, np)

	return
}

func TestCopy(t *testing.T) {
	c, m, _ := stageCombat(NP, ROWS, COLS)
	Ants := makeAnts(m, NP, NANT)

	c.Setup(Ants)
	cc := c.Copy()

	/* No longer deep equal since we only copy the
	 * stuff necessary for playing combat forward in the
	 * engine...

		if !reflect.DeepEqual(c, cc) {
			t.Error("Copy not deep equal")
		}
	*/
	eq, diff := CombatCheck(c, cc)
	if !eq {
		t.Error("Unaltered copy not deep equal", diff)
	}

	cc.Threat[0] = 9
	eq, diff = CombatCheck(c, cc)
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

func TestRiskMark(t *testing.T) {
	_, m, mask := stageCombat(NP, ROWS, COLS)
	c := NewCombat(m, mask, NP)
	Ants := makeAnts(m, NP, NANT)
	c.Setup(Ants)
	rm := RiskMark(m, &mask.Offsets, Ants, c.Ants1, c.Threat1, c.PThreat1)
	for i := 0; i < NP; i++ {
		n := make([]int, MaxRiskStat)
		for _, r := range rm[i] {
			n[r]++
		}
		if chatty {
			log.Print("Player ", i, " ants ", len(Ants[i]), " nrisk: ", len(rm[i]), " totals ", n)
		}
	}
}
func BenchmarkRiskMark(b *testing.B) {
	b.StopTimer()
	_, m, mask := stageCombat(NP, ROWS, COLS)
	c := NewCombat(m, mask, NP)
	Ants := makeAnts(m, NP, NANT)
	c.Setup(Ants)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		RiskMark(m, &mask.Offsets, Ants, c.Ants1, c.Threat1, c.PThreat1)
	}
}

func BenchmarkSetup(b *testing.B) {
	_, m, mask := stageCombat(NP, ROWS, COLS)
	c := NewCombat(m, mask, NP)
	Ants := makeAnts(m, NP, NANT)
	for i := 0; i < b.N; i++ {
		c.Reset()
		c.Setup(Ants)
	}
}

func BenchmarkSetupAlloc(b *testing.B) {
	_, m, mask := stageCombat(NP, ROWS, COLS)
	c := NewCombat(m, mask, NP)
	Ants := makeAnts(m, NP, NANT)
	for i := 0; i < b.N; i++ {
		c.ResetAlloc()
		c.Setup(Ants)
	}
}
