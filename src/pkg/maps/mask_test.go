package maps

import (
	"testing"
	"sort"
	"json"
	"log"
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
	m := MakeMask(5, 10, 100)
	o, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Err %v\nMask: %s", err, o)
	}
}

func TestMaskChange(t *testing.T) {
	for _, r2 := range []int{0, 1, 5, 55, 100} {
		v := maskCircle(r2)
		add, remove, union := maskChange(r2, v)
		if false {
			// Output the combat circle and variants.
			if r2 == 5 {
				log.Printf("%d v %d len=%d\n add: %v\nremove: %v\nunion: %v\n", r2, len(v), len(union), add, remove, union)
			}
		}

		if len(add) != 4 || len(remove) != 4 {
			t.Errorf("maskChange sizes are wrong add: %v remove: %v union: %v", add, remove, union)
			break
		}

		for i := range add {
			if len(add[i]) != len(remove[i]) {
				t.Errorf("maskChange sizes are wrong add: %v remove: %v", add, remove)
				break
			}
			if i > 1 && len(add[i-1]) != len(add[i]) {
				t.Errorf("maskChange sizes are wrong add: %v remove: %v", add, remove)
				break
			}
		}
	}

}

func TestMoveMask(t *testing.T) {
	mm := MakeMoveMask(5, 100)
	if len(mm) != 16 {
		t.Errorf("MoveMask size wrong")
	}
	if false {
		// dump masks
		for i := 0; i < 16; i++ {
			log.Printf("%v", mm[i])
		}
	}
}
