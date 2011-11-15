package maps

type Torus struct {
	Rows int
	Cols int
}

func (t *Torus) Size() int {
	return t.Rows * t.Cols
}

func (t *Torus) ToLocation(p Point) Location {
	p = t.Donut(p)
	return Location(p.r*t.Cols + p.c)
}

func (t *Torus) ToPoint(l Location) (p Point) {
	p = Point{r: int(l) / t.Cols, c: int(l) % t.Cols}

	return
}

func (t *Torus) Donut(p Point) Point {
	if p.r < 0 {
		p.r += t.Rows
	}
	if p.r >= t.Rows {
		p.r -= t.Rows
	}
	if p.c < 0 {
		p.c += t.Cols
	}
	if p.c >= t.Cols {
		p.c -= t.Cols
	}

	return p
}

func (t *Torus) PointEqual(p1, p2 Point) bool {
	// todo donuts
	return p1.c == p2.c && p1.r == p2.r
}

func (t *Torus) PointAdd(p1, p2 Point) Point {
	return t.Donut(Point{r: p1.r + p2.r, c: p1.c + p2.c})
}

//Take a slice of Point and return a slice of Location
//Used for offsets so it does not donut things.
func (t *Torus) ToLocations(pv []Point) []Location {
	lv := make([]Location, len(pv), len(pv)) // maybe use cap(pv)
	for i, p := range pv {
		lv[i] = Location(p.r*t.Cols + p.c)
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
}

func (t *Torus) EDist2(l1, l2 Location) int {
}

func (t *Torus) Reduce(off Point) Point {
}

func (t *Torus) LocSort(origin Location, locs []Location) {
}
