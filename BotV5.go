package main
// The v4 Bot -- Marginally less Terrible!!!!

import (
	"fmt"
	"os"
	"rand"
	"sort"
	_ "log"
)

type IntSlice []int
	
func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }


type BotV5 struct {
	
}



//NewBot creates a new instance of your bot
func NewBotV5(s *State) Bot {
	mb := &BotV5{
		//do any necessary initialization here
	}

	return mb
}

type Target struct {
	item Item
	p Point             // Target Point
	loc Location        // Target Location
	arrivals []int      // Inbound Arrival time
	ant      []Location // Inbound Source
}	

type TargetSet map[Location]Target


func (* TargetSet) Add(item, loc, count)

func (bot *BotV5) DoTurn(s *State) os.Error {
	sv := []Point{{-1, 0}, {1, 0}, {0, 1}, {0, -1}}
	
	tset := make(TargetSet, 0)

	// Generate list of food and enemy hill points.
	for _, loc := range s.FoodLocations() {
		tset.Add(FOOD, loc, 1)
	}

	for _, loc := range s.EnemyHillLocations() {
		tar = append(tar, TargetSet{
		item: FOOD, 
		p: s.ToPoint(loc),
		loc: loc,
		arrivals: make([]int),
		and: make([]Location),
	})
	}

	// List of available ants
	ants := make([]Locations, len(s.Ants[0]), len(s.Ants[0]))
	ants := copy(ants, s.Ants[0])
	// TODO remove dedicated ants eg sentinel, capture, defense guys

	for iter := 0; iter < 5 && len(ants) > 0; iter++ {
		f, _, _ := MapFill(s.Map, s.Map.ToPoints(targets))


		// Build list of locations sorted by depth
		ll := make(map[int][]Location )
		var dl []int
		for loc, _ := range ants {
			depth := int(f.Depth[loc])
			if _, ok := ll[depth]; !ok {
				ll[depth] = make([]Location, 0)
				dl = append(dl, int(depth))
			}
			ll[depth] = append(ll[depth], loc)
		}
		
		sort.Sort(IntSlice(dl))
		
		for _,depth := range dl {
			for _, loc := range ll[depth] {
				p := s.Map.ToPoint(loc)
				dir := rand.Perm(4)
				for _, d := range dir {
					np := s.Map.PointAdd(s.Map.ToPoint(loc), sv[d])
					nl := s.Map.ToLocation(np)
					// log.Printf("Turn %d %d %v to %v depth %d to %d", s.Turn, d, p, np, depth, f.Depth[nl])
					
					if f.Depth[nl] < uint16(depth) && 
						( s.Map.Grid[nl] == LAND || s.Map.Grid[nl].IsEnemyHill()) {
						// We have a valid next step, path in to dest and see if 
						// We should remove ant and possibly target
						s.Map.Grid[nl] = MY_ANT
						s.Map.Grid[loc] = LAND
						fmt.Fprintf(os.Stdout, "o %d %d %c\n", p.r, p.c, ([4]byte{'n', 's', 'e', 'w'})[d])	
						break
					}
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

