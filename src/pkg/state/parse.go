package state

import (
	"os"
	"log"
	"bufio"
	"strconv"
	"strings"
	. "bugnuts/maps"
	. "bugnuts/torus"
)

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

type PlayerLoc struct {
	Player int
	Loc    Location
}

type Turn struct {
	Turn int
	W    []Location
	F    []Location
	A    []PlayerLoc
	H    []PlayerLoc
	D    []PlayerLoc
	End  bool
	*Map
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

// Take the settings from the state string and emit the header for input.
func (g *GameInfo) String() string {
	str := ""

	str += "turn 0\n"
	str += "loadtime " + strconv.Itoa(g.LoadTime) + "\n"
	str += "turntime " + strconv.Itoa(g.TurnTime) + "\n"
	str += "rows " + strconv.Itoa(g.Rows) + "\n"
	str += "cols " + strconv.Itoa(g.Cols) + "\n"
	str += "turns " + strconv.Itoa(g.Turns) + "\n"
	str += "viewradius2 " + strconv.Itoa(g.ViewRadius2) + "\n"
	str += "attackradius2 " + strconv.Itoa(g.AttackRadius2) + "\n"
	str += "spawnradius2 " + strconv.Itoa(g.SpawnRadius2) + "\n"
	str += "player_seed " + strconv.Itoa64(g.PlayerSeed) + "\n"

	return str
}

func stringPlayerLoc(t *Turn, pl *PlayerLoc) string {
	str := ""
	p := t.ToPoint(pl.Loc)
	str += strconv.Itoa(p.R) + " " + strconv.Itoa(p.C) + " " + strconv.Itoa(pl.Player)
	return str
}

func (t *Turn) String() string {
	str := ""
	for i := range t.W {
		p := t.ToPoint(t.W[i])
		str += "\nw " + strconv.Itoa(p.R) + " " + strconv.Itoa(p.C)
	}
	for i := range t.F {
		p := t.ToPoint(t.F[i])
		str += "\nf " + strconv.Itoa(p.R) + " " + strconv.Itoa(p.C)
	}
	for _, pl := range t.H {
		str += "\nh " + stringPlayerLoc(t, &pl)
	}
	for _, pl := range t.A {
		str += "\na " + stringPlayerLoc(t, &pl)
	}
	for _, pl := range t.D {
		str += "\nd " + stringPlayerLoc(t, &pl)
	}

	// we just chop first \n here
	if str != "" {
		str = str[1:]
	}

	return str
}

// Providing tl allows for better sizing of turn slices.
func (s *State) TurnScan(in *bufio.Reader, tl *Turn) (*Turn, os.Error) {
	var t *Turn
	if tl == nil {
		t = &Turn{
			Map: s.Map,
			W:   make([]Location, 0, 200),
			F:   make([]Location, 0, 10),
			A:   make([]PlayerLoc, 0, 10),
			H:   make([]PlayerLoc, 0, 10),
			D:   make([]PlayerLoc, 0, 10),
		}
	} else {
		t = &Turn{
			Map: s.Map,
			W:   make([]Location, 0, len(tl.W)*2),
			F:   make([]Location, 0, len(tl.F)*3/2),
			A:   make([]PlayerLoc, 0, len(tl.A)*3/2),
			H:   make([]PlayerLoc, 0, len(tl.H)*3/2),
			D:   make([]PlayerLoc, 0, len(tl.D)*2),
		}
	}

	var err os.Error
	var line string
	for {
		line, err = in.ReadString('\n')
		if err != nil {
			break
		}

		line = line[:len(line)-1] // remove the delimiter

		if line == "" {
			continue
		}

		if line == "go" {
			break // EXIT
		}

		if line == "end" {
			t.End = true
			break
		}

		words := strings.SplitN(line, " ", 5)
		if words[0] == "turn" {
			if len(words) != 2 {
				log.Printf("Invalid command format: \"%s\"", line)
			}
			t.Turn, err = strconv.Atoi(words[1])
			if err != nil {
				log.Printf("Atoi error %s \"%v\"", line, err)
			}
			continue
		}

		if len(words) < 3 || len(words) > 4 {
			log.Printf("Invalid command format: \"%s\"", line)
			continue
		}

		// Here we have parsed the turn lines and any terminating line like go or end
		// so now just points and players.
		var Row, Col, Player int
		Row, err = strconv.Atoi(words[1])
		if err != nil {
			log.Printf("Atoi error %s \"%v\"", line, err)
			continue
		}
		Col, err = strconv.Atoi(words[2])
		if err != nil {
			log.Printf("Atoi error %s \"%v\"", line, err)
			continue
		}

		loc := s.Map.ToLocation(Point{R: Row, C: Col})

		if len(words) > 3 {
			Player, err = strconv.Atoi(words[3])
			if err != nil {
				log.Printf("Atoi error %s \"%v\"", line, err)
				continue
			}
		}

		// Now handle items

		switch words[0] {
		case "f":
			t.F = append(t.F, loc)
		case "w":
			t.W = append(t.W, loc)
		case "a":
			t.A = append(t.A, PlayerLoc{Player: Player, Loc: loc})
		case "h":
			t.H = append(t.H, PlayerLoc{Player: Player, Loc: loc})
		case "d":
			t.D = append(t.D, PlayerLoc{Player: Player, Loc: loc})
		default:
			log.Printf("Unknown turn data \"%s\"", line)
		}
	}

	if err == os.EOF {
		t.End = true
		err = nil
	}

	return t, err
}
