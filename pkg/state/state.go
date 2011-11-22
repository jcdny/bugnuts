package state

import (
	"log"
	"rand"
	. "bugnuts/maps"
)

const (
	MaxPlayers int = 10
)

type Statistics struct {
	Dead []map[Location]int // Death count per location by player
	Died [MaxPlayers]int    // Chronicle of a Death Foretold
}

type Hill struct {
	Location Location //
	Player   int      // The owner of the hill
	// hill state info
	Found  int  // Turn we first saw it
	Seen   int  // Last turn we saw it
	Killed int  // First Turn we no longer saw it
	Killer int  // Who we think killed it, may be a guess
	guess  bool // Are we guessing location
	ssid   int  // was it a sym guess...
}

type State struct {
	// Game parameter set
	*GameInfo
	*Map
	SSID int // The sym id at turn start
	Turn int //current turn number

	AttackMask *Mask
	ViewMask   *Mask

	Ants         []map[Location]int // Ant lists List by playerid value is turn seen
	Food         map[Location]int   // Food Seen
	InitialHills int                // Number of hills per player from Turn 1
	NHills       [MaxPlayers]int    // Count of live hills per player
	Hills        map[Location]*Hill // Hill list

	Stats *Statistics

	// Map Metrics
	Met *Metrics
}

func (g *GameInfo) NewState() *State {
	// Initialize Maps and cache some precalculated results
	m := NewMap(g.Rows, g.Cols, -1)
	s := &State{
		GameInfo: g,
		Map:      m,
		Met:      NewMetrics(m),
	}
	// Mask Cache

	s.ViewMask = MakeMask(s.ViewRadius2, s.Rows, s.Cols)
	s.AttackMask = MakeMask(s.AttackRadius2, s.Rows, s.Cols)

	// From every point on the map we know nothing.
	for i := range s.Met.Unknown {
		s.Met.Unknown[i] = len(s.ViewMask.P)
	}

	// Food and Ant things
	s.Food = make(map[Location]int)
	s.Ants = make([]map[Location]int, MaxPlayers)
	s.Hills = make(map[Location]*Hill)
	s.Stats = &Statistics{
		Dead: make([]map[Location]int, MaxPlayers),
	}

	return s
}

func (s *State) Turn1() {
	s.InitialHills = len(s.HillLocations(0))
	seed := s.hillHash(s.PlayerSeed)
	log.Printf("TURN 1 Seed %v (Player Seed %v and hillHash)", seed, s.PlayerSeed)
	rand.Seed(seed)
}

// Compute the PlayerSeed XORed with the hill locations so 
// if the bot plays itself it has 2 different seeds but the game is 
// still deterministic.
func (s *State) hillHash(seed int64) int64 {
	for _, loc := range s.HillLocations(0) {
		seed ^= int64(loc) * 1327 // 10 less than 1337
	}

	return seed
}

func (s *State) SetBlock(l Location) {
	s.Map.Grid[l] = BLOCK
}

func (s *State) SetOccupied(l Location) {
	s.Map.Grid[l] = OCCUPIED
}
