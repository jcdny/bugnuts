package main

import (
	"flag"
	"log"
	"bufio"
	"os"
	"fmt"
	"time"
	"runtime"
	"strings"
	. "bugnuts/maps"
	. "bugnuts/viz"
	. "bugnuts/watcher"
	. "bugnuts/debug"
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

func init() {
	//	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.SetFlags(log.Lshortfile)

	vizList := ""
	vizHelp := "Visualize: all,none,useful"
	for flag := range Viz {
		vizHelp += "," + flag
	}

	flag.StringVar(&runBot, "b", "v7", "Which bot to run\n\t"+strings.Join(BotList(), "\n\t"))
	flag.StringVar(&vizList, "V", "", vizHelp)
	flag.IntVar(&debugLevel, "d", 0, "Debug level")
	flag.StringVar(&mapName, "m", "", "Map file -- Used to validate generated map, hill guessing etc.")
	flag.StringVar(&watchPoints, "w", "", "Watch points \"T1:T2@R,C,N[;T1:T2...]\", \":\" will watch everything")
	flag.Parse()

	if BotGet(runBot) == nil {
		log.Printf("Unrecognized bot \"%s\", Registered bots:\n\t%s\n", runBot, strings.Join(BotList(), "\n\t"))
		return
	}

	log.SetPrefix(runBot + ":")

	SetDebugLevel(debugLevel)
	SetViz(vizList, Viz)
}

func main() {
	//TurnTimer()

	var refmap *Map
	if mapName != "" {
		refmap, _ = MapLoadFile(MapFile(mapName))
	}

	in := bufio.NewReader(os.Stdin)

	// T0 START
	btime := time.Nanoseconds()
	etime, stime := btime, btime
	egc := runtime.MemStats.PauseTotalNs
	sgc := egc

	// Load game definition
	g, err := GameScan(in)
	if err != nil {
		log.Panicf("Start(in) failed (%s)", err)
	} else if Debug[DBG_Start] {
		log.Printf("Game Info:\n%v\n", g)
	}

	s := NewState(g)
	turns := make([]*Turn, 1, s.Turns+2)

	// Create watch points
	WS = NewWatches(s.Rows, s.Cols, s.Turns)
	if len(watchPoints) > 0 {
		wlist := strings.Split(watchPoints, ";")
		WS.Load(wlist)
	}

	bot := NewBot(runBot, s)
	if bot == nil {
		log.Printf("Failed to create bot \"%s\"", runBot)
		return
	}

	// TODO this is a hack
	if runBot == "v7" {
		s.Testing = false
	}

	// Send go to tell server we are ready to process turns
	fmt.Fprintf(os.Stdout, "go\n")

	// Timing for turn 0 
	// Timing hoohah
	btime = time.Nanoseconds() // btime is reset since we want per turn time excluding steup.
	stime, etime = etime, btime
	sgc, egc = egc, runtime.MemStats.PauseTotalNs
	if Debug[DBG_TurnTime] {
		log.Printf("TURN %d %.2fms %.2fms GC",
			0,
			float64(etime-stime)/1e6,
			float64(egc-sgc)/1e6)
	}

	for {
		// READ TURN INFO FROM SERVER
		var t *Turn
		t, _ = TurnScan(s.Map, in, t)
		if t.End {
			break
		}
		if t.Turn != s.Turn+1 {
			log.Printf("Turns out of order Turn parse is %d expected %d", t.Turn, s.Turn+1)
		}
		turns = append(turns, t)

		s.ProcessTurn(t)
		s.UpdateStatistics(t)

		if refmap != nil {
			count, out := MapValidate(refmap, s.Map)
			if count > 0 {
				log.Print(out)
			}
		}

		// Generate order list
		bot.DoTurn(s)

		// Timing hoohah
		stime, etime = etime, time.Nanoseconds()
		sgc, egc = egc, runtime.MemStats.PauseTotalNs
		if Debug[DBG_TurnTime] {
			log.Printf("TURN %d %.2fms %.2fms GC %.2f SGap",
				s.Turn,
				float64(etime-t.Started)/1e6,
				float64(egc-sgc)/1e6,
				float64(t.Started-stime)/1e6)
		}
	}

	if Debug[DBG_TurnTime] {
		etime = time.Nanoseconds()
		log.Printf("TOTAL TIME %.2fms/turn for %d Turns",
			float64(etime-btime)/1e6/float64(s.Turn), s.Turn)
	}
}
