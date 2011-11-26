package engine

import (
	"log"
	"bugnuts/maps"
	"bugnuts/state"
	"bugnuts/replay"
)

type Player struct {
	Hills   []maps.Location
	Seen    []bool
	Visible []bool
	IdMap   []maps.Item
}

type Game struct {
	*maps.Map
	*state.GameInfo
	Turn    int
	Players []*Player
}

func NewGame(gi *state.GameInfo, m *maps.Map) *Game {
	g := &Game{
		Map:      m,
		GameInfo: gi,
		Players:  make([]*Player, m.Players),
	}
	// Set Hill Locations and initialize the player map, each player is it's own 0th player.
	for i := range g.Players {
		log.Printf("Hills: %v", m.Hills(i))
		g.Players[i] = &Player{
			Hills:   m.Hills(i),
			IdMap:   make([]maps.Item, m.Players),
			Seen:    make([]bool, m.Size()),
			Visible: make([]bool, m.Size()),
		}
		g.Players[i].IdMap[i] = maps.MY_ANT
	}

	return g
}

func GenerateTurn(g *Game, r *replay.Replay) []*state.Turn {
	al := r.GetAnts(g.Turn)
	// per player update seen and visible masks
	3for i, p := range g.Players {

	}
}
