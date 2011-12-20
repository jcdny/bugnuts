// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

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
