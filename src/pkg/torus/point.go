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
