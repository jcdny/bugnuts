package maps

import (
	"testing"
	"log"
	"fmt"
	"os"
)

func TestMapNew(t *testing.T) {
	r, c := 5, 7
	m := NewMap(r, c, 1)

	if m.Size() != r*c {
		t.Errorf("Map size error")
	}

	for i, d := range m.BorderDist {
		if m.BorderDist[m.Size()-i-1] != d {
			t.Errorf("Border Distance Error, Not Symmetric: %v", m.BorderDist)
			break
		}
	}

	if m.BorderDist[0] != 1 ||
		m.BorderDist[m.Size()-1] != 1 ||
		m.BorderDist[c+2] != 2 {
		t.Errorf("Border Distance Error bounds wrong: %v", m.BorderDist)
	}
}

func TestMapLoad(t *testing.T) {

	file := "testdata/big.map"
	m, err := MapLoadFile(file)

	if err != nil {
		t.Errorf("Invalid load of map error os.Error == %v", err)
	}
	if m == nil {
		t.Errorf("Invalid load of map m == nil")
	}

	//m.WriteDebugImage("_maptest", 0, func(c, r int) image.NRGBAColor { return m.At(r, c) })
}

func TestMapId(t *testing.T) {
	out, err := os.Create("testdata/maps.csv.tmp")
	if err != nil {
		log.Panic("Open of testdata/maps.csv.tmp failed ", err)
	}
	defer out.Close()

	for _, name := range AllMaps {
		file := MapFile(name)
		m, err := MapLoadFile(file)
		if m == nil || err != nil {
			log.Panicf("Error reading %s: err %v map: %v", file, err, m)
		}
		mapid := m.MapId()
		fmt.Fprintf(out, "\"%s\",\"%s\",%d,%d,%d,%d,\"%v\"\n", mapid, name, m.Rows, m.Cols, m.Players, len(m.Hills(-1)), m)
		log.Printf("\"%s\",\"%s\",%d,%d,%d,%d\n", mapid, name, m.Rows, m.Cols, m.Players, len(m.Hills(-1)))
	}
}
