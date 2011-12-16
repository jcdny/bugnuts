package maps

import (
	"log"
	"os"
	"fmt"
	. "bugnuts/torus"
	. "bugnuts/debug"
	. "bugnuts/util"
)

const (
	SYMN        = 7  // Neighboorhood size.  needs to be odd and < 8
	SYMMAXCELLS = 32 // maximum number of cells for tranlations...
)

type SymHash int64 // SymHash needs to be int64 if SYMN = 7, int32 otherwise.

type SymTile struct {
	Hash      SymHash    // the minhash
	Locs      []Location // encountered tiles with this minhash
	Bits      uint8      // bits of info Min(SYMN*SYMN - N*Water, N*Water)
	Self      uint8      // number of matching self rotations
	Ignore    bool       // Ignore this tile for symmetry stuff.
	Origin    Point      // Origin for discovered symmetries
	Gen       int        // the generator for the origin
	Translate Point      // The offset for translation symmetry {0,0} for non translation
	Subtile   Torus      // The dimensions for the subtile == the map dim if none
	EquivSet  []Location // the location list for the identified symmetry of this tile.
}

type SymData struct {
	*Map                         // The associated map for the Symmetry data.
	Offsets                      // The offsets cache
	MinBits uint8                // Ignore hashes with less than MinBits bits of different info
	NLen    [16]int              // Number of equiv group for a given N
	MinHash []SymHash            // Sym data for a given point.
	Hashes  []*[8]SymHash        // Map from the location to all rotations of the given location
	Tiles   map[SymHash]*SymTile // Map from minhash to location list.
}

