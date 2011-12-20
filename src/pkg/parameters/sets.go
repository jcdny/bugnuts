// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package parameters

var defaultPriMap = map[string]int{"DEFEND": 15, "HILL": 10, "RALLY": 15, "FOOD": 20, "EXPLORE": 30, "WAYPOINT": 30}

var ParameterSets = map[string]*Parameters{
	"v5": &Parameters{
		ExpireFood:  12,
		PriorityMap: &defaultPriMap,
		Explore:     true,
	},
	"v6": &Parameters{
		ExpireFood:       12,
		PriorityMap:      &defaultPriMap,
		Outbound:         80,
		MinHorizon:       20,
		DefendDistance:   16,
		RiskOffThreshold: .33,
		Explore:          true,
	},
	"v7": &Parameters{
		ExpireFood:       -1,
		PriorityMap:      &defaultPriMap,
		Outbound:         80,
		MinHorizon:       20,
		DefendDistance:   16,
		RiskOffThreshold: .3,
		Explore:          true,
	},
	"v8": &Parameters{
		ExpireFood:       -1,
		PriorityMap:      &defaultPriMap,
		Outbound:         80,
		MinHorizon:       20,
		DefendDistance:   16,
		RiskOffThreshold: .3,
		Explore:          true,
	},
}
