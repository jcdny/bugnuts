package maps

import (
	"testing"
	"os"
)

func TestMapNew(t *testing.T) {
	r, c := 5, 7
	m := NewMap(r, c, 1)

	if m.Size() != r*c {
		t.Errorf("Map size error")
	}

	for i, d := range m.borderDist {
		if m.borderDist[m.Size()-i-1] != d {
			t.Errorf("Border Distance Error, Not Symmetric: %v", m.borderDist)
			break
		}
	}

	if m.borderDist[0] != 1 ||
		m.borderDist[m.Size()-1] != 1 ||
		m.borderDist[c+2] != 2 {
		t.Errorf("Border Distance Error bounds wrong: %v", m.borderDist)
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

	//m.WriteDebugImage("_maptest", 0, func(c, r int) image.NRGBAColor { return m.At(r, c) })
}
