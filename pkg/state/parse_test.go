package state

import (
	"log"
	"testing"
	"os"
	"bufio"
	"io/ioutil"
	"bytes"
)

func TestParse(t *testing.T) {
	fnm := "testdata/test.input"

	file, _ := os.Open(fnm)
	defer file.Close()

	in := bufio.NewReader(file)
	g, err := GameScan(in)
	log.Printf("Game:\n%v", g)

	if err != nil {
		t.Errorf("Reading %s error %v", fnm, err)
	}

	s := g.NewState()

	turns := make([]*Turn, 1, s.Turns)

	for {
		var t *Turn
		t = s.TurnScan(in, t)
		if t.End {
			log.Printf("End received at turn %d", len(turns))
			break
		}
		if t.Err != nil {
			log.Printf("Error on turn %d(%d): %v", len(turns), t.Turn, t.Err)
		}
		turns = append(turns, t)
	}
}

func BenchmarkParse(b *testing.B) {
	fnm := "testdata/test.input"

	data, _ := ioutil.ReadFile(fnm)

	for i := 0; i < b.N; i++ {
		in := bufio.NewReader(bytes.NewBuffer(data))
		g, err := GameScan(in)
		if err != nil {
			log.Panicf("Reading %s error %v", fnm, err)
		}
		s := g.NewState()
		turns := make([]*Turn, 1, s.Turns)
		for {
			var t *Turn
			t = s.TurnScan(in, t)
			if t.End {
				break
			}
			if t.Err != nil {
				log.Printf("Error on turn %d(%d): %v", len(turns), t.Turn, t.Err)
			}
			turns = append(turns, t)
		}
	}
}
