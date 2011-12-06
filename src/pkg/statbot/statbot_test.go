package statbot

import (
	"log"
	"testing"
	"strings"
	"os"
	"io/ioutil"
	"bugnuts/replay"
	. "bugnuts/state"
	. "bugnuts/game"
	. "bugnuts/engine"
	. "bugnuts/MyBot"
)

func statMe(in, out, user string) {
	match, err := replay.Load(in)
	if err != nil {
		log.Print("Read error for ", in, ":", err)
		return
	}
	unum := -1
	for i, pn := range match.PlayerNames {
		if user == pn {
			unum = i
			break
		}
	}
	if unum < 0 {
		log.Print("User name ", user, " not found in ", in)
		return
	}

	fd, err := os.Create(out)
	if err != nil {
		log.Print("Create error for ", out, ":", err)
		return
	}
	defer fd.Close()

	m := match.GetMap()
	g := NewGame(&match.GameInfo, m)
	g.Replay(match.Replay, 0, match.GameLength, true)
	replay := g.PlayerInput

	s := NewState(&match.GameInfo)
	bot := NewBot("sb", s)

	sb := bot.(*StatBot)
	sb.SetTrueState(g, unum)

	if bot == nil {
		log.Print("Unkown bot SB")
		log.Printf("Bots:\n%s\n", strings.Join(BotList(), "\n"))
		return
	}
	turns := make([]*Turn, 1, s.Turns+2)

	for i := 0; i < len(replay) && len(replay[i][unum].A) > 0; i++ {
		turns = append(turns, replay[i][unum])
		s.ProcessTurn(replay[i][unum])
		s.UpdateStatistics(replay[i][unum])
		sb.Report(s, fd)
	}
}

func TestStatBot(t *testing.T) {
	replays, _ := ioutil.ReadFile("testdata/bugnuts.tmp")
	for _, file := range strings.Split(string(replays), "\n")[:10] {
		in := "testdata/" + file
		t := strings.Split(file, "/")
		out := "testdata/tmp/" + t[len(t)-1] + ".csv"
		user := "bugnutsv5"
		log.Print("Read ", in, " write ", out, " for ", user)
		statMe(in, out, user)
	}
}
