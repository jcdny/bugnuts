package parameters

import (
	"testing"
	"log"
)

func TestParametersLoadSave(t *testing.T) {
	for key, p := range ParameterSets {
		p.Save("tmp/" + key + ".param")
	}
	for key := range ParameterSets {
		pnew := &Parameters{}
		pnew.Load("tmp/" + key + ".param")
		log.Printf("Load %s got %#v", key, pnew)
	}
}
