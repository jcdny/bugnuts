// MyBot is the original interface definition from the starter package.
package MyBot

import (
	"os"
	. "bugnuts/state"
)

// Bot interface
type Bot interface {
	DoTurn(*State) os.Error
}
