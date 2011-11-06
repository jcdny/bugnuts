package main

import (
	"testing"
	"os"
	"image"
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

	if err != os.EOF {
		t.Errorf("Invalid load of map error os.Error == %v", err)
	}
	if m == nil {
		t.Errorf("Invalid load of map m == nil")
	}

	m.WriteDebugImage("_maptest", 0, func(c, r int) image.NRGBAColor { return m.At(r, c) })
}
