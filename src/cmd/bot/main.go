package main

import (
	"flag"
	"log"
	"bufio"
	"os"
	"fmt"
	"strings"
	. "bugnuts/maps"
	. "bugnuts/viz"
	. "bugnuts/watcher"
	. "bugnuts/state"
	. "bugnuts/game"
	. "bugnuts/MyBot"
	// import bots to get their init() to register them.
	_ "bugnuts/statbot"
	_ "bugnuts/bot3"
	_ "bugnuts/bot5"
	_ "bugnuts/bot6"
	_ "bugnuts/bot7"
	_ "bugnuts/bot8"
)

var runBot string
var mapName string
var watchPoints string
var debugLevel int
var maxTurn int

func init() {
	//	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.SetFlags(log.Lshortfile)

	vizList := ""
	vizHelp := "Visualize: all,none,useful"
	for flag := range Viz {
		vizHelp += "," + flag
	}

	flag.StringVar(&runBot, "b", "v8", "Which bot to run\n\t"+strings.Join(BotList(), "\n\t"))
	flag.StringVar(&vizList, "V", "", vizHelp)
	flag.IntVar(&debugLevel, "d", 0, "Debug level")
	flag.StringVar(&mapName, "m", "", "Map file -- Used to validate generated map, hill guessing etc.")
	flag.StringVar(&watchPoints, "w", "", "Watch points \"T1:T2@R,C,N[;T1:T2...]\", \":\" will watch everything")
	flag.IntVar(&maxTurn, "T", 65535, "Max Turn")
	flag.Parse()

	if BotGet(runBot) == nil {
		log.Printf("Unrecognized bot \"%s\", Registered bots:\n\t%s\n", runBot, strings.Join(BotList(), "\n\t"))
		return
	}

	SetWatcherPrefix(runBot)

	SetDebugLevel(debugLevel)
	SetViz(vizList, Viz)
}

func main() {
	//TurnTimer()
	wd, _ := os.Getwd()
	log.Print("Running bot in ", wd)

	var refmap *Map
	if mapName != "" {
		refmap, _ = MapLoadFile(MapFile(mapName))
	}

	in := bufio.NewReader(os.Stdin)

	// Load game definition
	g, err := GameScan(in)
	if err != nil {
		log.Panicf("Start(in) failed (%s)", err)
	} else if Debug[DBG_Start] {
		log.Printf("Game Info:\n%v\n", g)
	}
	TurnSet(0)
	// Create watch points
	WS = NewWatches(g.Rows, g.Cols, g.Turns)
	if len(watchPoints) > 0 {
		wlist := strings.Split(watchPoints, ";")
		WS.Load(wlist)
	}

	TPush("NewState")
	s := NewState(g)
	TPop()

	turns := make([]*Turn, 1, s.Turns+2)

	TPush("NewBot")
	bot := NewBot(runBot, s)
	TPop()
	if bot == nil {
		log.Printf("Failed to create bot \"%s\"", runBot)
		return
	}

	// Send go to tell server we are ready to process turns
	fmt.Fprintf(os.Stdout, "go\n")

	for {
		// READ TURN INFO FROM SERVER
		TPush("@turnscan")
		var t *Turn
		t, err = TurnScan(s.Map, in, t)
		TPop()
		if t.End || err != nil {
			break
		}

		if t.Turn != s.Turn+1 {
			log.Printf("Turns out of order Turn parse is %d expected %d", t.Turn, s.Turn+1)
		}
		if s.Turn > maxTurn {
			log.Printf("Reached MaxTurn...")
			break
		}
		turns = append(turns, t)

		TPush("@process")
		s.ProcessTurn(t)
		TPop()
		TPush("@statistics")
		s.UpdateStatistics(t)
		TPop()

		if refmap != nil {
			TPush("@validatemap")
			count, out := MapValidate(refmap, s.Map)
			TPop()
			if count > 0 {
				log.Print(out)
			}
		}

		// Generate order list
		TPush("@turn")
		bot.DoTurn(s)
		TPop()
	}

	// If we are running on tcp dump to file
	if Debug[DBG_GatherTime] {
		file := "/tmp/" + runBot + ".csv"
		if os.Getenv("BHOST") != "" {
			file = os.Getenv("BHOST") + "-" + os.Getenv("GAME") + "-" + runBot + ".csv"
		}
		TDump(file)
	}
}
