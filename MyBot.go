package main

import (
	"os"
)

//Bot interface defines what we need from a bot
type Bot interface {
	DoTurn(s *State) os.Error
}
