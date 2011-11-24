package replay

import (
	"testing"
	"json"
	"log"
	"io/ioutil"
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
