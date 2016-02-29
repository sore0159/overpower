package db

import (
	"database/sql"
	"mule/hexagon"
	"mule/mydb"
	"mule/overpower"
)

var sourceT overpower.Source = &Source{}

// Usage Example:
// f := func(d DB) error {
//		source := d.NewSource(gid)
//		breakE, logE := overpower.RunGameTurn(source)
//		if logE != nil {
//			Log(logE)
//		}
//		if my, bad := Check(breakE, "rungameturn failed", "gid", gid); bad {
//			return my
//		}
//		err := source.Commit()
//		if my, bad := Check(err, "rungame turn commit failed", "gid", gid); bad {
//			return my
//		}
//		return nil
// }
// err := db.Transact(OPDB, f)
type Source struct {
	Gid int
	db  DB
	// ------ CHANGE ------ //
	game        *Game
	planets     []*Planet
	factions    []*Faction
	planetViews map[[3]int]*PlanetView
	// ------- MAKE ------- //
	MadePlanets       []*Planet
	MadePlanetViews   map[[3]int]*PlanetView
	MadeShips         []*Ship
	MadeShipViews     map[overpower.Ship][]*ShipView
	MadeMapViews      []*MapView
	MadeLaunchRecords []*LaunchRecord
	MadeBattleRecords map[int][]*BattleRecord
	// ------- DROP ------ //
	DroppedShips    []overpower.Ship
	DroppedTruces   []overpower.Truce
	dropOrders      bool
	dropPowerOrders bool
}

func (d DB) NewSource(gid int) *Source {
	return &Source{
		Gid:               gid,
		db:                d,
		planetViews:       map[[3]int]*PlanetView{},
		MadePlanets:       []*Planet{},
		MadePlanetViews:   map[[3]int]*PlanetView{},
		MadeShips:         []*Ship{},
		MadeShipViews:     map[overpower.Ship][]*ShipView{},
		MadeMapViews:      []*MapView{},
		MadeLaunchRecords: []*LaunchRecord{},
		MadeBattleRecords: map[int][]*BattleRecord{},
		DroppedShips:      []overpower.Ship{},
		DroppedTruces:     []overpower.Truce{},
	}
}

func (s *Source) Commit() error {
	var err error
	// ------- DROP ------ //
	if len(s.DroppedTruces) > 0 {
		err = s.db.DropOPTruces(s.DroppedTruces...)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "drop", "truces", "item", s.DroppedTruces); bad {
			return my
		}
	}
	if s.dropOrders {
		err = s.db.dropItemsIf("orders", C{"gid", s.Gid})
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "drop", "orders"); bad {
			return my
		}
	}
	if s.dropPowerOrders {
		err = s.db.dropItemsIf("powerorders", C{"gid", s.Gid})
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "drop", "powerorders"); bad {
			return my
		}
	}
	if len(s.DroppedShips) > 0 {
		err = s.db.dropTheseShips(s.DroppedShips)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "drop", "ships", "list", s.DroppedShips); bad {
			return my
		}
	}
	// ------- MAKE ------- //
	if len(s.MadePlanets) > 0 {
		group := &PlanetGroup{s.MadePlanets}
		err = s.db.makeGroup(group)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "make", "Planets", "list", s.MadePlanets); bad {
			return my
		}
	}
	if len(s.MadePlanetViews) > 0 {
		allPV := make([]*PlanetView, 0, len(s.MadePlanetViews))
		for _, pv := range s.MadePlanetViews {
			allPV = append(allPV, pv)
		}

		group := &PlanetViewGroup{allPV}
		err = s.db.makeGroup(group)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "make", "PlanetViews", "list", allPV); bad {
			return my
		}
	}
	if len(s.MadeShips) > 0 {
		group := &ShipGroup{s.MadeShips}
		err = s.db.makeGroup(group)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "make", "Ships", "list", s.MadeShips); bad {
			return my
		}
	}
	if len(s.MadeLaunchRecords) > 0 {
		group := &LaunchRecordGroup{s.MadeLaunchRecords}
		err = s.db.makeGroup(group)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "make", "LaunchRecords", "list", s.MadeLaunchRecords); bad {
			return my
		}
	}
	if len(s.MadeBattleRecords) > 0 {
		allLR := make([]*BattleRecord, 0, len(s.MadeBattleRecords))
		for _, list := range s.MadeBattleRecords {
			allLR = append(allLR, list...)
		}
		group := &BattleRecordGroup{allLR}
		err = s.db.makeGroup(group)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "make", "landingrecords", "list", allLR); bad {
			return my
		}
	}
	if len(s.MadeShipViews) > 0 {
		allSVs := make([]*ShipView, 0, len(s.MadeShipViews))
		for sh, list := range s.MadeShipViews {
			sid := sh.Sid()
			for _, shV := range list {
				// let's hope any new ships have thier sids by now
				shV.sid = sid
				allSVs = append(allSVs, shV)
			}
		}
		group := &ShipViewGroup{allSVs}
		err = s.db.makeGroup(group)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "make", "shipviews", "list", allSVs); bad {
			return my
		}
	}
	if len(s.MadeMapViews) > 0 {
		group := &MapViewGroup{s.MadeMapViews}
		err = s.db.makeGroup(group)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "make", "mapviews", "list", s.MadeMapViews); bad {
			return my
		}
	}
	// ------ CHANGE ------ //
	plVs := make([]*PlanetView, 0, len(s.planetViews))
	for _, item := range s.planetViews {
		plVs = append(plVs, item)
	}
	for _, group := range []mydb.UpdateGrouper{
		&GameGroup{[]*Game{s.game}},
		&PlanetGroup{s.planets},
		&FactionGroup{s.factions},
		&PlanetViewGroup{plVs},
	} {
		err = s.db.updateGroup(group)
		if my, bad := Check(err, "source commit error", "gid", s.Gid, "group", group); bad {
			return my
		}
	}
	return nil
}

