package replay

import (
	"strconv"
	"bugnuts/maps"
	"bugnuts/torus"
	"bugnuts/state"
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

func (r *Replay) GetGameInfo() *state.GameInfo {
	return &r.GameInfo
}

func (r *Replay) AntCount(turns int) [][]int {
	// count the ants per turn
	nants := make([][]int, r.Players)
	for _, a := range r.Ants {
		if len(nants[a.Player]) == 0 {
			nants[a.Player] = make([]int, turns+1)
		}
		for i := a.Start; i < a.End; i++ {
			nants[a.Player][i]++
		}
	}

	return nants
}

// Return ant locations l[turn][player][ant]
func (r *Replay) AntLocations(m *maps.Map, turns int) [][][]torus.Location {
	nants := r.AntCount(turns)

	// Allocate the slices
	al := make([][][]torus.Location, turns+1)
	for turn := 0; turn <= turns; turn++ {
		al[turn] = make([][]torus.Location, r.Players)
		for np := 0; np < r.Players; np++ {
			if nants[np][turn] > 0 {
				al[turn][np] = make([]torus.Location, 0, nants[np][turn])
			}
		}
	}
	// Populate the array
	for _, a := range r.Ants {
		turn := a.Start
		loc := m.ToLocation(torus.Point{a.Row, a.Col})
		for _, move := range a.Moves {
			al[turn][a.Player] = append(al[turn][a.Player], loc)
			turn++
			d := maps.ByteToDirection[move]
			switch d {
			case maps.NoMovement:
			case maps.InvalidMove:
			default:
				loc = m.LocStep[loc][d]
			}
		}
		al[turn][a.Player] = append(al[turn][a.Player], loc)
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
