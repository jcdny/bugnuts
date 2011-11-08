package main
// The v3 Bot -- Terrible!!!!

import (
	"fmt"
	"log"
	"math"
	"os"
	"rand"
)

type BotV3 struct {

}

//NewBot creates a new instance of your bot
func NewBotV3(s *State) Bot {
	mb := &BotV3{
	//do any necessary initialization here
	}

	return mb
}

func (bot *BotV3) Priority(i Item) int {
	return 1
}

func (bot *BotV3) DoTurn(s *State) os.Error {
	for loc, _ := range s.Ants[0] {
		p := s.Map.ToPoint(loc)
		best := math.MinInt32
		var score [4]int
		for d, op := range Steps {
			tp := s.PointAdd(p, op)
			if bot.validPoint(s, tp) {
				if false && rand.Intn(8) == 0 {
					score[d] = 500
				} else {
					score[d] = bot.Score(s, p, tp, s.viewMask.Add[d])
				}
				if score[d] > best {
					best = score[d]
				}
			} else {
				score[d] = -9999
			}
		}

		if Debug > 2 {
			log.Printf("TURN %d point %v score %v best %v", s.Turn, p, score, best)
		}

		if best > math.MinInt32 {
			var bestd []int
			for d, try := range score {
				if try == best {
					bestd = append(bestd, d)
				}
			}
			pp := rand.Perm(len(bestd))[0]
			// Swap the current and target cells
			tp := s.PointAdd(p, Steps[bestd[pp]])
			s.Map.Grid[s.ToLocation(tp)] = MY_ANT
			s.Map.Grid[s.ToLocation(p)] = LAND
			fmt.Fprintf(os.Stdout, "o %d %d %c\n", p.r, p.c, ([4]byte{'n', 's', 'e', 'w'})[bestd[pp]])
		}
	}
	fmt.Fprintf(os.Stdout, "go\n")

	return nil
}

func (bot *BotV3) Score(s *State, p, tp Point, pv []Point) int {
	score := 0

	// Score for explore
	for _, op := range pv {
		seen := s.Map.Seen[s.ToLocation(s.PointAdd(p, op))]
		switch {
		case seen < 1:
			score += 2
		case seen > s.Turn-2:
			score -= 1
		}
	}
	score = score * 17 / len(pv)

	if Debug > 3 {
		log.Printf("p %v tp %v explore score %d", p, tp, score)
	}

	// Score for nearby items
	for _, op := range s.viewMask.P {
		item := s.Map.Grid[s.ToLocation(s.PointAdd(tp, op))]
		inc := 0
		iname := ""
		if item != LAND && item != WATER {
			//log.Printf("%v %v %d %d",p, tp, op, d, item)
			d := Abs(op.c) + Abs(op.r)

			if item == MY_HILL {
				iname = "my hill"
				inc = -32 + 4*Min([]int{d, 8})
			}
			if item.IsEnemyHill() {
				iname = "enemy hill"
				inc = 1500 - 100*Min([]int{d, 10})
			}
			if item == FOOD {
				iname = "food"
				if d == 1 {
					inc = 1000
				} else {
					inc = 120 - 12*Min([]int{d, 10})
				}
			}
			if item == MY_ANT && d > 1 {
				iname = "my ant"
				inc = -30 + 5*Min([]int{d, 6})
			}
		}
		score += inc
		if Debug > 3 && iname != "" {
			log.Printf("tp %v (at %v) %s worth %d",
				tp, op, iname, inc)
		}
	}
	if Debug > 3 {
		log.Printf("p %v tp %v total score %d",
			p, tp, score)
	}
	return score
}

func (bot *BotV3) validPoint(s *State, p Point) bool {
	tgt := s.Map.Grid[s.ToLocation(p)]
	if tgt == FOOD || tgt == LAND || tgt.IsEnemyHill() {
		for _, op := range Steps {
			//make sure there is an exit
			ep := s.PointAdd(p, op)
			tgt := s.Map.Grid[s.ToLocation(ep)]
			if tgt == FOOD || tgt == LAND {
				return true
			}
		}
	}
	return false
}
