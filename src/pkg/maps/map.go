// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package maps

import (
	"log"
	"strconv"
	"os"
	"bufio"
	"strings"
	"fmt"
	"crypto/md5"
	"io"
	. "bugnuts/util"
	. "bugnuts/torus"
)

type Map struct {
	Torus                 // Defines the geometry of the map
	Players  int          // This is stored in the map file
	Grid     []Item       // The map data possibly run through the symmetry mapper
	TGrid    []Item       // The true map data - updated via ProcessTurn
	*SymData              // Symmetry data
	SID      int          // The integer id for symmetry so we can invalidate hill guesses etc.
	SMap     [][]Location // Symmetry map

	// internal cache data
	BorderDist []uint8       // border distance
	LocStep    [][5]Location // adjecent tile pointer
}

const MAXMAPSIZE = 65536 // this dictated by using uint16 etc for depth

func NewMap(rows, cols, players int) *Map {
	if rows < 1 || cols < 1 {
		log.Panicf("Invalid map size %d %d", rows, cols)
	}

	m := &Map{
		Torus: Torus{
			Rows: rows,
			Cols: cols,
		},
		Players: players,
		Grid:    make([]Item, rows*cols),
		TGrid:   make([]Item, rows*cols),
		// cache data
		BorderDist: borderDistance(rows, cols),
		LocStep:    locationStep(rows, cols),
		SMap:       make([][]Location, 0, rows*cols),
	}
	m.SymData = m.NewSymData(4)

	return m
}

// MapId computes an md5 hash of the string representation of the map data
func (m *Map) MapId() string {
	c := md5.New()
	s := m.String()
	io.WriteString(c, s)
	out := fmt.Sprintf("%x", c.Sum())

	return out
}

func (m *Map) String() string {
	s := "rows " + strconv.Itoa(m.Rows) + "\n"
	s += "cols " + strconv.Itoa(m.Rows) + "\n"
	s += "players " + strconv.Itoa(m.Players)
	for r := 0; r < m.Rows; r++ {
		s += "\nm "
		for _, item := range m.Grid[r*m.Cols : (r+1)*m.Cols] {
			s += string(item.ToSymbol())
		}
	}

	return s
}

func MapLoad(in *bufio.Reader) (*Map, os.Error) {
	var m *Map = nil
	var err os.Error

	lines := 0
	loc := 0
	nrow := 0
	rows := -1
	cols := -1
	players := -1

	for {
		var line string

		line, err = in.ReadString('\n')
		lines++
		if err != nil {
			break
		}

		line = line[:len(line)-1] //remove the delimiter
		if line == "" {
			continue
		}

		words := strings.SplitN(line, " ", 2)
		if len(words) != 2 {
			log.Printf("Invaid param line \"%s\"", line)
			continue
		}

		switch words[0] {
		case "rows":
			rows, _ = strconv.Atoi(words[1])
		case "cols":
			cols, _ = strconv.Atoi(words[1])
		case "players":
			players, _ = strconv.Atoi(words[1])
		case "m":
			if m == nil {
				m = NewMap(rows, cols, players)
			}

			if nrow > rows {
				log.Panicf("Map rows mismatch row %d expected %d", nrow, rows)
			}

			line = line[2:] // remove "m "
			if len(line) != cols {
				log.Panicf("Map line length mismatch line %d, got %d, expected %d", lines, len(words[1]), cols)
			}

			for _, c := range line {
				m.Grid[loc] = ToItem(byte(c))
				loc++
			}
		}
	}

	return m, err
}

func MapLoadFile(file string) (m *Map, err os.Error) {
	f, err := os.Open(file)
	if err != nil {
		return
	} else {
		defer f.Close()

		in := bufio.NewReader(f)
		m, err = MapLoad(in)
		if err == os.EOF {
			err = nil
		}
	}

	return
}

func MapValidate(ref *Map, gen *Map) (int, string) {
	out := ""
	count := 0

	if gen.Rows != ref.Rows || gen.Cols != ref.Cols {
		out += fmt.Sprintf("Map size mismatch: refence map r:%d C:%d and generated map r:%d C:%d\n",
			ref.Rows, ref.Cols, gen.Rows, gen.Cols)
		count++
	} else {
		for i, item := range gen.TGrid {
			if item != UNKNOWN && item != ref.Grid[i] &&
				(item == WATER || ref.Grid[i] == WATER ||
					item.IsHill() != gen.Grid[i].IsHill()) {
				out += fmt.Sprintf("%v ref %s gen %s\n", gen.ToPoint(Location(i)), ref.Grid[i], item)
				count++
			}
		}
	}

	return count, out
}

// Ruturn a uint8 array with distance to border in each cell
func borderDistance(rows, cols int) (out []uint8) {
	if rows > 255 || cols > 255 {
		log.Panic("Rows or cols > 255 in BorderDist")
	}
	out = make([]uint8, rows*cols, rows*cols)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			out[r*cols+c] = uint8(MinV(r+1, c+1, Abs(r-rows), Abs(c-cols)))
		}
	}

	return
}

// Generate the cache of one step moves from current cell
func locationStep(rows, cols int) (out [][5]Location) {
	out = make([][5]Location, rows*cols)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			loc := r*cols + c
			for i, step := range DirectionOffset {
				rstep := r + step.R
				cstep := c + step.C
				// Wrap if we need to
				if rstep < 0 {
					rstep += rows
				}
				if rstep >= rows {
					rstep -= rows
				}
				if cstep < 0 {
					cstep += cols
				}
				if cstep >= cols {
					cstep -= cols
				}
				nloc := rstep*cols + cstep
				out[loc][i] = Location(nloc)
			}
		}
	}

	return
}

// Return list of Hill Locations, player -1 is all players
func (m *Map) Hills(player int) []Location {
	// find a hill for start
	hills := make([]Location, 0)

	for i, item := range m.Grid {
		if item.IsHill() && (player < 0 || item == MY_HILL+Item(player)) {
			hills = append(hills, Location(i))
		}
	}

	return hills
}

func (m *Map) TApply() {
	// Takes the current TGrid and SMap and updates Grid
	for i := 0; i < m.Size(); i++ {
		m.Grid[i] = UNKNOWN
	}
	for i := 0; i < m.Size(); i++ {
		if m.TGrid[i] != UNKNOWN {
			for _, loc := range m.SMap[i] {
				if m.Grid[loc] != UNKNOWN {
					break
				}
				m.Grid[loc] = m.TGrid[i]
			}
		}
	}
}

// Return false if the symmetry mapping gives a result inconsistent with TGrid
func (m *Map) TSet(i Item, locs ...Location) bool {
	valid := len(m.SMap) > 0
	for _, loc := range locs {
		if m.TGrid[loc] != i {
			m.TGrid[loc] = i
			if valid && m.Grid[loc] != UNKNOWN && m.Grid[loc] != i {
				valid = false
			}
			if !valid {
				// with no sym map or and invalid one
				// just copy to Grid
				m.Grid[loc] = i
			} else {
				for _, lsym := range m.SMap[loc] {
					m.Grid[lsym] = i
				}
			}
		}

	}
	// ending smap is valid or we never had one
	return valid || len(m.SMap) == 0
}

func (m *Map) DumpMap() string {
	mout := make([]byte, len(m.Grid))
	for i, o := range m.Grid {
		mout[i] = o.ToSymbol()
	}

	str := ""
	for r := 0; r < m.Rows; r++ {
		str += string(mout[r*m.Cols : (r+1)*m.Cols-1])
		str += "\n"
	}
	return str
}
