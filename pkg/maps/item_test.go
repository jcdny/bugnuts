package maps

import (
	"testing"
)

func TestItemMap(t *testing.T) {
	for i := UNKNOWN; i < MAX_ITEM; i++ {
		if i != ToItem(i.ToSymbol()) {
			t.Errorf("Map from %v to %c to %v not reflexive", i, i.ToSymbol(), ToItem(i.ToSymbol()))
		}
	}

	if ToItem('}') != INVALID_ITEM {
		t.Errorf("Map from '}' returns %v, should be INVALID_ITEM %v", ToItem('}'), INVALID_ITEM)
	}
}
