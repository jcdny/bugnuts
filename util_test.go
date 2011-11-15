package main

import (
	"testing"
	"math"
)

type Lists struct {
	list       []int
	emin, emax int
}

func TestMinMax(t *testing.T) {
	L := []Lists{
		{[]int{1}, 1, 1},
		{[]int{1, 1, 1, 1, 1, 1}, 1, 1},
		{[]int{1, 0, -1}, -1, 1},
		{[]int{}, math.MaxInt32, math.MinInt32},
	}

	for _, l := range L {
		if amin := Min(l.list); amin != l.emin {
			t.Errorf("For min(%v) got %d expected %d", l.list, amin, l.emin)
		}
		if amax := Max(l.list); amax != l.emax {
			t.Errorf("For max(%v) got %d expected %d", l.list, amax, l.emax)
		}
	}
}

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
