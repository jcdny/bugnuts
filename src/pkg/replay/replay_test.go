package replay

import (
	"testing"
	"json"
	"log"
	"io/ioutil"
	"bugnuts/maps"
)

func TestExtractMetadata(t *testing.T) {

	files := []string{
		"testdata/replay.0.json",
		"testdata/replay.1.json",
	}
	for _, file := range files {
		m, err := Load(file)
		if err != nil {
			t.Errorf("Load of %s failed %v", file, err)
		}

		g, p := m.ExtractMetadata()
		if g.Challenge != "ants" || g.Location == "" {
			t.Errorf("Load of %s bad %#v", file, g)
		}
		if len(p) != len(m.PlayerNames) {
			t.Errorf("Player len mismatch %s %d != %d ", file, len(m.PlayerNames), len(p))
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
				t.Errorf("Map data mismatch %v", m2.ToPoint(maps.Location(j)))
			}
		}
	}
}
func BenchmarkReplayUnmarshall(b *testing.B) {

	buf, err := ioutil.ReadFile("testdata/replay.1.json")
	if err != nil {
		log.Panicf("Readfile error %v", err)
	}

	// Do the actual parse here
	for i := 0; i < b.N; i++ {
		m := &Match{}
		err = json.Unmarshal(buf, m)
	}
}
