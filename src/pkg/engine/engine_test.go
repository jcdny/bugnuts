package engine

import (
	"log"
	"testing"
	"bugnuts/replay"
	"fmt"
	"os"
)

func TestEngine(t *testing.T) {

	match, err := replay.Load("testdata/0.replay")

	if err != nil || match == nil {
		log.Panicf("Error loading replay: %v", err)
	}
	m := match.GetMap()

	g := NewGame(&match.GameInfo, m)

	tout := g.Replay(match.Replay, 0, match.GameLength)

	out, err := os.Create("testdata/0.input.tmp")
	if err != nil {
		log.Panic("open failed for testdata/0.input.tmp:", err)
	}
	defer out.Close()

	fmt.Fprintf(out, "turn 0\n%v\nready\n", g.GameInfo)
	for i := range tout {
		fmt.Fprint(out, "turn ", i+1, "\n", tout[i][0], "\ngo\n")
	}
}
