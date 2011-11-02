package main

import (
	"math"
)

type Mask struct {
	R      uint8 // Radius
	Stride int // cols

	P   []Point
	Loc []Location

	// Locations added or removed for a step in each direction
	Add    [][]Point
	Remove [][]Point
	// Union of points in all directions
	Union    []Point
	UnionLoc []Location
}

func maskCircle(r2 int) []Point {
	if r2 < 0 {
		return nil
	}

	d := int(math.Sqrt(float64(r2)))
	v := make([]Point, 0, (r2*22)/7+5)

	// Make the origin the first element so you can easily skip it.
	p := Point{r: 0, c: 0}
	v = append(v, p)

	for r := -d; r <= d; r++ {
		for c := -d; c <= d; c++ {
			if c != 0 || r != 0 {
				if c*c+r*r <= r2 {
					p = Point{r: int(r), c: int(c)}
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
			m[(p.c+off)+(p.r+off)*size]++
			m[(p.c+s.c+off)+(p.r+s.r+off)*size]--
		}

		for c := 0; c < size; c++ {
			for r := 0; r < size; r++ {
				switch {
				case m[c+r*size] > 0:
					rv = append(rv, Point{r: r - off, c: c - off})
				case m[c+r*size] < 0:
					av = append(av, Point{r: r - off, c: c - off})
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

func makeMask(r2, rows, cols int) *Mask {
	p := maskCircle(r2)
	add, rem, union := maskChange(r2, p)

	r := uint8(math.Sqrt(float64(r2)))
	m := &Mask{
		R:        r,
		P:        p,
		Add:      add,
		Remove:   rem,
		Union:    union,
		Loc:      ToOffsets(p, cols),
		UnionLoc: ToOffsets(union, cols),
	}

	return m
}

