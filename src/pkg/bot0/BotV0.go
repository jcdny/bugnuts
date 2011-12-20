// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

// The v0 Bot -- NOOP Bot; It does nothing, much like the goggles.
package bot0

import (
	"os"
	. "bugnuts/state"
)

type BotV0 struct {

}

//NewBot creates a new instance of your bot
func NewBotV0(s *State) Bot {
	mb := &BotV0{
	//do any necessary initialization here
	}

	return mb
}

func (bot *BotV0) Priority(i Item) int {
	return 1
}

func (bot *BotV0) DoTurn(s *State) os.Error {
	return nil
}
