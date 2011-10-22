package main

import (
	"log"
)

type Map struct {
	Rows int
	Cols int

	Grid    []Item
	Visible []bool
	Seen    []Seen
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

type Seen byte

const (
	UNSEEN Seen = iota
	SEEN
	VISITED
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
