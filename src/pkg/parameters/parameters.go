// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

// Parameters holds the per bot parameter files.
package parameters

import (
	"io/ioutil"
	"log"
	"json"
	. "bugnuts/maps"
)

type Parameters struct {
	ExpireFood       int // If we have not seen the food in this many turns then assume it's gone.
	PriorityMap      *map[string]int
	Outbound         int     // Radius inside which we step away from hill by default.
	MinHorizon       int     // minimum horizon of mystery to our hill
	DefendDistance   int     // How early to we consider an ant a threat to the hill
	RiskOffThreshold float64 // if static ants > than this turn off risk aversion
	Explore          bool    // Create initial explore points
}

func init() {
	ParameterSets["default"] = ParameterSets["V6"]
}

func (p *Parameters) Load(filename string) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("Parameter Readfile %s error %v", filename, err)
	}
	json.Unmarshal(buf, &p)
}

func (p *Parameters) Save(filename string) {
	o, err := json.Marshal(p)
	if err != nil {
		log.Panicf("Parameter Save %s failed error %v", filename, err)
	}
	ioutil.WriteFile(filename, o, 0666)
}

func (p *Parameters) Priority(i Item) int {
	key := ""
	if i.IsHill() {
		i = HILL1
	}

	switch i {
	case DEFEND:
		key = "DEFEND"
	case HILL1:
		key = "HILL"
	case RALLY:
		key = "RALLY"
	case FOOD:
		key = "FOOD"
	case WAYPOINT:
		key = "WAYPOINT"
	case EXPLORE:
		key = "EXPLORE"
	}

	if key == "" {
		return 0
	}

	val, ok := (*p.PriorityMap)[key]

	if !ok {
		log.Printf("Priority missing for key \"%s\"", key)
		return 0
	}
	return val
}

func (p *Parameters) MakePriMap() *[256]int {
	out := [256]int{}
	for i := Item(0); i < MAX_ITEM; i++ {
		out[i] = p.Priority(i)
	}

	return &out
}
