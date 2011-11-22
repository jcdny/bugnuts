package maps

import (
	. "bugnuts/util"
)

// Neighboorhood size, if > 5 then SymHash needs to be int64
// Also the tests are all coded to expect 5x5
const (
	SYMN       = 7
	SYMCELLMAX = 32 // maximum number of cells for tranlations...
)

type SymHash int64

type SymTile struct {
	Hash     SymHash    // the minhash
	Locs     []Location // encountered tiles with this minhash
	Bits     uint8      // bits of info Min(SYMN*SYMN - N*Water, N*Water)
	Self     uint8      // number of matching self rotations
	Ignore   bool       // Ignore this tile for symmetry stuff.
	Symmetry []uint8    // The list symmetries present
	Origin   Point      // Origin for the currently accepted Symmetry, {0,0} for translation
	Offset   Point      // The offset for translation symmetry {0,0} for non translation
	EquivSet []Location // the location imputed to
}

type SymData struct {
	*Map                         // The associated map for the Symmetry data.
	MinBits uint8                // Ignore hashes with less than MinBits bits of different info
	NLen    [16]int              // Number of equiv group for a given N
	MinHash []SymHash            // Sym data for a given point.
	Hashes  []*[8]SymHash        // Map from the location to all rotations of the given location
	Tiles   map[SymHash]*SymTile // Map from minhash to location list.
}

// The bit shuffle for the 8 symmetries a SYMNxSYMN neighborhood
var symMask [SYMN * SYMN][8]SymHash

// Map {r, c} -> {r*rr + c*cr, c*cc+ r*rc}
type symOffsets struct {
	RR int
	CR int
	RC int
	CC int
}

var symOffsetMap = [8]symOffsets{
	{1, 0, 0, 1},   // translation
	{1, 0, 0, -1},  // mirror vert
	{-1, 0, 0, 1},  // mirror horiz
	{0, -1, 1, 0},  // ccw 90
	{-1, 0, 0, -1}, // ccw 180
	{0, 1, -1, 0},  // ccw 270
	{0, 1, 1, 0},   // rot/mirror
	{0, -1, -1, 0}, // rot/mirror
}

const (
	SYMTRANS = iota
	SYMMIRRORC
	SYMMIRRORR
	SYMROT90
	SYMROT180
	SYMROT270
	SYMRM1
	SYMRM2
	SYMNONE
)

// Number of symmetry axes
type symAxes struct {
	Name string
	Id   uint8
	N    int
	R    bool
	C    bool
	D    bool
}

var symAxesMap = [8]symAxes{
	{"Trans", SYMTRANS, 0, false, false, false},
	{"mir-C", SYMMIRRORC, 1, false, true, false},
	{"mir-R", SYMMIRRORR, 1, true, false, false},
	{"rt_90", SYMROT90, 2, true, true, false},
	{"rt180", SYMROT180, 2, true, true, false},
	{"rt270", SYMROT270, 2, true, true, false},
	{"-rm1-", SYMRM1, 1, false, false, true},
	{"-rm2-", SYMRM2, 1, false, false, true},
}

func init() {
	// Generate the shuffle masks for the symmetries as defined by symOffsetMap

	// first the stepping arrays so the logic is clearer below 
	steps := [3][SYMN]int{}
	for i := 0; i < SYMN; i++ {
		steps[0][i] = SYMN - i - 1
		steps[2][i] = i
	}

	for sym, om := range symOffsetMap {
		bit := uint8(0)
		if om.RR != 0 {
			for _, r := range steps[om.RR+1] {
				for _, c := range steps[om.CC+1] {
					symMask[r*SYMN+c][sym] ^= 1 << bit
					bit++
				}
			}
		} else {
			for _, c := range steps[om.CR+1] {
				for _, r := range steps[om.RC+1] {
					symMask[r*SYMN+c][sym] ^= 1 << bit
					bit++
				}
			}
		}
	}
}

func (m *Map) NewSymData(minBits uint8) *SymData {
	s := SymData{
		Map:     m,
		MinBits: minBits,
		MinHash: make([]SymHash, m.Size()),
		Hashes:  make([]*[8]SymHash, m.Size()),
		Tiles:   make(map[SymHash]*SymTile, m.Size()/4),
	}
	s.NLen[0] = s.Size()

	return &s
}

// Tiles an entire map.
func (m *Map) Tile(minBits uint8) *SymData {
	s := m.NewSymData(minBits)

	for loc, _ := range m.Grid {
		s.Update(Location(loc))
	}

	return s
}

func (s *SymData) UpdateSymmetryData() {
	check := make(map[SymHash]bool, 100)
	for l, item := range s.TGrid {
		loc := Location(l)
		if item != UNKNOWN && s.Hashes[loc] == nil {
			hash, found := s.Update(loc)
			if found {
				check[hash] = true
			}
		}
	}
	for minhash, _ := range check {
		symset, origin, offset := s.SymAnalyze(minhash)
		s.Tiles[minhash].Symmetry = symset
		s.Tiles[minhash].Origin = origin
		s.Tiles[minhash].Offset = offset
	}
}

