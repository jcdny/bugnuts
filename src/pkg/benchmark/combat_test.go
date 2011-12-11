package bench

import (
	"log"
	"testing"
	"bugnuts/engine"
	. "bugnuts/util"
	. "bugnuts/torus"
	. "bugnuts/combat"
)

type testData struct {
	in, user string
	turn     int
	part     Point
}

var tests = []testData{
	//{"testdata/test1/0.replay", "bot8", 30, Point{28, 52}},  // 6 me 2 other
	{"../monte/testdata/test1/0.replay", "bot8", 30, Point{49, 19}}, // 5 me, 3, 1
	{"../pathing/testdata/replay.big.json.gz", "MomoBot", 300, Point{75, 102}},
}

func TestPartitions(t *testing.T) {
	for _, td := range tests {
		log.Print("Running test for ", td.in, "(", td.user, ") turn ", td.turn, " partition ", td.part)
		combatMe(td)
	}
}

func BenchmarkCombatSetup(b *testing.B) {
	b.StopTimer()
	t := tests[1]
	s := engine.LoadState(t.in, t.user, t.turn)
	if s != nil {
		c := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c.Setup(s.Ants)
		}
	}
}

func combatMe(t testData) {
	s := engine.LoadState(t.in, t.user, t.turn)
	if s != nil {
		c := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		c.Setup(s.Ants)

		TPush("** Partition " + t.in)
		ap, pmap := CombatPartition(s)
		TPop()

		log.Printf("Found %d partitions with %d mappings", len(ap), len(pmap))
		pkey := s.ToLocation(t.part)
		p, ok := ap[pmap[pkey][0]]
		if !ok {
			log.Print("Could not find partition for ", t.part)
		} else {
			ps := NewPartitionState(c, p)
			log.Print("Partition ", pmap[pkey][0], ":", DumpPartitionState(ps))
		}
	}
}
