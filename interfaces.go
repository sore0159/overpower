package overpower

import (
	"mule/hexagon"
	"mule/mydb"
)

type Game interface {
	mydb.Updater
	//
	Turn() int
	IncTurn()
	Gid() int
	Name() string
	Owner() string
	HasPW() bool
	IsPwd(string) bool
	//
}

type Faction interface {
	Gid() int
	Fid() int
	Owner() string
	Name() string
	Done() bool
	SetDone(bool)
	//
	mydb.Updater
}

type Planet interface {
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
	//
	mydb.Updater
}

type PlanetView interface {
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
	//
	mydb.Updater
}

type FactionView interface {
	Gid() int
	Fid() int
	Center() hexagon.Pixel
	Zoom() int
	//
	mydb.Updater
}
