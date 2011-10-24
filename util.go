package main

import (
	"math"
)

// Utility to implement manhattan distance sorting on a slice of Point Offsets
// TODO think about how to do this in context of torus.  Assumes Offset < side/2
type OffsetSlice []Point

func (p OffsetSlice) Len() int           { return len(p) }
func (p OffsetSlice) Less(i, j int) bool { return Abs(p[i].r)+Abs(p[i].c) < Abs(p[j].r)+Abs(p[j].c) }
func (p OffsetSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func Max(x []int) int {
	xm := math.MinInt32
	for _, y := range x {
		if y > xm {
			xm = y
		}
	}

	return xm
}

func Min(x []int) int {
	xm := math.MaxInt32
	for _, y := range x {
		if y < xm {
			xm = y
		}
	}

	return xm
}

// Precompute circle points for lookup for a given r2 and number of map columns.
func GenCircleTable(r2 int) []Point {
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
func moveChangeCache(r2 int, v []Point) (add [][]Point, remove [][]Point) {
	// compute the size of the array we need to hold shifted circle
	d := int(math.Sqrt(float64(r2)))
	//TODO compute d from v rather than r2 so we can use different masks

	off := d + 1    // offset to get origin
	size := 2*d + 3 // one on either side + origin

	// Ordinal moves
	// TODO pass in
	sv := []Point{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
	for _, s := range sv {
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
	}

	return
}
