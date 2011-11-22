package viz

import (
	"os"
	"fmt"
	"strings"
	"log"
	. "bugnuts/maps"
	. "bugnuts/util"
	. "bugnuts/state"
	. "bugnuts/pathing"
)

var Viz = map[string]bool{
	"path":    false,
	"vcount":  false,
	"horizon": false,
	"threat":  false,
	"error":   false,
	"targets": false,
	"monte":   false,
	"sym":     false,
}

func SetViz(vizList string, Viz map[string]bool) {
	if vizList != "" {
		for _, word := range strings.Split(strings.ToLower(vizList), ",") {
			switch word {
			case "all":
				for flag, _ := range Viz {
					Viz[flag] = true
				}
			case "none":
				for flag, _ := range Viz {
					Viz[flag] = false
				}
			case "useful":
				Viz["path"] = true
				Viz["horizon"] = true
				Viz["targets"] = true
				Viz["error"] = true
				Viz["monte"] = true
			default:
				_, ok := Viz[word]
				if !ok {
					log.Printf("Visualization flag %s not known", word)
				} else {
					Viz[word] = true
				}
			}
		}
	}
}

func VizLine(m *Map, p1, p2 Point, arrow bool) {
	ltype := "line"
	if arrow {
		ltype = "arrow"
	}

	if Abs(p1.R-p2.R) > m.Rows/2 {
		if p1.R < m.Rows/2 {
			p2.R -= m.Rows
		} else {
			p2.R += m.Rows
		}
	}
	if Abs(p2.C-p1.C) > m.Cols/2 {
		if p1.C < m.Cols/2 {
			p2.C -= m.Cols
		} else {
			p2.C += m.Cols
		}
	}

	fmt.Fprintf(os.Stdout, "v %s %d %d %d %d\n", ltype, p1.R, p1.C, p2.R, p2.C)
}

func Visualize(s *State) {
	if Viz["horizon"] {
		for _, loc := range s.Met.HBorder {
			p := s.ToPoint(Location(loc))
			fmt.Fprintf(os.Stdout, "v tileBorder %d %d MM\n", p.R, p.C)
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
				p := s.ToPoint(Location(i))
				fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.R, p.C)
			}
		}
		fmt.Fprintf(os.Stdout, "v setFillColor 0 0 0 1.0\n")
	}

	if Viz["vcount"] {
		lnvis := -1
		for i, nvis := range s.Met.VisCount {
			if nvis > 1 {
				if nvis > 8 {
					nvis = 8
				}
				if nvis != lnvis {
					fmt.Fprintf(os.Stdout, "v setFillColor 255 255 255 %.1f\n", float64(nvis)*.1)
					lnvis = nvis
				}

				p := s.ToPoint(Location(i))
				fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.R, p.C)
			}
		}
		fmt.Fprintf(os.Stdout, "v setFillColor 0 0 0 1.0\n")
	}

	if Viz["monte"] {
		VizMCPaths(s)
	}
	if Viz["sym"] {
		log.Printf("Visalizing symmetry")
		m := s.Map
		if len(m.SMap) > 0 {
			for _, item := range []Item{WATER, LAND} {
				if item == WATER {
					fmt.Fprintf(os.Stdout, "v setFillColor 0 0 128 .3\n")
				} else {
					fmt.Fprintf(os.Stdout, "v setFillColor 0 128 0 .3\n")
				}
				for i, gitem := range m.Grid {
					if item == gitem && m.TGrid[i] != gitem {
						p := s.ToPoint(Location(i))
						fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.R, p.C)
					}
				}
			}
		}
	}
}

func VizTargets(s *State, tset *TargetSet) {
	for loc, target := range *tset {
		p := s.ToPoint(loc)
		fmt.Fprintf(os.Stdout, "v star %d %d .3 1 %d true\n", p.R, p.C, target.Count+2)
	}
}

func VizMCPaths(s *State) {
	if s.Met.MCPaths < 1 {
		return
	}

	for i, val := range s.Met.MCDist {
		if val > 0 {
			vout := val * 64 / (s.Met.MCDistMax + 1)
			if val == s.Met.MCDistMax {
				fmt.Fprintf(os.Stdout, "v setFillColor %d %d %d %.1f\n",
					0, 0, 255, .75)
			} else {
				fmt.Fprintf(os.Stdout, "v setFillColor %d %d %d %.1f\n",
					heat64[vout].R, heat64[vout].G, heat64[vout].B, .4)
			}
			p := s.ToPoint(Location(i))
			fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.R, p.C)
		}
	}
}

func VizMCHillIn(s *State) {
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
			dist, _ := f.MontePathIn(s.Map, ants, paths, 1)
			maxdist := Max(dist)
			for i, val := range dist {
				if val > 0 {
					vout := val * 64 / (maxdist + 1)
					if val == maxdist {
						fmt.Fprintf(os.Stdout, "v setFillColor %d %d %d %.1f\n",
							0, 0, 255, .75)
					} else {
						fmt.Fprintf(os.Stdout, "v setFillColor %d %d %d %.1f\n",
							heat64[vout].R, heat64[vout].G, heat64[vout].B, .5)
					}
					p := s.ToPoint(Location(i))
					fmt.Fprintf(os.Stdout, "v tile %d %d\n", p.R, p.C)
				}
			}
		}
	}
}
