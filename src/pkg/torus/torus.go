// Package torus implements geometry on a torus.
// It supports conversion between Locations (indexes into linear
// arrays) and Points ({row, col} coordinates)
package torus

import (
	"log"
)

// A Torus is defined by its dimensions.
// By convention {0,0} is the top right.
type Torus struct {
	Rows int
	Cols int
}

// A Location is a linear index to the torus position.
type Location int

type LocationSlice []Location

func (p LocationSlice) Len() int           { return len(p) }
func (p LocationSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p LocationSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Size returns the number of cells in a given Torus.
func (t *Torus) Size() int {
	return t.Rows * t.Cols
}

// ToLocation converts a Point into a Location
func (t *Torus) ToLocation(p Point) Location {
	p = t.Donut(p)
	return Location(p.R*t.Cols + p.C)
}

// ToPoint converts a Location into a Point
func (t *Torus) ToPoint(l Location) (p Point) {
	p = Point{R: int(l) / t.Cols, C: int(l) % t.Cols}

	return
}

// Donut converts a Point which may be outside the canonical torus to 
// a Point with 0 <= R < Rows and 0 <= C < Cols.
func (t *Torus) Donut(p Point) Point {
	if p.R < 0 {
		p.R += t.Rows
	}
	if p.R >= t.Rows {
		p.R -= t.Rows
	}
	if p.C < 0 {
		p.C += t.Cols
	}
	if p.C >= t.Cols {
		p.C -= t.Cols
	}

	return p
}

// PointEqual compares two Points for equality and assumes they are already in standard coordinates.
func (t *Torus) PointEqual(p1, p2 Point) bool {
	// todo donuts
	return p1.C == p2.C && p1.R == p2.R
}

// PointAdd adds two Points and converts them to standard coordinates.
func (t *Torus) PointAdd(p1, p2 Point) Point {
	return t.Donut(Point{R: p1.R + p2.R, C: p1.C + p2.C})
}

// ToLocations takes a slice of Point and return a slice of Location
// It is used for offsets so it does not convert to standard coordinates.
func (t *Torus) ToLocations(pv []Point) []Location {
	lv := make([]Location, len(pv), len(pv)) // maybe use cap(pv)
	for i, p := range pv {
		lv[i] = Location(p.R*t.Cols + p.C)
	}

	return lv
}

// Take a slice of Location and return slice of Point, Does not
// convert to standard coordinates since the locations could be
// Offsets.
func (t *Torus) ToPoints(lv []Location) []Point {
	pv := make([]Point, len(lv))
	for i, l := range lv {
		pv[i] = t.ToPoint(l)
	}

	return pv
}

// MDist computes manhattan distance between two Locations.
func (t *Torus) MDist(l1, l2 Location) int {
	// handle odd rows/cols...
	log.Panicf("unimplemented")
	return 0
}

// Edist2 computes the euclidian distance between two Locations.
func (t *Torus) EDist2(l1, l2 Location) int {
	log.Panicf("unimplemented")
	return 0
}
