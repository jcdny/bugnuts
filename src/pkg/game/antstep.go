package game

import (
	. "bugnuts/torus"
	. "bugnuts/maps"
)

type Neighborhood struct {
	//TODO add hill distance step
	Valid    bool
	Threat   int
	PrThreat int
	Goal     int
	PrFood   int
	Combat   int
	//Vis     int
	//Unknown int
	//Land    int
	Perm   int // permuter
	D      Direction
	Safest bool
	Item   Item
}

type AntStep struct {
	Source  Location   // our original location
	Move    Direction  // the next step
	Dest    []Location // track routing
	Steps   []int      // and distance
	Steptot int        // and sum total distance
	N       []*Neighborhood
	Foodp   bool
	Goalp   bool
	Combat  Location // pointer to the combat partition
	Perm    int      // to randomize ants when sorting
	NFree   int
}

// Order ants for trying to move.
type AntSlice []*AntStep

func (p AntSlice) Len() int      { return len(p) }
func (p AntSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p AntSlice) Less(i, j int) bool {
	if p[i].Goalp != p[j].Goalp {
		return p[i].Goalp
	}
	if p[i].NFree != p[j].NFree {
		// order by min degree of freedom but 0 degree last.
		return p[i].NFree < p[j].NFree && p[i].NFree != 0
	}
	if p[i].Goalp && p[i].Steps[0] != p[j].Steps[0] {
		return p[i].Steps[0] < p[j].Steps[0]
	}
	return p[i].Perm > p[j].Perm
}

// For ordering perspective moves...
type ENSlice []*Neighborhood

func (p ENSlice) Len() int      { return len(p) }
func (p ENSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ENSlice) Less(i, j int) bool {
	if p[i].Valid != p[j].Valid {
		return p[i].Valid
	}
	if p[i].Combat != p[j].Combat {
		return p[i].Combat > p[j].Combat
	}
	if p[i].Threat != p[j].Threat {
		return p[i].Threat < p[j].Threat
	}
	if p[i].PrThreat != p[j].PrThreat {
		return p[i].PrThreat < p[j].PrThreat
	}
	if p[i].Goal != p[j].Goal {
		return p[i].Goal > p[j].Goal
	}
	if p[i].PrFood != p[j].PrFood {
		return p[i].PrFood > p[j].PrFood
	}
	return p[i].Perm < p[j].Perm
}
