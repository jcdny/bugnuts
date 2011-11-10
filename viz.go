package main

import (
	"os"
	"fmt"
)

func VizLine(m *Map, p1, p2 Point, arrow bool) {
	ltype := "line"
	if arrow {
		ltype = "arrow"
	}

	if Abs(p1.r-p2.r) > m.Rows/2 {
		if p1.r < m.Rows/2 {
			p2.r -= m.Rows
		} else {
			p2.r += m.Rows
		}
	}
	if Abs(p2.c-p1.c) > m.Cols/2 {
		if p1.c < m.Cols/2 {
			p2.c -= m.Cols
		} else {
			p2.c += m.Cols
		}
	}

	fmt.Fprintf(os.Stdout, "v %s %d %d %d %d\n", ltype, p1.r, p1.c, p2.r, p2.c)
}

func (s *State) Viz() {
	if Viz["horizon"] {
		for _, loc := range s.Map.HBorder {
			p := s.Map.ToPoint(Location(loc))
			fmt.Fprintf(os.Stdout, "v tileBorder %d %d MM\n", p.r, p.c)
		}
	}

	if Viz["threat"] {
		lthreat := int8(10)
		for i, threat := range s.ThreatMap(s.Turn) {
			if threat > 0 {
				if lthreat != threat {
					fmt.Fprintf(os.Stdout, "v setFillColor 255 0 0 %.1f\n", float64(threat)*.2)
					lthreat = threat
				}
				p := s.Map.ToPoint(Location(i))
				fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.r, p.c)
			}
		}
		fmt.Fprintf(os.Stdout, "v setFillColor 0 0 0 1.0\n")
	}

	if Viz["vcount"] {
		lnvis := -1
		for i, nvis := range s.Map.VisCount {
			if nvis > 1 {
				if nvis > 8 {
					nvis = 8
				}
				if nvis != lnvis {
					fmt.Fprintf(os.Stdout, "v setFillColor 255 255 255 %.1f\n", float64(nvis)*.1)
					lnvis = nvis
				}

				p := s.Map.ToPoint(Location(i))
				fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.r, p.c)
			}
		}
		fmt.Fprintf(os.Stdout, "v setFillColor 0 0 0 1.0\n")
	}

	if Viz["monte"] {
		s.VizMCPaths()
	}
}

func (s *State) VizTargets(tset *TargetSet) {
	for loc, target := range *tset {
		p := s.Map.ToPoint(loc)
		fmt.Fprintf(os.Stdout, "v star %d %d .5 1.2 %d true\n", p.r, p.c, target.Count+3)
	}
}

func (s *State) VizMCPaths() {
	if s.Map.MCPaths < 1 {
		return
	}

	for i, val := range s.Map.MCDist {
		if val > 0 {
			vout := val * 64 / (s.Map.MCDistMax + 1)
			fmt.Fprintf(os.Stdout, "v setFillColor %d %d %d %.1f\n",
				heat64[vout].R, heat64[vout].G, heat64[vout].B, .5)
			p := s.Map.ToPoint(Location(i))
			fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.r, p.c)
		}
	}
}

func (s *State) VizMCHillIn() {
	hills := make(map[Location]int, 6)
	for _, loc := range s.HillLocations(0) {
		hills[loc] = 1
	}

	if len(hills) > 0 {
		ants := make([]Location, 0, 100)
		f, _, _ := MapFillSeed(s.Map, hills, 0)

		for i := 1; i < len(s.Ants); i++ {
			for loc, _ := range s.Ants[i] {
				steps := f.Depth[loc] - f.Depth[f.Seed[loc]]
				if steps < 64 {
					ants = append(ants, Location(loc))
				}
			}
		}
		if len(ants) > 0 {
			// do up to 512 paths, but no more than 32 per ant
			paths := 512 / len(ants)
			if paths > 32 {
				paths = 32
			}
			dist := f.MontePathIn(s.Map, ants, paths, 1)
			maxdist := Max(dist)
			for i, val := range dist {
				if val > 0 {
					vout := val * 64 / (maxdist + 1)
					fmt.Fprintf(os.Stdout, "v setFillColor %d %d %d %.1f\n",
						heat64[vout].R, heat64[vout].G, heat64[vout].B, .5)
					p := s.Map.ToPoint(Location(i))
					fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.r, p.c)
				}
			}
		}
	}
}
