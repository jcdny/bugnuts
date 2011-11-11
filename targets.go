package main

import (
	"log"
	"math"
)

type Target struct {
	Item     Item
	Loc      Location // Target Location
	Count    int      // how many do we want at this location
	Pri      int      // target priority.
	Terminal bool     // Is it a terminating target 

	arrivals []int      // Inbound Arrival time
	player   []int      // Inbound player
	ant      []Location // Inbound Source
}

type TargetSet map[Location]*Target

func (tset *TargetSet) Add(item Item, loc Location, count, pri int) {
	if Debug > 3 {
		log.Printf("Adding target: item %v loc %v count %d pri %d", item, loc, count, pri)
	}
	t, found := (*tset)[loc]
	if pri < 1 {
		log.Panicf("Target pri must be > 1")
	}
	if !found || t.Pri < pri {
		// We already have this point in the target set, replace if pri is higher
		(*tset)[loc] = &Target{
			Item:     item,
			Loc:      loc,
			Count:    count,
			Pri:      pri,
			Terminal: item.IsTerminal(),
		}
	}

}

func (tset *TargetSet) Remove(loc Location) {
	t, ok := (*tset)[loc]
	if Debug > 3 {
		if ok {
			log.Printf("Removing target: item %v loc %v count %d pri %d", t.Item, t.Loc, t.Count, t.Pri)
		} else {
			log.Printf("Removing target: not found Loc %d", loc)
		}
	}
	if ok {
		// We already have this point in the target set, replace if pri is higher
		(*tset)[loc] = t, false
	}
}

func (tset TargetSet) Merge(src *TargetSet) {
	if src == nil {
		return
	}
	// Run through explore targets
	for loc, tgt := range *src {
		nt, found := tset[loc]
		if !found || nt.Pri > tgt.Pri {
			tset[loc] = tgt
		}
	}
}

func (tset *TargetSet) Pending() int {
	n := 0
	for _, t := range *tset {
		if t.Count > 0 {
			n++
		}
	}

	return n
}

func (tset *TargetSet) Active() map[Location]int {
	tp := make(map[Location]int, tset.Pending())
	for _, t := range *tset {
		if t.Count > 0 {
			tp[t.Loc] = t.Pri
		}
	}

	return tp
}

func MakeExplorers(s *State, scale float64, count, pri int) *TargetSet {

	// Set an initial group of targetpoints
	if scale <= 0 {
		scale = 1.0
	}

	stride := int(math.Sqrt(float64(s.ViewRadius2)) * scale)

	tset := make(TargetSet, (s.Rows * s.Cols / (stride * stride)))

	for r := 5; r < s.Rows; r += stride {
		for c := 5; c < s.Cols; c += stride {
			loc := s.Map.ToLocation(Point{r: r, c: c})
			tset.Add(EXPLORE, loc, count, pri)
		}
	}

	return &tset
}
