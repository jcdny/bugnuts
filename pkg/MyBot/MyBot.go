package MyBot

import (
	"os"
	. "bugnuts/state"
	. "bugnuts/maps"
)

// Bot interface
type Bot interface {
	DoTurn(*State) os.Error
	Priority(Item)
}
