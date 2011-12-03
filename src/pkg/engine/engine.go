package engine

import (
	"log"
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
}

type Game struct {
	*maps.Map
	*game.GameInfo
	ViewMask   *maps.Mask
	AttackMask *maps.Mask
	Turn       int
	Players    []*Player
}

const CanonicalOrder bool = true

func NewGame(gi *game.GameInfo, m *maps.Map) *Game {
	g := &Game{
		Map:      m,
		GameInfo: gi,
		Players:  make([]*Player, m.Players),
	}

	g.ViewMask = maps.MakeMask(g.ViewRadius2, g.Rows, g.Cols)
	g.AttackMask = maps.MakeMask(g.AttackRadius2, g.Rows, g.Cols)

	// Set Hill Locations and initialize the player map, each player is it's own 0th player.
	for i := range g.Players {
		log.Printf("Hills: %v", m.Hills(i))
		g.Players[i] = &Player{
			Hills:   m.Hills(i),
			IdMap:   make([]int, m.Players),
			Seen:    make([]bool, m.Size()),
			Visible: make([]bool, m.Size()),
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
func (g *Game) Replay(r *replay.Replay, tmin, tmax int) [][]*game.Turn {
	tout := make([][]*game.Turn, 0, tmax-tmin+1)

	// Extract the location data from the replay.
	// We have to run from 0 since we need to update Player.Seen from start
	ants := r.AntLocations(g.Map, 0, tmax)
	food := r.FoodLocations(g.Map, 0, tmax)
	hills := r.HillLocations(g.Map, 0, tmax)

	for i := 0; i <= tmax; i++ {
		tset := g.GenerateTurn(ants[i], hills[i], food[i])
		if i >= tmin {
			tout = append(tout, tset)
		}
	}

	return tout
}

var _false [maps.MAXMAPSIZE]bool

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

	if CanonicalOrder {
		sort.Sort(torus.LocationSlice(seen))
	}
	return seen
}

// Generate the Turn output for each player given a collection of ant locations
func (g *Game) GenerateTurn(ants [][]torus.Location, hills []game.PlayerLoc, food []torus.Location) []*game.Turn {
	turns := make([]*game.Turn, len(g.Players))

	// Handle Combat for the passed locations.

	// Handle Razes

	// Handle Spawns

	// Handle Gather

	// Update visibility, generating new water, all ants (updating IdMap), hills, and food seen
	for i, p := range g.Players {
		t := &game.Turn{Map: g.Map}
		seen := p.UpdateVisibility(g, ants[i])
		for _, loc := range seen {
			if g.Map.Grid[loc] == maps.WATER {
				t.W = append(t.W, loc)
			}
		}
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

		if CanonicalOrder {
			sort.Sort(game.PlayerLocSlice(t.A))
			sort.Sort(game.PlayerLocSlice(hills))

		}

		for _, h := range hills {
			if p.Visible[h.Loc] {
				if p.IdMap[h.Player] < 0 {
					p.IdMap[h.Player] = util.Max(p.IdMap) + 1
				}
				t.H = append(t.H, game.PlayerLoc{Loc: h.Loc, Player: p.IdMap[h.Player]})
			}
		}
		for _, loc := range food {
			if p.Visible[loc] {
				t.F = append(t.F, loc)
			}
		}

		turns[i] = t
	}
	return turns
}
