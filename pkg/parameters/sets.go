package parameters

var defaultPriMap = map[string]int{"DEFEND": 15, "HILL": 10, "RALLY": 15, "FOOD": 20, "EXPLORE": 30, "WAYPOINT": 30}

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
