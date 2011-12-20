// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package pathing

import (
	"testing"
	. "bugnuts/torus"
)

func TestQ(t *testing.T) {
	q := NewQueue(100)

	q.Q(Point{C: 1})
	q.Q(Point{C: 2})

	if n := q.Size(); n != 2 {
		t.Errorf("Queue size not 2: %#v", q)
	}

	if p := q.DQ(); p.C != 1 {
		t.Errorf("Expected %v got %v", Point{C: 1}, p)
	}

	for i := 3; i < 10; i++ {
		q.Q(Point{C: i})
	}

	if qpos := q.Position(Point{C: 8}); qpos != 6 {
		t.Errorf("Expected position 6 got %d", qpos)
	}
	if qpos := q.Position(Point{C: 99}); qpos != -1 {
		t.Errorf("Expected qpos -1 got %d", qpos)
	}

	for i := 2; i < 10; i++ {
		if p := q.DQ(); p.C != i {
			t.Errorf("Expected %v got %v", Point{C: i}, p)
		}
	}

	if n := q.Size(); n != 0 {
		t.Errorf("Queue Size should be 0 got %#v", q)
	}

	if !q.Empty() {
		t.Errorf("Queue q.Empty should be true")
	}
}
