package main

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

const (
	Offset          = map[string]Point{"n": {-1, 0}, "s": {1, 0}, "e": {0, -1}, "w": {0, 1}, "-": {0, 0}}
	DirectionString = [5]string{"n", "e", "s", "w", "-"}
)
