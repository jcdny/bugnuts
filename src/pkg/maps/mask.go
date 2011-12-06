package maps

import (
	"math"
	"fmt"
	"log"
	. "bugnuts/torus"
)

type Offsets struct {
	P []Point
	L []Location
}

type Mask struct {
	R      uint8 // Radius
	Stride int   // Cols

	Offsets
	// Locations added or removed for a step in each direction
	Union  Offsets     // Union is the one step out in any direction
	Add    []Offsets   // Points added for step in dir d
	Remove []Offsets   // Points removed for step in dir d
	MM     []*MoveMask // See FreedomKey for how to index into this 
}

// MoveMask is generated for a given number of degrees of freedom
type MoveMask struct {
	R      uint8    // Radius
	NFree  int      // Degrees of freedom
	PStep  int      // Probability denom
	Stride int      // Cols used to generate loc offsets.
	MinPr  []uint16 // Pr numerator, matches Offsets order
	MaxPr  []uint16 // Pr numerator, matches Offsets order
	Offsets
}

const MoveMaskPStep = 60

func maskCircle(r2 int) []Point {
	if r2 < 0 {
		return nil
	}

	d := int(math.Sqrt(float64(r2)))
	v := make([]Point, 0, (r2*22)/7+5)

	// Make the origin the first element so you can easily skip it.
	p := Point{R: 0, C: 0}
	v = append(v, p)

	for r := -d; r <= d; r++ {
		for c := -d; c <= d; c++ {
			if c != 0 || r != 0 {
				if c*c+r*r <= r2 {
					p = Point{R: int(r), C: int(c)}
					v = append(v, p)
				}
			}
		}
	}

	return v
}

// Given a []Point vector, compute the change from stepping north, south, east, west
// Useful for updating visibility, ranking move values.
func maskChange(r2 int, v []Point) (add, remove [][]Point, union []Point) {
	// compute the size of the array we need to hold shifted circle
	d := int(math.Sqrt(float64(r2)))

	//TODO compute d from v rather than r2 so we can use different masks
	off := d + 1    // offset to get origin
	size := 2*d + 3 // one on either side + origin

	union = make([]Point, len(v), len(v)+4*size)
	copy(union, v)

	// Ordinal moves
	for _, s := range Steps {
		m := make([]int, size*size)

		av := []Point{}
		rv := []Point{}

		for _, p := range v {
			m[(p.C+off)+(p.R+off)*size]++
			m[(p.C+s.C+off)+(p.R+s.R+off)*size]--
		}

		for r := 0; r < size; r++ {
			for c := 0; c < size; c++ {
				switch {
				case m[c+r*size] > 0:
					rv = append(rv, Point{R: r - off, C: c - off})
				case m[c+r*size] < 0:
					av = append(av, Point{R: r - off, C: c - off})
				}
			}
		}

		add = append(add, av)
		remove = append(remove, rv)
		union = union[0 : len(union)+len(av)]
		copy(union[len(union)-len(av):len(union)], av)

	}

	return
}

// Generate a mask for a circle, including the added/removed list for
// steps in any directions plus a union of the 1step move There is
// also the move mask which includes probabilities of a cell falling
// within the mask given available moves.
func MakeMask(r2, rows, cols int) *Mask {
	p := maskCircle(r2)
	add, rem, union := maskChange(r2, p)
	addo := make([]Offsets, 0, len(add))
	for _, pv := range add {
		addo = append(addo, PointsToOffsets(pv, cols))
	}
	remo := make([]Offsets, 0, len(rem))
	for _, pv := range rem {
		remo = append(remo, PointsToOffsets(pv, cols))
	}
	uniono := PointsToOffsets(union, cols)

	r := uint8(math.Sqrt(float64(r2)))
	m := &Mask{
		R:      r,
		Stride: cols,
		Add:    addo,
		Remove: remo,
		Union:  uniono,

		MM: MakeMoveMask(r2, cols),
	}

	m.Offsets = PointsToOffsets(p, cols)

	return m
}

