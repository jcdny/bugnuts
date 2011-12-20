// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package util

func Permutations(n, states int) [][]int {
	nl := make([]int, n)
	for i := range nl {
		nl[i] = states
	}
	return PermuteList(nl)
}

// PermuteList takes a list of state counts and generates the permutations...
func PermuteList(N []int) [][]int {

	n := len(N)
	cnt := make([]int, n)

	nperm := 1
	for i := range N {
		nperm *= N[i]
	}

	out := make([][]int, nperm)
	buf := make([]int, nperm*n)
	for c, i := 0, 0; i < nperm; i++ {
		out[i] = buf[i*n : (i+1)*n]
		for j := 0; j < n; j++ {
			if c > 0 {
				cnt[j] += 1
				c--
			}
			if cnt[j] == N[j] {
				cnt[j] = 0
				c = 1
			}
		}
		c = 1
		copy(out[i], cnt)
	}

	return out
}
