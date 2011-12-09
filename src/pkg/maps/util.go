package maps

import (
	"rand"
	. "bugnuts/torus"
	. "bugnuts/util"
)

var PermStepD5 [4][6][5]Direction

func init() {
	permdstep := map[string][]string{
		"ns": []string{"ew-", "e-w", "we-", "w-e", "-ew", "-we"},
		"ew": []string{"ns-", "n-s", "sn-", "s-n", "-ns", "-sn"},
	}
	for k, v := range permdstep {
		da := ByteToDirection[k[0]]
		db := ByteToDirection[k[1]]
		for n, s := range v {
			PermStepD5[da][n][0] = da
			PermStepD5[da][n][4] = db
			PermStepD5[db][n][0] = db
			PermStepD5[db][n][4] = da
			for j, d := range s {
				PermStepD5[da][n][j+1] = ByteToDirection[d]
				PermStepD5[db][n][j+1] = ByteToDirection[d]
			}
		}
	}
}

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

// PointsToOffsets convert a set of points to location offsets
func PointsToOffsets(pv []Point, cols int) Offsets {
	out := Offsets{
		P: make([]Point, len(pv), len(pv)),
		L: make([]Location, len(pv), len(pv)),
	}
	copy(out.P, pv)
	R := 0
	for i, p := range pv {
		R = MaxV(Abs(p.R), Abs(p.C), R)
		out.L[i] = Location(p.R*cols + p.C)
	}
	out.R = uint8(R)
	out.cacheL = make(map[Location][]Location, 4*(R*cols-R*R))
	return out
}

//  LocationsToOffsets takes a list of location offsets and returns a Points vector
func LocationsToOffsets(locs []Location, cols int) Offsets {
	out := Offsets{
		P: make([]Point, len(locs), len(locs)),
		L: make([]Location, len(locs), len(locs)),
	}

	copy(out.L, locs)
	for i, l := range locs {
		out.P[i] = Point{R: int(l) / cols, C: int(l) % cols}
	}

	return out
}

// Draw from the arrays with rand.Intn(24) for a random permutation of directions.
func Permute4(r *rand.Rand) *[4]Direction {
	return &Perm4[r.Intn(24)]
}

// Returns a permute4 + guard used in eg NPathIn
func Permute4G(r *rand.Rand) *[5]Direction {
	return &Perm4G[r.Intn(24)]
}

// Draw from the arrays with rand.Intn(120) for a random permutation of directions including NoMovement
func Permute5(r *rand.Rand) *[5]Direction {
	return &Perm5[r.Intn(120)]
}

// Return a set of directions with d first, opposite last an the other 3 permuted.
func Permute5D(d Direction, r *rand.Rand) *[5]Direction {
	return &PermStepD5[d][r.Intn(6)]
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
