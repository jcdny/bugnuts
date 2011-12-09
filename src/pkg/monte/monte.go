package monte

import (
	"log"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/game"
	. "bugnuts/combat"
)

type AntState struct {
	Start    Location
	End      Location
	NStep    int
	Steps    [8]Direction
	Prefered Direction
}

func MonteDraw(c *Combat, am []AntMove, N int) {
	log.Print("Drawing ", N)
}
