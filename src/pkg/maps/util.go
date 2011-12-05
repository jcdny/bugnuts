package maps

import (
	"rand"
	. "bugnuts/torus"
	. "bugnuts/util"
)

// Utility to implement manhattan distance sorting on a slice of Point Offsets
// TODO think about how to do this in context of torus.  Assumes Offset < side/2
type OffsetSlice []Point

func (p OffsetSlice) Len() int { return len(p) }
// Sort metric is Manhattan distance from origin.
func (p OffsetSlice) Less(i, j int) bool { return Abs(p[i].R)+Abs(p[i].C) < Abs(p[j].R)+Abs(p[j].C) }
func (p OffsetSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type IntSlice []int

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Convert a set of points to location offsets
func ToOffsets(pv []Point, cols int) []Location {
	out := make([]Location, len(pv), len(pv))

	for i, p := range pv {
		out[i] = Location(p.R*cols + p.C)
	}

	return out
}

// Convert a list of location offsets back to signed Points.
func ToOffsetPoints(loc []Location, cols int) (out []Point) {
	out = make([]Point, len(loc))
	for i, l := range loc {
		out[i] = Point{R: int(l) / cols, C: int(l) % cols}
	}

	return out
}

// Draw from the arrays with rand.Intn(24) for a random permutation of directions.
func Permute4() *[4]Direction {
	return &Perm4[rand.Intn(24)]
}

// Returns a permute4 + guard used in eg NPathIn
func Permute4G() *[5]Direction {
	return &Perm4G[rand.Intn(24)]
}

// Draw from the arrays with rand.Intn(120) for a random permutation of directions including NoMovement
func Permute5() *[5]Direction {
	return &Perm5[rand.Intn(120)]
}

// Perm4G is the directions permuted with the NoMovement guard.
// The ordering here is relevant since NPathInString uses the perm flag to index into this 
// in order to generate bounding paths.
var Perm4G = [...][5]Direction{
	{0, 1, 2, 3, 4},
	{3, 2, 1, 0, 4},
	{0, 1, 3, 2, 4},
	{0, 2, 1, 3, 4},
	{0, 2, 3, 1, 4},
	{0, 3, 1, 2, 4},
	{0, 3, 2, 1, 4},
	{1, 0, 2, 3, 4},
	{1, 0, 3, 2, 4},
	{1, 2, 0, 3, 4},
	{1, 2, 3, 0, 4},
	{1, 3, 0, 2, 4},
	{1, 3, 2, 0, 4},
	{2, 0, 1, 3, 4},
	{2, 0, 3, 1, 4},
	{2, 1, 0, 3, 4},
	{2, 1, 3, 0, 4},
	{2, 3, 0, 1, 4},
	{2, 3, 1, 0, 4},
	{3, 0, 1, 2, 4},
	{3, 0, 2, 1, 4},
	{3, 1, 0, 2, 4},
	{3, 1, 2, 0, 4},
	{3, 2, 0, 1, 4},
}
var Perm4 = [...][4]Direction{
	{0, 1, 2, 3},
	{0, 1, 3, 2},
	{0, 2, 1, 3},
	{0, 2, 3, 1},
	{0, 3, 1, 2},
	{0, 3, 2, 1},
	{1, 0, 2, 3},
	{1, 0, 3, 2},
	{1, 2, 0, 3},
	{1, 2, 3, 0},
	{1, 3, 0, 2},
	{1, 3, 2, 0},
	{2, 0, 1, 3},
	{2, 0, 3, 1},
	{2, 1, 0, 3},
	{2, 1, 3, 0},
	{2, 3, 0, 1},
	{2, 3, 1, 0},
	{3, 0, 1, 2},
	{3, 0, 2, 1},
	{3, 1, 0, 2},
	{3, 1, 2, 0},
	{3, 2, 0, 1},
	{3, 2, 1, 0},
}

var Perm5 = [...][5]Direction{
	{0, 1, 2, 3, 4},
	{0, 1, 2, 4, 3},
	{0, 1, 3, 2, 4},
	{0, 1, 3, 4, 2},
	{0, 1, 4, 2, 3},
	{0, 1, 4, 3, 2},
	{0, 2, 1, 3, 4},
	{0, 2, 1, 4, 3},
	{0, 2, 3, 1, 4},
	{0, 2, 3, 4, 1},
	{0, 2, 4, 1, 3},
	{0, 2, 4, 3, 1},
	{0, 3, 1, 2, 4},
	{0, 3, 1, 4, 2},
	{0, 3, 2, 1, 4},
	{0, 3, 2, 4, 1},
	{0, 3, 4, 1, 2},
	{0, 3, 4, 2, 1},
	{0, 4, 1, 2, 3},
	{0, 4, 1, 3, 2},
	{0, 4, 2, 1, 3},
	{0, 4, 2, 3, 1},
	{0, 4, 3, 1, 2},
	{0, 4, 3, 2, 1},
	{1, 0, 2, 3, 4},
	{1, 0, 2, 4, 3},
	{1, 0, 3, 2, 4},
	{1, 0, 3, 4, 2},
	{1, 0, 4, 2, 3},
	{1, 0, 4, 3, 2},
	{1, 2, 0, 3, 4},
	{1, 2, 0, 4, 3},
	{1, 2, 3, 0, 4},
	{1, 2, 3, 4, 0},
	{1, 2, 4, 0, 3},
	{1, 2, 4, 3, 0},
	{1, 3, 0, 2, 4},
	{1, 3, 0, 4, 2},
	{1, 3, 2, 0, 4},
	{1, 3, 2, 4, 0},
	{1, 3, 4, 0, 2},
	{1, 3, 4, 2, 0},
	{1, 4, 0, 2, 3},
	{1, 4, 0, 3, 2},
	{1, 4, 2, 0, 3},
	{1, 4, 2, 3, 0},
	{1, 4, 3, 0, 2},
	{1, 4, 3, 2, 0},
	{2, 0, 1, 3, 4},
	{2, 0, 1, 4, 3},
	{2, 0, 3, 1, 4},
	{2, 0, 3, 4, 1},
	{2, 0, 4, 1, 3},
	{2, 0, 4, 3, 1},
	{2, 1, 0, 3, 4},
	{2, 1, 0, 4, 3},
	{2, 1, 3, 0, 4},
	{2, 1, 3, 4, 0},
	{2, 1, 4, 0, 3},
	{2, 1, 4, 3, 0},
	{2, 3, 0, 1, 4},
	{2, 3, 0, 4, 1},
	{2, 3, 1, 0, 4},
	{2, 3, 1, 4, 0},
	{2, 3, 4, 0, 1},
	{2, 3, 4, 1, 0},
	{2, 4, 0, 1, 3},
	{2, 4, 0, 3, 1},
	{2, 4, 1, 0, 3},
	{2, 4, 1, 3, 0},
	{2, 4, 3, 0, 1},
	{2, 4, 3, 1, 0},
	{3, 0, 1, 2, 4},
	{3, 0, 1, 4, 2},
	{3, 0, 2, 1, 4},
	{3, 0, 2, 4, 1},
	{3, 0, 4, 1, 2},
	{3, 0, 4, 2, 1},
	{3, 1, 0, 2, 4},
	{3, 1, 0, 4, 2},
	{3, 1, 2, 0, 4},
	{3, 1, 2, 4, 0},
	{3, 1, 4, 0, 2},
	{3, 1, 4, 2, 0},
	{3, 2, 0, 1, 4},
	{3, 2, 0, 4, 1},
	{3, 2, 1, 0, 4},
	{3, 2, 1, 4, 0},
	{3, 2, 4, 0, 1},
	{3, 2, 4, 1, 0},
	{3, 4, 0, 1, 2},
	{3, 4, 0, 2, 1},
	{3, 4, 1, 0, 2},
	{3, 4, 1, 2, 0},
	{3, 4, 2, 0, 1},
	{3, 4, 2, 1, 0},
	{4, 0, 1, 2, 3},
	{4, 0, 1, 3, 2},
	{4, 0, 2, 1, 3},
	{4, 0, 2, 3, 1},
	{4, 0, 3, 1, 2},
	{4, 0, 3, 2, 1},
	{4, 1, 0, 2, 3},
	{4, 1, 0, 3, 2},
	{4, 1, 2, 0, 3},
	{4, 1, 2, 3, 0},
	{4, 1, 3, 0, 2},
	{4, 1, 3, 2, 0},
	{4, 2, 0, 1, 3},
	{4, 2, 0, 3, 1},
	{4, 2, 1, 0, 3},
	{4, 2, 1, 3, 0},
	{4, 2, 3, 0, 1},
	{4, 2, 3, 1, 0},
	{4, 3, 0, 1, 2},
	{4, 3, 0, 2, 1},
	{4, 3, 1, 0, 2},
	{4, 3, 1, 2, 0},
	{4, 3, 2, 0, 1},
	{4, 3, 2, 1, 0},
}
