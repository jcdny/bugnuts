package engine

import (
	//"log"
	"sort"
	"bugnuts/torus"
	"bugnuts/maps"
	"bugnuts/game"
	"bugnuts/replay"
	"bugnuts/util"
)

type Player struct {
	Hills   []torus.Location
	Seen    []bool
	Visible []bool
	IdMap   []int
	Threat  []int
}

type Game struct {
	*maps.Map
	*game.GameInfo
	ViewMask   *maps.Mask
	AttackMask *maps.Mask
	Turn       int
	Players    []*Player
	Threat     []int
	AntCount   []int
	PlayerMap  []int
}

var _minusone [maps.MAXMAPSIZE]int
var _zero [maps.MAXMAPSIZE]int
var _false [maps.MAXMAPSIZE]bool

var dloc torus.Location

func init() {
	for i := range _minusone {
		_minusone[i] = -1
	}
}

func NewGame(gi *game.GameInfo, m *maps.Map) *Game {
	// allocate all the threat arrays at the same time.
	threat := make([]int, m.Size()*(m.Players))

	g := &Game{
		Map:       m,
		GameInfo:  gi,
		Players:   make([]*Player, m.Players),
		Threat:    make([]int, m.Size()),
		PlayerMap: make([]int, m.Size()),
		AntCount:  make([]int, m.Size()),
	}
	dloc = g.ToLocation(torus.Point{R: 32, C: 14})

	g.ViewMask = maps.MakeMask(g.ViewRadius2, g.Rows, g.Cols)
	g.AttackMask = maps.MakeMask(g.AttackRadius2, g.Rows, g.Cols)

	// Set Hill Locations and initialize the player map, each player is it's own 0th player.
	for i := range g.Players {
		g.Players[i] = &Player{
			Hills:   m.Hills(i),
			IdMap:   make([]int, m.Players),
			Seen:    make([]bool, m.Size()),
			Visible: make([]bool, m.Size()),
			Threat:  threat[i*m.Size() : (i+1)*m.Size()],
		}
		for j := range g.Players[i].IdMap {
			g.Players[i].IdMap[j] = -1
		}
		g.Players[i].IdMap[i] = 0
	}

	return g
}

// Replay a game, returns turns for all players between turn tmin and tmax inclusive.
// Assumes Game is in initial state.
func (g *Game) Replay(r *replay.Replay, tmin, tmax int, canonicalorder bool) [][]*game.Turn {
	tout := make([][]*game.Turn, 0, tmax-tmin+1)

	// Extract the location data from the replay.
	// We have to run from 0 since we need to update Player.Seen from start
	ants, spawn := r.AntLocations(g.Map, 0, tmax)
	food := r.FoodLocations(g.Map, 0, tmax)
	hills := r.HillLocations(g.Map, 0, tmax)

	for i := 0; i <= tmax; i++ {
		tset := g.GenerateTurn(ants[i], spawn[i], hills[i], food[i], canonicalorder)
		if i >= tmin {
			tout = append(tout, tset)
		}
	}

	return tout
}

func (p *Player) UpdateVisibility(g *Game, ants []torus.Location) []torus.Location {
	copy(p.Visible, _false[:len(p.Visible)])

	seen := make([]torus.Location, 0, 100)
	for _, loc := range ants {
		// apply vis mask
		ap := g.ToPoint(loc)
		for _, op := range g.ViewMask.P {
			l := g.ToLocation(g.PointAdd(ap, op))
			p.Visible[l] = true
			if !p.Seen[l] {
				p.Seen[l] = true
				seen = append(seen, l)
			}
		}
	}

	return seen
}

func (g *Game) ComputeThreat(ants [][]torus.Location) {
	copy(g.PlayerMap, _minusone[:len(g.PlayerMap)])
	copy(g.AntCount, _zero[:len(g.AntCount)])
	copy(g.Threat, _zero[:len(g.Threat)])

	for np := range ants {
		copy(g.Players[np].Threat, _zero[:len(g.Threat)])

		for _, loc := range ants[np] {
			g.AntCount[loc]++

			// If we encounter suicides on the first one we remove the original threat from the
			// first ant encountered then ignore them subsequently
			inc := 1
			if g.AntCount[loc] == 1 {
				g.PlayerMap[loc] = np
			} else if g.AntCount[loc] == 2 {
				g.PlayerMap[loc] = -1
				inc = -1
			}

			if g.AntCount[loc] < 2 {
				ap := g.ToPoint(loc)
				for _, op := range g.AttackMask.P {
					nloc := g.ToLocation(g.PointAdd(ap, op))
					g.Threat[nloc] += inc
					g.Players[np].Threat[nloc] += inc
				}
			}
		}

		//log.Printf("p %d threat at dloc %d", np, g.Players[np].Threat[dloc])
		//log.Printf("p %d g threat at dloc %d", np, g.Threat[dloc])
	}
}

