package monte

import (
	"log"
	"testing"
	"bugnuts/engine"
	. "bugnuts/maps"
	. "bugnuts/game"
	. "bugnuts/torus"
	. "bugnuts/combat"
)

type testData struct {
	in, user string
	turn     int
	part     Point
}

var tests = []testData{
	{"testdata/test1/0.replay", "bot8", 30, Point{28, 52}},
}

func TestMonteCarlo(t *testing.T) {
	for _, t := range tests {
		log.Print("Running test for ", t.in, "(", t.user, ") turn ", t.turn, " partition ", t.part)
		combatMe(t)
	}
}

func combatMe(t testData) {
	s := engine.LoadState(t.in, t.user, t.turn)
	if s != nil {
		ap, pmap := CombatPartition(s)
		log.Printf("%v", ap)
		log.Printf("%v", pmap)
		pkey := s.ToLocation(t.part)
		p, ok := ap[pmap[pkey][0]]
		if !ok {
			log.Print("Could not find partition for ", t.part)
		} else {
			log.Print("Simulate for ", p)

			c := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
			c.Setup(s.Ants)
			// players
			am, _ := c.PartitionMoves(p)
			DumpAntMove(s.Map, am, -1, t.turn)
		}
	}

	log.Print(genperm(4))
}

// generate the list of permuted directions for n players
func genperm(n uint) [][]Direction {
	nperm := uint(4) << (2 * (n - 1))
	log.Print("nperm ", nperm)
	dl := make([]Direction, nperm*n)
	out := make([][]Direction, nperm)
	for i := uint(0); i < nperm; i++ {
		for s := uint(0); s < n; s++ {
			dl[i*n+s] = Direction((i >> (2 * s)) & 3)
		}
		out[i] = dl[i*n : (i+1)*n]
	}

	return out
}
