// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package torus

import (
	"log"
	. "bugnuts/util"
)

// Rot determines the rotation origin for a pair of points where p2 is a 90 degree rotation from p1.
func (t *Torus) Rot(p1, p2 Point, sym int) (orig Point) {
	var r0, c0 int
	if sym == 3 {
		c0 = (p1.C + p2.C + p1.R - p2.R + 1) / 2
		r0 = (p2.C - p1.C + p1.R + p2.R + 1) / 2
		orig = t.Donut(Point{r0, c0})
		if orig.R >= t.Rows/2 {
			orig.R -= t.Rows / 2
			orig.C -= t.Cols / 2
			if orig.C < 0 {
				orig.C += t.Cols
			}
		}
	} else {
		p := t.Mirror(p1, p2, 0)
		r0 = p.R
		p = t.Mirror(p1, p2, 1)
		c0 = p.C
		orig = Point{r0, c0}
	}
	//log.Print("ROT", t, p1, p2, r0, c0)

	return
}

// ReflectRM1 computes the reflection of p1 across origin o.
func (t Torus) ReflectRM1(p1 Point, o int) Point {
	return t.Donut(Point{p1.C - o, p1.R + o})
}

// ReflectRM2 computes the reflection of p1 across origin o.
func (t Torus) ReflectRM2(p1 Point, o int) Point {
	return t.Donut(Point{o - p1.C, o - p1.R})
}

// Diag computes the origin of a pair of points which are diagonal mirrorings.
func (t *Torus) Diag(p1, p2 Point, sym int) (orig int) {
	orig = -1
	if sym == 6 {
		// RM1
		orig = (p1.C + p2.C - p1.R - p2.R) / 2
		for orig < 0 {
			orig += t.Cols
		}
	} else if sym == 7 {
		// RM2
		// solve for R = 0
		orig = (p1.C + p2.C + p1.R + p2.R) / 2
		for orig-t.Cols > 0 {
			orig -= t.Cols
		}
	}
	//log.Print("sym,p1,p2,orig:", sym, p1, p2, orig)

	return
}

func (t *Torus) Mirrors(p Point, mr, mc int) []Point {
	pv := make([]Point, 1, 4)
	pv[0] = p
	if mr > -1 {
		p.R = 2*mr - p.R
		pv = append(pv, t.Donut(p))
	}
	if mc > -1 {
		for _, pp := range pv {
			pp.C = 2*mc - p.C
			pv = append(pv, t.Donut(pp))
		}
	}
	return pv
}

