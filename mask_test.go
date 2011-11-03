package main

import (
	"testing"
	"sort"
	"json"
)

func TestMaskCircle(t *testing.T) {
	v := maskCircle(77)
	if len(v) != 241 {
		t.Errorf("maskCircle(77) expected 241 got %d: %v", len(v), v)
	}

	sort.Sort(OffsetSlice(v))
	//log.Printf("%T %v", v, v)

	v = maskCircle(1)
	if len(v) != 5 {
		t.Errorf("maskCircle(1) expected len=5 got %v", v)
	}

	v = maskCircle(0)
	if len(v) != 1 {
		t.Errorf("maskCircle(0) expected len=1 got %v", v)
	}

	for i := 1000; i > 0; i-- {
		v = maskCircle(i)
		// masks should have 4 way symmetry
		if (len(v)-1)%4 != 0 {
			t.Errorf("maskCircle(i) expected len to be 1+4*n got %d+4*%d", (len(v)-1)%4, len(v)/4)
			break
		}
	}
}

func TestMakeMask(t *testing.T) {
	m := makeMask(5, 10, 100)
	o, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Err %v\nMask: %s", err, o)
	}
}
