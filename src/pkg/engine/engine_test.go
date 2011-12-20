// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package engine

import (
	"log"
	"testing"
	"bugnuts/replay"
	"bugnuts/maps"
	"fmt"
	"os"
)

func TestEngine(t *testing.T) {
	match, err := replay.Load("testdata/0.replay")

	if err != nil || match == nil {
		log.Panicf("Error loading replay: %v", err)
	}
	m := match.GetMap()

	g := NewGame(&match.GameInfo, m)

	g.Replay(match.Replay, 0, match.GameLength, true)

	for i := range g.PlayerInput[0] {
		filein := fmt.Sprint("testdata/0.bot", i, ".input")
		fileout := filein + ".tmp"
		out, err := os.Create(fileout)
		defer out.Close()
		if err != nil {
			log.Panic("open failed for ", fileout, ":", err)
		}
		fmt.Fprintf(out, "turn 0\n%v\nready\n", g.GameInfo)
		for turn := range g.PlayerInput {
			if len(g.PlayerInput[turn][i].A) > 0 && turn < match.GameLength {
				fmt.Fprint(out, "turn ", turn+1, "\n")
			} else {
				fmt.Fprint(out, "end\n")
			}
			fmt.Fprint(out, g.PlayerInput[turn][i], "\ngo\n")
			if len(g.PlayerInput[turn][i].A) == 0 {
				break
			}
		}
	}
}

func getMatch(file string) (match *replay.Match, m *maps.Map, vm, am *maps.Mask) {
	var err os.Error
	match, err = replay.Load(file)
	if err != nil || match == nil {
		log.Panic("Error loading replay", file, ":", err)
	}
	m = match.GetMap()
	gi := &match.GameInfo
	vm = maps.MakeMask(gi.ViewRadius2, gi.Rows, gi.Cols)
	am = maps.MakeMask(gi.AttackRadius2, gi.Rows, gi.Cols)
	m.OffsetsCachePopulateAll(vm)
	m.OffsetsCachePopulateAll(am)

	return
}

// BenchmarkEngine times turn generation for a replay file.
func BenchmarkEngine(b *testing.B) {
	match, m, vm, am := getMatch("testdata/0.replay")

	for i := 0; i < b.N; i++ {
		g := NewGameMasks(&match.GameInfo, m, vm, am)
		g.Replay(match.Replay, 0, match.GameLength, false)
	}
}

// BenchmarkEngineOrdered times the generation in cannonical order to reproduce the 
// python input files.
func BenchmarkEngineOrdered(b *testing.B) {
	match, m, vm, am := getMatch("testdata/0.replay")

	for i := 0; i < b.N; i++ {
		g := NewGameMasks(&match.GameInfo, m, vm, am)
		g.Replay(match.Replay, 0, match.GameLength, true)
	}
}
