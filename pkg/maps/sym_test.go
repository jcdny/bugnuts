package maps

import (
	"testing"
	"log"
	"fmt"
	"os"
)

var expect = []int{64, 2, 0, 0, 0, 4, 128, 4096, 128, 1536, 64, 128, 2, 1, 0, 0, 0, 2, 64, 128, 64, 48, 2, 4, 0, 0, 0, 0, 0, 1, 2, 4, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 1, 0, 1, 2, 4, 2, 1, 0, 1, 2, 4, 64, 2, 0, 2, 64, 128, 64, 2, 0, 2, 64, 128, 128, 4, 0, 4, 128, 4096, 128, 4, 0, 4, 128, 4096, 64, 2, 0, 2, 64, 128, 64, 2, 0, 2, 64, 128, 2, 1, 0, 1, 2, 4, 2, 1, 0, 1, 2, 4, 2, 1, 0, 0, 0, 0, 0, 0, 0, 1, 2, 4, 64, 2, 0, 0, 0, 1, 2, 4, 2, 48, 64, 128, 128, 4, 0, 0, 0, 2, 64, 128, 64, 1536, 128, 4096}
var expectmap = map[int][8]int{
	48:   [8]int{1572864, 16777224, 48, 2097153, 16809984, 513, 8388624, 1048578},
	128:  [8]int{8192, 128, 2048, 131072, 2048, 8192, 131072, 128},
	1536: [8]int{49152, 8388612, 1536, 4194306, 525312, 16416, 4194312, 2097156},
	0:    [8]int{0, 0, 0, 0, 0, 0, 0, 0},
	4:    [8]int{1024, 4194304, 16384, 4, 16384, 1024, 4, 4194304},
	4096: [8]int{4096, 4096, 4096, 4096, 4096, 4096, 4096, 4096},
	1:    [8]int{16777216, 16, 1, 1048576, 1048576, 16, 16777216, 1},
	64:   [8]int{65536, 262144, 256, 64, 262144, 64, 256, 65536},
	2:    [8]int{524288, 8, 32, 2097152, 32768, 512, 8388608, 2},
}

func TestTile(t *testing.T) {
	m := mapMeBaby("testdata/sym.map")
	sym, smap, _ := m.tile()

	str := ""
	mismatch := false
	for loc, _ := range sym {
		if loc%m.Cols == 0 {
			str += fmt.Sprintf("\n%2d :", loc/m.Cols)
		}
		if sym[loc] != expect[loc] {
			mismatch = true
			str += fmt.Sprintf("%5d*", sym[loc])
		} else {
			str += fmt.Sprintf("%6d", sym[loc])
		}
	}

	if mismatch {
		t.Error("Invalid symmetry point")
		log.Printf("got:\n%s", str)
	}

	if len(expectmap) != len(smap) {
		t.Error("Sym Map length mismatch %d vs %d", len(smap), len(expectmap))
	} else {
		for k, v := range smap {
			ve, ok := expectmap[k]
			if !ok {
				t.Errorf("Key missing from smap %d", k)
			} else {
				if len(v) != len(ve) {
					t.Errorf("Values different length key %d", k)
				} else {
					//log.Printf("%d:%#v", k, *v)
					for i, _ := range v {
						if v[i] != ve[i] {
							t.Errorf("Value mismatch key %d, %v vs %v", k, *v, ve)
						}
					}
				}
			}
		}
	}
}

func (m *Map) tile() ([]int, map[int]*[8]int, map[int][]Location) {
	sym := make([]int, len(m.Grid))
	smap := make(map[int]*[8]int)
	sloc := make(map[int][]Location)
	for loc, _ := range sym {
		sval, i8 := m.SymCompute(Location(loc))
		sym[loc] = sval
		smap[sval] = i8
		_, ok := sloc[sval]
		if !ok {
			sloc[sval] = make([]Location, 0)
		}
		sloc[sval] = append(sloc[sval], Location(loc))
	}

	return sym, smap, sloc
}

type SymData struct {
	Sym       []int           // Sym data for a given point.
	SymRotate map[int]*[8]int // Map from the min in to matching rotations
	SymTiles  map[int][]Location
}

func mapMeBaby(file string) *Map {
	m, err := MapLoadFile(file)
	if m == nil || err != os.EOF {
		log.Panicf("Error reading %s: err %v map: %v", file, err, m)
	}

	return m
}

func BenchmarkTile(b *testing.B) {
	m := mapMeBaby("testdata/maps/mmaze_05p_01.map")
	for i := 0; i < b.N; i++ {
		m.tile()
	}
}

func (m *Map) Symmie(sym *SymData) *SymData {
	// Initial setup
	if sym == nil {
		sym = &SymData{}
		sym.Sym, sym.SymRotate, sym.SymTiles = m.tile()
	}

	return sym
}

func TestSymmie(t *testing.T) {
	for _, name := range AllMaps {
		m := mapMeBaby("testdata/maps/" + name + ".map")
		if m == nil {
			t.Error("Map nil")
		}
		sym := m.Symmie(nil)
		if sym == nil {
			t.Error("Sym nil")
		}

		log.Printf("MAP %s Sym: SymRotate: %d entries SymTiles: %d entries", name, len(sym.SymRotate), len(sym.SymTiles))
		done := false
		for tile, llist := range sym.SymTiles {
			if !done && len(llist) < 20 {
				done = true
				log.Printf("Tile %d: %v", tile, m.ToPoints(llist))
				str := "\n"
				for _, l1 := range llist {
					for _, l2 := range llist {
						pd := m.SymDiff(l1, l2)
						str += fmt.Sprintf(" [%3d%4d]", pd.R, pd.C)
					}
					str += "\n"
				}
				log.Print(str)
			}
		}
	}
}

func (m *Map) SymDiff(l1, l2 Location) Point {
	p1 := m.ToPoint(l1)
	p2 := m.ToPoint(l2)

	r := p2.R - p1.R
	c := p2.C - p1.C

	if r > m.Rows/2 {
		r -= m.Rows
	}
	if r < -m.Rows/2 {
		r += m.Rows
	}
	if c > m.Cols/2 {
		c -= m.Cols
	}
	if c < -m.Cols/2 {
		c += m.Cols
	}
	if r < 0 {
		r = -r
		c = -c
	}
	return Point{R: r, C: c}
}
