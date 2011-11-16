package pathing

import (
	"log"
	. "bugnuts/maps"
)

type Queue struct {
	c []Point
}

func QNew(cap int) *Queue {
	q := &Queue{
		c: make([]Point, 0, cap),
	}

	return q
}

func (q *Queue) DQ() Point {
	if q.Empty() {
		log.Panicf("Queue.DQ() empty queue")
	}
	p := q.c[0]
	q.c = q.c[1:]
	//log.Printf("DQ %v", p)
	return p
}

func (q *Queue) Q(p Point) {
	//log.Printf("Q %v", p)
	q.c = append(q.c, p)
}

func (q *Queue) Empty() bool {
	if len(q.c) < 1 {
		return true
	}
	return false
}

func (q *Queue) Size() int {
	return len(q.c)
}

func (q *Queue) Cap() int {
	return cap(q.c)
}

func (q *Queue) Position(p Point) int {
	pos := -1

	for i, qp := range q.c {
		// assumes we have aleady wrapped.
		if p.C == qp.C && p.R == qp.R {
			pos = i
			break
		}
	}

	return pos
}
