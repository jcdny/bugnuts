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
		var turn *Turn
		turn, err = s.TurnScan(in, turn)
		if turn.End {
			log.Printf("End received at turn %d", len(turns))
			break
		}
		if err != nil {
			log.Printf("Error on turn %d(%d): %v", len(turns), turn.Turn, err)
		}
		turns = append(turns, turn)
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
			var turn *Turn
			turn, err = s.TurnScan(in, turn)
			if turn.End {
				break
			}
			if err != nil {
				log.Printf("Error on turn %d(%d): %v", len(turns), turn.Turn, err)
			}
			turns = append(turns, turn)
		}
	}
}
