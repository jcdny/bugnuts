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
type Direction byte

const (
	North      Direction = 'n'
	East                 = 'e'
	South                = 's'
	West                 = 'w'
	NoMovement           = '-'
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
	log.Panicf("%v is not a valid direction", d)
	return ""
}

func max(x []int) int {
	xm := -1 << 31
	for _, y := range x {
		if y > xm {
			xm = y
		}
	}

	return xm
}

func min(x []int) int {
	xm := int(^uint(0) >> 1)
	for _, y := range x {
		if y < xm {
			xm = y
		}
	}

	return xm
}

func (s *State) DumpSeen() {
	mseen := max(s.Map.Seen)
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
