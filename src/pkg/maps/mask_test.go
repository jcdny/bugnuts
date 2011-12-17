package maps

import (
	"testing"
	"sort"
	"json"
	"log"
	. "bugnuts/torus"
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
	if true {
		// dump masks
		for i := 0; i < 16; i++ {
			log.Printf("%v", mm[i])
			log.Printf("Add: %v", mm[i].Add)
		}
	}
}

func TestMoveMask2(t *testing.T) {
	mm := MakeMoveMask2(5, 100)
	if len(mm) != 4096 {
		t.Error("MoveMask size wrong:", len(mm))
	}
	if true {
		// dump masks
		for i := 0; i < 2; i++ {
			log.Printf("%v", mm[i])
		}
	}
}

const (
	bCols = 80
	bRows = 80
	bMask = 77
)

func TestApplyCacheCreate(t *testing.T) {
	bRows := 23
	bCols := 17
	m := NewMap(bRows, bCols, 1)

	j := 0
	mm := MakeMask(bMask, bRows, bCols)
	for loc := range m.Grid {
		m.ApplyOffsets(Location(loc), &mm.Offsets, func(l Location) { j++ })
	}
	for loc := range m.Grid {
		m.ApplyOffsets(Location(loc), &mm.Offsets, func(l Location) { j++ })
	}

	if j != 2*bRows*bCols*len(mm.Offsets.P) {
		t.Error("Invalid ApplyOffsets expected j=", 2*bRows*bCols*len(mm.Offsets.P), " got j=", j)
	}
}

func TestCacheAll(t *testing.T) {
	m := NewMap(bRows, bCols, 1)
	mm := MakeMask(bMask, bRows, bCols)
	m.OffsetsCachePopulateAll(mm, 0)

	r := int(mm.Offsets.R)
	e := 2*r*(bRows+bCols) - 4*r*r

	if len(mm.Offsets.cacheL) != e {
		t.Error("Cache size error expect ", e, " got ", len(mm.Offsets.cacheL))
	}
}

func BenchmarkApplyCached(b *testing.B) {
	m := NewMap(bRows, bCols, 1)
	mm := MakeMask(bMask, bRows, bCols)
	for loc := range m.Grid {
		m.ApplyOffsets(Location(loc), &mm.Offsets, func(l Location) {})
	}

	for i := 0; i < b.N; i++ {
		j := 0
		for loc := range m.Grid {
			m.ApplyOffsets(Location(loc), &mm.Offsets, func(l Location) { j++ })
		}
	}
}

func BenchmarkApplyNone(b *testing.B) {
	m := NewMap(bRows, bCols, 1)
	mm := MakeMask(bMask, bRows, bCols)
	o := &mm.Offsets

	for i := 0; i < b.N; i++ {
		j := 0
		for loc := range m.Grid {
			if m.BorderDist[loc] <= o.R {
				ap := m.ToPoint(Location(loc))
				for j, op := range o.P {
					j += int(m.ToLocation(m.PointAdd(ap, op)))
				}
			} else {
				for _, lo := range o.L {
					j += loc + int(lo)
				}
			}
		}
	}
}

func BenchmarkApplyNoCache(b *testing.B) {
	m := NewMap(bRows, bCols, 1)
	mm := MakeMask(bMask, bRows, bCols)
	mm.Offsets.nocache = true
	for i := 0; i < b.N; i++ {
		j := 0
		for loc := range m.Grid {
			m.ApplyOffsets(Location(loc), &mm.Offsets, func(l Location) { j++ })
		}
	}
}

func BenchmarkApplyCacheCreateA(b *testing.B) {
	m := NewMap(bRows, bCols, 1)
	for i := 0; i < b.N; i++ {
		mm := MakeMask(bMask, bRows, bCols)
		o := &mm.Offsets
		for loc := Location(0); int(loc) < len(m.Grid); loc++ {
			m.ApplyOffsets(loc, o, func(l Location) {})
		}
	}
}

func BenchmarkApplyCacheCreateB(b *testing.B) {
	m := NewMap(bRows, bCols, 1)
	for i := 0; i < b.N; i++ {
		mm := MakeMask(bMask, bRows, bCols)
		o := &mm.Offsets
		size := 2*int(o.R)*(m.Rows+m.Cols) - 4*int(o.R)*int(o.R)
		o.cacheL = make(map[Location][]Location, size)
		for loc := Location(0); int(loc) < len(m.Grid); loc++ {
			if m.BorderDist[loc] <= o.R {
				cl := make([]Location, len(o.P))
				ap := m.ToPoint(loc)
				for j, op := range o.P {
					cl[j] = m.ToLocation(m.PointAdd(ap, op))
				}
				mm.Offsets.cacheL[loc] = cl
			}
		}
	}
}

func BenchmarkCacheAll(b *testing.B) {
	m := NewMap(bRows, bCols, 1)

	for i := 0; i < b.N; i++ {
		mm := MakeMask(bMask, bRows, bCols)
		m.OffsetsCachePopulateAll(mm, 0)
	}
}
