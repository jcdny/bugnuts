// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package util

import (
	"math"
)

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func Max(x []int) int {
	xm := math.MinInt32
	for _, y := range x {
		if y > xm {
			xm = y
		}
	}

	return xm
}
func MaxF64(x []float64) float64 {
	xm := -math.MaxFloat64
	for _, y := range x {
		xm = math.Fmax(xm, y)
	}

	return xm
}

func Min(x []int) int {
	xm := math.MaxInt32
	for _, y := range x {
		if y < xm {
			xm = y
		}
	}

	return xm
}

func MinInt8(x []int8) int8 {
	xm := int8(math.MaxInt8)
	for _, y := range x {
		if y < xm {
			xm = y
		}
	}

	return int8(xm)
}

func MinV(v1 int, vn ...int) (m int) {
	m = v1
	for _, vi := range vn {
		if vi < m {
			m = vi
		}
	}
	return
}

func MaxV(v1 int, vn ...int) (m int) {
	m = v1
	for _, vi := range vn {
		if vi > m {
			m = vi
		}
	}
	return
}
