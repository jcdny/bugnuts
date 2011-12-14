package monte

import (
	"log"
	"testing"
	"bugnuts/engine"
	. "bugnuts/torus"
	. "bugnuts/combat"
	. "bugnuts/game"
)

type testData struct {
	in, user string
	turn     int
	part     Point
}

var tests = []testData{
	//{"testdata/test1/0.replay", "bot8", 30, Point{28, 52}},  // 6 me 2 other
	{"testdata/test1/0.replay", "bot8", 30, Point{49, 19}}, // 5 me, 3, 1
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestGenPerm(t *testing.T) {
	if len(genperm(2)) != 16 ||
		len(genperm(3)) != 64 ||
		len(genperm(4)) != 256 {
		t.Error("genperm length wrong")
	}
}

func BenchmarkMonteCarlo(b *testing.B) {
	for _, td := range tests {
		benchMe(td, b)
	}
}

func benchMe(t testData, b *testing.B) {
	s := engine.LoadState(t.in, t.user, t.turn)
	if s != nil {
		c := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		c.Setup(s.Ants)

		ap, pmap := CombatPartition(s)
		pkey := s.ToLocation(t.part)
		p, ok := ap[pmap[pkey][0]]
		if !ok {
			log.Print("Could not find partition for ", t.part)
		} else {
			ps := NewPartitionState(c, p)
			perm := genperm(uint(ps.PLive))
			move := make([][]AntMove, len(perm))
			fs := FirstStep(s.Map, c, ps)

			for ip, p := range perm {
				move[ip] = make([]AntMove, ps.ALive)
				ib, ie := 0, 0
				for np := 0; np < ps.PLive; np++ {
					ib, ie = ie, ie+len(fs[np][p[np]])
					copy(move[ip][ib:ie], fs[np][p[np]])
				}
			}
			for i := 0; i < b.N; i++ {
				for ip := range move {
					c.Resolve(move[ip])
					c.Unresolve(move[ip])
				}
			}
		}
	}
}