func (s *Source) Game() (overpower.Game, error) {
	group := NewGameGroup()
	err := s.db.getGroup(group, C{"gid", s.Gid})
	if my, bad := Check(err, "source get Games failure", "gid", s.Gid); bad {
		return nil, my
	}
	list := group.List
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	s.game = list[0]
	return list[0], nil
}

func (s *Source) Factions() ([]overpower.Faction, error) {
	group := NewFactionGroup()
	err := s.db.getGroup(group, C{"gid", s.Gid})
	if my, bad := Check(err, "source get factions failure", "gid", s.Gid); bad {
		return nil, my
	}
	list := group.List
	s.factions = list
	return convertFactions2OP(list...), nil
}

func (s *Source) Planets() ([]overpower.Planet, error) {
	group := NewPlanetGroup()
	err := s.db.getGroup(group, C{"gid", s.Gid})
	if my, bad := Check(err, "source get Planets failure", "gid", s.Gid); bad {
		return nil, my
	}
	list := group.List
	s.planets = list
	return convertPlanets2OP(list...), nil
}

func (s *Source) Orders() ([]overpower.Order, error) {
	return s.db.GetOrders("gid", s.Gid)
}
func (s *Source) Ships() ([]overpower.Ship, error) {
	return s.db.GetShips("gid", s.Gid)
}
func (s *Source) PowerOrders() ([]overpower.PowerOrder, error) {
	return s.db.GetPowerOrders("gid", s.Gid)
}
func (s *Source) Truces() ([]overpower.Truce, error) {
	return nil, nil
}

func (s *Source) NewMapView(fid int, center hexagon.Coord) (mapview overpower.MapView) {
	mv := NewMapView()
	mv.gid, mv.fid, mv.center = s.Gid, fid, center
	s.MadeMapViews = append(s.MadeMapViews, mv)
	return mv
}

func (s *Source) NewShipView(sh overpower.Ship, fid, turn int, loc, dest hexagon.NullCoord, trail hexagon.CoordList) (shipview overpower.ShipView) {
	sv := NewShipView()
	sv.gid = s.Gid
	sv.turn, sv.loc, sv.dest, sv.trail, sv.fid = turn, loc, dest, trail, fid
	sv.controller, sv.size, sv.sid = sh.Fid(), sh.Size(), sh.Sid()
	if list, ok := s.MadeShipViews[sh]; ok {
		s.MadeShipViews[sh] = append(list, sv)
	} else {
		s.MadeShipViews[sh] = []*ShipView{sv}
	}
	return sv
}

func (s *Source) NewShip(fid, size, turn int, path hexagon.CoordList) (ship overpower.Ship) {
	sh := NewShip()
	sh.gid, sh.fid, sh.size, sh.launched, sh.path = s.Gid, fid, size, turn, path
	s.MadeShips = append(s.MadeShips, sh)
	sh.justmade = true
	return sh
}

func (s *Source) DropShip(ship overpower.Ship) {
	s.DroppedShips = append(s.DroppedShips, ship)
}
func (s *Source) DropOrders() {
	s.dropOrders = true
}
func (s *Source) DropPowerOrders() {
	s.dropPowerOrders = true
}

func (s *Source) DropTruce(pl overpower.Planet, trucer, trucee int) {
	tr := &Truce{
		gid:    pl.Gid(),
		fid:    trucer,
		loc:    pl.Loc(),
		trucee: trucee,
	}
	s.DroppedTruces = append(s.DroppedTruces, tr)
}

// name primaryFac, primaryPres, prPower, seFac, secondaryPres, sePower
// antimatter, tachyons, loc
func (s *Source) NewPlanet(name string, pF, pPr, pPw, sF, sPr, sPw, antiM, tach int, loc hexagon.Coord) (planet overpower.Planet) {
	p := NewPlanet()
	p.gid = s.Gid
	p.loc = loc
	if pF != 0 {
		p.primaryfaction = sql.NullInt64{Valid: true, Int64: int64(pF)}
	}
	if sF != 0 {
		p.secondaryfaction = sql.NullInt64{Valid: true, Int64: int64(sF)}
	}
	p.primarypresence = pPr
	p.primarypower = pPw
	p.secondarypresence = sPr
	p.secondarypower = sPw
	p.antimatter = antiM
	p.tachyons = tach

	s.MadePlanets = append(s.MadePlanets, p)
	return p
}

