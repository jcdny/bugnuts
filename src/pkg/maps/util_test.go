// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package maps

import (
	"testing"
	"rand"
)

func TestPoints(t *testing.T) {
	v := maskCircle(10)

	vl := PointsToOffsets(v, 13)
	vp := LocationsToOffsets(vl.L, 13)
	vl2 := PointsToOffsets(vp.P, 13)

	if len(vl.L) != len(vl2.L) || len(v) != len(vp.P) {
		t.Error("Point length mismatch")
	} else {
		for i, l := range vl.L {
			if l != vl2.L[i] {
				t.Error("Point roundtrip failed %v, %v, %v, %v", v, vl, vp, vl2)
				break
			}
		}
	}
}

func sumDir(x []Direction) int {
	j := Direction(0)
	for _, d := range x {
		j += d
	}

	return int(j)
}

func TestPermutes(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	for i := 0; i < 256; i++ {
		for d := Direction(0); d < NoMovement; d++ {
			p5d := Permute5D(d, rng)
			if len(p5d) != 5 || sumDir(p5d[:]) != 10 {
				t.Error("Permute5D wrong ", p5d)
			}
		}

		p5 := Permute5(rng)
		p4g := Permute4G(rng)
		p4 := Permute4(rng)
		if len(p5) != 5 || sumDir(p5[:]) != 10 ||
			len(p4g) != 5 || sumDir(p4g[:]) != 10 ||
			len(p4) != 4 || sumDir(p4[:]) != 6 {
			t.Error("Permutes returning invalid array")
			break
		}
	}
}
