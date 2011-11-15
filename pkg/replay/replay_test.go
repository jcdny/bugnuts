package replay

import (
	"testing"
	"json"
	"io/ioutil"
	"log"
)

func TestReplayLoad(t *testing.T) {
	m := &Match{}

	buf, err := ioutil.ReadFile("testdata/replay.0.json")
	if err != nil {
		log.Panicf("Readfile error %v", err)
	}

	// Do the actual parse here
	err = json.Unmarshal(buf, m)

	if err != nil || m.ReplayFormat != "json" {
		t.Errorf("Error on Unmarshal %v, format found was %s", err, m.ReplayFormat)
	}

	log.Printf("%#v", m)
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
