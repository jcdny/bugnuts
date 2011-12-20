// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package util

func Gcd(n, m int) int {
	for m != 0 {
		n, m = m, n%m
	}
	return n
}

func Lcm(n, m int) int {
	return m / Gcd(n, m) * n
}
