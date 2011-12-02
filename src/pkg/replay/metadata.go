package replay

import (
	"strconv"
	"sort"
	"bugnuts/maps"
	"bugnuts/torus"
	"bugnuts/game"
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
		if a.Start <= tmax && a.End >= tmin {
			for i := MaxV(a.Start-tmin, 0); i <= MinV(tmax, a.End)-tmin; i++ {
				nants[a.Player][i]++
			}
		}
	}

	return nants
}

// Return ant locations in array [(turn-tmin)][player][ant] for turns tmin..tmax inclusive
func (r *Replay) AntLocations(m *maps.Map, tmin, tmax int) [][][]torus.Location {
	nants := r.AntCount(tmin, tmax)

	// Allocate the slices
	al := make([][][]torus.Location, tmax-tmin+1)
	for turn := 0; turn <= tmax-tmin; turn++ {
		al[turn] = make([][]torus.Location, r.Players)
		for np := 0; np < r.Players; np++ {
			if nants[np][turn] > 0 {
				al[turn][np] = make([]torus.Location, 0, nants[np][turn])
			}
		}
	}
	// Populate the array
	for _, a := range r.Ants {
		if a.Start <= tmax && a.End >= tmin {
			turn := a.Start
			loc := m.ToLocation(torus.Point{a.Row, a.Col})
			for _, move := range a.Moves {
				if turn+1 > tmax {
					break
				} else if turn >= tmin {
					al[turn-tmin][a.Player] = append(al[turn-tmin][a.Player], loc)
				}
				turn++
				d := maps.ByteToDirection[move]
				if d != maps.InvalidMove {
					loc = m.LocStep[loc][d]
				}
			}
			al[turn-tmin][a.Player] = append(al[turn-tmin][a.Player], loc)
		}
	}

	return al
}

type foodSlice []FoodHistory

func (p foodSlice) Len() int           { return len(p) }
func (p foodSlice) Less(i, j int) bool { return p[i].Row*256+p[i].Col < p[j].Row*256+p[j].Col }
func (p foodSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Return food locations per turn
func (r *Replay) FoodLocations(m *maps.Map, tmin, tmax int) [][]torus.Location {
	food := make([][]torus.Location, tmax-tmin+1)
	sort.Sort(foodSlice(r.Food))

	for _, f := range r.Food {
		gather := f.Gather
		if gather == 0 {
			gather = tmax
		}
		if f.Spawn <= tmax && gather >= tmin {
			loc := m.ToLocation(torus.Point{f.Row, f.Col})
			for i := MaxV(f.Spawn-tmin, 0); i <= MinV(tmax, gather)-tmin; i++ {
				food[i] = append(food[i], loc)
			}
		}
	}
	return food
}

type hillSlice []HillData

func (p hillSlice) Len() int           { return len(p) }
func (p hillSlice) Less(i, j int) bool { return p[i].Row*256+p[i].Col < p[j].Row*256+p[j].Col }
func (p hillSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Return hill locations for the given turns.
// [turn][hill]
func (r *Replay) HillLocations(m *maps.Map, tmin, tmax int) [][]game.PlayerLoc {
	hills := make([][]game.PlayerLoc, tmax-tmin+1)
	// sort here so we get them in location order so we match replay output.
	sort.Sort(hillSlice(r.Hills))
	for _, h := range r.Hills {
		loc := m.ToLocation(torus.Point{h.Row, h.Col})
		razed := h.Razed
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