func MakeMoveMask(r2 int, cols int) []*MoveMask {
	if r2 < 0 {
		log.Panic("Radius must be > 0")
	}
	rad := int(math.Sqrt(float64(r2)))
	stride := 2*rad + 3
	size := stride * stride
	center := Location(size / 2)

	// generate a mask for combat radius
	off := PointsToOffsets(maskCircle(r2), stride)

	mm := make([]*MoveMask, 16)
	// loop over possible states
	for i := 0; i < 16; i++ {
		pr := make([]int, size)
		nfree := 0

		// degrees of freedom
		for bit := uint(0); bit < 4; bit++ {
			if i&(1<<bit) > 0 {
				nfree++
			}
		}

		// pstep prob
		pstep := MoveMaskPStep / (nfree + 1)

		// now generate the actual probabilities
		for bit := uint(0); bit < 5; bit++ {
			if (i+16)&(1<<bit) > 0 {
				offset := Location(DirectionOffset[bit].R*stride + DirectionOffset[bit].C)
				for _, l := range off.L {
					loc := center + offset + l
					pr[loc] += pstep
				}
			}
		}

		// Given maxpr Generate the location offsets and point offsets for the masks
		mpt := make([]Point, 0, len(pr))
		minpr := make([]uint16, 0, len(pr))
		maxpr := make([]uint16, 0, len(pr))

		off := rad + 1
		for r := 0; r < stride; r++ {
			for c := 0; c < stride; c++ {
				p := pr[r*stride+c]
				if p > 0 {
					mpt = append(mpt, Point{R: r - off, C: c - off})
					minpr = append(minpr, uint16(60-p))
					maxpr = append(maxpr, uint16(p))
				}
			}
		}

		mask := MoveMask{
			R:      uint8(rad),
			NFree:  nfree,
			PStep:  pstep,
			Stride: cols, // This is for the Locations, not lstride we use internally here
			MinPr:  minpr,
			MaxPr:  maxpr,
		}
		mask.Offsets = PointsToOffsets(mpt, cols)
		mm[i] = &mask
	}

	return mm
}

// ApplyMask applies a precomputed mask centered on location loc
func (m *Map) ApplyMask(loc Location, mask *Mask, x func(loc Location)) {
	if m.BorderDist[loc] > mask.R {
		for _, loff := range mask.L {
			x(loc + loff)
		}
	} else {
		ap := m.ToPoint(loc)
		for _, op := range mask.P {
			l := m.ToLocation(m.PointAdd(ap, op))
			x(l)
		}
	}
}
// ApplyMaskBreak applies a precomputed mask centered on location loc via function x, if x returns false the apply is aborted.
func (m *Map) ApplyMaskBreak(loc Location, mask *Mask, x func(loc Location) bool) {
	if m.BorderDist[loc] > mask.R {
		for _, loff := range mask.L {
			if !x(loc + loff) {
				return
			}
		}
	} else {
		ap := m.ToPoint(loc)
		for _, op := range mask.P {
			l := m.ToLocation(m.PointAdd(ap, op))
			if !x(l) {
				return
			}
		}
	}
}

func (mm *MoveMask) String() string {
	s := fmt.Sprintf("free %d pstep: %d stride %d\n*** minpr:", mm.NFree, mm.PStep, mm.Stride)
	stride := int(2*mm.R + 3)

	minpr := make([]uint16, stride*stride)
	for i := range minpr {
		minpr[i] = MoveMaskPStep
	}

	maxpr := make([]uint16, stride*stride)
	off := stride * stride / 2
	for i, p := range mm.P {
		minpr[p.R*stride+p.C+off] = mm.MinPr[i]
		maxpr[p.R*stride+p.C+off] = mm.MaxPr[i]
	}

	for r := 0; r < stride; r++ {
		s += "\n"
		for c := 0; c < stride; c++ {
			s += fmt.Sprintf("%3d", minpr[r*stride+c])
		}
	}
	s += "\n*** maxpr"
	for r := 0; r < stride; r++ {
		s += "\n"
		for c := 0; c < stride; c++ {
			s += fmt.Sprintf("%3d", maxpr[r*stride+c])
		}
	}

	return s
}

func (m *Map) FreedomKey(loc Location) int {
	key := 0
	for i, l := range m.LocStep[loc] {
		// skip no move
		if l != loc && StepableItem[m.Grid[l]] {
			key += 1 << uint(i)
		}
	}

	return key
}

// Compute degrees of freedom taking into account threat, returned value can be used to index into mask.MM
func (m *Map) FreedomKeyThreat(loc Location, t []int8, nsup [4]int8) int {
	key := 0
	for i, l := range m.LocStep[loc] {
		if l != loc && StepableItem[m.Grid[l]] && (len(t) == 0 || t[l] < nsup[i]) {
			key += 1 << uint(i)
		}
	}

	return key
}