// ResolveCombat takes the ant location slices return list of dead ants and update the ant locations in place to remove dead ants.
// Does not preserve order on the per player list of ant locations.
func (g *Game) ResolveCombat(ants [][]torus.Location) []game.PlayerLoc {
	dead := make([]game.PlayerLoc, 0, 20)

	// suicides first.
	for np := range ants {
		// walk through ants swapping in the end of list for any dead ants
		lant := len(ants[np])
		for i := 0; i < lant; {
			if g.AntCount[ants[np][i]] > 1 {
				dead = append(dead, game.PlayerLoc{Loc: ants[np][i], Player: np})
				lant--
				ants[np][i] = ants[np][lant]
			} else {
				i++
			}
		}
		// truncate away the leftovers.
		ants[np] = ants[np][:lant]
	}

	// Now actual combat
	for np := range ants {
		lant := len(ants[np])
		for i := 0; i < lant; {
			loc := ants[np][i]

			t := g.Threat[loc] - g.Players[np].Threat[loc]
			if loc == dloc {
				// log.Printf("p %d net threat at dloc %d", np, t)
			}
			if t > 0 {
				ap := g.ToPoint(loc)
				for _, op := range g.AttackMask.P {
					nloc := g.ToLocation(g.PointAdd(ap, op))
					ntp := g.PlayerMap[nloc]
					if ntp >= 0 && ntp != np && t >= g.Threat[nloc]-g.Players[ntp].Threat[nloc] {
						dead = append(dead, game.PlayerLoc{Loc: ants[np][i], Player: np})
						lant--
						ants[np][i] = ants[np][lant]
						i-- // we just increment back to original i after the break
						break
					}
				}
			}
			i++
		}
		ants[np] = ants[np][:lant]
	}

	return dead
}

// Generate the Turn output for each player given a collection of ant locations
func (g *Game) GenerateTurn(ants [][]torus.Location, spawn, hills []game.PlayerLoc, food []torus.Location, canonicalorder bool) []*game.Turn {
	turns := make([]*game.Turn, len(g.Players))

	// Handle Combat for the passed locations.
	g.ComputeThreat(ants)
	dead := g.ResolveCombat(ants)

	if canonicalorder {
		sort.Sort(game.PlayerLocSlice(dead))
		sort.Sort(game.PlayerLocSlice(hills))
		sort.Sort(torus.LocationSlice(food))
	}

	// Handle Razes

	// Handle Spawns
	for _, s := range spawn {
		ants[s.Player] = append(ants[s.Player], s.Loc)
	}

	// Handle Gather

	// Update visibility, generating new water, all ants (updating IdMap), hills, and food seen
	for i, p := range g.Players {
		t := &game.Turn{Map: g.Map}
		seen := p.UpdateVisibility(g, ants[i])

		// newly visible water
		for _, loc := range seen {
			if g.Map.Grid[loc] == maps.WATER {
				t.W = append(t.W, loc)
			}
		}

		// visible ants
		for np := range ants {
			for i := range ants[np] {
				if p.Visible[ants[np][i]] {
					if p.IdMap[np] < 0 {
						p.IdMap[np] = util.Max(p.IdMap) + 1
					}
					t.A = append(t.A, game.PlayerLoc{Loc: ants[np][i], Player: p.IdMap[np]})
				}
			}
		}

		// visible dead
		for _, d := range dead {
			if p.Visible[d.Loc] || d.Player == i {
				if p.IdMap[d.Player] < 0 {
					p.IdMap[d.Player] = util.Max(p.IdMap) + 1
				}
				t.D = append(t.D, game.PlayerLoc{Loc: d.Loc, Player: p.IdMap[d.Player]})
			}
		}

		// visible hills
		for _, h := range hills {
			if p.Visible[h.Loc] {
				if p.IdMap[h.Player] < 0 {
					p.IdMap[h.Player] = util.Max(p.IdMap) + 1
				}
				t.H = append(t.H, game.PlayerLoc{Loc: h.Loc, Player: p.IdMap[h.Player]})
			}
		}

		// visible food
		for _, loc := range food {
			if p.Visible[loc] {
				t.F = append(t.F, loc)
			}
		}

		if canonicalorder {
			// got t.D, t.H and t.F by sorting inputs.
			sort.Sort(torus.LocationSlice(t.W))
			sort.Sort(game.PlayerLocSlice(t.A))
		}
		turns[i] = t
	}
	return turns
}
