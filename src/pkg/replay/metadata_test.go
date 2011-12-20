// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package replay

import (
	"testing"
	"log"
	"fmt"
	"os"
	"reflect"
	"bugnuts/maps"
	"bugnuts/torus"
	"bugnuts/combat"
)

func replayGet(file string) (*Match, *maps.Map) {
	match, err := Load(file)
	if err != nil {
		log.Panicf("Load of %s failed: %v", file, err)
	}
	m := match.GetMap()

	return match, m
}

func TestFoodLocations(t *testing.T) {
	files := []string{
		"testdata/replay.0.json",
	}

	for _, file := range files {
		match, m := replayGet(file)

		f := match.FoodLocations(m, 5, 10)
		log.Print("food:", f)
	}
}

func TestHillLocations(t *testing.T) {
	files := []string{
		"testdata/replay.0.json",
	}

	for _, file := range files {
		match, m := replayGet(file)

		f := match.HillLocations(m, 5, 10)
		log.Print("hills:", f)
	}
}

func DumpMove(m *maps.Map, am [][][]combat.AntMove, p int) {
	for turn := range am {
		log.Print("TURN ", turn)
		for np := range am[turn] {
			if p == np || p < 0 {
				log.Print("PLAYER ", np)
				for _, move := range am[turn][np] {
					log.Printf("Move t=%d p=%d %#v %v %#v", turn, np, m.ToPoint(move.From), move.D, m.ToPoint(move.To))
				}
			}
		}
	}
}

func TestAntMoves(t *testing.T) {
	files := []string{
		"testdata/replay.0.json",
	}

	for _, file := range files {
		match, m := replayGet(file)

		ac := match.AntCount(0, match.GameLength)
		ac9 := match.AntCount(9, 9)
		ac910 := match.AntCount(9, 10)

		log.Print("am:**********************************************************")
		am := match.AntMoves(m, 0, match.GameLength)
		log.Print("am9:************************************************************")
		am9 := match.AntMoves(m, 9, 9)
		log.Print("am910:************************************************************")
		am910 := match.AntMoves(m, 9, 10)
		//DumpMove(m, am, -1)
	OUT:
		for p := range ac {
			am0am9 := reflect.DeepEqual(am[9][p], am9[0][p])
			am0am910 := reflect.DeepEqual(am[9][p], am910[0][p])
			am1am910 := reflect.DeepEqual(am[10][p], am910[1][p])
			if !am0am9 || !am0am910 || !am1am910 {
				log.Printf("Match: am0am9: %v am0am910 %v; am1am910 %v", am0am9, am0am910, am1am910)
				log.Print("am[9]:**********************************************************")
				DumpMove(m, am[9:11], p)
				log.Print("am9:************************************************************")
				DumpMove(m, am9, p)
				log.Print("am910[0]:************************************************************")
				DumpMove(m, am910, p)
				t.Error("AntMove move location mismatch for subsets player ", p)
				break OUT
			}

			for turn := range ac[p] {
				if ac[p][turn] != len(am[turn][p]) {
					t.Errorf("AntCount, AntMove count mismatch player %d turn %d: count(%d) != locs(%d)", p, turn, ac[p][turn], len(am[turn][p]))
					break OUT
				}
			}

			if ac[p][9] != ac9[p][0] ||
				ac[p][9] != ac910[p][0] ||
				ac[p][10] != ac910[p][1] {
				t.Errorf("AntCount mismatch for full versus subset player %d; ac %d ac9 %d, ac %d ac910[0] %d, ac %d ac910[1] %d", p,
					ac[p][9], ac9[p][0],
					ac[p][9], ac910[p][0],
					ac[p][10], ac910[p][1])
			}

		}
		if false {
			p := 2
			for turn := range am {
				log.Printf("%d: %v", turn, am[turn][p])
			}
		}
	}
}

type testData struct {
	file string
	g    *GameResult
	p    []PlayerResult
}

