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
	. "bugnuts/watcher"
	. "bugnuts/viz"
	. "bugnuts/debug"
	. "bugnuts/state"
	. "bugnuts/bot6"
)

var runBot string
var mapFile string
var paramKey string
var watchPoints string
var debugLevel int

var WS *Watches

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	vizList := ""
	vizHelp := "Visualize: all,none,useful,"
	for flag, _ := range Viz {
		vizHelp += flag
	}

	flag.StringVar(&vizList, "V", "", vizHelp)

	flag.IntVar(&debugLevel, "d", 0, "Debug level 0 none 1 game 2 per turn 3 per ant 4 excessive")
	flag.StringVar(&runBot, "b", "CUR", "Which bot to run")
	flag.StringVar(&mapFile, "m", "", "Map file, used to validate generated map, hill guessing etc.")
	flag.StringVar(&paramKey, "p", "", "Parameter set, defaults to default.BOT")
	flag.StringVar(&watchPoints, "w", "", "Watch points \"T1:T2@R,C,N[;T1:T2...]\", \":\" will watch all")

	flag.Parse()

	SetDebugLevel(debugLevel)
	SetViz(vizList, Viz)
}

var Times = make(map[string]int64, 30)

func main() {
	var s State
	var bot Bot

	in := bufio.NewReader(os.Stdin)

	err := s.Start(in)
	WS = NewWatches(s.Rows, s.Cols, s.Turns)
	if len(watchPoints) > 0 {
		wlist := strings.Split(watchPoints, ";")
		WS.Load(wlist)
	}

	if err != nil {
		log.Panicf("Start(in) failed (%s)", err)
	} else if Debug[DBG_Start] {
		log.Printf("State:\n%v\n", &s)
	}

	// Set up bot
	switch runBot {
	/* 
		case "v0":
			bot = NewBotV0(&s)
		case "v3":
			bot = NewBotV3(&s)
		case "v4":
			bot = NewBotV4(&s)
		case "v5":
			bot = NewBotV5(&s)
	*/
	case "CUR":
		fallthrough // no flag given run latest defined bot...
	case "v6":
		bot = NewBotV6(&s)
	default:
		log.Printf("Unkown bot %s", runBot)
		return
	}

	// some of the state updating like treatment of non-visible food 
	// depends on the bot parameters.
	//s.bot = &bot

	var refmap *Map
	if mapFile != "" {
		refmap, _ = MapLoadFile("testdata/maps/" + mapFile)
	}

	// TODO Time from load to measure other bots calc time in preload.
	// Send go to tell server we are ready to process turns
	fmt.Fprintf(os.Stdout, "go\n")

	stime := time.Nanoseconds()
	ntime := stime
	ptime := stime
	ngctime := runtime.MemStats.PauseTotalNs
	ngc := runtime.MemStats.NumGC
	pgc := ngc
	pgctime := ngctime
	for {

		ptime, ntime = ntime, time.Nanoseconds()
		pgctime, ngctime = ngctime, runtime.MemStats.PauseTotalNs
		pgc, ngc = ngc, runtime.MemStats.NumGC
		runtime.GC()
		if Debug[DBG_TurnTime] || Debug[DBG_AllTime] {
			log.Printf("TURN %d TOOK %.2fms gc %.3fms Ngc %d", s.Turn,
				float64(ntime-ptime)/1000000,
				float64(ngctime-pgctime)/1000000, ngc-pgc)
		}

		// READ TURN INFO FROM SERVER
		Times["preparse"] = time.Nanoseconds()
		turn, err := s.ParseTurn()
		Times["postparse"] = time.Nanoseconds()

		if refmap != nil {
			count, out := MapValidate(refmap, s.Map)
			if count > 0 {
				log.Print(out)
			}
		}

		if err == os.EOF || turn == "end" {
			break
		}

		if Debug[DBG_Turns] {
			log.Printf("TURN %d Generating orders turn", s.Turn)
		}
		// Generate order list
		Times["preturn"] = time.Nanoseconds()
		bot.DoTurn(&s)
		Times["postturn"] = time.Nanoseconds()

		// additional thinking til near timeout
		if Debug[DBG_AllTime] {
			log.Printf("%d Parse %.2fms, Turn %.2fms", s.Turn,
				float64(Times["postparse"]-Times["preparse"])/1000000,
				float64(Times["postturn"]-Times["preturn"])/1000000)
		}

	}

	// Read end of game data.

	// Do end of game diagnostics

	//s.DumpSeen()
	//s.DumpMap()

	if Debug[DBG_TurnTime] {
		ntime = time.Nanoseconds()
		log.Printf("TOTAL TIME %.2fms/turn for %d Turns",
			float64(ntime-stime)/1000000/float64(s.Turn), s.Turn)
	}
	if Debug[DBG_Results] {
		log.Printf("Bot Result %v", bot)
	}
}
