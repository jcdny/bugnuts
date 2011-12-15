package combat

import (
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/pathing"
)

// Threatfill returns a fill with the distance to the threat surface
func ThreatFill(m *Map, gthreat []int, pthreat []int, maxdepth uint16, size int) *Fill {
	tb, interior := ThreatBorder(m, gthreat, pthreat, size)
	torig := make(map[Location]int, len(tb))
	for _, loc := range tb {
		torig[loc] = 2
	}
	ft := NewFill(m)
	ft.MapFillSeedMD(torig, 2, maxdepth)
	for _, loc := range interior {
		ft.Depth[loc] = 1
	}
	return ft
}

// ThreatBorder returns a list of locations where the threat as adjacent to a no threat location.
func ThreatBorder(m *Map, gthreat []int, pthreat []int, size int) (surf []Location, interior []Location) {
	if size < 100 {
		size = 100
	}
	surf = make([]Location, 0, size)
	interior = make([]Location, 0, size)
	for loc, t := range gthreat {
		if t > 0 && t > pthreat[loc] && m.Grid[loc] != WATER {
			var d int
			for d = 0; d < 4; d++ {
				nl := m.LocStep[loc][d]
				if gthreat[nl] == pthreat[nl] {
					surf = append(surf, Location(loc))
					break
				}
			}
			if d == 4 {
				interior = append(interior, Location(loc))
			}
		}
	}

	return
}

const (
	RiskSafe = iota
	RiskAverse
	RiskNeutral
	Suicidal
	MaxRiskStat
)

func (c *Combat) Riskly(Ants []map[Location]int) []map[Location]int {
	return RiskMark(c.Map, &c.AttackMask.Offsets, Ants, c.Ants1, c.Threat1, c.PThreat1)
}

// Generate a list of risk differentiators
func RiskMark(m *Map, o *Offsets, Ants []map[Location]int, amask, Tg []int, Tp [][]int) (rm []map[Location]int) {
	rm = make([]map[Location]int, len(Ants))
	for np := range Ants {
		rm[np] = make(map[Location]int, len(Ants[np])*5)
		for aloc := range Ants[np] {
			for d := 0; d < 5; d++ {
				loc := m.LocStep[aloc][d]
				myt := Tg[loc] - Tp[np][loc]
				if myt != 0 && amask[loc]&PlayerFlag[np] != 0 {
					if _, ok := rm[np][loc]; !ok {
						// only locations where there is any 1 step risk
						mint := 999
						m.ApplyOffsetsBreak(loc, o, func(nl Location) bool {
							if amask[nl]&PlayerMask[np] != 0 {
								for _, tp := range PlayerList[amask[nl]&PlayerMask[np]] {
									t := Tg[nl] - Tp[tp][nl]
									if t < mint {
										mint = t
									}
								}
							}
							return mint >= myt
						})
						switch {
						case mint < myt:
							rm[np][loc] = Suicidal
						case mint == myt:
							rm[np][loc] = RiskNeutral
						case mint > myt:
							rm[np][loc] = RiskAverse
						}
					}
				}
			}
		}
	}

	return
}
