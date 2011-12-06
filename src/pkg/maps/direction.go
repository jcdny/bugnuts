package maps

import (
	"log"
	. "bugnuts/torus"
)

//Direction represents the direction concept for issuing orders.
type Direction uint8

const (
	North Direction = iota
	East
	South
	West
	NoMovement
	InvalidMove
)

var ByteToDirection [256]Direction

var DirectionMap = map[string]Point{"n": {-1, 0}, "e": {0, -1}, "s": {1, 0}, "w": {0, 1}, "-": {0, 0}}
var DirectionOffset = [5]Point{North: {-1, 0}, East: {0, 1}, South: {1, 0}, West: {0, -1}, NoMovement: {0, 0}}
var Steps = [4]Point{North: {-1, 0}, East: {0, 1}, South: {1, 0}, West: {0, -1}} // exclude no move
var Steps2 = []Point{{-1, 0}, {0, 1}, {1, 0}, {0, -1},                           // one step moves
	{-2, 0}, {-1, 1}, {0, 2}, {1, 1}, {2, 0}, {1, -1}, {0, -2}, {-1, -1}, // two step moves
}
var Steps3 = []Point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}, // one step moves
	{-2, 0}, {-1, 1}, {0, 2}, {1, 1}, {2, 0}, {1, -1}, {0, -2}, {-1, -1}, // two step moves
	{-3, 0}, {-2, 1}, {-1, 2}, {0, 3}, {1, 2}, {2, 1}, // 3 steps
	{3, 0}, {2, -1}, {1, -2}, {0, -3}, {-1, -2}, {-2, -1},
}

var DirectionChar = [5]string{"n", "e", "s", "w", "-"}

func init() {
	for i := range ByteToDirection {
		ByteToDirection[i] = InvalidMove
	}
	for i, char := range "nesw-" {
		ByteToDirection[char] = Direction(i)
	}
	for i, char := range "NESW" {
		ByteToDirection[char] = Direction(i)
	}
}

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

	log.Printf("%v : invalid direction", d)
	return "-"
}

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