// Mirror returns the midpoint between two points.  Axis 0 is R and 1
// is C.  By convention it will return either the axis closer to {0,0}
// if the side dimension is even or the axis which exists between the
// even spread of p1/p2.
func (t *Torus) Mirror(p1, p2 Point, axis int) Point {

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
		return Point{}
	}

	//log.Print("mirror", axis, *t, p1, p2, odd, o, s, (s+1)/2)
	if !odd {
		if o >= (s+1)/2 {
			o -= (s + 1) / 2
		}
	} else {
		o = (o + (s+1)/2) % s
	}

	out := Point{}
	if axis == 0 {
		out.R = o
	} else {
		out.C = o
	}

	return out
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
func (t *Torus) ShiftReduce(p1, p2 Point, maxcells int) (Point, bool) {
	if t.PointEqual(p1, p2) {
		return Point{}, false
	}

	r := p2.R - p1.R
	c := p2.C - p1.C

	if r < 0 {
		r += t.Rows
	}
	if c < 0 {
		c += t.Cols
	}

	if r == 0 || c == 0 {
		return Point{}, false
	}

	n := Lcm(Lcm(t.Rows, r)/r, Lcm(t.Cols, c)/c)
	if n > maxcells {
		return Point{}, false
	}

	l := 65535
	rm, cm := r, c
	coff := [3]int{0, 0, -t.Cols}
	roff := [3]int{0, -t.Rows, 0}

	for i := 0; i < n; i++ {
		cs := (c + i*c) % t.Cols
		rs := (r + i*r) % t.Rows
		for j := 0; j < 3; j++ {
			css := cs + coff[j]
			rss := rs + roff[j]
			if Abs(css)+Abs(rss) < l && (css != 0 && rss != 0) {
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

	//log.Print(r, c, n, rm, cm)

	return Point{R: rm, C: cm}, true
}

// point to the list
func (t *Torus) SymAddPoint(in []Point, p Point) []Point {
	for p.C-t.Cols > 0 {
		p.C -= t.Cols
	}
	for p.R-t.Rows > 0 {
		p.R -= t.Rows
	}

	for _, ip := range in {
		if ip.C == p.C && ip.R == p.R {
			return in
		}
	}

	return append(in, p)
}

func SetAddInt(in []int, p int) []int {
	for _, ip := range in {
		if ip == p {
			return in
		}
	}

	return append(in, p)
}

func SetAddPoint(in []Point, p Point) []Point {
	for _, ip := range in {
		if ip.R == p.R && ip.C == p.C {
			return in
		}
	}

	return append(in, p)
}

// Given a set of translation offsets returns the tiling dimension
func (t *Torus) BlockDim(in []Point) (dim Torus) {
	dim = *t
	for _, p := range in {
		if p.R != p.C && (p.R == 0 || p.C == 0) {
			if p.C == 0 {
				dim.Rows = p.R
			}
			if p.R == 0 {
				dim.Cols = p.C
			}
		}
	}

	return dim
}

func (t *Torus) Translation(in []Point) Point {
	tr := Point{}
	n := 0
	for _, p := range in {
		if p.R != 0 && p.C != 0 {
			tr = p
			n++
		}
	}
	if n == 1 {
		return tr
	}
	return Point{}
}
func (t *Torus) TranslationLen(p Point) int {
	if p.R == 0 && p.C == 0 {
		return 1
	}
	if p.R == 0 {
		return Abs(t.Cols / p.C)
	}
	if p.C == 0 {
		return Abs(t.Rows / p.R)
	}
	result := Abs(Lcm(Lcm(p.R, t.Rows)/p.R, Lcm(p.C, t.Cols)/p.C))
	//log.Print(t, p, " LEN: ", result, " lcm(r,R)/r ", Lcm(p.R, t.Rows)/p.R, " lcm(c,C)/c ", Lcm(p.C, t.Cols)/p.C)

	return result
}

// ReduceReduce takes a list of translation offsets and generates the shortest spanning set
// of offsets, will remove subtiles.
func (t *Torus) ReduceReduce(in []Point) []Point {
	out := make([]Point, 0)
	left := make([]Point, 0)

	n := 0
	for i, p := range in {
		if Abs(p.R) < 5 || Abs(p.C) < 5 {
			// drop nonsensical translates
			continue
		} else {
			if n > 0 && t.PointEqual(in[n-1], p) {
				continue
			} else {
				if n < i {
					in[n] = p
				}
				n++
			}
		}
	}
	//log.Print(n, in[:n])
	in = in[:n]

	if len(in) == 0 {
		return []Point{}
	} else if len(in) == 1 {
		out = append(out, in[0])
		return out
	}

	// figure out shortest line in set
	l := 65535
	min := -1
	//log.Print("using in", in)
	for i, p := range in {
		if Abs(p.R)+Abs(p.C) < l {
			l = Abs(p.R) + Abs(p.C)
			min = i
		}
	}
	if min < 0 {
		log.Print("no min")
		return []Point{}
	} else {
		//log.Print("min ", in[min])
	}

	pm := in[min]
	for _, p := range in {
		if t.PointEqual(pm, p) || t.TranslationEquiv(pm, p) {
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
		n := pm.C * p.R / pm.R
		for n < p.C {
			n += t.Cols
		}
		for n > p.C {
			n -= t.Cols
		}
		return n == p.C
	}

	if pm.C != 0 && Abs(pm.C) < Abs(p.C) && Abs(p.C)%Abs(pm.C) == 0 {
		n := pm.R * p.C / pm.C
		for n < p.R {
			n += t.Rows
		}
		for n > p.R {
			n -= t.Rows
		}
		return n == p.R
	}
	return false
}

// Translations produces the list of locations generated from a given
// Location and translation.  It will return an empty slice in the
// event that the translation is not periodic in maxcells steps.
func (t *Torus) Translations(m Torus, l1 Location, o Point, ll []Location, maxcells int) []Location {
	// NB t is the subtile and m the larger map -- jump through hoops to make 
	// locations correct...
	p1 := t.Donut(m.ToPoint(l1))
	ll = append(ll, m.ToLocation(p1))
	// log.Print(t, m, l1, m.ToPoint(l1), m.ToLocation(p1), len(ll), maxcells)
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
		ll = append(ll, Location(p.R*m.Cols+p.C))
	}
	return []Location{}
}

// TransMap returns a mapping from a given location to the its equiv set.
func (t *Torus) TransMap(p Point, maxcells int) [][]Location {
	size := t.Size()
	smap := make([][]Location, size)
	marr := make([]Location, 0, size)

	n := 0
	for i := range smap {
		if smap[i] == nil {
			marr = t.Translations(*t, Location(i), p, marr, maxcells)
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
