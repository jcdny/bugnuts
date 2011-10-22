package main

import (
	"os"
)

type MyBot struct {

}

//Bot interface defines what we need from a bot
type Bot interface {
	DoTurn(s *State) os.Error
}

//NewBot creates a new instance of your bot
func NewBot(s *State) Bot {
	mb := &MyBot{
	//do any necessary initialization here
	}

	return mb
}

//DoTurn is where you should do your bot's actual work.
func (mb *MyBot) DoTurn(s *State) os.Error {
	//returning an error will halt the whole program!
	return nil
}
