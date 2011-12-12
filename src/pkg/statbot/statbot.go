// StatBot - Dummy bot that simply outputs the statistics generated out of
// state.UpdateStatistics().
package statbot

import (
	"os"
	"fmt"
	"bugnuts/engine"
	. "bugnuts/game"
	. "bugnuts/state"
	. "bugnuts/parameters"
	. "bugnuts/MyBot"
)

type StatBot struct {
	P    *Parameters
	G    *engine.Game
	PMax int
	NP   int
}

func init() {
	RegisterABot(ABot{Key: "sb", Desc: "Statbot - noop bot that collects statistics", PKey: "sb", NewBot: NewStatBot})
}

//NewBot creates a new instance of your bot
func NewStatBot(s *State, pset *Parameters) Bot {
	mb := &StatBot{
		P: pset,
	}
	//mb.PriMap = mb.P.MakePriMap()

	return mb
}
func (bot *StatBot) SetTrueState(g *engine.Game, np int) {
	bot.G = g
	bot.NP = np
	bot.PMax = len(bot.G.Players)

	// have to add an inverse map for players we never saw during the actual game
	for i, invp := range bot.G.Players[np].InvMap {
		if invp < 0 {
			bot.G.Players[np].AddIdMap(i)
		}
	}

}

func (bot *StatBot) StatHeader() string {
	s := "Turn,"
	if bot.G != nil {
		for i := range bot.G.Players {
			s += fmt.Sprint("Ntrue", i, ",")
		}
	} else {
		bot.NP = 0
		bot.PMax = MaxPlayers
	}

	s += "Unknown,Horizon,HorizonMax,HorizonMaxTurn,DiedTotAll,Food,"

	s += "PSeen,"

	for i := 0; i < bot.PMax; i++ {
		s += fmt.Sprint("Nseen", i, ",")
	}

	return s
}

func (bot *StatBot) StatLine(turn int, s *Statistics) string {
	ts := &s.TStats[turn]

	out := ""
	out = fmt.Sprint(turn, ",")
	if bot.G != nil {
		for i := range bot.G.Players {
			nant := 0
			pnum := bot.G.Players[bot.NP].InvMap[i]
			for _, pl := range bot.G.PlayerInput[turn-1][pnum].A {
				if pl.Player == 0 {
					nant++
				}
			}
			out += fmt.Sprint(nant) + ","
		}
	}

	out += fmt.Sprint(ts.Unknown, ",",
		ts.Horizon, ",",
		s.HorizonMax, ",",
		s.HorizonMaxTurn, ",",
		s.DiedTotAll, ",",
		ts.Food, ",")
	out += fmt.Sprint(-1, ",")
	for i := 0; i < bot.PMax; i++ {
		out += fmt.Sprint(ts.Seen[i], ",")
	}
	return out
}

func (bot *StatBot) Report(s *State, fd *os.File) {
	if s.Turn == 1 {
		fmt.Fprint(fd, bot.StatHeader(), "\n")
	}
	fmt.Fprint(fd, bot.StatLine(s.Turn, s.Stats), "\n")
}

func (bot *StatBot) DoTurn(s *State) os.Error {
	bot.Report(s, os.Stdout)

	return nil
}
