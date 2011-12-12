package viz

import (
	"os"
	"fmt"
	"strings"
	"log"
	. "bugnuts/maps"
	. "bugnuts/torus"
	. "bugnuts/util"
	. "bugnuts/state"
	. "bugnuts/pathing"
	. "bugnuts/combat"
)

var Viz = map[string]bool{
	"path":    false,
	"goals":   false,
	"vcount":  false,
	"horizon": false,
	"threat":  false,
	"error":   false,
	"targets": false,
	"monte":   false,
	"sym":     false,
	"combat":  false,
	"tborder": false,
}

func SetViz(vizList string, Viz map[string]bool) {
	if vizList != "" {
		for _, word := range strings.Split(strings.ToLower(vizList), ",") {
			switch word {
			case "all":
				for flag := range Viz {
					Viz[flag] = true
				}
			case "none":
				for flag := range Viz {
					Viz[flag] = false
				}
			case "useful":
				Viz["path"] = true
				Viz["goals"] = true
				Viz["combat"] = true
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

func VizPath(m *Map, p Point, steps string, color int) {
	if color == 1 {
		slc(cRed, .5)
	} else if color == 2 {
		slc(cGreen, .5)
	}
	fmt.Fprintf(os.Stdout, "v rp %d %d %s\n", p.R, p.C, steps)
	if color > 0 {
		slc(cBlack, 1.0)
	}
}

func VizLine(m *Map, p1, p2 Point, arrow bool) {
	ltype := "l"
	if arrow {
		ltype = "a"
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
			fmt.Fprintf(os.Stdout, "v tb %d %d MM\n", p.R, p.C)
		}
	}

	if Viz["threat"] {
		lthreat := 10
		for i, threat := range s.C.PThreat[0] {
			if threat > 0 {
				if lthreat != threat {
					sfc(cRed, float64(threat)*.2)
					lthreat = threat
				}
				p := s.ToPoint(Location(i))
				fmt.Fprintf(os.Stdout, "v t %d %d\n", p.R, p.C)
			}
		}
		sfc(cBlack, 1.0)
	}

	if Viz["tborder"] {
		slc(cRed, 1.0)
		for _, loc := range ThreatBorder(s.C.Map, s.C.Threat1, s.C.PThreat1[0], 0) {
			p := s.ToPoint(Location(loc))
			fmt.Fprintf(os.Stdout, "v tb %d %d MM\n", p.R, p.C)
		}
		slc(cBlack, 1.0)
	}

	if Viz["vcount"] {
		lnvis := -1
		for i, nvis := range s.Met.VisCount {
			if nvis > 1 {
				if nvis > 8 {
					nvis = 8
				}
				if nvis != lnvis {
					sfc(cWhite, float64(nvis)*.1)
					lnvis = nvis
				}

				p := s.ToPoint(Location(i))
				fmt.Fprintf(os.Stdout, "v t %d %d\n", p.R, p.C)
			}
		}
		sfc(cBlack, 1.0)
	}

	if Viz["monte"] {
		VizMCPaths(s)
	}
	if Viz["sym"] {
		m := s.Map
		if len(m.SMap) > 0 {
			for _, item := range []Item{WATER, LAND} {
				if item == WATER {
					sfc(cBlue, .2)
				} else {
					sfc(cGreen, .2)
				}
				for i, gitem := range m.Grid {
					if item == gitem && m.TGrid[i] != gitem {
						p := s.ToPoint(Location(i))
						fmt.Fprintf(os.Stdout, "v t %d %d\n", p.R, p.C)
					}
				}
			}
		}
	}
}

func VizTargets(s *State, tset *TargetSet) {
	for loc, target := range *tset {
		p := s.ToPoint(loc)
		fmt.Fprintf(os.Stdout, "v s %d %d .3 1 %d true\n", p.R, p.C, target.Count+2)
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
				sfc(cBlue, .75)
			} else {
				sfc(heat64[vout], .4)
			}
			p := s.ToPoint(Location(i))
			fmt.Fprintf(os.Stdout, "v t %d %d\n", p.R, p.C)
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
			for loc := range s.Ants[i] {
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
			dist, _ := f.MontePathIn(s.Rand, ants, paths, 1)
			maxdist := Max(dist)
			for i, val := range dist {
				if val > 0 {
					vout := val * 64 / (maxdist + 1)
					if val == maxdist {
						sfc(cBlue, .75)
					} else {
						sfc(heat64[vout], .5)
					}
					p := s.ToPoint(Location(i))
					fmt.Fprintf(os.Stdout, "v t %d %d\n", p.R, p.C)
				}
			}
		}
	}
}

func vizCircle(p Point, r float64, fill bool) {
	fmt.Fprintf(os.Stdout, "v c %d %d %f %v\n",
		p.R, p.C, r, fill)
}

func VizFrenemies(s *State, ap Partitions, pmap map[Location]map[Location]struct{}) {
	i := 0
	for ploc, p := range ap {
		for _, loc := range p.Ants {
			//log.Printf("ploc %v loc %v pmap %v", ploc, loc, pmap[loc])
			slc(qual6[i%6], 1)
			p := s.ToPoint(loc)
			if loc == ploc {
				sfc(qual6[i%6], .5)
				vizCircle(p, 1, true)
				sfc(cWhite, 1)
			} else {
				vizCircle(p, 1, false)
			}
			pp := s.ToPoint(ploc)
			VizLine(s.Map, p, pp, false)
		}
		i++
	}
	if i > 0 {
		slc(cBlack, 1.0)
	}
}
