package maps

import (
	"log"
)

type Location int

type Torus struct {
	Rows int
	Cols int
}

func (t *Torus) Size() int {
	return t.Rows * t.Cols
}

func (t *Torus) ToLocation(p Point) Location {
	p = t.Donut(p)
	return Location(p.R*t.Cols + p.C)
}

func (t *Torus) ToPoint(l Location) (p Point) {
	p = Point{R: int(l) / t.Cols, C: int(l) % t.Cols}

	return
}

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

func (t *Torus) PointEqual(p1, p2 Point) bool {
	// todo donuts
	return p1.C == p2.C && p1.R == p2.R
}

func (t *Torus) PointAdd(p1, p2 Point) Point {
	return t.Donut(Point{R: p1.R + p2.R, C: p1.C + p2.C})
}

//Take a slice of Point and return a slice of Location
//Used for offsets so it does not donut things.
func (t *Torus) ToLocations(pv []Point) []Location {
	lv := make([]Location, len(pv), len(pv)) // maybe use cap(pv)
	for i, p := range pv {
		lv[i] = Location(p.R*t.Cols + p.C)
	}

	return lv
}

// Take a slice of locations and return points, Does not Donut 
// since the locations could be Offsets
func (t *Torus) ToPoints(lv []Location) []Point {
	pv := make([]Point, len(lv))
	for i, l := range lv {
		pv[i] = t.ToPoint(l)
	}

	return pv
}

func (t *Torus) MDist(l1, l2 Location) int {
	// handle odd rows/cols...
	log.Panicf("unimplemented")
	return 0
}

func (t *Torus) EDist2(l1, l2 Location) int {
	log.Panicf("unimplemented")
	return 0
}

// Take an offset and reduce it to the minimum magnitude.
func (t *Torus) Reduce(off Point) Point {
	log.Panicf("unimplemented")
	return Point{0, 0}
}

func (t *Torus) LocSort(origin Location, locs []Location) {
	log.Panicf("unimplemented")
}
