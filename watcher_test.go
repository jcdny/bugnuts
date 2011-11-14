package main

import (
	"testing"
	"math"
	"log"
)

func TestWatchParse(t *testing.T) {
	ws := Watches{}
	wtest := map[string]*Watch{
		"1:2@5,7,3": &Watch{Start: 1, End: 2, R: 5, C: 7, N: 3, Player: -1},
		"1:2@5,7":   &Watch{Start: 1, End: 2, R: 5, C: 7, N: 0, Player: -1},
		"1:@5,7":    &Watch{Start: 1, End: math.MaxInt32, R: 5, C: 7, N: 0, Player: -1},
		"@5,7,3":    &Watch{Start: 0, End: math.MaxInt32, R: 5, C: 7, N: 3, Player: -1},
		"@5,7":      &Watch{Start: 0, End: math.MaxInt32, R: 5, C: 7, N: 0, Player: -1},
		"1:10":      &Watch{Start: 1, End: 10, N: -1, Player: -1},
		":10":       &Watch{Start: 0, End: 10, N: -1, Player: -1},
		"5:":        &Watch{Start: 5, End: math.MaxInt32, N: -1, Player: -1},
		"5":         &Watch{Start: 5, End: 5, N: -1, Player: -1},
		"5@1,2":     &Watch{Start: 5, End: 5, R: 1, C: 2, N: 0, Player: -1},
		"10:1":      &Watch{Start: 1, End: 10, N: -1, Player: -1},
		"10:1:2":    &Watch{Start: 1, End: 10, N: -1, Player: 2},
		"10:1:2:5":  nil,
		"":          nil,
		":@1,2,3,4": nil,
		"@1:2":      nil,
		"@1":        nil,
	}

	for s, w := range wtest {
		wp, err := ws.Parse(s)
		if false && err != nil {
			log.Printf("For \"%s\" Error: %v", err, s)
		}
		if w == nil && (wp != nil || err == nil) {
			t.Errorf("Parse fail expected for \"%s\" got %v %#v", s, err, wp)
		} else if w != nil && (wp == nil || err != nil) {
			t.Errorf("Parse failed for \"%s\" %v expected %#v", s, err, w)
		} else if w != nil {
			if wp.Start != w.Start ||
				wp.End != w.End ||
				wp.R != w.R ||
				wp.C != w.C ||
				wp.N != w.N ||
				wp.Player != w.Player {
				t.Errorf("Parse mismatch \"%s\" -> %#v not %#v", s, wp, w)
			}
		}
	}
}

type WTest struct {
	Wlist []string
	N     []int
}

const (
	rows  = 8
	cols  = 11
	all   = rows * cols
	turns = 5
)

func TestWatches(t *testing.T) {

	watches := []WTest{
		{Wlist: []string{"2:2@4,4", "@6,6,1", "4"}, N: []int{9, 9, 10, 9, all, 9}},
		{Wlist: []string{":2", "4:"}, N: []int{all, all, all, 0, all, all}},
		{Wlist: []string{"4:@6,10,2"}, N: []int{0, 0, 0, 0, 25, 25}},
		{Wlist: []string{"4:@6,10,2", "5@5,9,2"}, N: []int{0, 0, 0, 0, 25, 34}},
		{Wlist: []string{"4:@6,10,2", "5@5,9,2"}, N: []int{0, 0, 0, 0, 25, 34}},
		{Wlist: []string{"3:@6,10,2", "3@5,9,2", "1:2@0,0"}, N: []int{0, 1, 1, 34, 25, 25}},
	}

	for ntest, wtest := range watches {
		ws := NewWatches(rows, cols, turns)
		for _, s := range wtest.Wlist {
			w, _ := ws.Parse(s)
			ws.Add(w)
		}
		for turn := 1; turn <= 5; turn++ {
			n := 0
			s := ""
			for r := 0; r < ws.Rows; r++ {
				for c := 0; c < ws.Cols; c++ {
					if ws.Watched(Location(r*ws.Cols+c), turn, 0) {
						n++
						s += "T"
					} else {
						s += "F"
					}
				}
				s += "\n"
			}

			if wtest.N[turn] != n {
				t.Errorf("Watch count mismatch for test %d turn %d got %d expected %d:\n%s", ntest, turn, n, wtest.N[turn], s)
			}
		}
	}
}
