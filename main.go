package main

import (
	"flag"
	"log"
	"bufio"
	"os"
	"fmt"
	"time"
	"strings"
)

var Debug int = 0
var runBot string
var mapFile string
var paramKey string

var Viz = map[string]bool{
	"path":    false,
	"vcount":  false,
	"horizon": false,
	"threat":  false,
	"error":   false,
	"targets": false,
	"monte":   false,
}

func init() {
	vizList := ""
	vizHelp := "Visualize: all,none,useful,"
	for flag, _ := range Viz {
		vizHelp += flag
	}
	flag.StringVar(&vizList, "V", "", vizHelp)

	flag.IntVar(&Debug, "d", 0, "Debug level 0 none 1 game 2 per turn 3 per ant 4 excessive")
	flag.StringVar(&runBot, "b", "CUR", "Which bot to run")
	flag.StringVar(&mapFile, "m", "", "Map file, if provided will be used to validate generated map, hill guessing etc.")
	flag.StringVar(&paramKey, "p", "", "Parameter set, defaults to default.BOT")

	flag.Parse()

	if vizList != "" {
		for _, word := range strings.Split(strings.ToLower(vizList), ",") {
			switch word {
			case "all":
				for flag, _ := range Viz {
					Viz[flag] = true
				}
			case "none":
				for flag, _ := range Viz {
					Viz[flag] = false
				}
			case "useful":
				Viz["path"] = true
				Viz["horizon"] = true
				Viz["targets"] = true
				Viz["error"] = true
				Viz["monte"] = true
			default:
				_, ok := Viz[word]
				if !ok {
					log.Printf("Visualization flag %s not known", word)
				} else {
					Viz[word] = true
				}
			}
		}
	}
}

var Times = make(map[string]int64, 30)

func main() {
	var s State
	var bot Bot

	in := bufio.NewReader(os.Stdin)

	err := s.Start(in)

	if err != nil {
		log.Panicf("Start(in) failed (%s)", err)
	} else if Debug > 1 {
		log.Printf("State:\n%v\n", &s)
	}

	// Set up bot
	switch runBot {
	case "v0":
		bot = NewBotV0(&s)
	case "v3":
		bot = NewBotV3(&s)
	case "v4":
		bot = NewBotV4(&s)
	case "v5":
		bot = NewBotV5(&s)
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
	s.bot = &bot

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
	for {

		ptime, ntime = ntime, time.Nanoseconds()
		log.Printf("%d TURN TOOK %.2fms", s.Turn,
			float64(ntime-ptime)/1000000)

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

		if Debug > 1 {
			log.Printf("TURN %d Generating orders turn", s.Turn)
		}
		// Generate order list
		Times["preturn"] = time.Nanoseconds()
		bot.DoTurn(&s)
		Times["postturn"] = time.Nanoseconds()

		// additional thinking til near timeout

		log.Printf("%d Parse %.2fms, Turn %.2fms", s.Turn,
			float64(Times["postparse"]-Times["preparse"])/1000000,
			float64(Times["postturn"]-Times["preturn"])/1000000)

	}

	// Read end of game data.

	// Do end of game diagnostics

	//s.DumpSeen()
	//s.DumpMap()

	if Debug > 0 {
		log.Printf("Bot Result %v", bot)
	}
	ntime = time.Nanoseconds()
	log.Printf("TOTAL TIME %.2fms for %d Turns",
		float64(ntime-stime)/1000000, s.Turn)

}
