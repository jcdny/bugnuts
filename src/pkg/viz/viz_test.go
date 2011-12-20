// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package main

import (
	"log"
	"testing"
)

func TestVizLine(t *testing.T) {
	file := "testdata/fill2.map" // fill.2 Point{r:4, c:5}
	m, _ := MapLoadFile(file)
	log.Printf("test")
	VizLine(m, Point{2, 2}, Point{6, 10}, false)
	VizLine(m, Point{1, 8}, Point{10, 2}, false)
}
