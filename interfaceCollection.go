package overpower

import (
	"mule/hexagon"
)

type Source interface {
	// -------- GET ------- //
	Game() (GameDat, error)
	Factions() ([]FactionDat, error)
	Planets() ([]PlanetDat, error)
	LaunchOrders() ([]LaunchOrderDat, error)
	Ships() ([]ShipDat, error)
	Truces() ([]TruceDat, error)
	PowerOrders() ([]PowerOrderDat, error)
	// ------- MAKE ------- //
	NewPlanet(name string,
		primaryFac, prPres, prPower,
		secondaryFac, sePres, sePower,
		antimatter, tachyons int,
		loc hexagon.Coord,
	) PlanetDat
	NewPlanetView(turn int, planet PlanetDat, exodus bool) PlanetViewDat
	NewMapView(fac int, center hexagon.Coord) MapViewDat
	NewShip(fid, sid, size, turn int, path hexagon.CoordList) ShipDat
	NewShipView(
		ship ShipDat, fid, turn int,
		loc, dest hexagon.NullCoord, trail hexagon.CoordList) ShipViewDat
	NewLaunchRecord(turn int, order LaunchOrderDat, ship ShipDat)
	NewBattleRecord(ship ShipDat, fid, turn,
		initPrimaryFac, initPrPres,
		initSecondaryFac, initSePres int,
		result PlanetDat,
		betrayals [][2]int,
	)
	NewPowerOrder(fid int, planet PlanetDat) PowerOrderDat
	// ------ CHANGE ----- //
	UpdatePlanetView(fid, turn int, planet PlanetDat) PlanetViewDat
	// ------- DROP ------ //
	ClearLaunchOrders() error
	//ClearPowerOrders() error
}

type BattleRecordGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	FID() int
	Loc() hexagon.Coord
	Turn() int
	Index() int
	PrimaryFaction() int
	PrimaryPresence() int
	SecondaryFaction() int
	SecondaryPresence() int
	ShipFaction() int
	ShipSize() int
	InitPrimaryFaction() int
	InitPrimaryPresence() int
	InitSecondaryFaction() int
	InitSecondaryPresence() int
	Betrayals() [][2]int
}
type BattleRecordSet interface {
	UnmarshalJSON([]byte) error
	DELETE()
}

type BattleRecordDat interface {
	BattleRecordGet
	BattleRecordSet
}

type FactionGet interface {
	MarshalJSON() ([]byte, error)
	SetFullJSON() // They say to never go full JSON...

	GID() int
	FID() int
	Owner() string
	Name() string
	DoneBuffer() int
	Score() int
	IsDone() bool
}
type FactionSet interface {
	UnmarshalJSON([]byte) error
	DELETE()

	SetDoneBuffer(int)
	SetScore(int)
}

type FactionDat interface {
	FactionGet
	FactionSet
}

type GameGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	Owner() string
	Name() string
	Turn() int
	//Autoturn() int
	AutoDays() [7]bool
	FreeAutos() int
	//Password() sql.NullString
	HasPassword() bool
	IsPassword(string) bool
	ToWin() int
	HighScore() int
	Winner() string
}
type GameSet interface {
	UnmarshalJSON([]byte) error
	DELETE()

	SetOwner(string)
	SetName(string)
	SetTurn(int)
	IncTurn()
	//SetAutoturn(int)
	SetAutoDays([7]bool)
	SetFreeAutos(int)
	//SetPassword(sql.NullString)
	SetToWin(int)
	SetHighScore(int)
}

type GameDat interface {
	GameGet
	GameSet
}

type LaunchRecordGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	FID() int
	Turn() int
	Source() hexagon.Coord
	Target() hexagon.Coord
	OrderSize() int
	Size() int
}
type LaunchRecordSet interface {
	UnmarshalJSON([]byte) error
	DELETE()
}

type LaunchRecordDat interface {
	LaunchRecordGet
	LaunchRecordSet
}

type MapViewGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	FID() int
	Center() hexagon.Coord
}
type MapViewSet interface {
	UnmarshalJSON([]byte) error
	DELETE()

	SetCenter(hexagon.Coord)
}

type MapViewDat interface {
	MapViewGet
	MapViewSet
}

type LaunchOrderGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	FID() int
	Source() hexagon.Coord
	Target() hexagon.Coord
	Size() int
}
type LaunchOrderSet interface {
	UnmarshalJSON([]byte) error
	DELETE()

	SetSize(int)
}

type LaunchOrderDat interface {
	LaunchOrderGet
	LaunchOrderSet
}

type PlanetGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	Loc() hexagon.Coord
	Name() string
	PrimaryFaction() int
	PrimaryPresence() int
	PrimaryPower() int
	SecondaryFaction() int
	SecondaryPresence() int
	SecondaryPower() int
	Antimatter() int
	Tachyons() int
}
type PlanetSet interface {
	UnmarshalJSON([]byte) error
	DELETE()

	SetPrimaryFaction(int)
	SetPrimaryPresence(int)
	SetPrimaryPower(int)
	SetSecondaryFaction(int)
	SetSecondaryPresence(int)
	SetSecondaryPower(int)
	SetAntimatter(int)
	SetTachyons(int)
}

type PlanetDat interface {
	PlanetGet
	PlanetSet
}

type PlanetViewGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	FID() int
	Loc() hexagon.Coord
	Name() string
	Turn() int
	PrimaryFaction() int
	PrimaryPresence() int
	PrimaryPower() int
	SecondaryFaction() int
	SecondaryPresence() int
	SecondaryPower() int
	Antimatter() int
	Tachyons() int
}
type PlanetViewSet interface {
	UnmarshalJSON([]byte) error
	DELETE()

	SetTurn(int)
	SetPrimaryFaction(int)
	SetPrimaryPresence(int)
	SetPrimaryPower(int)
	SetSecondaryFaction(int)
	SetSecondaryPresence(int)
	SetSecondaryPower(int)
	SetAntimatter(int)
	SetTachyons(int)
}

type PlanetViewDat interface {
	PlanetViewGet
	PlanetViewSet
}

type PowerOrderGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	FID() int
	Loc() hexagon.Coord
	UpPower() int
}

type PowerOrderSet interface {
	UnmarshalJSON([]byte) error
	DELETE()
}

type PowerOrderDat interface {
	PowerOrderGet
	PowerOrderSet
}

type ShipGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	FID() int
	SID() int
	Size() int
	Launched() int
	Path() hexagon.CoordList
}
type ShipSet interface {
	UnmarshalJSON([]byte) error
	DELETE()
}

type ShipDat interface {
	ShipGet
	ShipSet
}

type ShipViewGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	Turn() int
	FID() int
	SID() int
	Controller() int
	Size() int
	Loc() hexagon.NullCoord
	Dest() hexagon.NullCoord
	Trail() hexagon.CoordList
}
type ShipViewSet interface {
	UnmarshalJSON([]byte) error
	DELETE()
}

type ShipViewDat interface {
	ShipViewGet
	ShipViewSet
}

type TruceGet interface {
	MarshalJSON() ([]byte, error)

	GID() int
	FID() int
	Loc() hexagon.Coord
	Trucee() int
}
type TruceSet interface {
	UnmarshalJSON([]byte) error
	DELETE()
}

type TruceDat interface {
	TruceGet
	TruceSet
}
