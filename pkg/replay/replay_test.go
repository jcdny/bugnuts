package replay

import (
	"testing"
	"json"
	"io/ioutil"
	"log"
)

func TestReplayLoad(t *testing.T) {

	files := []string{
		"testdata/replay.0.json",
		"testdata/replay.1.json",
	}
	for _, s := range files {
		m := &Match{}

		buf, err := ioutil.ReadFile(s)
		if err != nil {
			log.Panicf("Readfile error %v", err)
		}

		// Do the actual parse here
		err = json.Unmarshal(buf, m)

		if err != nil || m.ReplayFormat != "json" {
			t.Errorf("Error on Unmarshal: %v (format found was %s)", err, m.ReplayFormat)
		}

		log.Printf("\n******************** %s ********************", s)

		log.Printf("%#v", m)
		// log.Printf("%#v", m.Replay.Hills)

		g, p := m.ExtractMetaData()
		log.Printf("%#v", g)
		for _, pd := range p {
			s, _ := json.Marshal(pd)
			log.Printf("%s", string(s))
		}
	}
}

func BenchmarkReplayLoad(b *testing.B) {

	buf, err := ioutil.ReadFile("testdata/replay.0.json")
	if err != nil {
		log.Panicf("Readfile error %v", err)
	}

	// Do the actual parse here
	for i := 0; i < b.N; i++ {
		m := &Match{}
		err = json.Unmarshal(buf, m)
	}
}
