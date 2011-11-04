package main

import (
	"flag"
	"log"
	"bufio"
	"os"
	"fmt"
)

var Debug int = 0
var runBot string
var mapFile string
var Viz bool = false

func init() {
	flag.BoolVar(&Viz, "V", false, "Debug level 0 none 1 game 2 per turn 3 per ant 4 excessive")
	flag.IntVar(&Debug, "d", 0, "Debug level 0 none 1 game 2 per turn 3 per ant 4 excessive")
	flag.StringVar(&runBot, "b", "CUR", "Which bot to run")
	flag.StringVar(&mapFile, "m", "", "Map file, if provided will be used to validate generated map, hill guessing etc.")

	flag.Parse()
}

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

	for {
		// RESET FOR NEXT PARSE

		// READ TURN INFO FROM SERVER
		turn, err := s.ParseTurn()

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
		bot.DoTurn(&s)

		// additional thinking til near timeout

		// emit orders
	}

	// Read end of game data.

	// Do end of game diagnostics

	//s.DumpSeen()
	//s.DumpMap()

	if Debug > 0 {
		log.Printf("Bot Result %v", bot)
	}
}
