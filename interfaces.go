package overpower

import (
	"mule/hexagon"
	"mule/mydb"
)

type Source interface {
	NewPlanetView(int, int, Planet) (PlanetView, bool)
	NewPlanet(string, int, int, int, int, int, int, hexagon.Coord) (Planet, bool)
}

type Game interface {
	mydb.Updater
	//
	Turn() int
	SetTurn(int)
	IncTurn()
	Gid() int
	Name() string
	Owner() string
	HasPW() bool
	IsPwd(string) bool
}

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
