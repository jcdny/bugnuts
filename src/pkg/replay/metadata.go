package replay

import (
	"log"
	"strconv"
	"bugnuts/maps"
	"bugnuts/torus"
	"bugnuts/game"
	"bugnuts/combat"
	. "bugnuts/util"
)

type GameResult struct {
	GameId     int
	Date       string
	GameLength int
	Challenge  string
	MatchupId  *int
	PostId     *int
	WorkerId   string
	Location   string
	MapId      string
}

type PlayerResult struct {
	Game           *GameResult
	PlayerName     string
	PlayerTurns    int
	UserId         *int
	SubmissionId   *int
	Score          int
	Rank           int
	Bonus          int
	Status         string
	ChallengeRank  *int
	ChallengeSkill *float64
}

func (r *Replay) GetMap() *maps.Map {
	m := maps.NewMap(r.Rows, r.Cols, r.Players)
	for r, rdat := range r.Map.Data {
		for c, item := range rdat {
			if maps.ToItem(byte(item)) == maps.WATER {
				m.Grid[r*m.Cols+c] = maps.WATER
			} else {
				m.Grid[r*m.Cols+c] = maps.LAND
			}
		}
	}
	for _, h := range r.Hills {
		loc := m.ToLocation(torus.Point{h.Row, h.Col})
		m.Grid[loc] = maps.MY_HILL + maps.Item(h.Player)
	}

	return m
}

func (r *Replay) GetGameInfo() *game.GameInfo {
	return &r.GameInfo
}

func (r *Replay) AntCount(tmin, tmax int) [][]int {
	// count the ants per turn
	nants := make([][]int, r.Players)
	for i := 0; i < r.Players; i++ {
		nants[i] = make([]int, tmax-tmin+1)
	}
	for _, a := range r.Ants {
		nm := len(a.Moves)
		if false {
			log.Print("a.End, a.Start+nm", a.End, " ", a.Start+nm)
		}
		if a.Start <= tmax && a.Start+nm+1 >= tmin {
			for i := MaxV(a.Start-tmin, 0); i <= MinV(tmax, a.End+1)-tmin; i++ {
				nants[a.Player][i]++
			}
		}
	}

	return nants
}

// Return ant locations in array [(turn-tmin)][player][ant] for turns tmin..tmax inclusive
// Return the spawns [turn][]PlayerLoc
func (r *Replay) AntMoves(m *maps.Map, tmin, tmax int) [][][]combat.AntMove {
	nants := r.AntCount(tmin, tmax)

	// Allocate the slices
	al := make([][][]combat.AntMove, tmax-tmin+1)
	for turn := 0; turn <= tmax-tmin; turn++ {
		al[turn] = make([][]combat.AntMove, r.Players)
		for np := 0; np < r.Players; np++ {
			if nants[np][turn] > 0 {
				al[turn][np] = make([]combat.AntMove, 0, nants[np][turn])
			}
		}
	}
	// Populate the array
	for _, a := range r.Ants {
		if a.Start <= tmax && a.End+1 >= tmin {
			var nloc torus.Location
			turn := a.Start
			loc := m.ToLocation(torus.Point{a.Row, a.Col})
			nloc = -1
			if turn >= tmin {
				// a spawned ant has From = -1, D: InvalidMove
				al[turn-tmin][a.Player] = append(al[turn-tmin][a.Player],
					combat.AntMove{From: -1, To: loc, D: maps.InvalidMove, Player: a.Player})
			}

			for _, move := range a.Moves {
				turn++
				d := maps.ByteToDirection[move]
				if d != maps.InvalidMove {
					nloc = m.LocStep[loc][d]
				} else {
					nloc = -1
				}
				if turn < 5 && a.Player < 2 {
					log.Printf("pl %d turn %d end %d loc %v nloc %v %v", a.Player, turn, a.End, m.ToPoint(loc), m.ToPoint(nloc), d)
				}
				if turn >= tmin && turn <= tmax {
					al[turn-tmin][a.Player] = append(al[turn-tmin][a.Player],
						combat.AntMove{From: loc, To: nloc, D: d, Player: a.Player})
				}
				loc = nloc
				if turn > tmax {
					break
				}
			}

			if turn+1 <= tmax && turn == a.End {
				// An ant that has no further moves.
				al[turn-tmin+1][a.Player] = append(al[turn-tmin+1][a.Player],
					combat.AntMove{From: loc, To: -1, D: maps.InvalidMove, Player: a.Player})
			}
		}
	}

	return al
}

