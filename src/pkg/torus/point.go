// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package torus

import (
	"fmt"
)
// Point is the coordinate on the torus.  It can be signed in the case of offset arrays 
// or points in non standard form.
type Point struct {
	R, C int
}

func (p *Point) String() string {
	return fmt.Sprintf("r:%d C:%d", p.R, p.C)
}
