package main

import (
	"testing"
	"math"
	"sort"
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

func TestGenCircleTable(t *testing.T) {
	exp7100 := []int{0, -201, -200, -199, -102, -101, -100, -99, -98, -2, -1, 1, 2, 98, 99, 100, 101, 102, 199, 200, 201}

	v := GenCircleTable(7)
	if len(v) != len(exp7100) {
		t.Errorf("GenCircleTable(7) expected %v got %v", exp7100, v)
	}

	sort.Sort(OffsetSlice(v))
	//log.Printf("%T %v", v, v)

	v = GenCircleTable(1)
	if len(v) != 5 {
		t.Errorf("GenCircleTable(1) expected len=5 got %v", v)
	}

	v = GenCircleTable(0)
	if len(v) != 1 {
		t.Errorf("GenCircleTable(0) expected len=1 got %v", v)
	}

	/* // test range of i
	 for i := 0; i <= 10; i += 1 {
		v = GenCircleTable(i)
		log.Printf("r2=%5d %5d %v", i, len(v), v)
	 }
	*/
}

func TestMoveChangeCache(t *testing.T) {
	for _, r2 := range []int{0, 1, 7, 55, 100} {
		v := GenCircleTable(r2)
		add, remove := moveChangeCache(r2, v)

		if len(add) != 4 || len(remove) != 4 {
			t.Errorf("moveChangeCache sizes are wrong add: %v remove: %v", add, remove)
			break
		}

		for i, _ := range add {
			if len(add[i]) != len(remove[i]) {
				t.Errorf("moveChangeCache sizes are wrong add: %v remove: %v", add, remove)
				break
			}
			if i > 1 && len(add[i-1]) != len(add[i]) {
				t.Errorf("moveChangeCache sizes are wrong add: %v remove: %v", add, remove)
				break
			}
		}
	}

}