// The bit shuffle for the 8 symmetries a SYMNxSYMN neighborhood
var symMask [SYMN * SYMN][8]SymHash
var symPointOffsets []Point

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
	TPush("Init sym")
	defer TPop()

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

	// populate and cache the Offseterator
	symPointOffsets = make([]Point, SYMN*SYMN)
	i := 0
	for r := -SYMN / 2; r < (SYMN+1)/2; r++ {
		for c := -SYMN / 2; c < (SYMN+1)/2; c++ {
			symPointOffsets[i] = Point{R: r, C: c}
			i++
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
	s.Offsets = PointsToOffsets(symPointOffsets, m.Cols)
	// Don't prepopulate since it's not used much....
	// s.offsetsCachePopulate(&s.Offsets)

	return &s
}

// Tiles an entire map.
func (m *Map) Tile(minBits uint8) *SymData {
	s := m.NewSymData(minBits)

	for loc := range m.Grid {
		s.Update(Location(loc))
	}

	return s
}

func (s *SymData) UpdateSymmetryData() {
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
		maxlen = len(s.Map.SMap[0])
	}

	for minhash := range check {
		tile := s.Tiles[minhash]
		eqlen := s.SymAnalyze(tile)
		//log.Print("symset, origin, offset, equiv:", symset, origin, offset, equiv)
		if !tile.Ignore && eqlen > maxlen {
			smap, valid := s.SymMapValidate(tile)
			if valid {
				if Debug[DBG_Symmetry] {
					log.Printf("Valid symmetry map len %d found", len(smap[0]))
				}
				maxlen = eqlen
				s.Map.SMap = smap
				updated = true
			} else {
				tile.Ignore = true
			}
			if false {
				VizSymTile(s.ToPoints(tile.Locs), valid)
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

func VizSymTile(pv []Point, valid bool) {
	if valid {
		fmt.Fprintf(os.Stdout, "v slc %d %d %d %.2f\n",
			0, 255, 0, 1.0)
	} else {
		fmt.Fprintf(os.Stdout, "v slc %d %d %d %.2f\n",
			255, 0, 0, 1.0)
	}
	for _, p := range pv {
		fmt.Fprintf(os.Stdout, "v r %d %d 6 6 false\n", p.R-3, p.C-3)
		fmt.Fprintf(os.Stdout, "v t %d %d \n", p.R, p.C)
	}
	fmt.Fprintf(os.Stdout, "v slc %d %d %d %.2f\n",
		0, 0, 0, 1.0)
}

// Returns the minhash, true if there is a potential new symmetry
func (s *SymData) Update(loc Location) (SymHash, bool) {
	newsym := false

	minhash, hashes, bits, self := s.SymCompute(Location(loc))

	if Debug[DBG_Symmetry] && WS.Watched(loc, 0) {
		log.Print("Minhash point, minhash, bits, self", s.ToPoint(loc), minhash, bits, self)
	}
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
	g := s.TGrid
	// TODO this might be faster...
	// TODO also might be worth discarding all land tiles quickly
	for r := -N; r < N+1; r++ {
		for c := -N; c < N+1; c++ {
			if p.R < N || p.R > s.Rows-N-1 || p.C < N || p.C > s.Cols-N-1 {
				nl = s.ToLocation(s.PointAdd(p, Point{R: r, C: c}))
			} else {
				nl = loc + Location(r*s.Cols+c)
			}

			if g[nl] == UNKNOWN {
				return -1, nil, SYMNONE, 0
			}

			if g[nl] == WATER {
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

// Compute the minhash for a given location, returning the bits of data, the minHash and all
// 8 hashes.  It returns (0, -1, nil) in the event it encounters an unknown tile...
// This is slower than the version above although I find it surprising that it is...
func (s *SymData) slowSymCompute(loc Location) (SymHash, *[8]SymHash, uint8, uint8) {
	id := [8]SymHash{}

	i := 0
	bits := 0
	g := s.TGrid
	s.ApplyOffsetsBreak(loc, &s.Offsets, func(nl Location) bool {
		if g[nl] == UNKNOWN {
			bits = -1
			return false
		}
		if g[nl] == WATER {
			bits++
			for rot, mask := range symMask[i] {
				id[rot] ^= mask
			}
		}
		i++
		return true
	})

	if bits < 0 {
		return -1, nil, SYMNONE, 0
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

func (s *SymData) Tiling(tile *SymTile) Torus {
	dim := s.Torus
	ndim := dim
	for i, l1 := range tile.Locs {
		p1 := dim.ToPoint(l1)
		for _, l2 := range tile.Locs[i+1:] {
			if s.Hashes[l1][0] == s.Hashes[l2][SYMTRANS] {
				p2 := dim.ToPoint(l2)
				if p1.C == p2.C {
					s := Abs(p2.R - p1.R)
					if s < ndim.Rows {
						if dim.Rows == Lcm(dim.Rows, s) {
							ndim.Rows = s
						}
					}
				}
				if p1.R == p2.R {
					s := Abs(p2.C - p1.C)
					if s < ndim.Cols {
						if dim.Cols == Lcm(dim.Cols, s) {
							ndim.Cols = s
						}
					}
				}
			}
		}
	}
	return ndim
}
// Update the analysis of a tile and return the length of the infered equiv set
func (s *SymData) SymAnalyze(tile *SymTile) (equivlen int) {
	dbg := false
	//dbg := true
	tile.Gen = SYMNONE
	equivlen = 0

	if tile == nil || len(tile.Locs) < 2 {
		return
	}

	llist := tile.Locs

	// Get the blocking for the map
	dim := s.Tiling(tile)
	if dbg {
		log.Print(s.Torus, dim)
	}

	// test for Translation symmetry
	redlist := make([]Point, 0, len(llist))
	n := 0
	for i, l1 := range llist {
		p1 := dim.Donut(s.ToPoint(l1))
		for _, l2 := range llist[i+1:] {
			if s.Hashes[l1][0] == s.Hashes[l2][SYMTRANS] {
				p2 := dim.Donut(s.ToPoint(l2))
				n++
				pd, good := dim.ShiftReduce(p1, p2, SYMMAXCELLS)
				if false && dbg {
					log.Print(p1, p2, pd, good)
				}
				if good {
					redlist = append(redlist, pd)
				}
			}
		}
	}
	if dbg {
		log.Print("pre reduce", redlist)
	}
	if len(redlist) > 0 {
		redlist = dim.ReduceReduce(redlist)
	}
	if dbg {
		log.Print("post reduce", redlist)
	}

	tile.Gen = SYMTRANS
	tile.Subtile = dim
	tile.Translate = s.Translation(redlist)
	if dbg {
		log.Print("Len is ", s.Size()/dim.Size(), " * ", dim.TranslationLen(tile.Translate))
	}
	equivlen = s.Size() / dim.Size() * dim.TranslationLen(tile.Translate)

	// If all we got was translations bail out.
	if n == len(llist)*(len(llist)-1)/2 {
		return
	}

	// Look for rotational symmetry, iff we have a square block
	rotorig := make([]Point, 0, 0)
	if dim.Rows == dim.Cols {
		for i, l1 := range llist {
			p1 := dim.Donut(s.ToPoint(l1))
			for _, l2 := range llist[i+1:] {
				p2 := dim.Donut(s.ToPoint(l2))
				if s.Hashes[l1][0] == s.Hashes[l2][SYMROT90] {
					porig := dim.Rot(p1, p2, SYMROT90)
					rotorig = dim.SymAddPoint(rotorig, porig)
				}
			}
		}
	}

	// can only have 1 rotation origin.
	if len(rotorig) > 1 {
		rotorig = rotorig[:0]
	} else if len(rotorig) == 1 {
		tile.Origin = rotorig[0]
		tile.Gen = SYMROT90
		equivlen *= 4

		return
	}

	// Look for mirror symmetry iff we did not find rot sym
	morigc := make([]Point, 0, 0)
	morigr := make([]Point, 0, 0)
	if len(rotorig) == 0 {
		for i, l1 := range llist {
			p1 := dim.Donut(s.ToPoint(l1))
			for _, l2 := range llist[i+1:] {
				p2 := dim.Donut(s.ToPoint(l2))
				if s.Hashes[l1][0] == s.Hashes[l2][SYMMIRRORC] {
					morigc = dim.SymAddPoint(morigc, dim.Mirror(p1, p2, 1))
				}
				if s.Hashes[l1][0] == s.Hashes[l2][SYMMIRRORR] {
					morigr = dim.SymAddPoint(morigr, dim.Mirror(p1, p2, 0))
				}
			}
		}
	}
	if len(morigc) == 1 {
		tile.Origin.C = morigc[0].C
		equivlen *= 2
		tile.Gen = SYMMIRRORC
	}
	if len(morigr) == 1 {
		tile.Origin.R = morigr[0].R
		equivlen *= 2
		tile.Gen = SYMMIRRORR
	}
	if tile.Gen == SYMMIRRORR || tile.Gen == SYMMIRRORC {
		return equivlen
	}

	mrotorig := make([]Point, 0, 0)
	for i, l1 := range llist {
		p1 := dim.Donut(s.ToPoint(l1))
		for _, l2 := range llist[i+1:] {
			p2 := dim.Donut(s.ToPoint(l2))
			//log.Print("ROT:", dim, s.ToPoint(l2), p2)
			if s.Hashes[l1][0] == s.Hashes[l2][SYMROT180] {
				mrotorig = dim.SymAddPoint(mrotorig, dim.Rot(p1, p2, SYMROT180))
			}
		}
	}

	if len(mrotorig) > 1 {
		mrotorig = mrotorig[:0]
	} else if len(mrotorig) == 1 {
		tile.Origin = mrotorig[0]
		tile.Gen = SYMROT180
		equivlen *= 2

		return
	}

	// Look for diagonal symmetry iff we did not find rot/mirror sym
	dorig := make([]Point, 0, 0)
	if false {
		for i, l1 := range llist {
			p1 := dim.Donut(s.ToPoint(l1))
			for _, l2 := range llist[i+1:] {
				p2 := dim.Donut(s.ToPoint(l2))
				if s.Hashes[l1][0] == s.Hashes[l2][SYMRM1] {
					dorig = dim.SymAddPoint(dorig, s.Diag(p1, p2, SYMRM1))
				}
				if s.Hashes[l1][0] == s.Hashes[l2][SYMRM2] {
					dorig = dim.SymAddPoint(dorig, s.Diag(p1, p2, SYMRM2))
				}
			}
		}
	}

	if len(dorig) > 0 {
		tile.Origin = dorig[0]
		tile.Gen = SYMRM1
		equivlen *= 2
	}

	return
}

func (s *SymData) TransMapValidate(p Point) ([][]Location, bool) {
	size := s.Size()
	smap := make([][]Location, size)
	marr := make([]Location, 0, size)

	n := 0
	for i := range smap {
		if smap[i] == nil {
			marr = s.Translations(Location(i), p, marr, SYMMAXCELLS)
			if false && n == 0 {
				log.Print("len(marr), size", len(marr), size)
			}
			if false && len(marr) == 0 || len(marr) > size {
				log.Print("len(marr), size", len(marr), size)
				return nil, false
			}
			item := UNKNOWN
			for _, loc := range marr[n:] {
				// Validate the equiv set is identical
				if item == UNKNOWN {
					item = s.TGrid[loc]
				} else if item != s.TGrid[loc] {
					// log.Print("i, n, loc, item, tgrid ", i, n, loc, int(item), s.TGrid[loc])
					return nil, false
				}
				smap[loc] = marr[n:]
			}
			n = len(marr)
		}
	}

	return smap, true
}

// Given an analyzed tile generate the map for loc, appending map to marr
func (tile *SymTile) Generate(t Torus, loc Location, marr []Location) []Location {
	return marr
}

// Takes a tile which has been analyzed and generates a sym map for it
// and simultaneously validates it.
func (s *SymData) SymMapValidate(tile *SymTile) ([][]Location, bool) {
	size := s.Size()
	smap := make([][]Location, size)
	marr := make([]Location, 0, size)

	n := 0
	// Take the first location we found as the starting point
	loc := int(tile.Locs[0])
	for i := range smap {
		if smap[loc] == nil {
			marr = tile.Generate(s.Torus, Location(loc), marr)
			if n == 0 {
				log.Print("len(marr), size", len(marr), size)
			}
			if len(marr) == 0 || len(marr) > size {
				log.Print("Invalid map len(marr), size, points", len(marr), size, marr[n:])
				return nil, false
			}
			item := UNKNOWN
			for _, mloc := range marr[n:] {
				// Validate the equiv set is identical
				if item == UNKNOWN {
					item = s.TGrid[mloc]
				} else if item != s.TGrid[mloc] {
					log.Print("Invalid point found: i, n, loc, item, tgrid ", i, n, loc, int(item), s.TGrid[loc])
					return nil, false
				}
				smap[mloc] = marr[n:]
			}
			n = len(marr)
		}
		// 1327 is 10 less than 1337 (and prime which is perhaps more important)
		// Do this to avoid worst case behavior where we are in eg center
		// with an invalid rotational symmtery and if we start at 0
		// we could potentially generate 70% of the map before encountering any data at all
		loc = (loc + 1327) % size
	}
	// Sanity check
	if len(marr) != size {
		log.Print("Tiling size mismatch ", len(marr), size)

		return nil, false
	}

	return smap, true
}

func (tile *SymTile) String() string {
	s := ""
	s += fmt.Sprintf("Hash: %d Bits: %d Self: %d Origin: %v Gen: %d Translate %v Subtile: %v",
		tile.Hash, tile.Bits, tile.Self, tile.Origin, tile.Gen, tile.Translate, tile.Subtile)
	return s
}
