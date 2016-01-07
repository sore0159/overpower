package json

import (
	"mule/overpower"
)

type Game struct {
	Turn      int     `json:"turn"`
	Gid       int     `json:"gid"`
	Name      string  `json:"name"`
	Owner     string  `json:"owner"`
	HasPW     bool    `json:"haspw"`
	AutoDays  [7]bool `json:"autodays"`
	FreeAutos int     `json:"freeautos"`
}

func LoadGame(g overpower.Game) *Game {
	return &Game{
		Turn:      g.Turn(),
		Gid:       g.Gid(),
		Name:      g.Name(),
		Owner:     g.Owner(),
		HasPW:     g.HasPW(),
		AutoDays:  g.AutoDays(),
		FreeAutos: g.FreeAutos(),
	}
}

func LoadGames(games []overpower.Game) []*Game {
	gs := make([]*Game, len(games))
	for i, g := range games {
		gs[i] = LoadGame(g)
	}
	return gs
}

/*

type Faction interface {
	mydb.Updater
	//
	Gid() int
	Fid() int
	Owner() string
	Name() string
	Done() bool
	SetDone(bool)
}

type Planet interface {
	mydb.Updater
	mydb.Inserter
	//
	Gid() int
	Pid() int
	Name() string
	Loc() hexagon.Coord
	Controller() int
	SetController(int)
	Inhabitants() int
	SetInhabitants(int)
	Resources() int
	SetResources(int)
	Parts() int
	SetParts(int)
}

type PlanetView interface {
	mydb.Updater
	mydb.Inserter
	//
	Gid() int
	Pid() int
	Fid() int
	Name() string
	Loc() hexagon.Coord
	Turn() int
	SetTurn(int)
	Controller() int
	SetController(int)
	Inhabitants() int
	SetInhabitants(int)
	Resources() int
	SetResources(int)
	Parts() int
	SetParts(int)
}

type FactionView interface {
	mydb.Updater
	//
	Gid() int
	Fid() int
	Center() hexagon.Pixel
	Zoom() int
}

type Order interface {
	Gid() int
	Fid() int
	Source() int
	Target() int
	Size() int
	SetSize(int)
}

type Ship interface {
	Gid() int
	Fid() int
	Sid() int
	Size() int
	Launched() int
	Path() []hexagon.Coord
}

type ShipView interface {
	Gid() int
	Fid() int
	Sid() int
	Turn() int
	Size() int
	Controller() int
	Loc() (hexagon.Coord, bool)
	Dest() (hexagon.Coord, bool)
	Trail() []hexagon.Coord
}

type MapView interface {
	mydb.Updater
	//
	Gid() int
	Fid() int
	Zoom() int
	Center() hexagon.Coord
	SetZoom(int)
	SetCenter(hexagon.Coord)
	Target1() hexagon.NullCoord
	SetTarget1(hexagon.Coord)
	DropTarget1()
	Target2() hexagon.NullCoord
	SetTarget2(hexagon.Coord)
	DropTarget2()
}

type Report interface {
	Gid() int
	Fid() int
	Turn() int
	Contents() []string
	AddContent(string)
}
*/
