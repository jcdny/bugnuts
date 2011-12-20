// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

// Bench is for functional testing of other packages without introducing circular dependencies.
package bench

import (
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}
