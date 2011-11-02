package main

import (
	"log"
	"strconv"
	"os"
	"bufio"
	"strings"
	"fmt"
)

type Map struct {
	Rows    int
	Cols    int
	Players int

	Grid   []Item   // Items seen
	Seen   []int    // Turn on which cell was last visible.
	Threat [][]int8 // how much threat is there on a given cell

	BDist []uint8 // cache of border distance
}

type Point struct {
	r, c int
}

func (p *Point) String() string {
	return fmt.Sprintf("r:%d s:%d", p.r, p.c)
}

func NewMap(rows, cols, players int) *Map {
	if rows < 1 || cols < 1 {
		log.Panicf("Invalid map size %d %d", rows, cols)
	}

	m := &Map{
		Rows:    rows,
		Cols:    cols,
		Players: players,
		Grid:    make([]Item, rows*cols),
		Seen:    make([]int, rows*cols),
		BDist:   BorderDistance(rows, cols),
	}

	return m
}

func (m *Map) Size() int {
	return m.Rows * m.Cols
}

func (m *Map) ToLocation(p Point) Location {
	p = m.Donut(p)
	return Location(p.r*m.Cols + p.c)
}

func (m *Map) ToPoint(l Location) (p Point) {
	p = Point{r: int(l) / m.Cols, c: int(l) % m.Cols}

	return
}

func (m *Map) Donut(p Point) Point {
	if p.r < 0 {
		p.r += m.Rows
	}
	if p.r >= m.Rows {
		p.r -= m.Rows
	}
	if p.c < 0 {
		p.c += m.Cols
	}
	if p.c >= m.Cols {
		p.c -= m.Cols
	}

	return p
}

func (m *Map) PointEqual(p1, p2 Point) bool {
	// todo donuts
	return p1.c == p2.c && p1.r == p2.r
}

func (m *Map) PointAdd(p1, p2 Point) Point {
	return m.Donut(Point{r: p1.r + p2.r, c: p1.c + p2.c})
}

func (m *Map) String() string {
	s := ""
	s += "rows " + strconv.Itoa(m.Rows) + "\n"
	s += "cols " + strconv.Itoa(m.Rows) + "\n"
	s += "players " + strconv.Itoa(m.Players) + "\n"
	for r := 0; r < m.Rows; r++ {
		s += "m "
		for _, item := range m.Grid[r*m.Cols : (r+1)*m.Cols] {
			s += string(item.ToSymbol())
		}
		s += "\n"
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

func (m *Map) HillLocations() []Location {
	// find a hill for start
	hills := []Location{}

	for i, item := range m.Grid {
		if item.IsHill() {
			hills = append(hills, Location(i))
		}
	}

	return hills
}

func MapLoadFile(file string) (*Map, os.Error) {
	var m *Map = nil

	f, err := os.Open(file)

	if err != nil {
		return nil, err
	} else {
		defer f.Close()

		in := bufio.NewReader(f)
		m, err = MapLoad(in)
	}

	return m, err
}

func MapValidate(ref *Map, gen *Map) (int, string) {
	out := ""
	count := 0

	if gen.Rows != ref.Rows || gen.Cols != ref.Cols {
		out += fmt.Sprintf("Map size mismatch: refence map r:%d c:%d and generated map r:%d c:%d\n",
			ref.Rows, ref.Cols, gen.Rows, gen.Cols)
		count++
	} else {
		for i, item := range gen.Grid {
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

//Take a slice of Point and return a slice of Location
//Used for offsets so it does not donut things.
func (m *Map) ToPoints(lv []Location) []Point {
	pv := make([]Point, len(lv))
	for i, l := range lv {
		pv[i] = m.ToPoint(l)
	}

	return pv
}


// Ruturn a uint8 array with distance to border in each cell
func BorderDistance(rows, cols int) (out []uint8) {
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

