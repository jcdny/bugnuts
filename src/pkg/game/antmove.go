// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package game

import (
	"log"
	. "bugnuts/maps"
	. "bugnuts/torus"
)

type AntMove struct {
	From   Location
	To     Location
	D      Direction
	Player int
}

// AntMove sorted by To then Player
type AntMoveSlice []AntMove

func (p AntMoveSlice) Len() int { return len(p) }
func (p AntMoveSlice) Less(i, j int) bool {
	return p[i].To < p[j].To || (p[i].To == p[j].To && p[i].Player < p[j].Player)
}
func (p AntMoveSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// AntMove sorted by Player then To
type AntMovePlayerSlice []AntMove

func (p AntMovePlayerSlice) Len() int { return len(p) }
func (p AntMovePlayerSlice) Less(i, j int) bool {
	return p[i].Player < p[j].Player || (p[i].Player == p[j].Player && p[i].To < p[j].To)
}
func (p AntMovePlayerSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func DumpAntMove(m *Map, am []AntMove, p int, turn int) {
	for _, a := range am {
		if p == a.Player || p < 0 {
			log.Printf("Move t=%d p=%d %#v %v %#v", turn, a.Player, m.ToPoint(a.From), a.D, m.ToPoint(a.To))
		}
	}
}
