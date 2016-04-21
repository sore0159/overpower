package models

import (
	"mule/mybad"
	"mule/mydb/db"
	sq "mule/mydb/sql"
)

var (
	Check = mybad.BuildCheck("package", "models")
)

type Manager struct {
	D                   db.DBer
	BattleRecordSession *BattleRecordSession
	FactionSession      *FactionSession
	GameSession         *GameSession
	LaunchRecordSession *LaunchRecordSession
	MapViewSession      *MapViewSession
	LaunchOrderSession  *LaunchOrderSession
	PlanetSession       *PlanetSession
	PlanetViewSession   *PlanetViewSession
	PowerOrderSession   *PowerOrderSession
	ShipSession         *ShipSession
	ShipViewSession     *ShipViewSession
	TruceSession        *TruceSession
}

func NewManager(d db.DBer) *Manager {
	return &Manager{D: d}
}

func (m *Manager) GID(gid int) sq.Condition {
	return sq.EQ("gid", gid)
}
func (m *Manager) FID(gid, fid int) sq.Condition {
	return sq.AND(sq.EQ("gid", gid), sq.EQ("fid", fid))
}
func (m *Manager) TURN(gid, fid, turn int) sq.Condition {
	return sq.AND(sq.EQ("gid", gid), sq.EQ("fid", fid), sq.EQ("turn", turn))
}

func (m *Manager) BattleRecord() *BattleRecordSession {
	s := NewBattleRecordSession(m.D)
	m.BattleRecordSession = s
	return s
}

func (m *Manager) CreateBattleRecord(item *BattleRecord) {
	if m.BattleRecordSession == nil {
		m.BattleRecordSession = NewBattleRecordSession(m.D)
	}
	item.sql.INSERT = true
	m.BattleRecordSession.List = append(m.BattleRecordSession.List, item)
}

func (m *Manager) Faction() *FactionSession {
	s := NewFactionSession(m.D)
	m.FactionSession = s
	return s
}

func (m *Manager) CreateFaction(item *Faction) {
	if m.FactionSession == nil {
		m.FactionSession = NewFactionSession(m.D)
	}
	item.sql.INSERT = true
	m.FactionSession.List = append(m.FactionSession.List, item)
}

func (m *Manager) Game() *GameSession {
	s := NewGameSession(m.D)
	m.GameSession = s
	return s
}

func (m *Manager) CreateGame(item *Game) {
	if m.GameSession == nil {
		m.GameSession = NewGameSession(m.D)
	}
	item.sql.INSERT = true
	m.GameSession.List = append(m.GameSession.List, item)
}

func (m *Manager) LaunchRecord() *LaunchRecordSession {
	s := NewLaunchRecordSession(m.D)
	m.LaunchRecordSession = s
	return s
}

func (m *Manager) CreateLaunchRecord(item *LaunchRecord) {
	if m.LaunchRecordSession == nil {
		m.LaunchRecordSession = NewLaunchRecordSession(m.D)
	}
	item.sql.INSERT = true
	m.LaunchRecordSession.List = append(m.LaunchRecordSession.List, item)
}

func (m *Manager) MapView() *MapViewSession {
	s := NewMapViewSession(m.D)
	m.MapViewSession = s
	return s
}

func (m *Manager) CreateMapView(item *MapView) {
	if m.MapViewSession == nil {
		m.MapViewSession = NewMapViewSession(m.D)
	}
	item.sql.INSERT = true
	m.MapViewSession.List = append(m.MapViewSession.List, item)
}

func (m *Manager) LaunchOrder() *LaunchOrderSession {
	s := NewLaunchOrderSession(m.D)
	m.LaunchOrderSession = s
	return s
}

func (m *Manager) CreateLaunchOrder(item *LaunchOrder) {
	if m.LaunchOrderSession == nil {
		m.LaunchOrderSession = NewLaunchOrderSession(m.D)
	}
	item.sql.INSERT = true
	m.LaunchOrderSession.List = append(m.LaunchOrderSession.List, item)
}

func (m *Manager) Planet() *PlanetSession {
	s := NewPlanetSession(m.D)
	m.PlanetSession = s
	return s
}

func (m *Manager) CreatePlanet(item *Planet) {
	if m.PlanetSession == nil {
		m.PlanetSession = NewPlanetSession(m.D)
	}
	item.sql.INSERT = true
	m.PlanetSession.List = append(m.PlanetSession.List, item)
}

func (m *Manager) PlanetView() *PlanetViewSession {
	s := NewPlanetViewSession(m.D)
	m.PlanetViewSession = s
	return s
}

func (m *Manager) CreatePlanetView(item *PlanetView) {
	if m.PlanetViewSession == nil {
		m.PlanetViewSession = NewPlanetViewSession(m.D)
	}
	item.sql.INSERT = true
	m.PlanetViewSession.List = append(m.PlanetViewSession.List, item)
}

func (m *Manager) PowerOrder() *PowerOrderSession {
	s := NewPowerOrderSession(m.D)
	m.PowerOrderSession = s
	return s
}

func (m *Manager) CreatePowerOrder(item *PowerOrder) {
	if m.PowerOrderSession == nil {
		m.PowerOrderSession = NewPowerOrderSession(m.D)
	}
	item.sql.INSERT = true
	m.PowerOrderSession.List = append(m.PowerOrderSession.List, item)
}

func (m *Manager) Ship() *ShipSession {
	s := NewShipSession(m.D)
	m.ShipSession = s
	return s
}

func (m *Manager) CreateShip(item *Ship) {
	if m.ShipSession == nil {
		m.ShipSession = NewShipSession(m.D)
	}
	item.sql.INSERT = true
	m.ShipSession.List = append(m.ShipSession.List, item)
}

