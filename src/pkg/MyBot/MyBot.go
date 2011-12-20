// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package MyBot

import (
	"os"
	. "bugnuts/state"
)

// Bot interface
type Bot interface {
	DoTurn(*State) os.Error
}
