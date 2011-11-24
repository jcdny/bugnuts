package replay

import (
	"testing"
	"reflect"
	"os"
	"log"
)

func TestReplayLoading(t *testing.T) {
	var err os.Error
	files := []string{
		"testdata/replay.0.html",
		"testdata/replay.0.html.gz",
		"testdata/replay.0.json",
		"testdata/replay.0.json.gz",
	}

	m := make([]*Match, len(files))
	for i, file := range files {
		m[i], err = Load(file)
		if err != nil {
			t.Errorf("Load of %s failed: %v", file, err)
		}
	}
	for i := 1; i < len(m); i++ {
		if !reflect.DeepEqual(m[i-1], m[i]) {
			t.Errorf("Replays %s and %s differ", files[i-1], files[i])
		} else if false {
			log.Printf("Files %s and %s match", files[i-1], files[i])
		}
	}
}
