package main

import ()

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
