package parameters

import ()

var defaultPriMap = map[string]int{"DEFEND": 10, "HILL": 10, "RALLY": 10, "FOOD": 10, "EXPLORE": 25, "WAYPOINT": 20}

var ParameterSets = map[string]*Parameters{
	"V5": &Parameters{
		ExpireFood:  12,
		PriorityMap: &defaultPriMap,
	},
	"V6": &Parameters{
		ExpireFood:     12,
		PriorityMap:    &defaultPriMap,
		Outbound:       80,
		MinHorizon:     20,
		DefendDistance: 16,
	},
}
