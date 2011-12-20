// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

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
	. "bugnuts/watcher"
)

func init() {
	SetWatcherPrefix("sb")
	log.SetFlags(log.Lshortfile)
	Debug[DBG_GatherTime] = true
	Debug[DBG_TurnTime] = true
}

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

	if bot == nil {
		log.Print("error creating bot", bot)
		return
	}
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
	var replaylist []string
	var user string
	if false {
		user = "a1k0n"
		replays, err := ioutil.ReadFile("testdata/" + user + ".tmp")
		if err != nil {
			t.Error(err)
		}
		replaylist = strings.Split(string(replays), "\n")
	} else {
		replaylist = []string{"fluxid.38928.replay"}
		user = "bugnutsv6"
	}

	for _, file := range replaylist {
		if file == "" {
			continue
		}
		in := "testdata/" + file
		t := strings.Split(file, "/")
		out := "testdata/" + t[len(t)-1] + ".csv"
		log.Print("Read ", in, " write ", out, " for ", user)
		statMe(in, out, user)
	}
}
