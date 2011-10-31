package main
// The v0 - The NOOP Bot.
//   Useful for testing forcast data from replays 


import (
	"os"
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

func (bot *BotV0) DoTurn(s *State) os.Error {

	return nil
}
