package main
// The v4 Bot -- Marginally less Terrible!!!!
//
// Entirely changed from v3 - now uses food and hill locations
// to set goals and does an iterated greedy BFS to path to goals.

import (
	"fmt"
	"os"
	"rand"
	"sort"
	_ "log"
)

type BotV4 struct {

}

//NewBot creates a new instance of your bot
func NewBotV4(s *State) Bot {
	mb := &BotV4{
	//do any necessary initialization here
	}

	return mb
}

func (bot *BotV4) DoTurn(s *State) os.Error {
	sv := []Point{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}

	// Generate list of food and enemy hill points.
	targets := s.FoodLocations()
	for _, loc := range s.EnemyHillLocations() {
		targets = append(targets, loc)
	}

	tmap := make(map[Location]int, len(targets))
	for _, loc := range targets {
		tmap[loc] = 1
	}
	// log.Printf("%v %v", targets, s.Map.ToPoints(targets))
	// Add search points

	f, _, _ := MapFill(s.Map, tmap, 0)

	// Build list of locations sorted by depth and attempt to go downhill
	ll := make(map[int][]Location)
	var dl []int
	for loc, _ := range s.Ants[0] {
		depth := int(f.Depth[loc])
		if _, ok := ll[depth]; !ok {
			ll[depth] = make([]Location, 0)
			dl = append(dl, int(depth))
		}
		ll[depth] = append(ll[depth], loc)
	}

	sort.Sort(IntSlice(dl))

	for _, depth := range dl {
		for _, loc := range ll[depth] {
			p := s.Map.ToPoint(loc)
			dir := rand.Perm(4)
			for _, d := range dir {
				np := s.Map.PointAdd(s.Map.ToPoint(loc), sv[d])
				nl := s.Map.ToLocation(np)
				// log.Printf("Turn %d %d %v to %v depth %d to %d", s.Turn, d, p, np, depth, f.Depth[nl])

				if f.Depth[nl] < uint16(depth) &&
					(s.Map.Grid[nl] == LAND || s.Map.Grid[nl].IsEnemyHill()) {
					s.Map.Grid[nl] = MY_ANT
					s.Map.Grid[loc] = LAND
					fmt.Fprintf(os.Stdout, "o %d %d %c\n", p.r, p.c, ([4]byte{'n', 's', 'e', 'w'})[d])
					break
				}
			}
		}
	}
	fmt.Fprintf(os.Stdout, "go\n")

	// refinements - path nearest ones remove food/ant pairs then regoal for spread and explore.
	// change depth of hill
	// tiebreak on global goals.

	return nil
}
