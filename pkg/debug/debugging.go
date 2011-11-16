package debug

import (
	"log"
	. "bugnuts/watcher"
)

const (
	DBG_TurnTime = iota + 1
	DBG_AllTime
	DBG_Iterations
	DBG_Threat
	DBG_Targets
	DBG_MoveErrors
	DBG_Stepping
	DBG_Targeting
	DBG_Combat
	DBG_PathIn
	DBG_Movement
	DBG_Sample
	DBG_BorderTargets
	DBG_PathInDefense
	DBG_Start
	DBG_Turns
	DBG_Results
	maxDBG
)

var Debug = make([]bool, maxDBG)
var WS *Watches

// DebugLevels are cumulative, but a 0 resets all to false and a -Level unsets.
var debugLevels = [][]int{
	0: []int{0},
	1: []int{DBG_TurnTime},
	2: []int{DBG_AllTime, DBG_Iterations},
	3: []int{DBG_Stepping},
}

func SetDebugLevel(dlev int) {
	if dlev > len(debugLevels) {
		log.Panicf("Max defined debug level is %d", len(debugLevels))
	}
	for d := 0; d <= dlev; d++ {
		for _, dbg := range debugLevels[d] {
			if dbg == 0 {
				for i, _ := range Debug {
					Debug[i] = false
				}
			} else if dbg > 0 {
				Debug[dbg] = true
			} else {
				Debug[dbg] = false
			}
		}
	}
}
