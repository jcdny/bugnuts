package maps

import (
	"fmt"
)

type Point struct {
	R, C int
}

func (p *Point) String() string {
	return fmt.Sprintf("r:%d C:%d", p.R, p.C)
}
