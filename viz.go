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
}

func (s *State) VizTargets(tset *TargetSet) {
	for loc, target := range *tset {
		p := s.Map.ToPoint(loc)
		fmt.Fprintf(os.Stdout, "v star %d %d .5 1.2 %d true\n", p.r, p.c, target.Count+3)
	}
}
