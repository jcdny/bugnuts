package engine

import (
	"log"
	"sort"
	"bugnuts/torus"
	"bugnuts/maps"
	"bugnuts/game"
	"bugnuts/replay"
	"bugnuts/util"
	"bugnuts/combat"
)

type Player struct {
	Hills   []torus.Location
	Seen    []bool
	Visible []bool
	IdMap   []int
	InvMap  []int
}

type Game struct {
	*maps.Map
	*game.GameInfo
	C          *combat.Combat
	ViewMask   *maps.Mask
	AttackMask *maps.Mask
	Turn       int
	Players    []*Player

	// Replay results
	PlayerInput [][]*game.Turn
	// PlayerOutput [][]*game.Moves
}

var _false [maps.MAXMAPSIZE]bool

func NewGame(gi *game.GameInfo, m *maps.Map) *Game {
	// allocate all the threat arrays at the same time.

	vm := maps.MakeMask(gi.ViewRadius2, gi.Rows, gi.Cols)
	am := maps.MakeMask(gi.AttackRadius2, gi.Rows, gi.Cols)
	return NewGameMasks(gi, m, vm, am)
}

func NewGameMasks(gi *game.GameInfo, m *maps.Map, vm *maps.Mask, am *maps.Mask) *Game {
	// allocate all the threat arrays at the same time.

	g := &Game{
		Map:        m,
		GameInfo:   gi,
		Players:    make([]*Player, m.Players),
		ViewMask:   vm,
		AttackMask: am,
	}

	// Set Hill Locations and initialize the player map, each player is it's own 0th player.
	for i := range g.Players {
		g.Players[i] = &Player{
			Hills:   m.Hills(i),
			IdMap:   make([]int, m.Players),
			InvMap:  make([]int, m.Players),
			Seen:    make([]bool, m.Size()),
			Visible: make([]bool, m.Size()),
		}
		for j := range g.Players[i].IdMap {
			g.Players[i].IdMap[j] = -1
			g.Players[i].InvMap[j] = -1
		}

		g.Players[i].AddIdMap(i)
	}

	return g
}

func (p *Player) AddIdMap(np int) {
	p.IdMap[np] = util.Max(p.IdMap) + 1
	p.InvMap[p.IdMap[np]] = np
}

// Replay a game, returns turns for all players between turn tmin and tmax inclusive.
// Assumes Game is in initial state.
func (g *Game) Replay(r *replay.Replay, tmin, tmax int, canonicalorder bool) {
	tout := make([][]*game.Turn, 0, tmax-tmin+1)

	// Extract the location data from the replay.
	// We have to run from 0 since we need to update Player.Seen from start
	ants := r.AntMoves(g.Map, 0, tmax)
	food := r.FoodLocations(g.Map, 0, tmax)
	hills := r.HillLocations(g.Map, 0, tmax)

	g.C = combat.NewCombat(g.Map, g.AttackMask, len(g.Players), nil)

	for i := 0; i <= tmax; i++ {
		g.Turn = i + tmin + 1
		tset := g.GenerateTurn(ants[i], hills[i], food[i], canonicalorder)
		for j := range tset {
			tset[j].Turn = g.Turn
		}
		if i >= tmin {
			tout = append(tout, tset)
		}
	}

	g.PlayerInput = tout
}

func (p *Player) UpdateVisibility(g *Game, ants []game.AntMove, np int, seen *[]torus.Location) {
	copy(p.Visible, _false[:len(p.Visible)])
	if false && np == 0 {
		game.DumpAntMove(g.Map, ants, np, g.Turn)
	}
	for i := range ants {
		if ants[i].Player == np && ants[i].To > -1 {
			g.ApplyOffsets(ants[i].To, &g.ViewMask.Offsets, func(l torus.Location) {
				p.Visible[l] = true
				if !p.Seen[l] {
					p.Seen[l] = true
					*seen = append(*seen, l)
				}
			})
		}
	}
	if false && np == 0 {
		sort.Sort(torus.LocationSlice(*seen))
		log.Printf("t %d Seen %v", g.Turn, g.ToPoints(*seen))
	}

	return
}

// Generate the Turn output for each player given a collection of ant locations
func (g *Game) GenerateTurn(ants [][]game.AntMove, hills []game.PlayerLoc, food []torus.Location, canonicalorder bool) []*game.Turn {
	turns := make([]*game.Turn, len(g.Players))

	// Handle Combat for the passed locations.
	g.C.Reset()

	moves, spawn := g.C.SetupReplay(ants)
	// NB: This is very tricky - we are avoiding allocating new slices so we
	// are cutting up ants in place into live and dead, but that means
	// appending to dead would overwrite live.
	live, dead := g.C.Resolve(moves)
	if canonicalorder {
		sort.Sort(game.AntMoveSlice(dead))
		sort.Sort(game.PlayerLocSlice(hills))
		sort.Sort(torus.LocationSlice(food))
	}

	// Handle Razes

	// Handle Spawns
	live = append(live, spawn...) // NB see Resolve above wrt live slice.

	// Handle Gather

	// Update visibility, generating new water, all ants (updating IdMap), hills, and food seen
	seen := make([]torus.Location, 0, 300)
	for np, p := range g.Players {
		t := &game.Turn{Map: g.Map}
		seen = seen[0:0]
		p.UpdateVisibility(g, live, np, &seen)
		// newly visible water
		for _, loc := range seen {
			if g.Map.Grid[loc] == maps.WATER {
				t.W = append(t.W, loc)
			}
		}
		// visible live ants
		for i := range live {
			if p.Visible[live[i].To] {
				if p.IdMap[live[i].Player] < 0 {
					p.AddIdMap(live[i].Player)
				}
				t.A = append(t.A, game.PlayerLoc{Loc: live[i].To, Player: p.IdMap[live[i].Player]})
			}
		}

		for i := range dead {
			if p.Visible[dead[i].To] || dead[i].Player == np {
				if p.IdMap[dead[i].Player] < 0 {
					p.AddIdMap(dead[i].Player)
				}
				t.D = append(t.D, game.PlayerLoc{Loc: dead[i].To, Player: p.IdMap[dead[i].Player]})
			}
		}

		// visible hills
		for _, h := range hills {
			if p.Visible[h.Loc] {
				if p.IdMap[h.Player] < 0 {
					p.AddIdMap(h.Player)
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
		turns[np] = t
	}
	return turns
}
