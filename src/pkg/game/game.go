// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

// Game has the game setup parser and various constants and datastructures.
package game

import (
	"os"
	"log"
	"bufio"
	"strconv"
	"strings"
)

const MaxPlayers = 10
const MaxMapDimension = 200

// Game definition
type GameInfo struct {
	Rows          int
	Cols          int
	LoadTime      int   //in milliseconds
	TurnTime      int   //in milliseconds
	Turns         int   //maximum number of turns in the game
	ViewRadius2   int   //view radius squared
	AttackRadius2 int   //battle radius squared
	SpawnRadius2  int   //spawn radius squared
	PlayerSeed    int64 `json:"player_seed"` //random player seed
}

func gameDefaults() *GameInfo {
	return &GameInfo{
		LoadTime:      3000,
		TurnTime:      500,
		Turns:         1000,
		ViewRadius2:   77,
		AttackRadius2: 5,
		SpawnRadius2:  1,
		PlayerSeed:    42,
	}
}

func NewGameInfo(rows, cols int) *GameInfo {
	g := gameDefaults()
	g.Rows = rows
	g.Cols = cols

	return g
}

// Take the settings from the state string and emit the header for input.
func (g *GameInfo) String() string {
	str := ""

	str += "loadtime " + strconv.Itoa(g.LoadTime) + "\n"
	str += "turntime " + strconv.Itoa(g.TurnTime) + "\n"
	str += "rows " + strconv.Itoa(g.Rows) + "\n"
	str += "cols " + strconv.Itoa(g.Cols) + "\n"
	str += "turns " + strconv.Itoa(g.Turns) + "\n"
	str += "viewradius2 " + strconv.Itoa(g.ViewRadius2) + "\n"
	str += "attackradius2 " + strconv.Itoa(g.AttackRadius2) + "\n"
	str += "spawnradius2 " + strconv.Itoa(g.SpawnRadius2) + "\n"
	str += "player_seed " + strconv.Itoa64(g.PlayerSeed)

	return str
}

func GameScan(in *bufio.Reader) (*GameInfo, os.Error) {
	g := gameDefaults()
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			return g, err
		}

		line = line[:len(line)-1] //remove the delimiter
		if line == "" {
			continue
		}

		if line == "ready" {
			break
		}

		words := strings.SplitN(line, " ", 2)
		if len(words) != 2 {
			log.Printf("Invaid param line \"%s\"", line)
			continue
		}

		if words[0] == "player_seed" {
			param64, err := strconv.Atoi64(words[1])
			if err != nil {
				log.Printf("Parse failed for \"%s\" (%v)", line, err)
				param64 = 42
			}
			g.PlayerSeed = param64
			continue
		}

		param, err := strconv.Atoi(words[1])
		if err != nil {
			log.Printf("Parse failed for \"%s\" (%v)", line, err)
			continue
		}

		switch words[0] {
		case "loadtime":
			g.LoadTime = param
		case "turntime":
			g.TurnTime = param
		case "rows":
			g.Rows = param
		case "cols":
			g.Cols = param
		case "turns":
			g.Turns = param
		case "viewradius2":
			g.ViewRadius2 = param
		case "attackradius2":
			g.AttackRadius2 = param
		case "spawnradius2":
			g.SpawnRadius2 = param
		case "turn":
			if param > 0 {
				log.Printf("Got turn %d before \"ready\", ignoring", param)
			}
		default:
			log.Printf("unknown command: %s", line)
		}
	}

	return g, nil
}
