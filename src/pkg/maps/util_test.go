package maps

import (
	"testing"
)

func TestPoints(t *testing.T) {
	v := maskCircle(10)

	vl := ToOffsets(v, 13)
	vp := ToOffsetPoints(vl, 13)
	vl2 := ToOffsets(vp, 13)

	if len(vl) != len(vl2) || len(v) != len(vp) {
		t.Error("Point length mismatch")
	} else {
		for i, l := range vl {
			if l != vl2[i] {
				t.Error("Point roundtrip failed %v, %v, %v, %v", v, vl, vp, vl2)
				break
			}
		}
	}
}
