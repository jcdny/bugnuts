package main

import (
	"math"
)

type Point struct {
	r, c int
}

type Location int

// Item codes from parsing turns
type Item byte

const (
	UNKNOWN Item = iota
	WATER
	FOOD
	LAND
	MY_ANT
	PLAYER1
	PLAYER2
	PLAYER3
	PLAYER4
	PLAYER5
	PLAYER6
	PLAYER7
	PLAYER8
	PLAYER9
	PLAYER10
	MY_HILL
	HILL1
	HILL2
	HILL3
	HILL4
	HILL5
	HILL6
	HILL7
	HILL8
	HILL9
	HILL10
	MY_DEAD
	DEAD1
	DEAD2
	DEAD3
	DEAD4
	DEAD5
	DEAD6
	DEAD7
	DEAD8
	DEAD9
	DEAD10
	MAX_ITEM
	INVALID_ITEM Item = 255
)

var itemToSym = [256]byte{' ', '%', '*', '.',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K',
	'!', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

var symToItem [256]Item

func init() {
	// Set up static symbol item mappings.
	for j := 0; j < 256; j++ {
		symToItem[j] = INVALID_ITEM
	}
	for i := UNKNOWN; i < MAX_ITEM; i++ {
		symToItem[itemToSym[i]] = i
	}
	for j := 0; j < 10; j++ {
		symToItem['0'+j] = MY_HILL + Item(j)
	}
}

// Map an Item code to a character
func (o Item) ToSymbol() byte {
	return itemToSym[o]
}

// Map a character to an Item code
func ToItem(c byte) Item {
	return symToItem[c]
}

// Precompute circle points for lookup for a given r2 and number of map columns.
func GenCircleTable(r2 int) []Point {
	if r2 < 0 {
		return nil
	}

	d := int(math.Sqrt(float64(r2)))
	v := make([]Point, 0, (r2*22)/7+5)

	// Make the origin the first element so you can easily skip it.
	p := Point{r: 0, c: 0}
	v = append(v, p)

	for r := -d; r <= d; r++ {
		for c := -d; c <= d; c++ {
			if c != 0 || r != 0 {
				if c*c+r*r <= r2 {
					p = Point{r: int(r), c: int(c)}
					v = append(v, p)
				}
			}
		}
	}

	return v
}

// Given a []Point vector, compute the change from stepping north, south, east, west
// Useful for updating visibility, ranking move values.
func moveChangeCache(r2 int, v []Point) (add [][]Point, remove [][]Point) {
	// compute the size of the array we need to hold shifted circle
	d := int(math.Sqrt(float64(r2)))
	//TODO compute d from v rather than r2 so we can use different masks

	off := d + 1    // offset to get origin
	size := 2*d + 3 // one on either side + origin

	// Ordinal moves
	// TODO pass in
	sv := []Point{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
	for _, s := range sv {
		m := make([]int, size*size)

		av := []Point{}
		rv := []Point{}

		for _, p := range v {
			m[(p.c+off)+(p.r+off)*size]++
			m[(p.c+s.c+off)+(p.r+s.r+off)*size]--
		}

		for c := 0; c < size; c++ {
			for r := 0; r < size; r++ {
				switch {
				case m[c+r*size] > 0:
					rv = append(rv, Point{r: r - off, c: c - off})
				case m[c+r*size] < 0:
					av = append(av, Point{r: r - off, c: c - off})
				}
			}
		}
		add = append(add, av)
		remove = append(remove, rv)
	}

	return
}
