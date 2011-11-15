package main

import (
	"math"
)

// Utility to implement manhattan distance sorting on a slice of Point Offsets
// TODO think about how to do this in context of torus.  Assumes Offset < side/2
type OffsetSlice []Point

func (p OffsetSlice) Len() int { return len(p) }
// Metric is Manhattan distance from origin.
func (p OffsetSlice) Less(i, j int) bool { return Abs(p[i].r)+Abs(p[i].c) < Abs(p[j].r)+Abs(p[j].c) }
func (p OffsetSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type IntSlice []int

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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

func MinInt8(x []int8) int8 {
	xm := int8(math.MaxInt8)
	for _, y := range x {
		if y < xm {
			xm = y
		}
	}

	return int8(xm)
}

func MinV(v1 int, vn ...int) (m int) {
	m = v1
	for _, vi := range vn {
		if vi < m {
			m = vi
		}
	}
	return
}

func MaxV(v1 int, vn ...int) (m int) {
	m = v1
	for _, vi := range vn {
		if vi > m {
			m = vi
		}
	}
	return
}

// Convert a set of points to location offsets
func ToOffsets(pv []Point, cols int) []Location {
	out := make([]Location, len(pv), len(pv))

	for i, p := range pv {
		out[i] = Location(p.r*cols + p.c)
	}

	return out
}

// Convert a list of location offsets back to signed Points.
func ToOffsetPoints(loc []Location, cols int) (out []Point) {
	out = make([]Point, len(loc))
	for i, l := range loc {
		out[i] = Point{r: int(l) / cols, c: int(l) % cols}
	}

	return out
}
