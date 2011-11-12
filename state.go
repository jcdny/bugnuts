package main

import (
	"os"
	"log"
	"rand"
	"bufio"
	"strconv"
	"strings"
)

const (
	MaxPlayers int = 10
)

type Statistics struct {
	Dead []map[Location]int // Death count per location by player
	Died [MaxPlayers]int    // Chronicle of deaths foretold
}

type Hill struct {
	Location Location //
	Player   int      // The owner of the hill
	// hill state info
	Found    int  // Turn we first saw it
	Seen     int  // Last turn we saw it
	Killed   int  // First Turn we no longer saw it
	Killer   int  // Who we think killed it, may be a guess
	guess    bool // Are we guessing location
	maxerror int  // the maximum steps to bound the unkown location
}

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

	attackMask *Mask
	viewMask   *Mask

	Ants   []map[Location]int // Ant lists List by playerid value is turn seen
	Hills  map[Location]*Hill // Hill list
	NHills int                // Number of hills per player from Turn 1
	Food   map[Location]int   // Food Seen

	Stats *Statistics

	// Map State
	Map *Map

	bot *Bot
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

	// Mask Cache
	s.viewMask = makeMask(s.ViewRadius2, s.Rows, s.Cols)
	s.attackMask = makeMask(s.AttackRadius2, s.Rows, s.Cols)

	// From every point on the map we know nothing.
	for i := range s.Map.Unknown {
		s.Map.Unknown[i] = len(s.viewMask.P)
	}

	// Food and Ant things
	s.Food = make(map[Location]int)
	s.Ants = make([]map[Location]int, MaxPlayers)
	s.Hills = make(map[Location]*Hill)
	s.Stats = &Statistics{
		Dead: make([]map[Location]int, MaxPlayers),
	}

	if s.PlayerSeed != 0 {
		rand.Seed(s.PlayerSeed)
	}

	return nil
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

func (s *State) PointAdd(p1, p2 Point) Point {
	return s.Donut(Point{r: p1.r + p2.r, c: p1.c + p2.c})
}

func (s *State) ResetGrid() {
	// Rotate threat maps and clear first.
	n := len(s.Map.Threat)
	if n > 1 {
		s.Map.Threat = append(s.Map.Threat[1:n], s.Map.Threat[0])
	}
	for i := range s.Map.Threat[0] {
		s.Map.Threat[0][i] = 0
	}

	// Set all seen map to land
	for i, t := range s.Map.Seen {
		s.Map.VisCount[i] = 0
		if t == s.Turn && s.Map.Grid[i] > LAND {
			s.Map.Grid[i] = LAND
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
			// should food clear any visibles, Remove if visible this turn


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

		loc := s.Map.ToLocation(Point{r: Row, c: Col})

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
			s.AddFood(loc)
		case "w":
			s.AddWater(loc)
		case "a":
			s.AddAnt(loc, Player)
		case "h":
			s.AddHill(loc, Player)
		case "d":
			s.AddDeadAnt(loc, Player)
		default:
			log.Printf("Unknown turn data \"%s\"", line)
		}
	}

	s.ProcessState() // Updater for all things visible

	// exit condition above is "go" or "end" or error on readline.
	return
}

func (s *State) AddWater(loc Location) {
	s.Map.Grid[loc] = WATER
}

func (s *State) AddFood(loc Location) {
	s.Food[loc] = s.Turn
}

func (s *State) AddAnt(loc Location, player int) {
	if s.Ants[player] == nil {
		s.Ants[player] = make(map[Location]int)
		// TODO New ant seen - start guessing hill loc
	}
	s.Ants[player][loc] = s.Turn

	if player == 0 {
		if s.Map.Unknown[loc] > 0 {
			s.UpdateLand(player, loc)
		}
		s.UpdateSeen(player, loc)
		s.UpdateVisCount(player, loc)
	}

	// TODO move this all to the inference step, should not be here!

	// Handle tracking Razes
	hill, found := s.Hills[loc]

	if found && !hill.guess {
		// for guesses we will update those when we validate
		// visibles - the guess location is by definition visible
		// since we got an ant in our location

		if hill.Seen == s.Turn {
			if hill.Player != player {
				// TODO work out how to infer raze - do we get the hill
				// sent with the player on the round after raze?
				// I assume not but need to check.
				//
				// If not this state should be treated as an error
				log.Printf("Error? player %d on hill player %d hill at %v",
					player, hill.Player, s.Map.ToPoint(loc))
			}
		} else if hill.Killed <= 0 {
			// we found a hill in the hash but its not marked killed
			// Mark it killed by the ant we found standing on it.
			// could be wrong ofc.
			hill.Killed = s.Turn
			hill.Killer = player
		}
	}
}

