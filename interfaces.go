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
	NewPlanet(string, int, int, int, int, int, hexagon.Coord) Planet
	NewPlanetView(int, Planet, bool) PlanetView
	NewMapView(int, hexagon.Coord) MapView
	NewShip(int, int, int, hexagon.CoordList) Ship
	NewShipView(Ship, int, int, hexagon.NullCoord, hexagon.NullCoord, hexagon.CoordList) ShipView
	NewReport(int, int) Report
	AddReportEvent(int, string) bool
	// ------ CHANGE ----- //
	UpdatePlanetView(int, int, Planet) PlanetView
	// ------- DROP ------ //
	DropShip(Ship)
	DropOrders()
	//	DropShipViews()
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
	Done() bool
	SetDone(bool)
	Score() int
	SetScore(int)
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
	Path() hexagon.CoordList
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
	Trail() hexagon.CoordList
}

type MapView interface {
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
