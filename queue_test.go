package main

import (
	"testing"
)

func TestQ(t *testing.T) {
	q := QNew(100)

	q.Q(Point{c: 1})
	q.Q(Point{c: 2})

	if n := q.Size(); n != 2 {
		t.Errorf("Queue size not 2: %#v", q)
	}

	if p := q.DQ(); p.c != 1 {
		t.Errorf("Expected %v got %v", Point{c: 1}, p)
	}

	for i := 3; i < 10; i++ {
		q.Q(Point{c: i})
	}

	if qpos := q.Position(Point{c: 8}); qpos != 6 {
		t.Errorf("Expected position 6 got %d", qpos)
	}
	if qpos := q.Position(Point{c: 99}); qpos != -1 {
		t.Errorf("Expected qpos -1 got %d", qpos)
	}

	for i := 2; i < 10; i++ {
		if p := q.DQ(); p.c != i {
			t.Errorf("Expected %v got %v", Point{c: i}, p)
		}
	}

	if n := q.Size(); n != 0 {
		t.Errorf("Queue Size should be 0 got %#v", q)
	}

	if !q.Empty() {
		t.Errorf("Queue q.Empty should be true")
	}
}