func (s *State) AddDeadAnt(loc Location, player int) {
	if s.Stats.Dead[player] == nil {
		s.Stats.Dead[player] = make(map[Location]int)
	}
	s.Stats.Dead[player][loc]++
	s.Stats.Died[player]++

	// TODO track suicides/sacrifices and who the killer was.
}

func (s *State) AddHill(loc Location, player int) {
	if hill, found := s.Hills[loc]; found {
		hill.Seen = s.Turn
		hill.guess = false
	} else {
		s.Hills[loc] = &Hill{
			Location: loc,
			Player:   player,
			Found:    s.Turn,
			Seen:     s.Turn,
			Killed:   0,
			Killer:   -1,
			guess:    false,
		}
	}
}

// Todo This could all be done in one step.  Also viewer count.
// Obvious optimizations: watch Adjacent Seen cells and do incremental updating.
func (s *State) UpdateLand(player int, loc Location) {
	nland := 0
	if s.Map.BorderDist[loc] > s.viewMask.R {
		// In interior of map so use loc offsets
		for _, offset := range s.viewMask.Loc {
			if s.Map.Grid[loc+offset] == UNKNOWN {
				s.Map.Grid[loc+offset] = LAND
			}
			if s.Map.Grid[loc+offset] != WATER {
				nland++
			}
		}
	} else {
		// non interior point lets go slow
		p := s.Map.ToPoint(loc)
		for _, op := range s.viewMask.P {
			l := s.ToLocation(s.PointAdd(p, op))
			if s.Map.Grid[l] == UNKNOWN {
				s.Map.Grid[l] = LAND
			}
			if s.Map.Grid[l] != WATER {
				nland++
			}
		}
	}

	s.Map.Land[loc] = nland
}

func (s *State) UpdateSeen(player int, loc Location) {
	s.Map.Unknown[loc] = 0
	if s.Map.BorderDist[loc] > s.viewMask.R {
		// In interior of map so use loc offsets
		for _, offset := range s.viewMask.Loc {
			s.Map.Seen[loc+offset] = s.Turn
		}
	} else {
		p := s.Map.ToPoint(loc)
		for _, op := range s.viewMask.P {
			s.Map.Seen[s.ToLocation(s.PointAdd(p, op))] = s.Turn
		}
	}
}

func (s *State) UpdateVisCount(player int, loc Location) {
	if s.Map.BorderDist[loc] > s.viewMask.R {
		// In interior of map so use loc offsets
		for _, offset := range s.viewMask.Loc {
			s.Map.VisCount[loc+offset]++
		}
	} else {
		p := s.Map.ToPoint(loc)
		for _, op := range s.viewMask.P {
			s.Map.VisCount[s.ToLocation(s.PointAdd(p, op))]++
		}
	}
}