// Return food locations per turn
func (r *Replay) FoodLocations(m *maps.Map, tmin, tmax int) [][]torus.Location {
	food := make([][]torus.Location, tmax-tmin+1)

	for _, f := range r.Food {
		gather := f.Gather - 1
		if f.Spawn <= tmax && gather >= tmin {
			loc := m.ToLocation(torus.Point{f.Row, f.Col})
			for i := MaxV(f.Spawn-tmin, 0); i <= MinV(tmax, gather)-tmin; i++ {
				food[i] = append(food[i], loc)
			}
		}
	}
	return food
}

// Return hill locations for the given turns.
// [turn][hill]
func (r *Replay) HillLocations(m *maps.Map, tmin, tmax int) [][]game.PlayerLoc {
	hills := make([][]game.PlayerLoc, tmax-tmin+1)
	for _, h := range r.Hills {
		loc := m.ToLocation(torus.Point{h.Row, h.Col})
		razed := h.Razed - 1
		if razed < 1 {
			razed = tmax
		}
		if razed >= tmin {
			for i := 0; i < MinV(razed, tmax)-tmin+1; i++ {
				hills[i] = append(hills[i], game.PlayerLoc{Loc: loc, Player: h.Player})
			}
		}
	}
	return hills
}

func (r *Match) ExtractMetadata() (g *GameResult, p []*PlayerResult) {

	g = &GameResult{
		GameId:     r.GameId,
		Date:       r.Date,
		GameLength: r.GameLength,
		Challenge:  r.Challenge,
		MatchupId:  r.MatchupId,
		PostId:     r.PostId,
		WorkerId:   r.WorkerId,
		Location:   r.Location,
		MapId:      r.GetMap().MapId(),
	}

	var uidp, subidp *int

	np := len(r.PlayerNames)
	p = make([]*PlayerResult, np)
	for i := 0; i < np; i++ {

		// Jump through hoops to denote absent fields
		if len(r.UserIds) == np {
			uid, err := strconv.Atoi(r.UserIds[i])
			if err == nil {
				uidp = new(int)
				*uidp = uid
			} else {
				uidp = nil
			}
		}
		if len(r.SubmissionIds) == np {
			subid, err := strconv.Atoi(r.SubmissionIds[i])
			if err == nil {
				subidp = new(int)
				*subidp = subid
			} else {
				subidp = nil
			}
		}

		p[i] = &PlayerResult{
			Game:         g,
			PlayerName:   r.PlayerNames[i],
			PlayerTurns:  r.PlayerTurns[i],
			UserId:       uidp,
			SubmissionId: subidp,
			Score:        r.Score[i],
			Rank:         r.Rank[i],
			Status:       r.Status[i],
		}

		// Again jump through hoops to denote absent fields
		cr := new(int)
		if len(r.ChallengeRank) == np {
			*cr, _ = strconv.Atoi(r.ChallengeRank[i])
		} else {
			cr = nil
		}
		p[i].ChallengeRank = cr

		var cs *float64 = new(float64)
		if len(r.ChallengeSkill) == np {
			*cs, _ = strconv.Atof64(r.ChallengeSkill[i])
		} else {
			cs = nil
		}
		p[i].ChallengeSkill = cs
	}

	return
}