func TestExtractMetadata(t *testing.T) {

	tests := []testData{
		{
			file: "testdata/replay.0.json",
			g:    &GameResult{GameId: 0, Date: "", GameLength: 288, Challenge: "ants", WorkerId: "", Location: "localhost", MapId: "c24acaf851f914c95f5686ae3a117691"},
			p: []PlayerResult{
				{PlayerName: "bot.sh", PlayerTurns: 288, Score: 34, Rank: 0, Bonus: 0, Status: "survived"},
				{PlayerName: "bugnutsv3", PlayerTurns: 203, Score: 4, Rank: 1, Bonus: 0, Status: "eliminated"},
				{PlayerName: "HunterBot.py", PlayerTurns: 288, Score: 2, Rank: 3, Bonus: 0, Status: "survived"},
				{PlayerName: "GreedyBot.py", PlayerTurns: 143, Score: 4, Rank: 1, Bonus: 0, Status: "eliminated"},
			},
		},
		{
			file: "testdata/replay.1.json",
			g:    &GameResult{GameId: 95405, Date: "2011-11-15T07:03:21+00:00", GameLength: 70, Challenge: "ants", MatchupId: new(int), PostId: new(int), WorkerId: "69", Location: "aichallenge.org", MapId: "e55cc99bdaf9a567b01c70bf21410e4d"},
			p: []PlayerResult{
				{PlayerName: "hohohoman", PlayerTurns: 70, UserId: new(int), SubmissionId: new(int), Score: 1, Rank: 0, Bonus: 0, Status: "eliminated", ChallengeRank: new(int), ChallengeSkill: new(float64)},
				{PlayerName: "amoore", PlayerTurns: 70, UserId: new(int), SubmissionId: new(int), Score: 1, Rank: 0, Bonus: 0, Status: "eliminated", ChallengeRank: new(int), ChallengeSkill: new(float64)},
			},
		},
	}
	// fill in pointers
	*tests[1].g.PostId = 61
	*tests[1].g.MatchupId = 98617
	*tests[1].p[0].UserId = 4780
	*tests[1].p[0].SubmissionId = 7545
	*tests[1].p[0].ChallengeRank = 3043
	*tests[1].p[0].ChallengeSkill = 43.1309840679169
	*tests[1].p[1].UserId = 10955
	*tests[1].p[1].SubmissionId = 21778
	*tests[1].p[1].ChallengeRank = 5435
	*tests[1].p[1].ChallengeSkill = 37.5642597675323

	for _, test := range tests {
		m, err := Load(test.file)
		if err != nil {
			t.Errorf("Load of %s failed %v", test.file, err)
		}

		g, p := m.ExtractMetadata()

		if false {
			fmt.Fprintf(os.Stdout, "%#v\n", g)
			for i := range p {
				fmt.Fprintf(os.Stdout, "%#v\n", p[i])
			}
		} else {
			if !reflect.DeepEqual(g, test.g) {
				t.Errorf("Game result mismatch %s:\nexpected: %#v\ngot: %#v", test.file, test.g, g)
			}
			if len(p) != len(test.p) {
				t.Errorf("Player len mismatch %s %d != %d ", test.file, len(test.p), len(p))
			} else {
				for i := range p {
					test.p[i].Game = g
					if !reflect.DeepEqual(p[i], &test.p[i]) {
						t.Errorf("Game result mismatch %s:\nexpected: %#v\ngot: %#v", test.file, test.p[i], p[i])
					}
				}
			}
			/*
				 log.Printf("%#v", g)
				 for _, pd := range p {
					 s, _ := json.Marshal(pd)
					 //log.Printf("%s", string(s))
				 }
			*/
		}
	}
}

func TestGetMap(t *testing.T) {
	files := []string{
		"testdata/replay.0.json",
		"testdata/replay.1.json",
	}
	mapfiles := []string{
		"testdata/replay.0.map",
		"testdata/replay.1.map",
	}

	for i, file := range files {
		match, err := Load(file)
		if err != nil {
			t.Errorf("Load of %s failed %v", file, err)
		}
		m := match.Replay.GetMap()

		m2, err := maps.MapLoadFile(mapfiles[i])
		if err != nil {
			t.Errorf("Load of %s failed %v", mapfiles[i], err)
		}
		if m.Players != m2.Players {
			t.Errorf("Player count mismatch for %s, %d and %d", file, m.Players, m2.Players)
		}
		for j, item := range m2.Grid {
			if item != m.Grid[j] {
				t.Errorf("Map data mismatch %v", m2.ToPoint(torus.Location(j)))
			}
		}
	}
}
