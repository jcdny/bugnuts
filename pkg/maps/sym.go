package maps

import (
	"log"
	. "bugnuts/util"
	. "bugnuts/debug"
)

const (
	SYMN       = 7  // Neighboorhood size.  needs to be odd and < 8
	SYMCELLMAX = 32 // maximum number of cells for tranlations...
)

type SymHash int64 // SymHash needs to be int64 if SYMN = 7, int32 otherwise.

type SymTile struct {
	Hash     SymHash    // the minhash
	Locs     []Location // encountered tiles with this minhash
	Bits     uint8      // bits of info Min(SYMN*SYMN - N*Water, N*Water)
	Self     uint8      // number of matching self rotations
	Ignore   bool       // Ignore this tile for symmetry stuff.
	Symmetry []uint8    // The list symmetries present
	Origin   Point      // Origin for the currently accepted Symmetry, {0,0} for translation
	Offset   Point      // The offset for translation symmetry {0,0} for non translation
	EquivSet []Location // the location list for the identified symmetry of this tile.
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
	RR, CR, RC, CC int
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
	Name    string
	Id      uint8
	N       int
	R, C, D bool
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

	steps := [3][SYMN]int{}
	for i := 0; i < SYMN; i++ {
		steps[0][i] = SYMN - i - 1 // negative steps
		steps[2][i] = i            // positive steps
	}

	for sym, omap := range symOffsetMap {
		bit := uint8(0)
		if omap.RR != 0 {
			// columns first
			for _, r := range steps[omap.RR+1] {
				for _, c := range steps[omap.CC+1] {
					symMask[r*SYMN+c][sym] ^= 1 << bit
					bit++
				}
			}
		} else {
			// rows first
			for _, c := range steps[omap.CR+1] {
				for _, r := range steps[omap.RC+1] {
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

	TPush("UpdateSymmetryData")
	defer TPop()

	check := make(map[SymHash]bool, 100)
	for l, item := range s.TGrid {
		loc := Location(l)
		if item != UNKNOWN && s.Hashes[loc] == nil {
			hash, newsym := s.Update(loc)
			if newsym {
				check[hash] = true
			}
		}
	}

	maxlen := 0
	updated := false
	if len(s.Map.SMap) > 0 {
		maxlen = len(s.Map.SMap)
	}

	TMark("hashed")

	for minhash, _ := range check {
		tile := s.Tiles[minhash]
		symset, origin, offset, equiv := s.SymAnalyze(minhash)
		tile.Symmetry = symset
		tile.Origin = origin
		tile.Offset = offset
		tile.EquivSet = equiv
		if !tile.Ignore && len(equiv) > maxlen {
			// TODO fix to handle other symmetries.
			smap, valid := s.TransMapValidate(tile.Offset)
			if valid {
				if Debug[DBG_Symmetry] {
					log.Printf("Valid symmetry map len %d found", len(s.Map.SMap[0]))
				}
				maxlen = len(equiv)
				s.Map.SMap = smap
				updated = true
			} else {
				tile.Ignore = true
			}
		}
	}

	if updated {
		if Debug[DBG_Symmetry] {
			log.Printf("Applying new symmetry map for sym len %d", len(s.Map.SMap[0]))
		}
		s.SID++
		s.TApply()
	}
}

// Returns the minhash, true if there is a potential new symmetry
func (s *SymData) Update(loc Location) (SymHash, bool) {
	newsym := false

	minhash, hashes, bits, self := s.SymCompute(Location(loc))
	if hashes != nil {
		s.MinHash[loc] = minhash
		s.Hashes[loc] = hashes

		tile, found := s.Tiles[minhash]
		if !found {
			// first time we have seen this minhash
			tile = &SymTile{
				Hash: minhash,
				Locs: make([]Location, 0, 4),
				Bits: bits,
				Self: self,
			}
			s.Tiles[minhash] = tile
			//log.Printf("New hash set")
		} else {
			var i int
			for i = 0; i < len(tile.EquivSet) && loc != tile.EquivSet[i]; i++ {
			}
			newsym = i < len(tile.EquivSet) || i == 0
		}

		// Keep track of number of equiv classes
		N := len(tile.Locs)
		s.NLen[0]--
		if N > 0 && N < len(s.NLen) {
			s.NLen[N]--
		}
		if N+1 < len(s.NLen) {
			s.NLen[N+1]++
		} else if N == len(s.NLen) {
			s.NLen[N-1]++
		}

		tile.Locs = append(tile.Locs, Location(loc))
	}

	return minhash, newsym
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

func (s *SymData) SymAnalyze(minhash SymHash) ([]uint8, Point, Point, []Location) {
	llist := s.Tiles[minhash].Locs
	if len(llist) < 2 {
		return []uint8{}, Point{0, 0}, Point{0, 0}, []Location{}
	}
	// test for Translation symmetry
	redlist := make([]Point, 0, 0)
	bad := 0
	for i, l1 := range llist {
		for _, l2 := range llist[i+1:] {
			if s.Hashes[l1][0] == s.Hashes[l2][SYMTRANS] {
				// TODO rectangular tilings need to be handled properly
				// Should add a test for that.
				pd, good := s.ShiftReduce(l1, l2)
				if !good {
					s.Tiles[minhash].Ignore = true
					// TODO most of the false positives happen with a mix of translate and
					// various rotations so currently I punt if I encounter a translation
					// which is not valid for tiling the torus.  Maybe wrong answer.
					return []uint8{}, Point{}, Point{}, []Location{}
				} else {
					redlist = append(redlist, pd)
				}
			} else {
				bad++
				break
			}
		}
	}

	if bad == 0 && len(redlist) > 0 {
		redlist = s.ReduceReduce(redlist)
		if len(redlist) == 1 {
			// Yay we got unambiguous translation...
			equiv := s.Translations(llist[0], redlist[0], []Location{})
			return []uint8{SYMTRANS}, Point{0, 0}, redlist[0], equiv
		}
	}

	// Test for mirroring
	found := make([]uint8, 0, 3)
	orig := Point{0, 0}
	for i, l1 := range llist {
		for _, l2 := range llist[i+1:] {
			//log.Printf("\n%#v\n%#v", s.Hashes[l1], s.Hashes[l2])
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
		return []uint8{SYMMIRRORC, SYMMIRRORR, SYMROT180}, orig, Point{0, 0}, []Location{}
	} else {
		return found, orig, Point{}, []Location{}
	}

	// TODO Test for rotations
	// For rotations and diagonal mirrorings the map needs to be square...
	return []uint8{}, Point{}, Point{}, []Location{}
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

func (s *SymData) TransMapValidate(p Point) ([][]Location, bool) {
	size := s.Size()
	smap := make([][]Location, size)
	marr := make([]Location, 0, size)

	n := 0
	for i, _ := range smap {
		if smap[i] == nil {
			marr = s.Translations(Location(i), p, marr)
			if len(marr) == 0 || len(marr) > size {
				return nil, false
			}
			item := UNKNOWN
			for _, loc := range marr[n:] {
				// Validate the equiv set only contains either land or water.
				if item != s.TGrid[loc] {
					if item == UNKNOWN {
						item = s.TGrid[loc]
					} else if s.TGrid[loc] != UNKNOWN {
						return nil, false
					}
				}
				smap[loc] = marr[n:]
			}
			n = len(marr)
		}
	}

	return smap, true
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
