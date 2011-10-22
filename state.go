package main

import (
	"os"
	"log"
	"rand"
	"strconv"
	"bufio"
	"strings"
)

//State keeps track of everything we need to know about the state of the game
type State struct {
	in *bufio.Reader
	// Game parameter set

	LoadTime      int   //in milliseconds
	TurnTime      int   //in milliseconds
	Rows          int   //number of rows in the map
	Cols          int   //number of columns in the map
	Turns         int   //maximum number of turns in the game
	ViewRadius2   int   //view radius squared
	AttackRadius2 int   //battle radius squared
	SpawnRadius2  int   //spawn radius squared
	PlayerSeed    int64 //random player seed
	Turn          int   //current turn number

	Map *Map
}

//Start takes the initial parameters from stdin
//Reads through the "ready" line.
func (s *State) Start(reader *bufio.Reader) os.Error {
	s.in = reader

	for {
		line, err := s.in.ReadString('\n')

		if err != nil {
			return err
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
				s.PlayerSeed = 42
			}
			s.PlayerSeed = param64
			continue
		}

		param, err := strconv.Atoi(words[1])
		if err != nil {
			log.Printf("Parse failed for \"%s\" (%v)", line, err)
			continue
		}

		switch words[0] {
		case "loadtime":
			s.LoadTime = param
		case "turntime":
			s.TurnTime = param
		case "rows":
			s.Rows = param
		case "cols":
			s.Cols = param
		case "turns":
			s.Turns = param
		case "viewradius2":
			s.ViewRadius2 = param
		case "attackradius2":
			s.AttackRadius2 = param
		case "spawnradius2":
			s.SpawnRadius2 = param
		case "turn":
			s.Turn = param
		default:
			log.Printf("unknown command: %s", line)
		}

	}

	s.Map = s.NewMap()

	if s.PlayerSeed != 0 {
		rand.Seed(s.PlayerSeed)
	}

	return nil
}

func (s *State) NewMap() *Map {
	size := s.Rows * s.Cols
	m := &Map{
		Grid:    make([]Item, size),
		Seen:    make([]Seen, size),
		Visible: make([]bool, size),
	}

	return m
}

func (s *State) String() string {
	str := ""

	str += "turn " + strconv.Itoa(s.Turn) + "\n"
	str += "loadtime " + strconv.Itoa(s.LoadTime) + "\n"
	str += "turntime " + strconv.Itoa(s.TurnTime) + "\n"
	str += "rows " + strconv.Itoa(s.Rows) + "\n"
	str += "cols " + strconv.Itoa(s.Cols) + "\n"
	str += "turns " + strconv.Itoa(s.Turns) + "\n"
	str += "viewradius2 " + strconv.Itoa(s.ViewRadius2) + "\n"
	str += "attackradius2 " + strconv.Itoa(s.AttackRadius2) + "\n"
	str += "spawnradius2 " + strconv.Itoa(s.SpawnRadius2) + "\n"
	str += "player_seed " + strconv.Itoa64(s.PlayerSeed) + "\n"
	str += "ready\n"

	return str
}

func (s *State) ToLocation(p Point) Location {
	return Location(p.r*s.Cols + p.c)
}

func (s *State) ToPoint(l Location) (p Point) {
	p = Point{r: int(l) / s.Cols, c: int(l) % s.Cols}

	return
}

func (s *State) ParseTurn() (line string, err os.Error) {
	for {
		line, err = s.in.ReadString('\n')

		if err != nil {
			break
		}

		line = line[:len(line)-1] // remove the delimiter

		if line == "" {
			continue
		}

		if line == "go" || line == "end" {
			break // EXIT
		}

		words := strings.SplitN(line, " ", 5)

		if words[0] == "turn" {
			if len(words) != 2 {
				log.Printf("Invalid command format: \"%s\"", line)
			}
			turn, err := strconv.Atoi(words[1])
			if err != nil {
				log.Printf("Atoi error %s \"%v\"", line, err)
			}
			if turn != s.Turn+1 {
				log.Printf("Turn number out of sync, expected %v got %v", s.Turn+1, turn)
			}
			s.Turn = turn

			continue
		}

		if len(words) < 3 || len(words) > 4 {
			log.Printf("Invalid command format: \"%s\"", line)
			continue
		}

		var Row, Col, Player int
		// Here we have parsed the turn lines and any terminating line like go or end
		// so now just points and players.
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
		p := Point{r: Row, c: Col}

		if len(words) > 3 {
			Player, err = strconv.Atoi(words[3])
			if err != nil {
				log.Printf("Atoi error %s \"%v\"", line, err)
				continue
			}
		}

		switch words[0] {
		case "f":
			s.AddFood(p)
		case "w":
			s.AddWater(p)
		case "a":
			s.AddAnt(p, Player)
		case "h":
			s.AddHill(p, Player)
		case "d":
			s.AddDeadAnt(p, Player)
		default:
			log.Printf("Unknown turn data \"%s\"", line)
		}
	}

	// exit condition above is "go" or "end" or error on readline.
	return
}

func (s *State) AddWater(p Point) {
	s.Map.Grid[s.ToLocation(p)] = WATER
}

func (s *State) AddFood(p Point) {
	s.Map.Grid[s.ToLocation(p)] = FOOD
}

func (s *State) AddAnt(p Point, Player int) {
	s.Map.Grid[s.ToLocation(p)] = MY_ANT + Item(Player)

}

func (s *State) AddDeadAnt(p Point, Player int) {
}

func (s *State) AddHill(p Point, Player int) {
}
