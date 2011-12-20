// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package bench

import (
	"log"
	"testing"
	"reflect"
	"bugnuts/engine"
	. "bugnuts/util"
	. "bugnuts/torus"
	. "bugnuts/combat"
)

type testData struct {
	in, user string
	turn     int
	part     Point
}

var tests = []testData{
	//{"testdata/test1/0.replay", "bot8", 30, Point{28, 52}},  // 6 me 2 other
	{"testdata/test1/0.replay", "bot8", 30, Point{49, 19}}, // 5 me, 3, 1
	{"../pathing/testdata/replay.big.json.gz", "MomoBot", 300, Point{75, 102}},
}

func TestPartitions(t *testing.T) {
	for _, td := range tests {
		log.Print("Running test for ", td.in, "(", td.user, ") turn ", td.turn, " partition ", td.part)
		combatMe(td)
	}
}

/*  see mystery.go for the what this tests...
 *

func TestCombatSetup(t *testing.T) {
	td := tests[1]
	s := engine.LoadState(td.in, td.user, td.turn)
	if s != nil {
		c := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		c.Setup(s.Ants)
		c1 := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		c1.SetupS1(s.Ants)
		c2 := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		c2.SetupS2(s.Ants)

		log.Print(reflect.DeepEqual(c, c1), reflect.DeepEqual(c1, c2), reflect.DeepEqual(c2, c))
	}
}

func BenchmarkCombatSetupS1(b *testing.B) {
	b.StopTimer()
	t := tests[1]
	s := engine.LoadState(t.in, t.user, t.turn)
	if s != nil {
		c := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		c.SetupS1(s.Ants)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c.SetupS1(s.Ants)
		}
	}
}
func BenchmarkCombatSetupS2(b *testing.B) {
	b.StopTimer()
	t := tests[1]
	s := engine.LoadState(t.in, t.user, t.turn)
	if s != nil {
		c := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		c.SetupS2(s.Ants)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c.SetupS2(s.Ants)
		}
	}
}
*/

func BenchmarkCombatSetup(b *testing.B) {
	b.StopTimer()
	t := tests[1]
	s := engine.LoadState(t.in, t.user, t.turn)
	if s != nil {
		c := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		c.Setup(s.Ants)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c.Setup(s.Ants)
		}
	}
}

func combatMe(t testData) {
	s := engine.LoadState(t.in, t.user, t.turn)
	if s != nil {
		c := NewCombat(s.Map, s.AttackMask, 10) // TODO player counts?
		c.SetupS2(s.Ants)

		TPush("** Partition " + t.in)
		ap, pmap := CombatPartition(s)
		TPop()

		log.Printf("Found %d partitions with %d mappings", len(ap), len(pmap))
		pkey := s.ToLocation(t.part)
		p, ok := ap[pmap[pkey][0]]
		if !ok {
			log.Print("Could not find partition for ", t.part)
		} else {
			ps := NewPartitionState(c, p)
			log.Print("Partition ", pmap[pkey][0], ":", DumpPartitionState(ps))
		}
	}
}
