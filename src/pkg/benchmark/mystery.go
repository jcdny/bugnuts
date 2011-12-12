package combat

// This is a mystery to me... the following routines are slower than the 
// existing setup function - even though they handle ant collisions.
// Some sort of oddity in the compiler code generation that I 
// do not have time to verify...
//
// bench.BenchmarkCombatSetupS1	    5000	    416106 ns/op
// bench.BenchmarkCombatSetupS2	    5000	    463537 ns/op
// bench.BenchmarkCombatSetup	   10000	    213160 ns/op

// Compute initial ant threat. returns a count of dead found.
// Should be 0 unless something has gone horribly wrong.
func (c *Combat) SetupS1(ants []map[Location]int) {
	for np := range ants {
		if len(ants[np]) > 0 && len(c.PThreat[np]) == 0 {
			c.PThreat[np] = make([]int, c.Map.Size())
		}
		for loc := range ants[np] {
			c.AntCount[loc]++
			c.PlayerMap[loc] = np
			tp := np
			c.ApplyOffsets(loc, &c.AttackMask.Offsets, func(nloc Location) {
				c.Threat[nloc]++
				c.PThreat[tp][nloc]++
			})
		}
	}
}

// Compute initial ant threat. returns a count of dead found.
// Should be 0 unless something has gone horribly wrong.
func (c *Combat) SetupS2(ants []map[Location]int) {
	for np := range ants {
		if len(ants[np]) > 0 && len(c.PThreat[np]) == 0 {
			c.PThreat[np] = make([]int, c.Map.Size())
		}
		for loc := range ants[np] {
			c.AddAntF(np, loc)
		}
	}
}

func (c *Combat) AddAntF(np int, loc Location) {
	c.AntCount[loc]++
	c.PlayerMap[loc] = np
	c.ApplyOffsets(loc, &c.AttackMask.Offsets, func(nloc Location) {
		c.Threat[nloc]++
		c.PThreat[np][nloc]++
	})
}
