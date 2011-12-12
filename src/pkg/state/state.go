package state

import (
	"rand"
	"sort"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/game"
)

const (
	MaxPlayers int = 10
)

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
	Rand         *rand.Rand         // the per bot stateful rng
	SSID         int                // The sym id at turn start
	Turn         int                // Current turn number, updated at start of ProcessTurn()
	Started      int64              // Nanoseconds since epoch at turn start
	Cutoff       int64              // Nanoseconds til EmitMove() required, see ProcessTurn() for details
	Ants         []map[Location]int // Ant lists List by playerid value is turn seen
	Food         map[Location]int   // Food Seen
	InitialHills int                // Number of hills per player from Turn 1
	NHills       [MaxPlayers]int    // Count of live hills per player
	Hills        map[Location]*Hill // Hill list
	Stats        *Statistics        //Computed statistics
	Met          *Metrics           // Map Metrics
	Testing      bool
	// Caches
	AttackMask *Mask
	ViewMask   *Mask
}

func NewState(g *GameInfo) *State {
	// Initialize Maps and cache some precalculated results
	m := NewMap(g.Rows, g.Cols, -1)
	s := &State{
		GameInfo: g,
		Map:      m,
		Met:      NewMetrics(m),
		Stats:    NewStatistics(g),
	}

	// Mask Cache

	s.ViewMask = MakeMask(s.ViewRadius2, s.Rows, s.Cols)
	s.AttackMask = MakeMask(s.AttackRadius2, s.Rows, s.Cols)

	// Populate calculated locations caches.
	m.OffsetsCachePopulateAll(s.ViewMask)
	m.OffsetsCachePopulateAll(s.AttackMask)

	// From every point on the map we know nothing.
	for i := range s.Met.Unknown {
		s.Met.Unknown[i] = len(s.ViewMask.P)
	}

	// Food and Ant things
	s.Food = make(map[Location]int)
	s.Ants = make([]map[Location]int, MaxPlayers)
	s.Hills = make(map[Location]*Hill)

	return s
}

func (s *State) Turn1() {
	s.InitialHills = len(s.HillLocations(0))
	seed := s.hillHash(s.PlayerSeed)
	s.Rand = rand.New(rand.NewSource(seed))
	// log.Printf("TURN 1 Seed used %v (Player Seed %v and hillHash)", seed, s.PlayerSeed)
}

// Compute the PlayerSeed XORed with the hill locations so 
// if the bot plays itself it has 2 different seeds but the game is 
// still deterministic.
func (s *State) hillHash(seed int64) int64 {
	seed = seed >> 16
	locs := s.HillLocations(0)
	sort.Sort(LocationSlice(locs))
	seed ^= int64(locs[0]) * 1327 // 10 less than 1337

	return seed
}

func (s *State) SetBlock(l Location) {
	s.Map.Grid[l] = BLOCK
}

func (s *State) SetOccupied(l Location) {
	s.Map.Grid[l] = OCCUPIED
}