// Returns the minhash, true if there is a potential new symmetry
func (s *SymData) Update(loc Location) (SymHash, bool) {
	var found bool

	minhash, hashes, bits, self := s.SymCompute(Location(loc))
	s.MinHash[loc] = minhash
	s.Hashes[loc] = hashes
	if hashes != nil {
		_, found := s.Tiles[minhash]
		if !found {
			// first time we have seen this minhash
			s.Tiles[minhash] = &SymTile{
				Hash: minhash,
				Locs: make([]Location, 0, 4),
				Bits: bits,
				Self: self,
			}
		}

		// Keep track of number of equiv classes
		N := len(s.Tiles[minhash].Locs)
		s.NLen[0]--
		if N > 0 && N < len(s.NLen) {
			s.NLen[N]--
		}
		if N+1 < len(s.NLen) {
			s.NLen[N+1]++
		} else if N == len(s.NLen) {
			s.NLen[N-1]++
		}

		s.Tiles[minhash].Locs = append(s.Tiles[minhash].Locs, Location(loc))
	}

	return minhash, found
}

// Compute the minhash for a given location, returning the bits of data, the minHash and all
// 8 hashes.  It returns (0, -1, nil) in the event it encounters an unknown tile...
func (s *SymData) SymCompute(loc Location) (SymHash, *[8]SymHash, uint8, uint8) {
	p := s.ToPoint(loc)
	id := [8]SymHash{}

	i := 0
	nl := loc
	N := SYMN / 2
	bits := 0
	// TODO this could be a lot faster.
	// TODO also might be worth discarding all land tiles quickly
	for r := -N; r < N+1; r++ {
		for c := -N; c < N+1; c++ {
			if p.R < N || p.R > s.Rows-N-1 || p.C < N || p.C > s.Cols-N-1 {
				nl = s.ToLocation(s.PointAdd(p, Point{R: r, C: c}))
			} else {
				nl = loc + Location(r*s.Cols+c)
			}

			if s.Grid[nl] == UNKNOWN {
				return -1, nil, SYMNONE, 0
			}

			if s.Grid[nl] == WATER {
				bits++
				for rot, mask := range symMask[i] {
					id[rot] ^= mask
				}
			}
			i++
		}
	}
	if bits > (SYMN*SYMN)/2 {
		bits = SYMN*SYMN - bits
	}

	self := 0
	for i := 1; i < 8; i++ {
		if id[0] == id[i] {
			self++
		}
	}

	return minSymHash(&id), &id, uint8(bits), uint8(self)
}

// annoying utility func.
func minSymHash(id *[8]SymHash) SymHash {
	// unrolled version was 300us faster over a 37ms tile of a full map...
	min := id[0]
	for i := 1; i < 8; i++ {
		if id[i] < min {
			min = id[i]
		}
	}
	return min
}

func (s *SymData) SymAnalyze(minhash SymHash) ([]uint8, Point, Point) {
	llist := s.Tiles[minhash].Locs

	// test for Translation symmetry
	redlist := make([]Point, 0, 0)
	bad := 0
	for i, l1 := range llist {
		for _, l2 := range llist[i+1:] {
			if s.Hashes[l1][0] == s.Hashes[l2][SYMTRANS] {
				pd, good := s.ShiftReduce(l1, l2)
				if !good {
					bad++
				} else {
					redlist = append(redlist, pd)
				}
			} else {
				bad++
			}
		}
	}
	if bad == 0 && len(redlist) > 0 {
		redlist = s.ReduceReduce(redlist)
		if len(redlist) == 1 {
			// Yay we got unambiguous translation...
			return []uint8{SYMTRANS}, Point{0, 0}, redlist[0]
		}
	}

	// Test for mirroring
	found := make([]uint8, 0, 3)
	orig := Point{0, 0}
	for i, l1 := range llist {
		for _, l2 := range llist[i+1:] {
			if s.Hashes[l1][0] == s.Hashes[l2][SYMMIRRORC] {
				orig.C = s.Mirror(l1, l2, 1)
				found = append(found, SYMMIRRORC)
			}
			if s.Hashes[l1][0] == s.Hashes[l2][SYMMIRRORR] {
				orig.R = s.Mirror(l1, l2, 0)
				found = append(found, SYMMIRRORR)
			}
			if s.Hashes[l1][0] == s.Hashes[l2][SYMROT180] {
				orig.C = s.Mirror(l1, l2, 1)
				orig.R = s.Mirror(l1, l2, 0)
				found = append(found, SYMROT180)
			}

		}
	}
	if len(found) > 1 {
		return []uint8{SYMMIRRORC, SYMMIRRORR, SYMROT180}, orig, Point{0, 0}
	} else {
		return found, orig, Point{0, 0}
	}

	// Test for rotations
	// TODO
	return []uint8{}, Point{0, 0}, Point{0, 0}
}

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

// Reduce a translation to its minumum length offset
// I should just do this with math but my head hurts.
func (t *Torus) ShiftReduce(l1, l2 Location) (Point, bool) {
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

	for i := 0; i < SYMCELLMAX+1; i++ {
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

// Take a list of translation offsets and generate list of shortest spanning set
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
		if i == min || t.EquivT(pm, p) {
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

func (t *Torus) EquivT(pm, p Point) bool {
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

// Given a point and a translation compute the list of locations
func (t *Torus) Translations(l1 Location, o Point, ll []Location) []Location {
	ll = append(ll, l1)
	p1 := t.ToPoint(l1)
	p := Point{}
	for i := 1; i < SYMCELLMAX+1; i++ {
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

func (t *Torus) TransMap(p Point) [][]Location {
	size := t.Size()
	smap := make([][]Location, size)
	marr := make([]Location, 0, size)

	n := 0
	for i, _ := range smap {
		if smap[i] == nil {
			marr = t.Translations(Location(i), p, marr)
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
