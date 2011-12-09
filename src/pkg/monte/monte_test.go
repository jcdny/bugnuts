package combat

import (
	"log"
	"testing"
	"bugnuts/replay"
	. "bugnuts/state"
	. "bugnuts/game"
	. "bugnuts/torus"
)

type testData struct {
	in, user string
	turn     int
	part     Point
}

var tests = []testData{
	{"testdata/test1/0.replay", "bot8", 30, Point{30, 53}},
}

func TestStatBot(t *testing.T) {
	for _, t := range tests {
		log.Print("Running test for ", t.in, "(", t.user, ") turn ", t.turn, " partition ", t.part)
		combatMe(t)
	}
}

func combatMe(t testData) {
	match, err := replay.Load(t.in)
	if err != nil {
		log.Print("Read error for ", t.in, ":", err)
		return
	}
	unum := -1
	for i, pn := range match.PlayerNames {
		if t.user == pn {
			unum = i
			break
		}
	}
	if unum < 0 {
		log.Print("User name ", t.user, " not found in ", t.in)
		return
	}

	m := match.GetMap()
	g := NewGame(&match.GameInfo, m)
	g.Replay(match.Replay, 0, match.GameLength, true)
	replay := g.PlayerInput

	s := NewState(&match.GameInfo)
	turns := make([]*Turn, 1, t.turn)

	var i int
	for i = 0; i < t.turn && len(replay[i][unum].A) > 0; i++ {
		turns = append(turns, replay[i][unum])
		s.ProcessTurn(replay[i][unum])
		s.UpdateStatistics(replay[i][unum])
	}
	if i != t.turn {
		log.Panic("turn not reached ", t.turn, " last ", i, " file ", t.in)
	}

	ap, pmap := CombatPartition(s)

	log.Print("%v", ap)
	log.Print("%v", pmap)
}
