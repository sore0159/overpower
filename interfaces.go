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
	Truces() ([]Truce, error)
	PowerOrders() ([]PowerOrder, error)
	// ------- MAKE ------- //

	// name primaryFac, primaryPres, prPower, seFac, secondaryPres, sePower
	// antimatter, tachyons, loc
	NewPlanet(string, int, int, int, int, int, int, int, int, hexagon.Coord) Planet
	// turn planet exodus
	NewPlanetView(int, Planet, bool) PlanetView
	NewMapView(int, hexagon.Coord) MapView
	NewShip(int, int, int, hexagon.CoordList) Ship
	NewShipView(Ship, int, int, hexagon.NullCoord, hexagon.NullCoord, hexagon.CoordList) ShipView
	NewLaunchRecord(Order, Ship)
	// ship, fid, turn, initPrFac, initPrPres, initSeFac, initSecPre, resPlanet
	// droppedtruces
	NewBattleRecord(Ship, int, int, int, int, int, int, Planet, [][2]int)
	// ------ CHANGE ----- //

	// viewer, turn, planet
	UpdatePlanetView(int, int, Planet) PlanetView
	// ------- DROP ------ //
	DropShip(Ship)
	DropOrders()
	DropPowerOrders()
	// planet trucer trucee
	DropTruce(Planet, int, int)
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

type MapView interface {
	Gid() int
	Fid() int
	Center() hexagon.Coord
}

type Planet interface {
	Gid() int
	Loc() hexagon.Coord
	Name() string

	PrimaryFaction() int
	SetPrimaryFaction(int)
	PrimaryPresence() int
	SetPrimaryPresence(int)
	PrimaryPower() int
	SetPrimaryPower(int)

	SecondaryFaction() int
	SetSecondaryFaction(int)
	SecondaryPresence() int
	SetSecondaryPresence(int)
	SecondaryPower() int
	SetSecondaryPower(int)

	Antimatter() int
	SetAntimatter(int)
	Tachyons() int
	SetTachyons(int)
}

type PlanetView interface {
	Gid() int
	Fid() int
	Loc() hexagon.Coord
	Name() string
	Turn() int
	SetTurn(int)

	PrimaryFaction() int
	PrimaryPresence() int
	PrimaryPower() int

	SecondaryFaction() int
	SecondaryPresence() int
	SecondaryPower() int

	Antimatter() int
	Tachyons() int
}

type PowerOrder interface {
	Gid() int
	Fid() int
	Loc() hexagon.Coord
	UpPower() bool
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

type LaunchRecord interface {
	Gid() int
	Fid() int
	Turn() int
	Source() hexagon.Coord
	Target() hexagon.Coord
	Size() int
}

type BattleRecord interface {
	Gid() int
	Fid() int
	Turn() int
	Index() int
	Loc() hexagon.Coord

	ShipFaction() int
	ShipSize() int

	InitPrimaryFaction() int
	InitSecondaryFaction() int
	InitPrimaryPresence() int
	InitSecondaryPresence() int

	ResultPrimaryFaction() int
	ResultSecondaryFaction() int
	ResultPrimaryPresence() int
	ResultSecondaryPresence() int

	Betrayals() [][2]int
}

type Truce interface {
	Gid() int
	Fid() int
	Loc() hexagon.Coord
	With() int
}
