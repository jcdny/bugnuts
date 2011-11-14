package main

import (
	"os"
	"log"
	"strings"
	"strconv"
)

type Watch struct {
	start int
	end   int
	locs  []Location
}

type Watch struct {
	Start int
	End   int
	R     int
	C     int
	N     int
}

type Watches struct {
	W     []*Watch
	Rows  int // the game characteristics
	Cols  int
	Turns int
	dirty bool
	// the cache of watches reduced to:
	turns  []bool     // turns for all locs
	locs   []bool     // locs for all turns
	wturns [][]*Watch // turn/region restricted watches slice of turns
}

func (ws *Watches) Watched(l Location, t Turn, p player) bool {
	if ws.dirty {
		ws.Update()
	}
	// Fast check for turns for all locations
	if ws.turns[t] {
		return true
	}
	// and locations for all turns

	if ws.locs[l] {
		return true
	}

	for _, w := range ws.wturns[t] {
		// only check for locations in a region since 
		// above we checked for locations with no turn restriction
		if (w.C+ws.cols-c)%ws.cols <= w.N ||
			(w.R+ws.rows-r)%ws.rows <= w.N {
			return true
		}
	}

	return false
}

func (ws *Watches) Update() {
	if !ws.dirty {
		return
	}

	// New cache
	ws.turns = make([]bool, ws.Turns+1)
	ws.locs = make([]bool, ws.Rows*ws.Cols)
	ws.wturns = make([][]*Watch, ws.Turns+1)

	for _, w := range ws.W {
		if w.N < 0 {
			for i := w.Start; i <= w.End; i++ {
				turns[i] = true
			}
		} else if w.Start == 0 && w.End > ws.Turns {
			for r := w.R - w.N; r <= w.R+w.N; r++ {
				for c := w.C - w.N; r <= w.C+w.N; c++ {
					locs[r*ws.Cols+c] = true
				}
			}
		} else {
			for i := w.Start; i <= w.End; i++ {
				ws.wturns[i] = append(ws.wturns[i], w)
			}
		}
	}

	ws.dirty = false
}

func (ws *Watches) Add(w *Watch) {
	ws.W = append(ws.W, w)
	ws.dirty = true
}

func Parse(s string) *Watch {
	w := Watch{s: s, N: -1}

	err := os.EOF

	s = strings.Replace(s, " ", "", -1)
	tok := strings.Split(s, "@")
	if len(tok) == 0 {
		return nil
	}

	turns := strings.Split(tok[0], ":")
	if len(turns) == 1 {
		if len(turns[0]) == 0 {
			w.start = 0
			w.end = math.MaxInt32
		} else {
			w.start, err = strconv.Atoi(turns[0])
			w.end = w.start
		}
	}
	if len(turns) > 1 {
		if len(turns[0]) == 0 {
			w.start = 0
		} else {
			w.start, err = strconv.Atoi(turns[0])
		}
		if len(turns[1]) == 0 {
			w.end = math.MaxInt32
		} else {
			w.start, err = strconv.Atoi(turns[1])
		}
	}
	if len(turns) > 2 {
		w.player, err = strconv.Atoi(turns[2])
	}
	if len(turns) > 3 {
		log.Printf("Invalid watch string %s", s)
	}

	if len(tok) > 1 {
		locs := strings.Split(tok[1], ",")
		if len(locs) < 2 {
			log.Printf("Invalid watch string %s", s)
		}
		ws.R, err = strconv.Atoi(locs[0])
		ws.C, err = strconv.Atoi(locs[1])
		if len(locs) > 2 {
			ws.N, err = strconv.Atoi(locs[2])
		}
	}

	return &w
}
