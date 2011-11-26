package replay

import (
	"testing"
	"log"
)

func TestAntLocations(t *testing.T) {
	files := []string{
		"testdata/replay.0.json",
	}

	for _, file := range files {
		match, err := Load(file)
		if err != nil {
			t.Errorf("Load of %s failed: %v", file, err)
		}
		m := match.GetMap()

		al := match.AntLocations(m, match.GameLength)
		for i := range al {
			log.Printf("%d: %v", i, m.ToPoints(al[i][2]))
		}
	}
}
