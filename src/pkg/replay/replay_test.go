package replay

import (
	"testing"
	"json"
	"log"
	"io/ioutil"
)

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
