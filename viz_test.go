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
