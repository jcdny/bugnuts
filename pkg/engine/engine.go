package engine

import (
	"log"
	"bugnuts/maps"
	"bugnuts/replay"
)

type Player struct {
	Hills   []Location
	Seen    []bool
	Visible []bool
}

type Game struct {
	*maps.Map
	Players []Player
}
