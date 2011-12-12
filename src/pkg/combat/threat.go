package combat

import (
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/pathing"
)

// Threatfill returns a fill with the distance to the threat surface
func ThreatFill(m *Map, gthreat []int, pthreat []int, maxdepth uint16, size int) *Fill {
	tb := ThreatBorder(m, gthreat, pthreat, size)
	torig := make(map[Location]int, len(tb))
	for _, loc := range tb {
		torig[loc] = 1
	}
	ft := NewFill(m)
	ft.MapFillSeedMD(torig, 1, maxdepth)

	return ft
}

// ThreatBorder returns a list of locations where the threat as adjacent to a no threat location.
func ThreatBorder(m *Map, gthreat []int, pthreat []int, size int) []Location {
	if size < 100 {
		size = 100
	}
	surf := make([]Location, 0, size)
	for loc, t := range gthreat {
		if t > 0 && t > pthreat[loc] && m.Grid[loc] != WATER {
			for d := 0; d < 4; d++ {
				nl := m.LocStep[loc][d]
				if gthreat[nl] == pthreat[nl] {
					surf = append(surf, Location(loc))
					break
				}
			}
		}
	}

	return surf
}
