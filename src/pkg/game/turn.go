package game

import (
	"os"
	"log"
	"bufio"
	"strconv"
	"strings"
	. "bugnuts/maps"
	. "bugnuts/torus"
)

type PlayerLoc struct {
	Player int
	Loc    Location
}

// PlayerLoc Sorting by Loc then Player
type PlayerLocSlice []PlayerLoc

func (p PlayerLocSlice) Len() int { return len(p) }
func (p PlayerLocSlice) Less(i, j int) bool {
	return p[i].Loc < p[j].Loc || (p[i].Loc == p[j].Loc && p[i].Player < p[j].Player)
}
func (p PlayerLocSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// PlayerLoc Sorting by player then loc
type LocPlayerSlice []PlayerLoc

func (p LocPlayerSlice) Len() int { return len(p) }
func (p LocPlayerSlice) Less(i, j int) bool {
	return p[i].Player < p[j].Player || (p[i].Player == p[j].Player && p[i].Loc < p[j].Loc)
}
func (p LocPlayerSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func stringPlayerLoc(t *Turn, pl *PlayerLoc) string {
	str := ""
	p := t.ToPoint(pl.Loc)
	str += strconv.Itoa(p.R) + " " + strconv.Itoa(p.C) + " " + strconv.Itoa(pl.Player)
	return str
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

func (t *Turn) String() string {
	str := ""
	for i := range t.W {
		p := t.ToPoint(t.W[i])
		str += "\nw " + strconv.Itoa(p.R) + " " + strconv.Itoa(p.C)
	}
	for _, pl := range t.H {
		str += "\nh " + stringPlayerLoc(t, &pl)
	}
	for _, pl := range t.A {
		str += "\na " + stringPlayerLoc(t, &pl)
	}
	for i := range t.F {
		p := t.ToPoint(t.F[i])
		str += "\nf " + strconv.Itoa(p.R) + " " + strconv.Itoa(p.C)
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
func TurnScan(m *Map, in *bufio.Reader, tl *Turn) (*Turn, os.Error) {
	var t *Turn
	if tl == nil {
		t = &Turn{
			Map: m,
			W:   make([]Location, 0, 200),
			F:   make([]Location, 0, 10),
			A:   make([]PlayerLoc, 0, 10),
			H:   make([]PlayerLoc, 0, 10),
			D:   make([]PlayerLoc, 0, 10),
		}
	} else {
		t = &Turn{
			Map: m,
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

		loc := m.ToLocation(Point{R: Row, C: Col})

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
