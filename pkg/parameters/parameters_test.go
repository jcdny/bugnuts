package main

import (
	"testing"
	"log"
)

func TestParametersLoadSave(t *testing.T) {
	for key, p := range ParameterSets {
		p.Save(key + ".param.tmp")
	}
	for key, _ := range ParameterSets {
		pnew := &Parameters{}
		pnew.Load(key + ".param.tmp")
		log.Printf("Load %s got %#v", key, pnew)
	}
}