func (s *Source) NewPlanetView(fid int, pl overpower.Planet, exodus bool) (planetview overpower.PlanetView) {
	pv := NewPlanetView()

	pv.gid = s.Gid
	pv.loc = pl.Loc()
	pv.name = pl.Name()
	pv.fid = fid

	if pl.PrimaryFaction() == fid || pl.SecondaryFaction() == fid {
		pv.turn = 1
		pF := pl.PrimaryFaction()
		sF := pl.SecondaryFaction()
		if pF != 0 {
			pv.primaryfaction = sql.NullInt64{Valid: true, Int64: int64(pF)}
		}
		if sF != 0 {
			pv.secondaryfaction = sql.NullInt64{Valid: true, Int64: int64(sF)}
		}
		pv.primarypresence = pl.PrimaryPresence()
		pv.primarypower = pl.PrimaryPower()
		pv.secondarypresence = pl.SecondaryPresence()
		pv.secondarypower = pl.SecondaryPower()
		pv.antimatter = pl.Antimatter()
		pv.tachyons = pl.Tachyons()
	}

	s.MadePlanetViews[[3]int{pv.loc[0], pv.loc[1], pv.fid}] = pv
	return pv
}

func (s *Source) UpdatePlanetView(fid, turn int, pl overpower.Planet) overpower.PlanetView {
	pv := NewPlanetView()
	pv.gid = s.Gid
	pv.fid = fid
	pv.loc = pl.Loc()
	pv.name = pl.Name()
	pv.turn = turn
	pF := pl.PrimaryFaction()
	sF := pl.SecondaryFaction()
	if pF != 0 {
		pv.primaryfaction = sql.NullInt64{Valid: true, Int64: int64(pF)}
	}
	if sF != 0 {
		pv.secondaryfaction = sql.NullInt64{Valid: true, Int64: int64(sF)}
	}
	pv.primarypresence = pl.PrimaryPresence()
	pv.primarypower = pl.PrimaryPower()
	pv.secondarypresence = pl.SecondaryPresence()
	pv.secondarypower = pl.SecondaryPower()
	pv.antimatter = pl.Antimatter()
	pv.tachyons = pl.Tachyons()

	s.planetViews[[3]int{pv.loc[0], pv.loc[1], pv.fid}] = pv
	return pv
}

func (s *Source) NewLaunchRecord(turn int, o overpower.Order, ship overpower.Ship) {
	lr := NewLaunchRecord()
	lr.gid = s.Gid
	lr.fid = o.Fid()
	lr.turn = turn
	lr.source = o.Source()
	lr.target = o.Target()
	lr.order = o.Size()
	if ship != nil {
		lr.size = ship.Size()
	}
	s.MadeLaunchRecords = append(s.MadeLaunchRecords, lr)
}

func (s *Source) NewBattleRecord(ship overpower.Ship, fid, turn, initPrFac, initPrPres, initSeFac, initSePres int, pl overpower.Planet, betrayals [][2]int) {
	lr := NewBattleRecord()
	list, ok := s.MadeBattleRecords[fid]

	//shipfaction           sql.NullInt64
	//shipsize              int
	//initprimaryfaction    sql.NullInt64
	//initprimarypresence   int
	//initsecondaryfaction  sql.NullInt64
	//initsecondarypresence int
	//betrayals mydb.IntList
	lr.gid = ship.Gid()
	lr.fid = fid
	lr.turn = turn
	lr.loc = pl.Loc()
	btr := make([]int, 0, len(betrayals)*2)
	for _, pt := range betrayals {
		btr = append(btr, pt[0], pt[1])
	}
	lr.betrayals = btr

	if ship != nil {
		lr.shipfaction = sql.NullInt64{Valid: true, Int64: int64(ship.Fid())}
		lr.shipsize = ship.Size()
	}
	if initPrFac != 0 {
		lr.initprimaryfaction = sql.NullInt64{Valid: true, Int64: int64(initPrFac)}
	}
	lr.initprimarypresence = initPrPres
	if initSeFac != 0 {
		lr.initsecondaryfaction = sql.NullInt64{Valid: true, Int64: int64(initSeFac)}
	}
	lr.initsecondarypresence = initSePres

	if resPrFac := pl.PrimaryFaction(); resPrFac != 0 {
		lr.primaryfaction = sql.NullInt64{Valid: true, Int64: int64(resPrFac)}
	}
	lr.primarypresence = pl.PrimaryPresence()
	if resSeFac := pl.SecondaryFaction(); resSeFac != 0 {
		lr.secondaryfaction = sql.NullInt64{Valid: true, Int64: int64(resSeFac)}
	}
	lr.secondarypresence = pl.SecondaryPresence()

	if !ok {
		s.MadeBattleRecords[fid] = []*BattleRecord{lr}
	} else {
		lr.index = len(list)
		s.MadeBattleRecords[fid] = append(list, lr)
	}
}
