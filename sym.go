package main

import (
	_ "log"
)

var symMask [25][8]int // the mask for 4 rotations and 2 mirrorings

func init() {
	pos := [5]int{0, 1, 2, 3, 4}
	neg := [5]int{4, 3, 2, 1, 0}
	// standard
	N, i := 0, uint(0)
	for _, r := range pos {
		for _, c := range pos {
			symMask[r*5+c][N] = 1 << i
			i++
		}
	}
	// ccw 90
	N++
	i = 0
	for _, c := range neg {
		for _, r := range pos {
			symMask[r*5+c][N] = 1 << i
			i++
		}
	}
	// ccw 180
	N++
	i = 0
	for _, r := range neg {
		for _, c := range neg {
			symMask[r*5+c][N] = 1 << i
			i++
		}
	}
	// ccw 270
	N++
	i = 0
	for _, c := range pos {
		for _, r := range neg {
			symMask[r*5+c][N] = 1 << i
			i++
		}
	}
	// mirror vertical
	N++
	i = 0
	for _, r := range pos {
		for _, c := range neg {
			symMask[r*5+c][N] = 1 << i
			i++
		}
	}
	// mirror horizontal
	N++
	i = 0
	for _, r := range neg {
		for _, c := range pos {
			symMask[r*5+c][N] = 1 << i
			i++
		}
	}
	// mirror+rotate
	N++
	i = 0
	for _, c := range pos {
		for _, r := range pos {
			symMask[r*5+c][N] = 1 << i
			i++
		}
	}
	// mirror+rotate
	N++
	i = 0
	for _, c := range neg {
		for _, r := range neg {
			symMask[r*5+c][N] = 1 << i
			i++
		}
	}
}

func (m *Map) SymCompute(loc Location) (int, *[8]int) {
	p := m.ToPoint(loc)
	id := &[8]int{}

	i := 0
	nl := loc
	for r := -2; r < 3; r++ {
		for c := -2; c < 3; c++ {
			if p.r < 2 || p.r > m.Rows-3 || p.c < 2 || p.c > m.Cols-3 {
				nl = m.ToLocation(m.PointAdd(p, Point{r: r, c: c}))
			} else {
				nl = loc + Location(r*m.Cols+c)
			}
			if m.Grid[nl] == UNKNOWN {
				return -1, nil
			}
			if m.Grid[nl] == WATER {
				for rot, mask := range symMask[i] {
					id[rot] ^= mask
				}
			}
			i++
		}
	}

	return Min(id[:]), id
}
