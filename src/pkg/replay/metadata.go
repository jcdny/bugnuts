package replay

import (
	"strconv"
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
	GameId         int
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
	m := maps.NewMap(r.Map.Rows, r.Map.Cols, r.Players)
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
			GameId:       r.GameId,
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
