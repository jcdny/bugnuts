package combat

import (
	"log"
	. "bugnuts/maps"
)

type ALoc uint8

type Arena struct {
	Grid    [256]Item
	LocStep [256][4]ALoc
	NFree   [256]uint8
}

// Extract an arena map from a full map.
func (m *Map) NewArena(loc Location) *Arena {
	a := &Arena{}

	p := m.ToPoint(loc)

	for r := 0; r < 16; r++ {
		for c := 0; c < 16; c++ {
			al := ALoc(r*16 + c)
			ml := m.ToLocation(m.PointAdd(p, Point{r, c}))
			a.Map[al] = m.Grid[ml]

			// fill in LocStep
			if r == 0 || c == 0 || r == 15 || c == 15 {
				// special case borders
			}
		}
	}
}

func (a *Arena) String() string {
	s := ""
	for r := 0; r < 16; r++ {
		for c := 0; c < 16; c++ {
			s += string(a[r*16+c].ToSymbol())
		}
		s += "\n"
	}

	return s
}