func (m *Manager) ShipView() *ShipViewSession {
	s := NewShipViewSession(m.D)
	m.ShipViewSession = s
	return s
}

func (m *Manager) CreateShipView(item *ShipView) {
	if m.ShipViewSession == nil {
		m.ShipViewSession = NewShipViewSession(m.D)
	}
	item.sql.INSERT = true
	m.ShipViewSession.List = append(m.ShipViewSession.List, item)
}

func (m *Manager) Truce() *TruceSession {
	s := NewTruceSession(m.D)
	m.TruceSession = s
	return s
}

func (m *Manager) CreateTruce(item *Truce) {
	if m.TruceSession == nil {
		m.TruceSession = NewTruceSession(m.D)
	}
	item.sql.INSERT = true
	m.TruceSession.List = append(m.TruceSession.List, item)
}

func (m *Manager) Close() error {
	var err error
	if m.BattleRecordSession != nil {
		err = m.BattleRecordSession.Close()
		if my, bad := Check(err, "manager close failure on BattleRecord Close"); bad {
			return my
		}
		m.BattleRecordSession = nil
	}

	if m.FactionSession != nil {
		err = m.FactionSession.Close()
		if my, bad := Check(err, "manager close failure on Faction Close"); bad {
			return my
		}
		m.FactionSession = nil
	}

	if m.GameSession != nil {
		err = m.GameSession.Close()
		if my, bad := Check(err, "manager close failure on Game Close"); bad {
			return my
		}
		m.GameSession = nil
	}

	if m.LaunchRecordSession != nil {
		err = m.LaunchRecordSession.Close()
		if my, bad := Check(err, "manager close failure on LaunchRecord Close"); bad {
			return my
		}
		m.LaunchRecordSession = nil
	}

	if m.MapViewSession != nil {
		err = m.MapViewSession.Close()
		if my, bad := Check(err, "manager close failure on MapView Close"); bad {
			return my
		}
		m.MapViewSession = nil
	}

	if m.LaunchOrderSession != nil {
		err = m.LaunchOrderSession.Close()
		if my, bad := Check(err, "manager close failure on LaunchOrder Close"); bad {
			return my
		}
		m.LaunchOrderSession = nil
	}

	if m.PlanetSession != nil {
		err = m.PlanetSession.Close()
		if my, bad := Check(err, "manager close failure on Planet Close"); bad {
			return my
		}
		m.PlanetSession = nil
	}

	if m.PlanetViewSession != nil {
		err = m.PlanetViewSession.Close()
		if my, bad := Check(err, "manager close failure on PlanetView Close"); bad {
			return my
		}
		m.PlanetViewSession = nil
	}

	if m.PowerOrderSession != nil {
		err = m.PowerOrderSession.Close()
		if my, bad := Check(err, "manager close failure on PowerOrder Close"); bad {
			return my
		}
		m.PowerOrderSession = nil
	}

	if m.ShipSession != nil {
		err = m.ShipSession.Close()
		if my, bad := Check(err, "manager close failure on Ship Close"); bad {
			return my
		}
		m.ShipSession = nil
	}

	if m.ShipViewSession != nil {
		err = m.ShipViewSession.Close()
		if my, bad := Check(err, "manager close failure on ShipView Close"); bad {
			return my
		}
		m.ShipViewSession = nil
	}

	if m.TruceSession != nil {
		err = m.TruceSession.Close()
		if my, bad := Check(err, "manager close failure on Truce Close"); bad {
			return my
		}
		m.TruceSession = nil
	}

	return nil
}

func CreateAllTables(d db.DBer) error {
	var err error

	err = GameTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table Game"); bad {
		return my
	}

	err = FactionTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table Faction"); bad {
		return my
	}

	err = PlanetTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table Planet"); bad {
		return my
	}

	err = PlanetViewTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table PlanetView"); bad {
		return my
	}

	err = BattleRecordTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table BattleRecord"); bad {
		return my
	}

	err = LaunchRecordTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table LaunchRecord"); bad {
		return my
	}

	err = MapViewTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table MapView"); bad {
		return my
	}

	err = LaunchOrderTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table LaunchOrder"); bad {
		return my
	}

	err = PowerOrderTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table PowerOrder"); bad {
		return my
	}

	err = ShipTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table Ship"); bad {
		return my
	}

	err = ShipViewTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table ShipView"); bad {
		return my
	}

	err = TruceTableCreate(d)
	if my, bad := Check(err, "Create all tables failure on table Truce"); bad {
		return my
	}

	return nil
}

func DropAllTables(d db.DBer) error {
	var err error
	err = BattleRecordTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table BattleRecord"); bad {
		return my
	}

	err = FactionTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table Faction"); bad {
		return my
	}

	err = GameTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table Game"); bad {
		return my
	}

	err = LaunchRecordTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table LaunchRecord"); bad {
		return my
	}

	err = MapViewTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table MapView"); bad {
		return my
	}

	err = LaunchOrderTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table LaunchOrder"); bad {
		return my
	}

	err = PlanetTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table Planet"); bad {
		return my
	}

	err = PlanetViewTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table PlanetView"); bad {
		return my
	}

	err = PowerOrderTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table PowerOrder"); bad {
		return my
	}

	err = ShipTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table Ship"); bad {
		return my
	}

	err = ShipViewTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table ShipView"); bad {
		return my
	}

	err = TruceTableDelete(d)
	if my, bad := Check(err, "Delete all tables failure on table Truce"); bad {
		return my
	}

	return nil
}
