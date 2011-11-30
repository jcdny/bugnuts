package combat

import (
	. "bugnuts/maps"
)

type ALoc uint8

type Arena struct {
	Grid    [256]Item
	LocStep [256][256]ALoc
}

// Extract an arena map from a full map.
func NewArena(m *Map, loc Location) *Arena {
	a := &Arena{}

	p := m.ToPoint(loc)

	for r := 0; r < 16; r++ {
		for c := 0; c < 16; c++ {
			al := ALoc(r*16 + c)
			ml := m.ToLocation(m.PointAdd(p, Point{r, c}))
			a.Grid[al] = m.Grid[ml]

			// Generate LocStep with special casing of broders.
			for d, step := range Steps {
				var alstep ALoc
				rstep := r + step.R
				cstep := c + step.C
				// Wrap if we need to
				if rstep < 0 || rstep >= 16 || cstep < 0 || cstep >= 16 {
					if StepableItem[m.Grid[m.LocStep[ml][d]]] {
						alstep = 0
					} else {
						alstep = 255
					}
				} else {
					alstep = ALoc(rstep*16 + cstep)
				}
				a.LocStep[al][d] = alstep
			}
		}
	}
	// above we mapped border stepable to 0 and blocked to 255.
	a.Grid[0] = LAND
	a.Grid[255] = WATER

	return a
}

func (a *Arena) String() string {
	s := ""
	for r := 0; r < 16; r++ {
		for c := 0; c < 16; c++ {
			s += string(a.Grid[r*16+c].ToSymbol())
		}
		s += "\n"
	}

	return s
}
