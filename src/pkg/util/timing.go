package util

import (
	"time"
	"log"
	//"os"
	//"os/signal"
)

var LMark []int64 = make([]int64, 0, 10)
var LStr []string = make([]string, 0, 10)

func TPush(s string) {
	mark := time.Nanoseconds()
	LMark = append(LMark, mark)
	LStr = append(LStr, s)
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

	log.Printf("** %.2fms %s", float64(diff)/1000000.0, s)

	return diff / 1000000
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
/*
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
