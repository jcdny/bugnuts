// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package state

import (
	"log"
	"fmt"
	"testing"
	. "bugnuts/maps"
)

func TestRuns(t *testing.T) {
	file := "testdata/run2.map"
	//file := "../maps/testdata/big.map"
	m, _ := MapLoadFile(file)
	met := NewMetrics(m)
	met.UpdateRuns()
	if false {
		log.Printf("%v", m)
		s := "\n"
		for d := 0; d < 4; d++ {
			s += fmt.Sprintf("%d ********\n", d)
			for r := 0; r < m.Rows; r++ {
				for c := 0; c < m.Cols; c++ {
					s += fmt.Sprintf("%d", met.Runs[r*m.Cols+c][d])
				}
				s += "\n"
			}
		}
		log.Print(s)
	}
}

func BenchmarkRuns(b *testing.B) {
	file := "../maps/testdata/maps/cell_maze_p04_01.map"
	m, _ := MapLoadFile(file)
	met := NewMetrics(m)
	for i := 0; i < b.N; i++ {
		met.UpdateRuns()
	}
}
