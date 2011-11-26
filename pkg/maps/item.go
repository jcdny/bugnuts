package maps

import ()

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
	PLAYERGUESS
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
	HILLGUESS // Not a real hill - our guess
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
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', '~', // Player, ~ is a player guess
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '?', // Hill ? is guess hill
	'!', 'z', 'y', 'x', 'w', 'v', 'u', 't', 's', 'r', // Dead ant
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', // Ant on hill
	'X', '+', '@', '=', '|', '&', // Operational things
}

var symToItem [256]Item

var TerminalItem [256]bool
var StepableItem [256]bool

func init() {
	// Set up static symbol item mappings.
	for j := 0; j < 256; j++ {
		symToItem[j] = INVALID_ITEM
	}
	for i := UNKNOWN; i < MAX_ITEM; i++ {
		symToItem[itemToSym[i]] = i
	}

	// Make Terminal
	for _, item := range []Item{EXPLORE, DEFEND, RALLY} {
		TerminalItem[item] = true
	}
	for item := Item(HILL1); item <= HILLGUESS; item++ {
		TerminalItem[item] = true
	}
	for item := Item(HILLANT1); item <= HILLANT9; item++ {
		TerminalItem[item] = true
	}

	// StepableItem
	for i := range StepableItem {
		StepableItem[i] = true
	}
	for _, item := range []Item{WATER, BLOCK, FOOD} {
		StepableItem[item] = false
	}

}

func (o Item) String() string {
	return string(o.ToSymbol())
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
	if (o >= MY_HILL && o <= HILLGUESS) || (o >= MY_HILLANT && o <= HILLANT9) {
		return true
	}

	return false
}

func (o Item) IsEnemyHill(player int) bool {
	if o == LAND || o == WATER {
		return false
	}
	if o >= MY_HILL && o <= HILLGUESS {
		return player != int(o-MY_HILL)
	}
	if o >= MY_HILLANT && o <= HILLANT9 {
		return player != int(o-MY_HILLANT)
	}

	return false
}

func (o Item) IsEnemyAnt(player int) bool {
	if o == LAND || o == WATER {
		return false
	}
	if o >= MY_ANT && o <= PLAYERGUESS {
		return player != int(o-MY_ANT)
	}
	if o >= MY_HILLANT && o <= HILLANT9 {
		return player != int(o-MY_HILLANT)
	}

	return false
}
