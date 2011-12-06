package maps

import (
	"testing"
)

func TestPoints(t *testing.T) {
	v := maskCircle(10)

	vl := PointsToOffsets(v, 13)
	vp := LocationsToOffsets(vl.L, 13)
	vl2 := PointsToOffsets(vp.P, 13)

	if len(vl.L) != len(vl2.L) || len(v) != len(vp.P) {
		t.Error("Point length mismatch")
	} else {
		for i, l := range vl.L {
			if l != vl2.L[i] {
				t.Error("Point roundtrip failed %v, %v, %v, %v", v, vl, vp, vl2)
				break
			}
		}
	}
}
