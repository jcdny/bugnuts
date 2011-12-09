package engine

import (
	"log"
	"bugnuts/replay"
	"bugnuts/state"
	"bugnuts/game"
)

func LoadState(file, player string, turn int) *state.State {
	match, err := replay.Load(file)
	if err != nil {
		log.Print("Read error for ", file, ":", err)
		return nil
	}
	unum := -1
	for i, pn := range match.PlayerNames {
		if player == pn {
			unum = i
			break
		}
	}
	if unum < 0 {
		log.Print("User name ", player, " not found in ", file)
		return nil
	}

	m := match.GetMap()
	g := NewGame(&match.GameInfo, m)
	g.Replay(match.Replay, 0, match.GameLength, true)
	replay := g.PlayerInput

	s := state.NewState(&match.GameInfo)
	turns := make([]*game.Turn, 1, turn)

	var i int
	for i = 0; i < turn && len(replay[i][unum].A) > 0; i++ {
		turns = append(turns, replay[i][unum])
		s.ProcessTurn(replay[i][unum])
		s.UpdateStatistics(replay[i][unum])
	}
	if i != turn {
		log.Print("turn ", turn, " not reached for player ", player, " last was ", i, " file ", file)
		return nil
	}

	return s
}
