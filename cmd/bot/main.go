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
	. "bugnuts/MyBot"
	. "bugnuts/bot6"
	. "bugnuts/parameters"
)

var runBot string
var mapFile string
var paramKey string
var watchPoints string
var debugLevel int

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	vizList := ""
	vizHelp := "Visualize: all,none,useful,"
	for flag, _ := range Viz {
		vizHelp += flag
	}

	flag.StringVar(&vizList, "V", "", vizHelp)

	flag.IntVar(&debugLevel, "d", 0, "Debug level 0 none 1 game 2 per turn 3 per ant 4 excessive")
	flag.StringVar(&runBot, "b", "V6", "Which bot to run")
	flag.StringVar(&mapFile, "m", "", "Map file, used to validate generated map, hill guessing etc.")
	flag.StringVar(&paramKey, "p", "", "Parameter set, defaults to default.BOT")
	flag.StringVar(&watchPoints, "w", "", "Watch points \"T1:T2@R,C,N[;T1:T2...]\", \":\" will watch all")

	flag.Parse()

	SetDebugLevel(debugLevel)
	SetViz(vizList, Viz)
}

var Times = make(map[string]int64, 30)

func main() {
	var refmap *Map
	if mapFile != "" {
		refmap, _ = MapLoadFile("testdata/maps/" + mapFile)
	}

	in := bufio.NewReader(os.Stdin)

	// Load game definition
	g, err := GameScan(in)
	if err != nil {
		log.Panicf("Start(in) failed (%s)", err)
	} else if Debug[DBG_Start] {
		log.Printf("Game Info:\n%v\n", g)
	}

	s := g.NewState()
	turns := make([]*Turn, 1, s.Turns+2)

	// Create watch points
	WS = NewWatches(s.Rows, s.Cols, s.Turns)
	if len(watchPoints) > 0 {
		wlist := strings.Split(watchPoints, ";")
		WS.Load(wlist)
	}

	// Set up bot
	if paramKey == "" {
		paramKey = runBot
	}

	var bot Bot
	switch runBot {
	case "CUR":
		fallthrough // no flag given run latest defined bot...
	case "V6":
		bot = NewBotV6(&s, ParameterSets[paramKey])
	default:
		log.Printf("Unkown bot %s", runBot)
		return
	}
	// Send go to tell server we are ready to process turns
	fmt.Fprintf(os.Stdout, "go\n")
	ttime := time.Nanoseconds()
	etime := time.Nanoseconds()
	egc := runtime.MemStats.PauseTotalNs
	for {
		// READ TURN INFO FROM SERVER]
		var t *Turn

		t, _ = s.TurnScan(in, t)
		if t.Turn != s.Turn+1 {
			log.Printf("Turns out of order new is %d expected %d", t.Turn, s.Turn+1)
		}
		turns = append(turns, t)
		s.ProcessTurn(t)

		if refmap != nil {
			count, out := MapValidate(refmap, s.Map)
			if count > 0 {
				log.Print(out)
			}
		}

		// Generate order list
		bot.DoTurn(&s)

		// Timing hohah
		stime, etime := etime, time.Nanoseconds()
		sgc, egc := egc, runtime.MemStats.PauseTotalNs
		if Debug[DBG_TurnTime] {
			log.Printf("TURN %d %.2fms %.2fms GC",
				s.Turn,
				float64(etime-stime)/1000000,
				float64(egc-sgc)/1000000)
		}
	}

	if Debug[DBG_TurnTime] {
		ntime = time.Nanoseconds()
		log.Printf("TOTAL TIME %.2fms/turn for %d Turns",
			float64(etime-ttime)/1000000/float64(s.Turn), s.Turn)
	}
}
