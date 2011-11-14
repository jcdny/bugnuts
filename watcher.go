package main

import (
	"os"
	"strings"
	"strconv"
	"math"
	"log"
)

type Watch struct {
	S      string
	Start  int
	End    int
	R      int
	C      int
	N      int
	Player int
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

func NewWatches(rows, cols, turns int) *Watches {
	ws := Watches{Rows: rows, Cols: cols, Turns: turns}
	return &ws
}

func (ws *Watches) Load(wlist []string) {
	for _, s := range wlist {
		w, err := ws.Parse(s)
		if err != nil {
			log.Printf("Problem loading watches: %v", err)
		} else {
			ws.Add(w)
		}
	}
}

func (ws *Watches) Watched(l Location, turn int, player int) bool {
	if len(ws.W) == 0 {
		return false
	}
	if ws.dirty {
		ws.update()
	}

	// Fast check for turns for all locations
	if turn >= 0 && ws.turns[turn] {
		return true
	}
	// and locations for all turns
	if ws.locs[l] {
		return true
	}

	if turn >= 0 {
		r, c := int(l)/ws.Cols, int(l)%ws.Cols
		for _, w := range ws.wturns[turn] {
			// only check for locations in a region since
			// above we checked for locations with no turn
			// restriction
			d := c - w.C
			if d < 0 {
				d = -d
			}
			if d > (ws.Cols+1)/2 {
				d = ws.Cols - d
			}
			if d > w.N {
				continue
			}
			d = r - w.R
			if d < 0 {
				d = -d
			}
			if d > (ws.Rows+1)/2 {
				d = ws.Rows - d
			}
			if d > w.N {
				continue
			}
			return true
		}
	}

	return false
}

func (ws *Watches) update() {
	if !ws.dirty {
		return
	}

	// New cache
	ws.turns = make([]bool, ws.Turns+1)
	ws.locs = make([]bool, ws.Rows*ws.Cols)
	ws.wturns = make([][]*Watch, ws.Turns+1)

	for _, w := range ws.W {
		if w.N < 0 {
			// Turn watches for all points
			for i := w.Start; i <= MinV(w.End, ws.Turns); i++ {
				ws.turns[i] = true
			}
		} else if w.Start == 0 && w.End > ws.Turns {
			// Location watches for all turns
			for r := w.R - w.N; r <= w.R+w.N; r++ {
				var ro int
				if r < 0 || r >= ws.Rows {
					ro = ((r + ws.Rows) % ws.Rows) * ws.Cols
				} else {
					ro = r * ws.Cols
				}
				for c := w.C - w.N; c <= w.C+w.N; c++ {
					var co int
					if c < 0 || c >= ws.Cols {
						co = (c + ws.Cols) % ws.Cols
					} else {
						co = c
					}
					ws.locs[ro+co] = true
				}
			}
		} else {
			// turn and location restricted watches.
			for i := w.Start; i <= MinV(ws.Turns, w.End); i++ {
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

func (ws *Watches) Parse(s string) (w *Watch, err os.Error) {
	w = &Watch{S: s, N: -1, Player: -1}

	s = strings.Replace(s, " ", "", -1)
	tok := strings.Split(s, "@")
	if len(tok) == 0 || len(s) == 0 {
		return nil, os.NewError("Empty watch string")
	}

	turns := strings.Split(tok[0], ":")
	if len(turns) == 1 {
		if len(turns[0]) == 0 {
			w.Start = 0
			w.End = math.MaxInt32
		} else {
			w.Start, err = strconv.Atoi(turns[0])
			if err != nil {
				return nil, err
			}
			w.End = w.Start
		}
	}
	if len(turns) > 1 {
		if len(turns[0]) == 0 {
			w.Start = 0
		} else {
			w.Start, err = strconv.Atoi(turns[0])
			if err != nil {
				return nil, err
			}
		}
		if len(turns[1]) == 0 {
			w.End = math.MaxInt32
		} else {
			w.End, err = strconv.Atoi(turns[1])
			if err != nil {
				return nil, err
			}
		}
	}
	if len(turns) > 2 {
		w.Player, err = strconv.Atoi(turns[2])
		if err != nil {
			return nil, err
		}
	}
	if len(turns) > 3 {
		return nil, os.NewError("Too many turn/player parameters \"" + s + "\"")
	}

	if len(tok) > 1 {
		locs := strings.Split(tok[1], ",")
		if len(locs) < 2 {
			return nil, os.NewError("Too few location parameters \"" + s + "\"")
		}
		w.R, err = strconv.Atoi(locs[0])
		w.C, err = strconv.Atoi(locs[1])
		if len(locs) > 2 {
			w.N, err = strconv.Atoi(locs[2])
			if err != nil {
				return nil, err
			}
		} else {
			w.N = 0
		}
		if len(locs) > 3 {
			return nil, os.NewError("Too many location parameters \"" + s + "\"")
		}
	}

	if w.End < w.Start {
		w.End, w.Start = w.Start, w.End
	}

	return w, nil
}
