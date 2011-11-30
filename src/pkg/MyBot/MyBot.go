package MyBot

import (
	"os"
	. "bugnuts/state"
)

// Bot interface
type Bot interface {
	DoTurn(*State) os.Error
}
