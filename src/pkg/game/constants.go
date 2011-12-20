// Copyright © 2011 Jeffrey Davis <jeff.davis@gmail.com>
// Use of this code is governed by the GPL version 2 or later.
// See the file LICENSE for details.

package game

// The risk characteristic cf combat/threat.go

const (
	RiskNone = iota
	RiskSafe
	RiskAverse
	RiskNeutral
	Suicidal
	MaxRiskStat
)

const (
	MS = 1000000
)
