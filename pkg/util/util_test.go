package util

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
