// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package maps

import (
	"testing"
	"log"
	"fmt"
	"time"
	. "bugnuts/torus"
	. "bugnuts/util"
	. "bugnuts/watcher"
)

func init() {
	//AllMaps = []string{"test"}
	log.SetFlags(log.Lshortfile)
}

func TestShiftReduce(t *testing.T) {
	T := Torus{Rows: 7, Cols: 7}

	p2 := Point{-1, 3}
	p, good := T.ShiftReduce(Point{}, p2, SYMMAXCELLS)
	if !good || p.R != 2 || p.C != 1 {
		t.Errorf("ShiftReduce: expected {2 1} got %v", good, p)
	}
}

func TestMirror(t *testing.T) {
	T := Torus{72, 72}
	p1 := Point{43, 0}
	p2 := Point{43, 65}
	m := T.Mirror(p1, p2, 1)
	if m.C != 33 {
		t.Errorf("Mirror: %v %v: %d", p1, p2, m)
	}
}

func TestTransMap(t *testing.T) {
	T := Torus{Rows: 7, Cols: 7}
	p := Point{-1, 3}
	m := T.TransMap(p, SYMMAXCELLS)
	log.Printf("TransMap %v", m)
}

func BenchmarkTransMap(b *testing.B) {
	// random_walk_07p_02
	T := Torus{Rows: 119, Cols: 147}
	p := Point{34, -21}
	for i := 0; i < b.N; i++ {
		T.TransMap(p, SYMMAXCELLS)
	}
}

func TestTile(t *testing.T) {
	// See end of file for expected data...
	m := mapMeBaby("./testdata/sym.map")
	s := m.Tile(0)
	sym := s.MinHash
	smap := s.Hashes
	str := ""

	test := true // set this to false to emit expected values...

	if !test {
		// Generate the expected data
		log.Printf("%#v", sym)
		for i := range smap {
			str += fmt.Sprintf("%#v,\n", smap[i])
		}
		log.Print("\n" + str)
	}
	if test {
		mismatch := false
		str = ""
		for loc := range sym {
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
			for loc, v := range smap {
				ve := expectmap[loc]
				//log.Printf("%d:%#v", k, *v)
				for i := range v {
					if v[i] != ve[i] {
						t.Errorf("Value mismatch loc %d, %v vs %v", loc, *v, ve)
					}
				}
			}
		}
	}
}

func mapMeBaby(name string) *Map {
	var file string
	if name[0] != '.' {
		file = MapFile(name)
	} else {
		file = name
	}
	m, err := MapLoadFile(file)
	// zap non land out.
	for i, item := range m.Grid {
		if !(item == LAND || item == WATER) {
			m.Grid[i] = LAND
		}
	}
	copy(m.TGrid, m.Grid) // need TGrid populated since this is all driven by the "observed grid"
	if m == nil || err != nil {
		log.Panicf("Error reading %s: err %v map: %v", file, err, m)
	}

	return m
}

var benchMap string = "mmaze_05p_01"

func BenchmarkTile0(b *testing.B) {
	b.StopTimer()
	m := mapMeBaby(benchMap)
	s := m.NewSymData(0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for loc := range m.Grid {
			s.Update(Location(loc))
		}
	}
}
func BenchmarkTile4(b *testing.B) {
	b.StopTimer()
	m := mapMeBaby(benchMap)
	s := m.NewSymData(4)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for loc := range m.Grid {
			s.Update(Location(loc))
		}
	}
}
func BenchmarkTile8(b *testing.B) {
	b.StopTimer()
	m := mapMeBaby(benchMap)
	s := m.NewSymData(8)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for loc := range m.Grid {
			s.Update(Location(loc))
		}
	}
}

