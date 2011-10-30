package replay

import (
	"testing"
	"os"
	"json"
)

func TestReplayLoad(t *testing.T) {
	m := &Match{}

	buf := make([]byte, 1000000)

	f, err := os.Open("testdata/replay.0.json")
	defer f.Close()
	n, err := f.Read(buf[:])
	buf = buf[0:n]

	// Do the actual parse here
	err = json.Unmarshal(buf, m)

	if err != nil || m.ReplayFormat != "json" {
		t.Errorf("Error on Unmarshal %v, format found was %s", err, m.ReplayFormat)
	}
}

