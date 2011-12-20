// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package util

import (
	"testing"
	"math"
	"rand"
	"log"
)

type Lists struct {
	list       []int
	emin, emax int
}

func TestMinMax(t *testing.T) {
	L := []Lists{
		{[]int{1}, 1, 1},
		{[]int{1, 1, 1, 1, 1, 1}, 1, 1},
		{[]int{1, 0, -1}, -1, 1},
		{[]int{}, math.MaxInt32, math.MinInt32},
	}

	for _, l := range L {
		if amin := Min(l.list); amin != l.emin {
			t.Errorf("For min(%v) got %d expected %d", l.list, amin, l.emin)
		}
		if amax := Max(l.list); amax != l.emax {
			t.Errorf("For max(%v) got %d expected %d", l.list, amax, l.emax)
		}
	}
}

func BenchmarkMinV(b *testing.B) {
	L := make([]int, 49)

	for i := range L {
		L[i] = rand.Intn(20)
	}

	for i := 0; i < b.N; i++ {
		MinV(1, L...)
	}
}

func BenchmarkMin(b *testing.B) {
	L := make([]int, 50)

	for i := range L {
		L[i] = rand.Intn(20)
	}

	for i := 0; i < b.N; i++ {
		Min(L)
	}
}

func TestP(t *testing.T) {
	p := Permutations(2, 4)
	log.Print(p)
}
