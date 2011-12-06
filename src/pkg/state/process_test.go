package state

import (
	"testing"
	"bugnuts/replay"
	. "bugnuts/game"
	. "bugnuts/engine"
)

func getGame(file string) (*GameInfo, [][]*Turn) {
	match, _ := replay.Load(file)
	m := match.GetMap()
	g := NewGame(&match.GameInfo, m)
	g.Replay(match.Replay, 0, match.GameLength, true)
	replay := g.PlayerInput

	return &match.GameInfo, replay
}

func BenchmarkProcess(b *testing.B) {
	gi, replay := getGame("testdata/bench.replay.gz")

	unum := 0
	turns := make([]*Turn, 1, gi.Turns+2)
	for i := 0; i < len(replay) && len(replay[i][unum].A) > 0; i++ {
		turns = append(turns, replay[i][unum])
	}

	for i := 0; i < b.N; i++ {
		s := NewState(gi)
		for turn := 0; turn < len(replay) && len(replay[turn][unum].A) > 0; turn++ {
			s.ProcessTurn(replay[turn][unum])
			s.UpdateStatistics(replay[turn][unum])
		}
	}
}
