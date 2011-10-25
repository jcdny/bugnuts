package main

import (
	"log"
	"strconv"
	"os"
	"bufio"
	"strings"
)

type Map struct {
	Rows    int
	Cols    int
	Players int

	Grid []Item // Items seen
	Seen []int  // Turn on which cell was last visible.
}

//Direction represents the direction concept for issuing orders.
type Direction int8

const (
	North Direction = iota
	East
	South
	West
	NoMovement
)

func (d Direction) String() string {
	switch d {
	case North:
		return "n"
	case South:
		return "s"
	case West:
		return "w"
	case East:
		return "e"
	case NoMovement:
		return "-"
	}

	log.Printf("%v is not a valid direction", d)
	return "-"
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
	}

	return m
}

func (s *State) DumpSeen() {
	mseen := Max(s.Map.Seen)
	str := ""

	for r := 0; r < s.Rows; r++ {
		for c := 0; c < s.Cols; c++ {
			str += strconv.Itoa(s.Map.Seen[r*s.Cols+c] * 10 / (mseen + 1))
		}
		str += "\n"
	}

	log.Printf("Turn %d\n%v\n", s.Turn, str)
}

func (s *State) DumpMap() {
	m := make([]byte, len(s.Map.Grid))
	str := ""

	for i, o := range s.Map.Grid {
		m[i] = o.ToSymbol()
	}

	for r := 0; r < s.Rows; r++ {
		for c := 0; c < s.Cols; c++ {
			str += string(m[r*s.Cols+c])
		}
		str += "\n"
	}

	log.Printf("Turn %d\n%v\n", s.Turn, str)
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

			if len(words[1]) != cols {
				log.Panicf("Map line length mismatch line %d, got %d, expected %d", lines, len(words[1]), cols)
			}

			for _, c := range words[1] {
				m.Grid[loc] = ToItem(byte(c))
				loc++
			}
		}
	}

	return m, err
}
