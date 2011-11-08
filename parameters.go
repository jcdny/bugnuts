package main

type Parameters struct {
	ExpireFood int // If we have not seen the food in this many turns then assume it's gone.
	Priority   map[Item]int
	Outbound   int // Radius inside which we step away from hill by default.
}

var ParameterSets = map[string]*Parameters{
	"V5": &Parameters{
		ExpireFood: 12,
		Priority:   map[Item]int{HILL1: 5, FOOD: 10, EXPLORE: 15},
	},
	"V6": &Parameters{
		ExpireFood: 12,
		Priority:   map[Item]int{DEFEND: 5, HILL1: 5, RALLY: 10, FOOD: 10, EXPLORE: 15, WAYPOINT: 15},
		Outbound:   50,
	},
}

func init() {
	ParameterSets["default"] = ParameterSets["V6"]
}
