// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package combat

import (
	"log"
	. "bugnuts/game"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/pathing"
	. "bugnuts/watcher"
)

// Threatfill returns a fill with the distance to the threat surface
func ThreatFill(m *Map, gthreat []int, pthreat []int, maxdepth uint16, size int) (ft *Fill, tb, interior []Location) {
	tb, interior = ThreatBorder(m, gthreat, pthreat, size)
	torig := make(map[Location]int, len(tb))
	for _, loc := range tb {
		torig[loc] = 2
	}
	ft = NewFill(m)
	ft.MapFillSeedMD(torig, 2, maxdepth)

	return
}

func ThreatPathin(f *Fill, ants []map[Location]int) (pathin []int, cutoff int) {
	pathin = make([]int, f.Size())

	for i := range ants {
		for loc := range ants[i] {
			if f.Depth[loc] != 0 {
				pathin[f.Seed[loc]]++
			} else {
				cutoff++
			}
		}
	}
	return
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
				if WS.Watched(loc, np) {
					log.Print("Compute for player ", np, " at ", m.ToPoint(loc), " Tg, Tp: ", Tg[loc], Tp[np][loc])
				}
				if myt != 0 && amask[loc]&PlayerFlag[np] != 0 {
					if _, ok := rm[np][loc]; !ok {
						// only locations where there is any 1 step risk
						mint := 999
						m.ApplyOffsetsBreak(loc, o, func(nl Location) bool {
							if WS.Watched(loc, np) {
								log.Print(m.ToPoint(nl), amask[nl], amask[nl]&PlayerMask[np], PlayerList[amask[nl]&PlayerMask[np]])
							}
							if amask[nl]&PlayerMask[np] != 0 {
								for _, tp := range PlayerList[amask[nl]&PlayerMask[np]] {
									if WS.Watched(loc, np) {
										log.Print(tp, Tg[nl], Tp[tp][nl])
									}
									t := Tg[nl] - Tp[tp][nl]
									if t < mint {
										mint = t
									}
								}
							}
							return mint >= myt
						})
						if WS.Watched(loc, np) {
							log.Print("Compute for player ", np, " at ", loc, " mint: ", mint)
						}
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
