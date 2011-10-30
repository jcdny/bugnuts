package replay

import (
	"json"
	"os"
	"strings"
	"strconv"
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
	var err os.Error

	str := strings.TrimRight(strings.TrimLeft(string(b), "[ "), " ]")
	item := strings.SplitN(str, ", ", 6)

	if len(item) != 6 {
		return os.NewError("Invalid AntHistory JSON:" + string(b))
	}

	a.Row, err = strconv.Atoi(item[0])
	if err != nil {
		return err
	}
	a.Col, err = strconv.Atoi(item[1])
	if err != nil {
		return err
	}
	a.Start, err = strconv.Atoi(item[2])
	if err != nil {
		return err
	}
	a.End, err = strconv.Atoi(item[3])
	if err != nil {
		return err
	}
	a.Player, err = strconv.Atoi(item[4])
	if err != nil {
		return err
	}
	a.Moves = strings.Trim(item[5], "\"")

	return err
}

func (h *HillData) UnmarshalJSON(b []byte) os.Error {
	var err os.Error

	str := strings.TrimRight(strings.TrimLeft(string(b), "[ "), " ]")
	item := strings.SplitN(str, ", ", 5)

	if len(item) != 4 {
		return os.NewError("Invalid HillData JSON:" + string(b))
	}

	h.Row, err = strconv.Atoi(item[0])
	if err != nil {
		return err
	}
	h.Col, err = strconv.Atoi(item[1])
	if err != nil {
		return err
	}
	h.Player, err = strconv.Atoi(item[2])
	if err != nil {
		return err
	}
	h.Razed, err = strconv.Atoi(item[3])
	if err != nil {
		return err
	}

	return err
}

func (f *FoodHistory) UnmarshalJSON(b []byte) os.Error {
	var err os.Error

	str := strings.TrimRight(strings.TrimLeft(string(b), "[ "), " ]")
	item := strings.SplitN(str, ", ", 5)

	if len(item) < 4 || len(item) > 5 {
		return os.NewError("Invalid FoodHistory JSON: " + string(b))
	}

	f.Row, err = strconv.Atoi(item[0])
	if err != nil {
		return err
	}
	f.Col, err = strconv.Atoi(item[1])
	if err != nil {
		return err
	}
	f.Spawn, err = strconv.Atoi(item[2])
	if err != nil {
		return err
	}
	f.Gather, err = strconv.Atoi(item[3])
	if err != nil {
		return err
	}
	if len(item) == 5 {
		// the player collecting can be omitted
		f.Player = new(int)
		*f.Player, err = strconv.Atoi(item[4])
		if err != nil {
			return err
		}
	}

	return err
}