// Take the settings from the state string and emit the header for input.
func (s *State) SettingsToString() string {
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

func (s *State) ProcessState() {
	// Assumes the loc data has all been read, and Seen/Land updated
	for loc, _ := range s.Food {
		s.Map.Grid[loc] = LAND
	}

	// Update hill data in map.
	for loc, hill := range s.Hills {
		if hill.Killed == 0 {
			if s.Map.Seen[loc] > hill.Seen {
				if hill.guess {
					// We just guessed at a location anyway, just remove for now
					s.Hills[loc] = &Hill{}, false

					// TODO: update the guess
				} else {
					// We don't see the hill to mark as killed by whoever we think was closest
					hill.Killed = s.Turn
				}
			}
		} else {
			if s.Map.Seen[loc] < s.Turn {
				// If the hill is not visible then set Horizon false
				// since it could be a source of ant.
				s.Map.Horizon[loc] = false
			}
			s.Map.Grid[loc] = MY_HILL + Item(hill.Player)
		}
	}

	if s.Turn == 1 {
		s.NHills = len(s.HillLocations(0))
	}

	for player, ants := range s.Ants {
		for loc, seen := range ants {
			if seen < s.Map.Seen[loc] || (player == 0 && seen < s.Turn) {
				ants[loc] = 0, false
			} else {
				if seen < s.Turn && player != 0 {
					s.Map.Horizon[loc] = false
				}
				if s.Map.Grid[loc].IsHill() {
					s.Map.Grid[loc] = MY_HILLANT + Item(player)
				} else {
					s.Map.Grid[loc] = MY_ANT + Item(player)
				}
			}
			// TODO Bug here since if an ant steps out of seen we don't assume it still exists
			// unless it was out move that remove it from vision

			// TODO if the ant was visble last turn, not now and there is an ant
			// we can see 1 step away from where it was assume the new ant is
			// the same ant.

			// TODO Think about this -- assuming appearing ants match missing ones,
			// tracking max ants in a region.
		}
	}

	for loc, _ := range s.Ants[0] {
		// Update the one step land count and unseen count for my ants
		s.Map.SumVisCount(loc, s.viewMask)
		for _, nloc := range s.Map.LocStep[loc] {
			s.Map.SumVisCount(nloc, s.viewMask)
			if s.Map.Unknown[nloc] > 0 {
				s.Map.UpdateCounts(nloc, s.viewMask)
			}
		}
	}

	s.Map.HBorder = s.StepHorizon(s.Map.HBorder)

	s.UpdateHillMaps()

	s.MonteCarloDensity()

	s.ComputeThreat(1, 0, s.attackMask, s.Map.Threat[len(s.Map.Threat)-1])

}

func (s *State) UpdateHillMaps() {
	// Generate the fill for all my hills.
	lend := make(map[Location]int)
	for _, hill := range s.HillLocations(0) {
		lend[hill] = 1
	}
	// log.Printf("Computing fill for %v", lend)
	s.Map.FHill, _, _ = MapFillSeed(s.Map, lend, 1)

	outbound := make(map[Location]int)
	samples, _ := s.Map.FHill.Sample(0, 60/s.NHills, 60/s.NHills)
	for _, loc := range samples {
		outbound[loc] = 1
	}
	if len(lend) > 0 && len(outbound) < 1 {
		log.Panicf("UpdateHillMaps no outside border")
	} else {
		s.Map.FDownhill, _, _ = MapFillSeed(s.Map, outbound, 1)
	}
}

func (s *State) FoodUpdate(age int) {
	// Nuke all stale food

	for loc, seen := range s.Food {
		// Better would be to compute expected pickup time give neighbors
		// in the pathing step and only revert to this when there were no
		// visible neighbors.
		//
		// Should track that anyway since does not make sense to run for 
		// food another bot will certainly get unless its to enter combat.

		if s.Map.Seen[loc] > seen || seen < s.Turn-age {
			s.Food[loc] = 0, false
			if s.Map.Grid[loc] == FOOD {
				s.Map.Grid[loc] = LAND
			}
		}
	}
}

func (s *State) FoodLocations() (l []Location) {
	for loc, _ := range s.Food {
		l = append(l, Location(loc))
	}

	return l
}

func (s *State) HillLocations(player int) (l []Location) {
	for loc, hill := range s.Hills {
		if hill.Player == player && hill.Killed == 0 {
			l = append(l, Location(loc))
		}
	}

	return l
}

func (s *State) EnemyHillLocations(player int) (l []Location) {
	for loc, hill := range s.Hills {
		if hill.Player != player && hill.Killed == 0 {
			l = append(l, Location(loc))
		}
	}

	return l
}

// Compute the threat for N turns out (currently only n = 0 or 1)
// if player > -1 then sum players not including player
func (s *State) ComputeThreat(turn, player int, mask *Mask, slice []int8) {
	var mv []Point
	switch turn {
	case 1:
		mv = mask.Union
	case 0:
		mv = mask.P
	default:
		log.Panicf("Illegal turns out = %d", turn)
	}

	if len(slice) != s.Rows*s.Cols {
		log.Panic("ComputeThreat slice size mismatch")
	}

	for i, _ := range s.Ants {
		if i != player {
			for loc, _ := range s.Ants[i] {
				p := s.Map.ToPoint(loc)
				for _, op := range mv {
					slice[s.ToLocation(s.PointAdd(p, op))]++
				}
			}
		}
	}

	return
}

func (s *State) Threat(turn int, l Location) int8 {
	i := len(s.Map.Threat) - turn + s.Turn - 1
	if i < 0 {
		log.Printf("Threat for turn %d on turn %d we only keep %d turns", turn, s.Turn, len(s.Map.Threat))
		return 0
	}
	return s.Map.Threat[i][l]
}

func (s *State) ThreatMap(turn int) []int8 {
	i := len(s.Map.Threat) - turn + s.Turn - 1
	if i < 0 {
		log.Printf("Threat for turn %d on turn %d we only keep %d turns", turn, s.Turn, len(s.Map.Threat))
		return nil
	}
	return s.Map.Threat[i]
}

func (s *State) SetBlock(l Location) {
	s.Map.Grid[l] = BLOCK
}
func (s *State) SetOccupied(l Location) {
	s.Map.Grid[l] = OCCUPIED
}

func (s *State) MoveAnt(from, to Location) bool {
	if from == to {
		return true
	}
	if s.ValidStep(to) {
		s.Map.Grid[from], s.Map.Grid[to] = LAND, OCCUPIED
		return true
	}
	return false
}

func (s *State) ValidStep(loc Location) bool {
	i := s.Map.Grid[loc]

	return i != WATER && i != BLOCK && i != OCCUPIED && i != FOOD && i != MY_ANT && i != MY_HILLANT
}

func (s *State) StepHorizon(hlist []Location) []Location {
	m := s.Map
	hlist = hlist[0:0]

	// Remove now visible cells; dont bother with water here.
	for loc, Seen := range m.Seen {
		if Seen >= s.Turn {
			m.Horizon[loc] = true
		}
	}

	// generate list of cells on border; exclude water cells here
	for loc, h := range m.Horizon {
		if !h && m.Grid[loc] != WATER {
			for _, nloc := range m.LocStep[loc] {
				if m.Horizon[nloc] && m.Grid[nloc] != WATER {
					// if the point has an adjacent non horizon point which is not water 
					// then add it to the border list.
					hlist = append(hlist, Location(loc))
					break
				}
			}
		}
	}

	// step one from all cells on border
	for _, loc := range hlist {
		for d := 0; d < 4; d++ {
			m.Horizon[m.LocStep[loc][d]] = false
		}
	}

	return hlist
}

// Update expected locations and flows for enemy ants
func (s *State) MonteCarloDensity() {
	tgt := make(map[Location]int, 32)
	for _, loc := range s.EnemyHillLocations(-1) {
		tgt[loc] = 1
	}

	for _, loc := range s.FoodLocations() {
		tgt[loc] = 20
	}
	if len(tgt) > 0 {
		ants := make([]Location, 0, 200)
		f, _, _ := MapFillSeed(s.Map, tgt, 0)

		for player := 1; player < len(s.Ants); player++ {
			for loc, _ := range s.Ants[player] {
				endloc := f.Seed[loc]
				steps := f.Depth[loc] - f.Depth[f.Seed[loc]]
				hill, ok := s.Hills[endloc]
				if !ok || hill.Player != player {
					if (ok && steps < 50) || steps < 16 {
						ants = append(ants, Location(loc))
					}
				}
			}
		}

		if len(ants) > 0 {
			// do up to 512 paths, but no more than 32 per ant
			paths := 1024 / len(ants)
			if paths > 32 {
				paths = 32
			}
			s.Map.MCDist = f.MontePathIn(s.Map, ants, paths, 1)
			s.Map.MCDistMax = Max(s.Map.MCDist)
			s.Map.MCPaths = paths

		} else {
			s.Map.MCPaths = 0
		}
	}
}
func (s *State) ToPoint(l Location) (p Point) {
	p = Point{r: int(l) / s.Cols, c: int(l) % s.Cols}
	return
}
