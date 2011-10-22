package main

import (
	"testing"
	"log"
)

func TestItemMap(t *testing.T) {
	for i := UNKNOWN; i < MAX_ITEM; i++ {
		if i != ToItem(i.ToSymbol()) {
			t.Errorf("Map from %v to %c to %v not reflexive",i, i.ToSymbol(), ToItem(i.ToSymbol()))
		}
	}

	if ToItem('~') != INVALID_ITEM {
		t.Errorf("Map from '~' returns %v, should be INVALID_ITEM %v", ToItem('~'), INVALID_ITEM)
	}
}

func TestGenCircleTable(t *testing.T) {
	exp7100 := []int{0, -201, -200, -199, -102, -101, -100, -99, -98, -2, -1, 1, 2, 98, 99, 100, 101, 102, 199, 200, 201}

	v := GenCircleTable(7)
	if len(v) != len(exp7100) {
		t.Errorf("GenCircleTable(7) expected %v got %v", exp7100, v)
	}

	v = GenCircleTable(1)
	if len(v) != 5 {
		t.Errorf("GenCircleTable(1) expected len=5 got %v", v)
	}

	v = GenCircleTable(0)
	if len(v) != 1 {
		t.Errorf("GenCircleTable(0) expected len=1 got %v", v)
	}

	for i := 0; i <= 10; i += 1 {
		v = GenCircleTable(i)
		log.Printf("r2=%5d %5d %v", i, len(v), v)
	}
}