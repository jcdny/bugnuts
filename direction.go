package main

import (
	"log"
)

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

var DirectionMap = map[string]Point{"n": {-1, 0}, "e": {0, -1}, "s": {1, 0}, "w": {0, 1}, "-": {0, 0}}
var DirectionOffset = [5]Point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}, {0, 0}}
var Steps = [4]Point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} // exclude no move
var DirectionChar = [5]string{"n", "e", "s", "w", "-"}

type Path struct {
	steps []Location

	// A path can take reference another path
	refpath *Path
	offset  int
}

func (m *Map) PointsToPath(points []Point) *Path {
	path := Path{
		steps: make([]Location, len(points), len(points)),
	}

	for i, p := range points {
		// check pts describe a real path
		path.steps[i] = m.ToLocation(p)
	}

	return &path
}
