package combat

import (
	"testing"
	"log"
	"bugnuts/maps"
	"bugnuts/replay"
)

var M *maps.Map
var A *Arena

func init() {
	file := "../replay/testdata/replay.0.json"
	match, err := replay.Load(file)
	if err != nil {
		log.Panicf("Load of %s failed %v", file, err)
	}
	M = match.Replay.GetMap()

	A = NewArena(M, maps.Location(9*96+11))
	log.Printf("Arena:\n%v", A)
}

func TestArena(t *testing.T) {
	files := []string{
		"../replay/testdata/replay.0.json",
	}

	for _, file := range files {
		match, err := replay.Load(file)
		if err != nil {
			t.Errorf("Load of %s failed %v", file, err)
		}
		m := match.Replay.GetMap()

		a := NewArena(m, maps.Location(9*96+11))
		log.Printf("Arena:\n%v", a)
	}
}

func BenchmarkMonte(b *testing.B) {
	start := ALoc(40)
	l := start
	n := 0
	tn := 0
	s := 0
	for i := 0; i < b.N*10000; i++ {
		var d maps.Direction
		for _, d = range maps.Permute5() {
			s++
			if d == maps.NoMovement {
				break
			} else if nl := A.LocStep[l][d]; maps.StepableItem[A.Grid[nl]] {
				l = nl
				break
			}
		}
		if l == 0 {
			l = start
			n++
			tn += s
			s = 0
		}
	}
	log.Printf("mean steps to exit %.2f exited %d Steps %d %d", float64(tn)/float64(n), n, tn, b.N)
}

func BenchmarkMontef(b *testing.B) {
	start := maps.Location(136)
	l := start
	n := 0
	tn := 0
	s := 0
	for i := 0; i < b.N*10000; i++ {
		s++
		var d maps.Direction
		for _, d = range maps.Permute5() {
			if d == maps.NoMovement {
				break
			} else if nl := M.LocStep[l][d]; maps.StepableItem[M.Grid[nl]] {
				l = nl
				break
			}
		}
		if l == 0 {
			l = start
			n++
			tn += s
			s = 0
		}
	}
	log.Printf("mean steps to exit %.2f", float64(tn)/float64(n))
}
