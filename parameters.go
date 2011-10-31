package main

type Parameters struct {
	ExpireFood int // If we have not seen the food in this many turns then assume it's gone.
	Priority   map[Item]int
}

var ParameterSets = map[string]*Parameters{
	"V5": &Parameters{
		ExpireFood: 12,
		Priority:   map[Item]int{HILL1: 5, FOOD: 10, EXPLORE: 15},
	},
}
