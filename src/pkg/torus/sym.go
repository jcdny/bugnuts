package torus

import (
	. "bugnuts/util"
)

// Mirror returns the midpoint between two points.  Axis 0 is R and 1
// is C.  By convention it will return either the axis closer to {0,0}
// if the side dimension is even or the axis which exists between the
// even spread of p1/p2.
func (t *Torus) Mirror(l1, l2 Location, axis int) int {
	p1 := t.ToPoint(l1)
	p2 := t.ToPoint(l2)

	var o, s int
	var odd bool
	if axis == 0 {
		o = p1.R + (p2.R-p1.R+1)/2
		odd = (p1.R-p2.R)%2 == 1
		s = t.Rows
	} else if axis == 1 {
		o = p1.C + (p2.C-p1.C+1)/2
		odd = (p1.C-p2.C)%2 == 1
		s = t.Cols
	}

	if o == 0 {
		return 0
	}

	if !odd {
		if o > (s+1)/2 {
			o -= (s + 1) / 2
		}
	} else {
		o = (o + (s+1)/2) % s
	}

	return o
}

// SymDiff returns the difference between two points.  By convention will return 
// r >= 0.
func (t *Torus) SymDiff(l1, l2 Location) Point {
	p1 := t.ToPoint(l1)
	p2 := t.ToPoint(l2)

	r := p2.R - p1.R
	c := p2.C - p1.C

	if r > t.Rows/2 {
		r -= t.Rows
	}
	if r < -t.Rows/2 {
		r += t.Rows
	}
	if c > t.Cols/2 {
		c -= t.Cols
	}
	if c < -t.Cols/2 {
		c += t.Cols
	}
	if r < 0 {
		r = -r
		c = -c
	}
	return Point{R: r, C: c}
}

// ShiftReduce will take a translation to its minumum length offset.
// I should just do this with math but my head hurts.
func (t *Torus) ShiftReduce(l1, l2 Location, maxcells int) (Point, bool) {
	p1 := t.ToPoint(l1)
	p2 := t.ToPoint(l2)

	r := p2.R - p1.R
	c := p2.C - p1.C

	if r < 0 {
		r += t.Rows
	}
	if c < 0 {
		c += t.Cols
	}

	l := 65535
	rm, cm := r, c
	coff := [3]int{0, 0, -t.Cols}
	roff := [3]int{0, -t.Rows, 0}

	for i := 0; i < maxcells+1; i++ {
		cs := (c + i*c) % t.Cols
		rs := (r + i*r) % t.Rows
		if cs == 0 && rs == 0 && i != 0 {
			return Point{R: rm, C: cm}, true
		}

		for j := 0; j < 3; j++ {
			css := cs + coff[j]
			rss := rs + roff[j]
			if Abs(css)+Abs(rss) < l && (css != 0 || rss != 0) {
				l = Abs(css) + Abs(rss)
				if rss < 0 {
					cm = -css
					rm = -rss
				} else {
					cm = css
					rm = rss
				}
			}
		}
	}

	return Point{R: 0, C: 0}, false
}

// ReduceReduce takes a list of translation offsets and generates the shortest spanning set
// of offsets
func (t *Torus) ReduceReduce(in []Point) []Point {
	out := make([]Point, 0)
	left := make([]Point, 0)

	if len(in) == 1 {
		out = append(out, in[0])
		return out
	}

	// figure out shortest line in set
	l := Abs(in[0].R) + Abs(in[0].C)
	min := 0
	for i, p := range in[1:] {
		if Abs(p.R)+Abs(p.C) < l {
			l = Abs(p.R) + Abs(p.C)
			min = i
		}
	}

	pm := in[min]
	for i, p := range in {
		if i == min || t.TranslationEquiv(pm, p) {
			continue
		}
		left = append(left, p)
	}
	if len(left) == 0 {
		out = append(out, pm)
		return out
	}

	return append(out, t.ReduceReduce(left)...)
}

// TranslationEquiv returns true if pm and p are equivalent translations
func (t *Torus) TranslationEquiv(pm, p Point) bool {
	if pm.R != 0 && Abs(pm.R) < Abs(p.R) && Abs(p.R)%Abs(pm.R) == 0 {
		if p.C == pm.C*(p.R/pm.R) || p.C == pm.C*(p.R/pm.R)-t.Cols {
			return true
		}
	}
	if pm.C != 0 && Abs(pm.C) < Abs(p.C) && Abs(p.C)%Abs(pm.C) == 0 {
		if p.R == pm.R*(p.C/pm.C) || p.R == pm.R*(p.C/pm.C)-t.Rows {
			return true
		}
	}
	return false
}

// Translations produces the list of locations generated from a given
// Location and translation.  It will return an empty slice in the
// event that the translation is not periodic in maxcells steps.
func (t *Torus) Translations(l1 Location, o Point, ll []Location, maxcells int) []Location {
	ll = append(ll, l1)
	p1 := t.ToPoint(l1)
	p := Point{}
	for i := 1; i < maxcells+1; i++ {
		p.C = (p1.C + i*o.C) % t.Cols
		p.R = (p1.R + i*o.R) % t.Rows
		if p.C < 0 {
			p.C += t.Cols
		}
		if p.R < 0 {
			p.R += t.Rows
		}
		if p.R == p1.R && p.C == p1.C {
			return ll
		}
		ll = append(ll, Location(p.R*t.Cols+p.C))
	}
	return []Location{}
}

// TransMap returns a mapping from a given location to the its equiv set.
func (t *Torus) TransMap(p Point, maxcells int) [][]Location {
	size := t.Size()
	smap := make([][]Location, size)
	marr := make([]Location, 0, size)

	n := 0
	for i, _ := range smap {
		if smap[i] == nil {
			marr = t.Translations(Location(i), p, marr, maxcells)
			if len(marr) == 0 || len(marr) > size {
				return nil
			}
			for _, loc := range marr[n:] {
				smap[loc] = marr[n:]
			}
			n = len(marr)
		}
	}

	return smap
}
