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
	HILL_GUESS // Not a real hill - our guess
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
	MY_HILLANT
	HILLANT1
	HILLANT2
	HILLANT3
	HILLANT4
	HILLANT5
	HILLANT6
	HILLANT7
	HILLANT8
	HILLANT9
	EXPLORE  // An explore goal - terminal
	DEFEND   // A defense spot - terminal
	RALLY    // rally point for future attack - terminal
	WAYPOINT // a place to go on the way somewhere - terminal
	BLOCK    // A moved ant or something else preventing stepping in
	OCCUPIED // An ant has moved here so it can't be moved into
	MAX_ITEM
	INVALID_ITEM Item = 255
)

var itemToSym = [256]byte{' ', '%', '*', '.',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '?', // ? is guess hill
	'!', 'z', 'y', 'x', 'w', 'v', 'u', 't', 's', 'r',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'X', '+', '@', '=', '|', '&',
}

var symToItem [256]Item

var Terminal = map[Item]bool{EXPLORE: true, DEFEND: true, WAYPOINT: false, FOOD: false, RALLY: true}

func init() {
	// Set up static symbol item mappings.
	for j := 0; j < 256; j++ {
		symToItem[j] = INVALID_ITEM
	}
	for i := UNKNOWN; i < MAX_ITEM; i++ {
		symToItem[itemToSym[i]] = i
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

func (o Item) IsHill() bool {
	if (o >= MY_HILL && o <= HILL_GUESS) || (o >= MY_HILLANT && o <= HILLANT9) {
		return true
	}

	return false
}
func (o Item) IsEnemyHill() bool {
	if (o > MY_HILL && o <= HILL_GUESS) || (o > MY_HILLANT && o <= HILLANT9) {
		return true
	}

	return false
}
func (o Item) IsTerminal() bool {
	termp, ok := Terminal[o]
	return termp || (!ok && o.IsEnemyHill())
}

func (o Item) String() string {
	return string(o.ToSymbol())
}
