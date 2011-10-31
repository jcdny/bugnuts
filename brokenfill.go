// An attempt to make flood fill faster by not using as much stack.
// not enough brain to make it work right now.  Also turns out to be
// premature opt. since the fill seems to take < 5ms 

// Generate a fill from Map m return fill slice, max Q size, max depth
func ExpMapFill(m *Map, origin Point) (*Fill, int, int) {

	newDepth := uint16(1) // dont start with 0 since 0 means not visited.
	safe := 0

	Directions := []Point{{0, -1}, {-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	Diagonals := []Point{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}, {-1, 1}}

	f := m.NewFill()

	q := QNew(100) // TODO think more about q cap

	q.Q(origin)
	f.Depth[m.ToLocation(origin)] = newDepth

	for !q.Empty() {
		p := q.DQ()

		Depth := f.Depth[m.ToLocation(p)]
		newDepth := Depth + 1

		validlast := false

		for i, s := range Diagonals {
			// on a given diagonal just go until we
			// stop finding same depth
			for {
				// Debug lets not infinite loop
				if safe++; safe > 1000 {
					log.Panicf("Oh No Crazytime")
				}

				fillp := m.PointAdd(p, Directions[i])
				nloc := m.ToLocation(fillp)

				// log.Printf("p: %v np: %v i: %d item: %c d: %d", p, fillp, i, m.Grid[nloc].ToSymbol(), f.Depth[nloc])
				//log.Printf("%s", PrettyFill(m, f, p, fillp, q, Depth))

				if m.Grid[nloc] != WATER && f.Depth[nloc] == 0 {
					f.Depth[nloc] = newDepth
					// Queue a new start point
					if !validlast {
						q.Q(fillp)
						//log.Printf("Q %v", fillp)
					}
					validlast = true
				} else {
					validlast = false
				}

				np := m.PointAdd(p, s)
				nloc = m.ToLocation(np)
				//log.Printf("p %v np %v fillp %v %v", p, np, fillp, f)

				if f.Depth[nloc] == Depth {
					p = np
				} else {
					break
				}
			}
		}
	}

	return f, 0, int(newDepth - 1)

}
