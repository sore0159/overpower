package overpower

import (
	"mule/hexagon"
)

type Source interface {
	// -------- GET ------- //
	Game() (Game, error)
	Factions() ([]Faction, error)
	Planets() ([]Planet, error)
	Orders() ([]Order, error)
	Ships() ([]Ship, error)
	// ------- MAKE ------- //
	NewPlanet(string, int, int, int, int, hexagon.Coord) Planet
	NewPlanetView(int, Planet, bool) PlanetView
	NewMapView(int, hexagon.Coord) MapView
	NewShip(int, int, int, hexagon.CoordList) Ship
	NewShipView(Ship, int, int, hexagon.NullCoord, hexagon.NullCoord, hexagon.CoordList) ShipView
	NewLaunchRecord(Ship)
	NewLandingRecord(int, int, Ship, [3]int)
	// ------ CHANGE ----- //
	UpdatePlanetView(int, int, Planet) PlanetView
	// ------- DROP ------ //
	DropShip(Ship)
	DropOrders()
}

type Game interface {
	Turn() int
	SetTurn(int)
	IncTurn()
	Gid() int
	Name() string
	Owner() string
	HasPW() bool
	IsPwd(string) bool
	AutoDays() [7]bool
	SetAutoDays([7]bool)
	FreeAutos() int
	SetFreeAutos(int)
	ToWin() int
	HighScore() int
	SetHighScore(int)
}

type Faction interface {
	Gid() int
	Fid() int
	Owner() string
	Name() string
	IsDone() bool
	DoneBuffer() int
	SetDoneBuffer(int)
	Score() int
	SetScore(int)
}

type Planet interface {
	Gid() int
	Loc() hexagon.Coord
	Name() string
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
	Gid() int
	Fid() int
	Loc() hexagon.Coord
	Turn() int
	Name() string
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

type Order interface {
	Gid() int
	Fid() int
	Source() hexagon.Coord
	Target() hexagon.Coord
	Size() int
	SetSize(int)
}

type Ship interface {
	Gid() int
	Fid() int
	Sid() int
	Size() int
	Launched() int
	Path() hexagon.CoordList
}

type ShipView interface {
	Gid() int
	Fid() int
	Sid() int
	Turn() int
	Size() int
	Controller() int
	Loc() hexagon.NullCoord
	Dest() hexagon.NullCoord
	Trail() hexagon.CoordList
}

type MapView interface {
	Gid() int
	Fid() int
	Center() hexagon.Coord
	SetCenter(hexagon.Coord)
}

type Report interface {
	Gid() int
	Fid() int
	Turn() int
	Contents() []string
	AddContent(string)
}

type LaunchRecord interface {
	Gid() int
	Fid() int
	Turn() int
	Source() hexagon.Coord
	Target() hexagon.Coord
	Size() int
}

type LandingRecord interface {
	Gid() int
	Fid() int
	Turn() int
	Index() int
	Target() hexagon.Coord
	Size() int
	FirstController() int
	ResultController() int
	ResultInhabitants() int
}
