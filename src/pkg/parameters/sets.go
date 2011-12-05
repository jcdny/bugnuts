package parameters

var defaultPriMap = map[string]int{"DEFEND": 15, "HILL": 10, "RALLY": 15, "FOOD": 20, "EXPLORE": 30, "WAYPOINT": 30}

var ParameterSets = map[string]*Parameters{
	"v5": &Parameters{
		ExpireFood:  12,
		PriorityMap: &defaultPriMap,
	},
	"v6": &Parameters{
		ExpireFood:       12,
		PriorityMap:      &defaultPriMap,
		Outbound:         80,
		MinHorizon:       20,
		DefendDistance:   16,
		RiskOffThreshold: .33,
	},
	"v7": &Parameters{
		ExpireFood:       -1,
		PriorityMap:      &defaultPriMap,
		Outbound:         80,
		MinHorizon:       20,
		DefendDistance:   16,
		RiskOffThreshold: .3,
	},
	"v8": &Parameters{
		ExpireFood:       -1,
		PriorityMap:      &defaultPriMap,
		Outbound:         80,
		MinHorizon:       20,
		DefendDistance:   16,
		RiskOffThreshold: .3,
	},
}
