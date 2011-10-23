package main

import (
	"log"
	"strconv"
)

type Map struct {
	Rows int
	Cols int

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
