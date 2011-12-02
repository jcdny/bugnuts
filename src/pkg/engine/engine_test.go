package engine

import (
	"log"
	"testing"
	"bugnuts/replay"
)

func TestEngine(t *testing.T) {

	match, err := replay.Load("testdata/replay.1.json")

	if err != nil || match == nil {
		log.Panicf("Error loading replay: %v", err)
	}
	m := match.GetMap()

	g := NewGame(&match.GameInfo, m)

	tout := g.Replay(match.Replay, 0, match.GameLength)

	log.Printf("turn 0\n%v\nready\n", g.GameInfo)
	for i := range tout {
		log.Print("\nturn ", i+1, "\n", tout[i][0], "\n")
	}
}
