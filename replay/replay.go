package replay

import (
	"json"
	"os"
)

// [57, 61, 246, 250, 2, "esss"]
// spawn at 57/61 turn 246
type AntHistory struct {
	Row    int
	Col    int
	Start  int
	End    int
	Player int
	Moves  string
}

type MapData struct {
	Cols int
	Rows int
	Data []string
}

type HillData struct {
	Row    int
	Col    int
	Player int // Hill Owner
	Razed  int // Razed on turn

}

type FoodHistory struct {
	Row    int
	Col    int
	Spawn  int
	Gather int
	Player *int
}

type ReplayData struct {
	Ants          []AntHistory
	AttackRadius2 int
	Bonus         []int
	Cutoff        string
	Food          []FoodHistory
	FoodRate      int `json:"food_rate"`
	FoodStart     int `json:"food_start"`
	FoodTurn      int `json:"food_turn"`
	Hills         []HillData
	HiveHistory   [][]int `json:"hive_history"`
	Loadtime      int
	Map           MapData
	PlayerSeed    int64 `json:"player_seed"`
	Players       int
	RankingTurn   int `json:"ranking_turn"`
	Revision      int
	Scores        [][]int
	SpawnRadius2  int
	Turns         int
	Turntime      int
	ViewRadius2   int
	WinningTurn   int `json:"winning_turn"`
}

type Match struct {
	Score        []int
	Challenge    string
	GameId       int `json:"game_id"`
	GameLength   int `json:"game_length"`
	Location     string
	PlayerNames  []string
	PlayerTurns  []int
	Rank         []int
	Replay       *ReplayData `json:"replaydata"`
	ReplayFormat string
	Status       []string
}

func ReplayUnmarshal(b []byte) (*Match, os.Error) {
	m := &Match{}
	err := json.Unmarshal(b[:], m)

	return m, err
}

func (a *AntHistory) UnmarshalJSON(b []byte) os.Error {
	var ah = make([]interface{}, 0, 6)

	err := json.Unmarshal(b, &ah)
	if err != nil {
		return err
	}
	if len(ah) != 6 {
		return os.NewError("Invalid AntHistory JSON:" + string(b))
	}

	a.Row = int(ah[0].(float64))
	a.Col = int(ah[1].(float64))
	a.Start = int(ah[2].(float64))
	a.End = int(ah[3].(float64))
	a.Player = int(ah[4].(float64))
	a.Moves = ah[5].(string)

	return err
}

func (h *HillData) UnmarshalJSON(b []byte) os.Error {
	var ah = make([]interface{}, 0, 5)

	err := json.Unmarshal(b, &ah)
	if err != nil {
		return err
	} else if len(ah) != 4 {
		return os.NewError("Invalid HillData JSON:" + string(b))
	}

	h.Row = int(ah[0].(float64))
	h.Col = int(ah[1].(float64))
	h.Player = int(ah[2].(float64))
	h.Razed = int(ah[3].(float64))

	return err
}

func (f *FoodHistory) UnmarshalJSON(b []byte) os.Error {
	var fa = make([]interface{}, 0, 5)

	err := json.Unmarshal(b, &fa)
	if err != nil {
		return err
	} else if len(fa) < 4 || len(fa) > 5 {
		return os.NewError("Invalid FoodHistory JSON: " + string(b))
	}

	f.Row = int(fa[0].(float64))
	f.Col = int(fa[1].(float64))
	f.Spawn = int(fa[1].(float64))
	f.Gather = int(fa[1].(float64))
	if len(fa) == 5 {
		f.Player = new(int)
		*f.Player = int(fa[1].(float64))
	}

	return err
}
