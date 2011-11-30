package engine

import (
	"log"
	"testing"
	"bugnuts/replay"
)

func TestEngine(t *testing.T) {

	match, err := replay.Load("testdata/replay.0.json")
	if err != nil {
		log.Panicf("Error loading replay: %v", err)
	}

	m := match.GetMap()

	g := NewGame(&match.GameInfo, m)

	log.Printf("turn 0\n%vready\n", g.GameInfo)

	for i := 0; i < match.GameLength; i++ {
		//turns := GenerateTurn(match.Replay)
		//log.Printf("%v", ts[i][0])
		log.Printf("%d", i)
	}
}
