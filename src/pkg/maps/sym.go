package maps

import (
	"log"
	"os"
	"fmt"
	. "bugnuts/torus"
	. "bugnuts/watcher"
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
	Subtile   Torus      // The dimensions for the subtile == the map dim if none
	Gen       int        // the generator for the origin
	Origin    Point      // Origin for discovered symmetries
	RM1       int
	RM2       int
	MR        int
	MC        int
	Translate Point      // The offset for translation symmetry {0,0} for non translation
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
	Fails   int
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
	{0, 1, 1, 0},   // rot/mirror, diagonal 
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
	SYMMIRBOTH
	SYMRMBOTH
	SYMMIR8
	SYMROT8
)

var symLabels = [...]string{
	SYMTRANS:   "Trans",
	SYMMIRRORC: "MirrC",
	SYMMIRRORR: "MirrR",
	SYMROT90:   "Rotat",
	SYMROT180:  "R_180",
	SYMROT270:  "R_270",
	SYMRM1:     "NDiag",
	SYMRM2:     "PDiag",
	SYMNONE:    "NONE",
	SYMMIRBOTH: "MBoth",
	SYMRMBOTH:  "DBoth",
	SYMMIR8:    "Mirr8",
	SYMROT8:    "Rota8",
}

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
	TPush("@updatesymmetry")
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

	maxlen := 1
	updated := false
	if len(s.Map.SMap) > 0 {
		maxlen = len(s.Map.SMap[0])
	}

	if s.Fails > 200 {
		// 200 fails with no new sym - lets give up
		return
	}

	// TODO do in order of len of locs... so we dont build a map for 2 then
	// immediately rebuild for 4
	for minhash := range check {
		tile := s.Tiles[minhash]
		if tile.Ignore || len(tile.Locs) > 16 {
			// len less than 16 mostly just to avoid catastrophe 
			// any map with more sym than that is a shit show anyway
			tile.Ignore = true
			continue
		}
		//log.Print("symset, origin, offset, equiv:", symset, origin, offset, equiv)
		eqlen := s.SymAnalyze(tile)
		if tile.Ignore {
			s.Fails++
			continue
		}
		if eqlen > maxlen {
			smap, valid := s.SymMapValidate(tile)
			if valid {
				if Debug[DBG_Symmetry] {
					log.Printf("Valid symmetry map len %d found", len(smap[0]))
				}
				maxlen = eqlen
				s.Map.SMap = smap
				updated = true
				s.Fails = 0
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
	if dim.Size()/ndim.Size() > 20 {
		return dim
	}
	return ndim
}
// Update the analysis of a tile and return the length of the infered equiv set
func (s *SymData) SymAnalyze(tile *SymTile) (equivlen int) {
	tile.Gen = SYMNONE
	tile.MC = -1
	tile.MR = -1
	tile.RM1 = -1
	tile.RM2 = -1

	equivlen = 0

	if tile == nil || len(tile.Locs) < 2 {
		return
	}

	llist := tile.Locs

	// Get the blocking for the map
	dim := s.Tiling(tile)

	// test for Translation symmetry
	redlist := make([]Point, 0, len(llist))
	n := 0
	for i, l1 := range llist {
		p1 := dim.Donut(s.ToPoint(l1))
		for _, l2 := range llist[i+1:] {
			if s.Hashes[l1][0] == s.Hashes[l2][SYMTRANS] {
				p2 := dim.Donut(s.ToPoint(l2))
				n++
				if pd, good := dim.ShiftReduce(p1, p2, SYMMAXCELLS); good {
					redlist = SetAddPoint(redlist, pd)
				}
			}
		}
	}

	if len(redlist) > 0 {
		redlist = dim.ReduceReduce(redlist)
	}

	tile.Subtile = dim
	tile.Translate = s.Translation(redlist)
	tlen := dim.TranslationLen(tile.Translate)

	if Debug[DBG_Symmetry] {
		log.Print("Eq Len is ", s.Size()/dim.Size(), " * BLOCKS(", dim.TranslationLen(tile.Translate), ")")
	}
	equivlen = s.Size() / dim.Size() * tlen
	if equivlen > 1 {
		tile.Gen = SYMTRANS // a tiling is classed as a symtrans.
	}
	// If all we got was translations bail out.
	if n == len(llist)*(len(llist)-1)/2 || tlen > 1 {
		return
	}

	// Look for rotational symmetry, iff we have a square block
	rotorig := make([]Point, 0, 0)
	ndiag := 0
	if dim.Rows == dim.Cols && len(tile.Locs) <= equivlen*8 {
		for i, l1 := range llist {
			p1 := dim.Donut(s.ToPoint(l1))
			for _, l2 := range llist[i+1:] {
				if s.Hashes[l1][0] == s.Hashes[l2][SYMROT90] {
					p2 := dim.Donut(s.ToPoint(l2))
					porig := dim.Rot(p1, p2, SYMROT90)
					rotorig = dim.SymAddPoint(rotorig, porig)
				}
				if s.Hashes[l1][0] == s.Hashes[l2][SYMRM1] ||
					s.Hashes[l1][0] == s.Hashes[l2][SYMRM2] {
					ndiag++
				}
			}
		}
		// can only have 1 rotation origin.
		if len(rotorig) > 1 {
			// log.Print("Multiple rotation origins ", rotorig, len(tile.Locs))
			rotorig = rotorig[:0]
		} else if len(rotorig) == 1 {
			tile.Origin = rotorig[0]
			tile.Gen = SYMROT90
			equivlen *= 4
			if false && ndiag > 2 {
				equivlen *= 2
				tile.Gen = SYMROT8
			}
			return
		}
	}

	// Look for mirror symmetry iff we did not find rot sym
	morigc := make([]Point, 0, 0)
	morigr := make([]Point, 0, 0)
	ndiag = 0
	if len(rotorig) == 0 {
		for i, l1 := range llist {
			p1 := dim.Donut(s.ToPoint(l1))
			for _, l2 := range llist[i+1:] {
				if s.Hashes[l1][0] == s.Hashes[l2][SYMMIRRORC] {
					p2 := dim.Donut(s.ToPoint(l2))
					morigc = dim.SymAddPoint(morigc, dim.Mirror(p1, p2, 1))
				}
				if s.Hashes[l1][0] == s.Hashes[l2][SYMMIRRORR] {
					p2 := dim.Donut(s.ToPoint(l2))
					morigr = dim.SymAddPoint(morigr, dim.Mirror(p1, p2, 0))
				}
				if s.Hashes[l1][0] == s.Hashes[l2][SYMRM1] ||
					s.Hashes[l1][0] == s.Hashes[l2][SYMRM2] {
					ndiag++
				}
			}
		}
	}
	if len(morigc) == 1 {
		tile.MC = morigc[0].C
		equivlen *= 2
		tile.Gen = SYMMIRRORC
	}
	if len(morigr) == 1 {
		tile.MR = morigr[0].R
		equivlen *= 2
		if tile.Gen == SYMMIRRORC {
			tile.Gen = SYMMIRBOTH
		} else {
			tile.Gen = SYMMIRRORR
		}
	}

	if ndiag > 2 && dim.Rows == dim.Cols &&
		tile.Gen == SYMMIRRORR || tile.Gen == SYMMIRRORC || tile.Gen == SYMMIRBOTH {
		equivlen *= 2
		tile.Gen = SYMMIR8
	}

	if tile.Gen != SYMNONE && tile.Gen != SYMTRANS {
		return equivlen
	}

	// Look for diagonal symmetry iff we did not find rot/mirror sym
	rm1orig := make([]int, 0, 0)
	rm2orig := make([]int, 0, 0)
OUT:
	for i, l1 := range llist {
		p1 := dim.Donut(s.ToPoint(l1))
		for _, l2 := range llist[i+1:] {
			p2 := dim.Donut(s.ToPoint(l2))
			if s.Hashes[l1][0] == s.Hashes[l2][SYMRM1] {
				rm1orig = SetAddInt(rm1orig, dim.Diag(p1, p2, SYMRM1))
			}
			if s.Hashes[l1][0] == s.Hashes[l2][SYMRM2] {
				rm2orig = SetAddInt(rm2orig, dim.Diag(p1, p2, SYMRM2))
			}
			if len(rm2orig) > 2 || len(rm1orig) > 2 {
				break OUT
			}
		}
	}

	if len(rm1orig) > 1 {
		//log.Print("RM1 dups: ", rm1orig)
		rm1orig = rm1orig[:0]
	}
	if len(rm1orig) == 1 {
		tile.Gen = SYMRM1
		tile.RM1 = rm1orig[0]
		equivlen *= 2
	}

	if len(rm2orig) > 1 {
		//log.Print("RM2 dups: ", rm2orig)
		rm2orig = rm2orig[:0]
	}
	if len(rm2orig) == 1 {
		tile.RM2 = rm2orig[0]
		equivlen *= 2
		if tile.Gen == SYMRM1 {
			tile.Gen = SYMRMBOTH
		} else {
			tile.Gen = SYMRM2
		}
	}

	if tile.Gen == SYMRM1 ||
		tile.Gen == SYMRM2 ||
		tile.Gen == SYMRMBOTH {
		return
	}

	mrotorig := make([]Point, 0, 0)
	for i, l1 := range llist {
		p1 := dim.Donut(s.ToPoint(l1))
		for _, l2 := range llist[i+1:] {
			p2 := dim.Donut(s.ToPoint(l2))
			if s.Hashes[l1][0] == s.Hashes[l2][SYMROT180] {
				mrotorig = dim.SymAddPoint(mrotorig, dim.Rot(p1, p2, SYMROT180))
			}
		}
	}

	if len(mrotorig) > 1 {
		//log.Print("Multiple rot180 points ", mrotorig, len(tile.Locs))
	} else if len(mrotorig) == 1 {
		tile.Origin = mrotorig[0]
		tile.Gen = SYMROT180
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
			marr = s.Torus.Translations(s.Torus, Location(i), p, marr, SYMMAXCELLS)
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

func (tile *SymTile) Rot180(t Torus, loc Location, marr []Location) []Location {
	var pr, pc int
	st := tile.Subtile
	p := st.Donut(t.ToPoint(loc))
	if st.Rows%2 == 0 {
		pr = p.R - tile.Origin.R + 1
		pc = p.C - tile.Origin.C + 1
	} else {
		pr = tile.Origin.R - p.R
		pc = tile.Origin.C - p.C
	}

	marr = append(marr, t.ToLocation(p))
	marr = append(marr, t.ToLocation(st.Donut(Point{tile.Origin.R - pr, tile.Origin.C - pc})))

	// log.Print(tile.Origin, loc, st, marr[len(marr)-2:], p, 
	// Point{tile.Origin.R - pr, tile.Origin.C - pc},
	// st.Donut(Point{tile.Origin.R - pr, tile.Origin.C - pc}))

	return marr
}
func (tile *SymTile) Mirrors(t Torus, loc Location, marr []Location) []Location {
	st := &tile.Subtile
	if tile.Origin.R == 0 && tile.Origin.C == 0 {
		return marr
	}
	marr = append(marr, t.ToLocations(st.Mirrors(st.Donut(t.ToPoint(loc)), tile.MR, tile.MC))...)

	return marr
}
func (tile *SymTile) Rotations(t Torus, loc Location, marr []Location) []Location {
	var pr, prs, pc, pcs int
	st := tile.Subtile
	p := st.Donut(t.ToPoint(loc))
	pr = p.R - tile.Origin.R
	pc = p.C - tile.Origin.C

	if st.Rows%2 == 0 {
		prs = -pr - 1
		pcs = -pc - 1
	} else {
		prs = -pr
		pcs = -pc
		if pr == 0 && pc == 0 {
			// odd square has fixed center point
			// TODO origin choice should pick odd one so 
			// we can ignore dups below.  No odd sized maps for 
			// now though.
			marr = append(marr, t.ToLocation(p))
			return marr
		}
	}

	marr = append(marr, t.ToLocations([]Point{
		p,
		st.Donut(Point{tile.Origin.R + pc, tile.Origin.C + prs}),
		st.Donut(Point{tile.Origin.R + prs, tile.Origin.C + pcs}),
		st.Donut(Point{tile.Origin.R + pcs, tile.Origin.C + pr}),
	})...)

	return marr
}

func (tile *SymTile) Diagonals(t Torus, loc Location, marr []Location) []Location {
	st := tile.Subtile
	if tile.Translate.R != 0 || tile.Translate.C != 0 {
		return marr
	}

	p := make([]Point, 0, 4)
	p = append(p, st.Donut(t.ToPoint(loc)))
	if tile.RM1 != -1 {
		p = SetAddPoint(p, st.ReflectRM1(p[0], tile.RM1))
	}
	if tile.RM2 != -1 {
		for _, pp := range p {
			p = SetAddPoint(p, st.ReflectRM2(pp, tile.RM2))
		}
	}

	marr = append(marr, t.ToLocations(p)...)

	return marr
}

// Given an analyzed tile generate the map for loc, appending map to marr
// t is the Main map.
func (tile *SymTile) Generate(t Torus, loc Location, marr []Location) []Location {
	//mm := make([]Location, 0)
	//log.Print(tile.Rotations(t, 6894, mm))
	n := len(marr)
	// Steps are generate the subtile points
	switch tile.Gen {
	case SYMTRANS:
		marr = tile.Subtile.Translations(t, loc, tile.Translate, marr, SYMMAXCELLS)
	case SYMMIRRORC, SYMMIRRORR, SYMMIRBOTH, SYMMIR8:
		marr = tile.Mirrors(t, loc, marr)
	case SYMROT90, SYMROT8:
		marr = tile.Rotations(t, loc, marr)
	case SYMRM1, SYMRM2, SYMRMBOTH:
		marr = tile.Diagonals(t, loc, marr)
	case SYMROT180:
		marr = tile.Rot180(t, loc, marr)
	case SYMNONE:
		return marr
	default:
		log.Panic("TODO - invalid tile.Gen", tile.Gen)
	}

	if t.Cols != tile.Subtile.Cols || t.Rows != tile.Subtile.Rows {
		m := len(marr)
		for _, l := range marr[n:m] {
			p := t.ToPoint(l)
			for c := 0; c < t.Cols/tile.Subtile.Cols; c++ {
				for r := 0; r < t.Rows/tile.Subtile.Rows; r++ {
					if r != 0 || c != 0 {
						np := t.PointAdd(p, Point{r * tile.Subtile.Rows, c * tile.Subtile.Cols})
						marr = append(marr, t.ToLocation(np))
					}
				}
			}
		}
	}
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
			if len(marr) == 0 { // || len(marr) > size {
				if len(marr) > size {
					log.Print("Invalid map len(marr), size, points ", len(marr), size, marr[n:])
					log.Print(tile)
				}
				return nil, false
			}

			found := false
			for _, mloc := range marr[n:] {
				if int(mloc) == loc {
					found = true
				}
				if len(smap[mloc]) != 0 {
					log.Print("Already seen ", loc, mloc, tile)
				}
			}
			if !found {
				log.Print("loc not returned in marr ", tile.Gen, loc, marr[n:], "\n", tile)
			}

			item := UNKNOWN
			for _, mloc := range marr[n:] {
				// Validate the equiv set is identical
				if item == UNKNOWN {
					item = s.TGrid[mloc]
				} else if item != s.TGrid[mloc] {
					if Debug[DBG_Symmetry] {
						log.Print("Invalid point found: i, n, loc, mloc, item, tgrid ",
							i, n, loc, mloc, int(item), int(s.TGrid[mloc]), marr[n:])
					}
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
		if Debug[DBG_Symmetry] {
			log.Print("Tiling size mismatch ", len(marr), size)
		}
		return nil, false
	}

	return smap, true
}

func (tile *SymTile) String() string {
	s := ""
	s += fmt.Sprintf("Hash: %d Bits: %d Self: %d Origin: %v Gen: %v Translate %v Subtile: %v\n",
		tile.Hash, tile.Bits, tile.Self, tile.Origin, symLabels[tile.Gen], tile.Translate, tile.Subtile)
	s += fmt.Sprintf("RM1: %v RM2 %v MR %v MC %v",
		tile.RM1, tile.RM2, tile.MR, tile.MC)
	return s
}