func TestZSym(t *testing.T) {
	//AllMaps := []string{"maze_p04_31"}
	for _, name := range AllMaps {
		log.Printf("***************************  %s ***************************************************", name)
		m := mapMeBaby(name)
		if m == nil {
			t.Error("Map nil")
		}
		TPush("Tile " + name)
		sym := m.Tile(0)
		TPop()
		if sym == nil {
			t.Error("Sym nil")
		}

		log.Printf("MAP %s Tiles: %d entries rows %d cols %d", name, len(sym.Tiles), m.Rows, m.Cols)
		log.Printf("NLen: %v", sym.NLen)

		peak := Max(sym.NLen[:])
		var cycle, n int
		for cycle, n = range sym.NLen {
			if n == peak {
				break
			}
		}
		if len(sym.Tiles) > 0 {
			bad := 0
			done := 0
			for _, tile := range sym.Tiles {
				if len(tile.Locs) == cycle {
					eqlen := sym.SymAnalyze(tile)
					if done == 0 && cycle == eqlen && tile.Self == 0 {
						log.Print("MMM ", tile.Gen, "::\"", name, "\":MapSymData{", cycle, ",", eqlen, ",", tile.Gen,
							",\"", symLabels[tile.Gen], "\",", sym.Size()/tile.Subtile.Size(),
							",", tile.Subtile.Rows, ",", tile.Subtile.Cols,
							",Point{", tile.Translate.R, ",", tile.Translate.C, "}",
							",", tile.RM1, ",", tile.RM2, ",", tile.MR, ",", tile.MC,
							",Point{", tile.Origin.R, ",", tile.Origin.C, "}",
							"},")

						done++
					}
					if false && cycle != eqlen && tile.Self == 0 {
						if bad < 4 {
							log.Print("Len: ", eqlen, " Tile::", tile, "\n\t", m.ToPoints(tile.Locs))
							if len(AllMaps) == 0 {
								sym.symdump(tile.Hash, m)
							}
						}
						bad++
						//sym.symdump(minhash, m)
					}
				}
			}
			log.Print("Mismatched: ", bad, " ", name, " ", peak, cycle)
		}
	}
}

func TestSymMatch(t *testing.T) {
	//Debug[DBG_Symmetry] = true
	WS = NewWatches(0, 0, 0)
	//AllMaps := []string{"test"}
	for name, sd := range MapSym {
		if false && sd.Gen != SYMRMBOTH {
			continue
		}
		m := mapMeBaby(name)
		if m == nil {
			t.Error("Map nil")
		}
		log.Print(name)
		TPush("UpdateSymmetryData:" + name)
		m.SymData.UpdateSymmetryData(time.Nanoseconds() + 5*1e12)
		TPop()

		peak := Max(m.SymData.NLen[:])
		var found, cycle, n int
		for cycle, n = range m.SymData.NLen {
			if n == peak {
				break
			}
		}

		if len(m.SMap) > 0 {
			found = len(m.SMap[0])
		}

		log.Print("MAP ", name, " ", m.Rows, " ", m.Cols, " Cycle: ", cycle, " found ", found, found == cycle, "\n", sd)
	}
}

// Fancy dump of Symmetry information including reduced translation
// offsets, matching symmetry,
func (sym *SymData) symdump(tile SymHash, m *Map) {
	llist := sym.Tiles[tile].Locs
	redlist := make([]Point, 0, 8)

	str := "\n"
	for _, l1 := range llist {
		// offset matrix
		p1 := sym.ToPoint(l1)
		for _, l2 := range llist {
			var pd Point
			p2 := sym.ToPoint(l2)
			good := false
			if l1 != l2 {
				if sym.Hashes[l1][0] == sym.Hashes[l2][0] {
					pd, good = m.ShiftReduce(p1, p2, SYMMAXCELLS)
					if good {
						redlist = append(redlist, pd)
					}
				} else {
					pd = m.SymDiff(l1, l2)
				}
				str += fmt.Sprintf("   [%3d%4d]   |", pd.R, pd.C)
			} else {
				str += "               |"
			}
		}
		str += "\n"

		// pairwise symmetry
		for _, l2 := range llist {
			sid := ""
			for i2 := uint8(0); i2 < 8; i2++ {
				if sym.Hashes[l1][0] == sym.Hashes[l2][i2] {
					if l1 == l2 {
						// translation sym
						sid += "IDENT"
						break
					} else {
						sid += fmt.Sprintf("%s", symAxesMap[i2].Name)
					}
				}
			}
			str += fmt.Sprintf("%15s|", sid)

		}
		str += "\n"

		if true {
			str += "    "
			// Symmetry matrix
			for _ = range llist {
				for i1 := uint8(0); i1 < 8; i1++ {
					str += fmt.Sprintf("%d", i1)
				}
				str += "   |    "
			}
			str += "\n"
			for i1 := uint8(0); i1 < 8; i1++ {
				str += fmt.Sprintf("%d___", i1)
				for _, l2 := range llist {
					for i2 := uint8(0); i2 < 8; i2++ {
						if sym.Hashes[l1][i1] == sym.Hashes[l2][i2] {
							str += "*"
						} else {
							str += " "
						}
					}
					str += "___|____"
				}
				str += "\n"
			}
		}
	}
	if len(redlist) > 0 {
		redlist = m.ReduceReduce(redlist)
	}
	log.Printf("Tile %d %v\nReduce: %v\n%s\n", tile, m.ToPoints(llist), redlist, str)
}
