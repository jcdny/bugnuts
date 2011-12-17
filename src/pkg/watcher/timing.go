package watcher

import (
	"time"
	"log"
	"fmt"
	"os"
)

type TimingData struct {
	Name       string
	TLast      int
	Total      int64
	Accumulate []int64
	Count      []int
	Started    []int64
	Stopped    []int64
	stack      int
}

var Times = make(map[string]*TimingData, 100)
var TurnTimer *TimingData
var LMark []int64 = make([]int64, 0, 10)
var LStr []string = make([]string, 0, 10)

//Todo channels and stuff

func init() {
	TurnTimer = NewTimingData("turntimer", 2000)
}

func NewTimingData(name string, n int) *TimingData {
	td := &TimingData{
		Name:       name,
		Accumulate: make([]int64, n),
		Count:      make([]int, n),
		Started:    make([]int64, n),
		Stopped:    make([]int64, n),
	}
	return td
}

func TPush(s string) {
	mark := time.Nanoseconds()
	LMark = append(LMark, mark)
	LStr = append(LStr, s)

	if Debug[DBG_GatherTime] && s[0] == '@' {
		turn := TurnGet()
		t, ok := Times[s]
		if !ok {
			t = NewTimingData(s[1:], 2500)
			Times[s] = t
		}
		if t.Started[turn] == 0 {
			t.Started[turn] = mark
		}
		t.stack++
		t.Count[turn]++
		t.TLast = turn
	}
}

func TPop() int64 {
	if len(LMark) < 1 {
		return 0
	}
	mark := time.Nanoseconds()
	diff := mark - LMark[len(LMark)-1]
	s := LStr[len(LStr)-1]

	LMark = LMark[:len(LMark)-1]
	LStr = LStr[:len(LStr)-1]

	if Debug[DBG_GatherTime] && s[0] == '@' {
		if t, ok := Times[s]; ok {
			t.Total += diff
			t.Accumulate[t.TLast] += diff
			if t.stack < 2 {
				t.Stopped[t.TLast] = mark
				t.stack = 0
			} else {
				t.stack--
			}
			if Debug[DBG_TurnTime] && s[0] == '@' {
				log.Printf("** %.2fms/%.2fms/%.2fms %s",
					float64(diff)/1e6,
					float64(t.Accumulate[t.TLast])/1e6,
					float64(t.Total)/1e6,
					s)
			}
			return diff / 1e6
		}
	}

	if Debug[DBG_TurnTime] && s[0] == '@' {
		log.Printf("** %.2fms %s", float64(diff)/1000000.0, s)
	}

	return diff / 1e6
}

func TMark(s string) int64 {
	if len(LMark) < 1 {
		return 0
	}
	mark := time.Nanoseconds()
	diff := mark - LMark[len(LMark)-1]
	ts := LStr[len(LStr)-1]

	log.Printf("** %.2fms %s: %s", float64(diff)/1000000.0, ts, s)

	return diff / 1000000
}

func TDump(file string) {
	if !Debug[DBG_GatherTime] {
		log.Print("Attempt to store turn time data not accumulated")
	}
	fd, err := os.Create(file)
	if err != nil {
		log.Print("Failed to create ", file, ": ", err)
		return
	}
	defer fd.Close()

	fmt.Fprintf(fd, "name,turn,count,accumulated,started,stopped\n")
	for i := 0; i < len(TurnTimer.Started); i++ {
		if TurnTimer.Started[i] != 0 && TurnTimer.Stopped[i] != 0 {
			fmt.Fprintf(fd, "%s,%d,%d,%.2f,%.2f,%.2f\n",
				TurnTimer.Name, i, TurnTimer.Count[i],
				0.0, 0.0,
				float64(TurnTimer.Stopped[i]-TurnTimer.Started[i])/1e6)
		}
	}

	for _, t := range Times {
		for i := 0; i <= t.TLast; i++ {
			if i > 0 || t.Count[i] > 0 {
				fmt.Fprintf(fd, "%s,%d,%d,%.2f,%.2f,%.2f\n",
					t.Name, i, t.Count[i],
					float64(t.Accumulate[i])/1e6,
					float64(t.Started[i]-TurnTimer.Started[i])/1e6,
					float64(t.Stopped[i]-TurnTimer.Started[i])/1e6)
			}
		}
	}
}

/*
	//"os"
	//"os/signal"
func TurnTimer() {
	log.Print("Starting timing")
	go func() {
		for isig := range signal.Incoming {
			sig := isig.(os.UnixSignal)
			switch sig {
			case os.SIGCONT:
				log.Printf("Got SigCont")
			default:
				log.Printf("Unexpected signal %v", sig)
			}
		}
	}()
}
*/
