package foo

func (s *State) Threat(turn int, l Location) int8 {
	i := len(s.Met.Threat) - turn + s.Turn - 1
	if i < 0 {
		log.Printf("Threat for turn %d on turn %d we only keep %d turns", turn, s.Turn, len(s.Met.Threat))
		return 0
	}
	return s.Met.Threat[i][l]
}

func (s *State) PrThreat(turn int, l Location) int {
	i := len(s.Met.PrThreat) - turn + s.Turn - 1
	if i < 0 {
		log.Printf("Threat for turn %d on turn %d we only keep %d turns", turn, s.Turn, len(s.Met.Threat))
		return 0
	}
	return s.Met.PrThreat[i][l]
}

func (s *State) ThreatMap(turn int) []int8 {
	i := len(s.Met.Threat) - turn + s.Turn - 1
	if i < 0 {
		log.Printf("Threat for turn %d on turn %d we only keep %d turns", turn, s.Turn, len(s.Met.Threat))
		return nil
	}
	return s.Met.Threat[i]
}

// Compute the threat for N turns out (currently only n = 0 or 1)
// if player > -1 then sum players not including player
func (s *State) ComputeThreat(turn, player int, mask []*MoveMask, threat []int8, pthreat []int) {
	if turn > 1 || turn < 0 {
		log.Panicf("Illegal turns out = %d", turn)
	}

	if len(threat) != s.Rows*s.Cols || len(threat) != len(pthreat) {
		log.Panic("ComputeThreat slice size mismatch")
	}

	var mythreat []int8
	if player >= 0 && turn > 0 && s.Testing {
		mythreat = make([]int8, s.Map.Size())
		for loc := range s.Ants[player] {
			p := s.ToPoint(loc)
			m := mask[s.FreedomKey(loc)]
			for _, op := range m.P {
				mythreat[s.ToLocation(s.PointAdd(p, op))]++
			}
		}
	}

	m := mask[0] // for 0 turns out we just use the 0 degree of freedom mask.
	for i := range s.Ants {
		if i != player {
			for loc := range s.Ants[i] {
				p := s.Map.ToPoint(loc)
				if turn > 0 {
					if false { // crazy or willing to sacrifice or other rules
						m = mask[s.Map.FreedomKey(loc)]
					} else {
						var nsup [4]int8
						for _, op := range mask[0].P {
							l := s.ToLocation(s.PointAdd(p, op))
							if _, ok := s.Ants[i][l]; ok {
								if op.R >= 0 {
									nsup[South]++
								}
								if op.R <= 0 {
									nsup[North]++
								}
								if op.C >= 0 {
									nsup[East]++
								}
								if op.C <= 0 {
									nsup[West]++
								}

							}
						}
						m = mask[s.Map.FreedomKeyThreat(loc, mythreat, nsup)]
					}
				}
				for i, op := range m.P {
					threat[s.ToLocation(s.PointAdd(p, op))]++
					pthreat[s.ToLocation(s.PointAdd(p, op))] += m.MaxPr[i]
				}
			}
		}
	}

	return
}

func (s *State) ResetGrid() {
	// Rotate threat maps and clear first.
	n := len(s.Met.Threat)
	if n > 1 {
		s.Met.Threat = append(s.Met.Threat[1:n], s.Met.Threat[0])
		s.Met.PrThreat = append(s.Met.PrThreat[1:n], s.Met.PrThreat[0])
	}
	for i := range s.Met.Threat[0] {
		s.Met.Threat[0][i] = 0
		s.Met.PrThreat[0][i] = 0
	}

	// Set all seen map to land
	for i, t := range s.Met.Seen {
		s.Met.VisCount[i] = 0
		if t == s.Turn && s.Map.Grid[i] > LAND {
			s.Map.Grid[i] = LAND
		}
	}
}
