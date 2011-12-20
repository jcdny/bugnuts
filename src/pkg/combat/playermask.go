// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package combat

import (
	"log"
	. "bugnuts/game"
)

var PlayerFlag [MaxPlayers]int
var PlayerMask [MaxPlayers]int
var PlayerList [][]int

func init() {
	if MaxPlayers > 31 {
		log.Panic("Unable to support more than 31 players with masks as int")
	}
	for i := range PlayerFlag {
		PlayerFlag[i] = 1 << uint(i)
		PlayerMask[i] = ^PlayerFlag[i]
	}
	nm := 1 << uint(len(PlayerFlag))
	PlayerList = make([][]int, nm)
	buf := make([]int, 0, len(PlayerFlag))
	for m := 0; m < nm; m++ {
		buf = buf[:0]
		for i := 0; i < len(PlayerFlag); i++ {
			if m&PlayerFlag[i] != 0 {
				buf = append(buf, i)
			}
		}
		PlayerList[m] = make([]int, len(buf))
		copy(PlayerList[m], buf)
	}
}
