package main

import (
	"os"
	"log"
	"rand"
	"math"
	"bufio"
	"fmt"
	"strconv"
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

	// Cached circles etc...
	Radius        int
	ViewPoints    []Point
	ViewLocations []Location
	ViewAdd       [][]Point // NSEW points added
	ViewRem       [][]Point // NSEW points removed

	Ants []Point          // My Ants
	Food map[Location]int // Food Seen

	// Map State
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

	// Initialize Maps and cache some precalculated results
	s.Map = NewMap(s.Rows, s.Cols, -1)

	// collection of viewpoints
	s.ViewPoints = GenCircleTable(s.ViewRadius2)
	s.ViewLocations = s.ToLocations(s.ViewPoints)
	s.Radius = int(math.Sqrt(float64(s.ViewRadius2)))

	// Step innovation cache
	add, remove := moveChangeCache(s.ViewRadius2, s.ViewPoints)
	s.ViewAdd = add
	s.ViewRem = remove

	// Food and Ant things
	s.Food = make(map[Location]int)
	s.Ants = make([]Point, 0, 500)

	if s.PlayerSeed != 0 {
		rand.Seed(s.PlayerSeed)
	}

	return nil
}

func (s *State) ParamsToString() string {
	str := ""

	str += "turn 0\n"
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

func (s *State) String() string {
	str := ""

	str += "turn " + strconv.Itoa(s.Turn) + "\n"
	str += "rows " + strconv.Itoa(s.Rows) + "\n"
	str += "cols " + strconv.Itoa(s.Cols) + "\n"
	str += "player_seed " + strconv.Itoa64(s.PlayerSeed) + "\n"
	return str
}

func (s *State) Donut(p Point) Point {
	if p.r < 0 {
		p.r += s.Rows
	}
	if p.r >= s.Rows {
		p.r -= s.Rows
	}
	if p.c < 0 {
		p.c += s.Cols
	}
	if p.c >= s.Cols {
		p.c -= s.Cols
	}

	return p
}

// Take a Point and return a Location
func (s *State) ToLocation(p Point) Location {
	p = s.Donut(p)
	return Location(p.r*s.Cols + p.c)
}

//Take a slice of Point and return a slice of Location
//Used for offsets so it does not donut things.
func (s *State) ToLocations(pv []Point) []Location {
	lv := make([]Location, len(pv), len(pv)) // maybe use cap(pv)
	for i, p := range pv {
		lv[i] = Location(p.r*s.Cols + p.c)
	}

	return lv
}

func (s *State) ToPoint(l Location) (p Point) {
	p = Point{r: int(l) / s.Cols, c: int(l) % s.Cols}

	return
}

func (s *State) PointAdd(p1, p2 Point) Point {
	return s.Donut(Point{r: p1.r + p2.r, c: p1.c + p2.c})
}

func (s *State) ResetGrid() {
	for i, t := range s.Map.Seen {
		if t == s.Turn {
			if s.Map.Grid[i] > LAND {
				s.Map.Grid[i] = LAND
			}
		}
	}
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

			s.ResetGrid() // TODO Mysterious to have it here...
			s.Ants = s.Ants[:0]

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
	l := s.ToLocation(p)
	s.Map.Grid[l] = FOOD
	s.Food[l] = s.Turn // Save food seen to check if its moved.
}

func (s *State) AddAnt(p Point, Player int) {
	if s.Map.Grid[s.ToLocation(p)] <= MY_HILL {
		s.Map.Grid[s.ToLocation(p)] = MY_ANT + Item(Player)
		if Player == 0 {
			s.UpdateLand(p)
			s.UpdateSeen(p)
			s.Ants = append(s.Ants, p) // Track my ants
		}
	}
}

func (s *State) AddDeadAnt(p Point, Player int) {
}

func (s *State) AddHill(p Point, Player int) {
	s.Map.Grid[s.ToLocation(p)] = MY_HILL + Item(Player)
}

func (s *State) UpdateLand(p Point) {
	if p.c > s.Radius && p.c+s.Radius < s.Cols && p.r > s.Radius && p.r+s.Radius < s.Rows {
		// In interior of map so use loc offsets
		l := s.ToLocation(p)
		for _, offset := range s.ViewLocations {
			if s.Map.Grid[l+offset] == UNKNOWN {
				s.Map.Grid[l+offset] = LAND
			}
		}
	} else {
		for _, op := range s.ViewPoints {
			l := s.ToLocation(s.PointAdd(p, op))
			if s.Map.Grid[l] == UNKNOWN {
				s.Map.Grid[l] = LAND
			}
		}
	}
}

func (s *State) UpdateSeen(p Point) {
	if p.c > s.Radius && p.c+s.Radius < s.Cols && p.r > s.Radius && p.r+s.Radius < s.Rows {
		// In interior of map so use loc offsets
		l := s.ToLocation(p)
		for _, offset := range s.ViewLocations {
			s.Map.Seen[l+offset] = s.Turn
		}
	} else {
		for _, op := range s.ViewPoints {
			s.Map.Seen[s.ToLocation(s.PointAdd(p, op))] = s.Turn
		}
	}
}

func (s *State) DoTurn() {
	sv := []Point{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
	for _, p := range s.Ants {
		best := math.MinInt32
		var score [4]int
		for d, op := range sv {
			tp := s.PointAdd(p, op)
			if s.validPoint(tp) {
				if false && rand.Intn(8) == 0 {
					score[d] = 500
				} else {
					score[d] = s.Score(p, tp, s.ViewAdd[d])
				}
				if score[d] > best {
					best = score[d]
				}
			} else {
				score[d] = -9999
			}
		}

		if Debug > 2 {
			log.Printf("TURN %d point %v score %v best %v", s.Turn, p, score, best)
		}

		if best > math.MinInt32 {
			var bestd []int
			for d, try := range score {
				if try == best {
					bestd = append(bestd, d)
				}
			}
			pp := rand.Perm(len(bestd))[0]
			// Swap the current and target cells
			tp := s.PointAdd(p, sv[bestd[pp]])
			s.Map.Grid[s.ToLocation(tp)] = MY_ANT
			s.Map.Grid[s.ToLocation(p)] = LAND
			fmt.Fprintf(os.Stdout, "o %d %d %c\n", p.r, p.c, ([4]byte{'n', 's', 'e', 'w'})[bestd[pp]])
		}
	}
	fmt.Fprintf(os.Stdout, "go\n")

}

func (s *State) Score(p, tp Point, pv []Point) int {
	score := 0

	// Score for explore
	for _, op := range pv {
		seen := s.Map.Seen[s.ToLocation(s.PointAdd(p, op))]
		switch {
		case seen < 1:
			score += 2
		case seen > s.Turn-2:
			score -= 1
		}
	}
	score = score * 17 / len(pv)

	if Debug > 3 {
		log.Printf("p %v tp %v explore score %d", p, tp, score)
	}

	// Score for nearby items
	for _, op := range s.ViewPoints {
		item := s.Map.Grid[s.ToLocation(s.PointAdd(tp, op))]
		inc := 0
		iname := ""
		if item != LAND && item != WATER {
			//log.Printf("%v %v %d %d",p, tp, op, d, item)
			d := Abs(op.c) + Abs(op.r)

			if item == MY_HILL {
				iname = "my hill"
				inc = -32 + 4*Min([]int{d, 8})
			}
			if item > MY_HILL && item < HILL10 {
				iname = "enemy hill"
				inc = 1500 - 100*Min([]int{d, 10})
			}
			if item == FOOD {
				iname = "food"
				if d == 1 {
					inc = 1000
				} else {
					inc = 120 - 12*Min([]int{d, 10})
				}
			}
			if item == MY_ANT && d > 1 {
				iname = "my ant"
				inc = -30 + 5*Min([]int{d, 6})
			}
		}
		score += inc
		if Debug > 3 && iname != "" {
			log.Printf("tp %v (at %v) %s worth %d",
				tp, op, iname, inc)
		}
	}
	if Debug > 3 {
		log.Printf("p %v tp %v total score %d",
			p, tp, score)
	}
	return score
}

func (s *State) validPoint(p Point) bool {
	sv := []Point{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
	tgt := s.Map.Grid[s.ToLocation(p)]
	if tgt == FOOD || tgt == LAND || (tgt > MY_HILL && tgt <= HILL10) {
		for _, op := range sv {
			//make sure there is an exit
			ep := s.PointAdd(p, op)
			tgt := s.Map.Grid[s.ToLocation(ep)]
			if tgt == FOOD || tgt == LAND {
				return true
			}
		}
	}
	return false
}
