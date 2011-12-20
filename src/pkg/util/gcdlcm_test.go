// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package util

import (
	"testing"
)

type GcdLcmLists struct {
	n, m, gcd, lcm int
}

var L = []GcdLcmLists{
	{1, 1, 1, 1},
	{12, 24, 12, 24},
	{2, 3, 1, 6},
	{25, 49, 1, 25 * 49},
	{15, 1, 1, 15},
	{20, 4, 4, 20},
}

func TestGcd(t *testing.T) {
	for _, l := range L {
		gcd := Gcd(l.n, l.m)
		if gcd != l.gcd {
			t.Error("Gcd ", gcd, " for ", l)
		}
	}
}
func TestLcm(t *testing.T) {
	for _, l := range L {
		lcm := Lcm(l.n, l.m)
		if lcm != l.lcm {
			t.Error("Lcm ", lcm, " for ", l)
		}
	}
}
