package models

import (
	"database/sql"
	"mule/hexagon"
	sq "mule/mydb/sql"
	"mule/overpower"
)

//var sourceTest overpower.Source = &Source{}

type Source struct {
	GID        int
	M          *Manager
	Where      sq.Condition
	updatedPVs map[[3]int]int
}

func NewSource(m *Manager, gid int) *Source {
	return &Source{
		GID:        gid,
		M:          m,
		Where:      sq.EQ("gid", gid),
		updatedPVs: map[[3]int]int{},
	}
}

func (s *Source) Game() (overpower.GameDat, error) {
	list, err := s.M.Game().SelectWhere(s.Where)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrNoneFound
	}
	return list[0], nil
}

func (s *Source) Factions() ([]overpower.FactionDat, error) {
	return s.M.Faction().SelectWhere(s.Where)
}
func (s *Source) Planets() ([]overpower.PlanetDat, error) {
	return s.M.Planet().SelectWhere(s.Where)
}
func (s *Source) LaunchOrders() ([]overpower.LaunchOrderDat, error) {
	return s.M.LaunchOrder().SelectWhere(s.Where)
}
func (s *Source) Ships() ([]overpower.ShipDat, error) {
	return s.M.Ship().SelectWhere(s.Where)
}
func (s *Source) Truces() ([]overpower.TruceDat, error) {
	return s.M.Truce().SelectWhere(s.Where)
}
func (s *Source) PowerOrders() ([]overpower.PowerOrderDat, error) {
	return s.M.PowerOrder().SelectWhere(s.Where)
}

func (s *Source) UpdatePlanetView(fid, turn int, planet overpower.PlanetDat) overpower.PlanetViewDat {
	pv := &PlanetView{
		GID:  s.GID,
		FID:  fid,
		Turn: turn,
		//
		Name:              planet.Name(),
		Loc:               planet.Loc(),
		PrimaryPresence:   planet.PrimaryPresence(),
		PrimaryPower:      planet.PrimaryPower(),
		SecondaryPresence: planet.SecondaryPresence(),
		SecondaryPower:    planet.SecondaryPower(),
		Antimatter:        planet.Antimatter(),
		Tachyons:          planet.Tachyons(),
	}
	pF := planet.PrimaryFaction()
	sF := planet.SecondaryFaction()
	if pF != 0 {
		pv.PrimaryFaction = sql.NullInt64{Valid: true, Int64: int64(pF)}
	}
	if sF != 0 {
		pv.SecondaryFaction = sql.NullInt64{Valid: true, Int64: int64(sF)}
	}
	pv.sql.UPDATE = true
	//
	pt := [3]int{pv.Loc[0], pv.Loc[1], fid}
	sess := s.M.PlanetView()
	if i, ok := s.updatedPVs[pt]; ok {
		sess.List[i] = pv
	} else {
		s.updatedPVs[pt] = len(sess.List)
		sess.List = append(sess.List, pv)
	}
	sess.List = append(sess.List, pv)
	return pv.Intf()
}

func (s *Source) ClearLaunchOrders() error {
	return s.M.LaunchOrder().Delete(s.Where)
}
func (s *Source) ClearPowerOrders() error {
	return s.M.PowerOrder().Delete(s.Where)
}

func (s *Source) NewPlanet(name string,
	primaryFac, prPres, prPower,
	secondaryFac, sePres, sePower,
	antimatter, tachyons int,
	loc hexagon.Coord,
) overpower.PlanetDat {
	pl := &Planet{
		GID:               s.GID,
		Name:              name,
		Loc:               loc,
		PrimaryPresence:   prPres,
		PrimaryPower:      prPower,
		SecondaryPresence: sePres,
		SecondaryPower:    sePower,
		Antimatter:        antimatter,
		Tachyons:          tachyons,
	}
	if primaryFac != 0 {
		pl.PrimaryFaction = sql.NullInt64{Valid: true, Int64: int64(primaryFac)}
	}
	if secondaryFac != 0 {
		pl.SecondaryFaction = sql.NullInt64{Valid: true, Int64: int64(secondaryFac)}
	}
	s.M.CreatePlanet(pl)
	return pl.Intf()
}

func (s *Source) NewPlanetView(fid int, pl overpower.PlanetDat, exodus bool) overpower.PlanetViewDat {
	pv := &PlanetView{
		GID:  s.GID,
		FID:  fid,
		Loc:  pl.Loc(),
		Name: pl.Name(),
	}

	if pl.PrimaryFaction() == fid || pl.SecondaryFaction() == fid {
		pv.Turn = 1
		pF := pl.PrimaryFaction()
		sF := pl.SecondaryFaction()
		if pF != 0 {
			pv.PrimaryFaction = sql.NullInt64{Valid: true, Int64: int64(pF)}
		}
		if sF != 0 {
			pv.SecondaryFaction = sql.NullInt64{Valid: true, Int64: int64(sF)}
		}
		pv.PrimaryPresence = pl.PrimaryPresence()
		pv.PrimaryPower = pl.PrimaryPower()
		pv.SecondaryPresence = pl.SecondaryPresence()
		pv.SecondaryPower = pl.SecondaryPower()
		pv.Antimatter = pl.Antimatter()
		pv.Tachyons = pl.Tachyons()
	}
	s.M.CreatePlanetView(pv)
	return pv.Intf()
}

func (s *Source) NewMapView(fid int, center hexagon.Coord) overpower.MapViewDat {
	mv := &MapView{
		GID:    s.GID,
		FID:    fid,
		Center: center,
	}
	s.M.CreateMapView(mv)
	return mv.Intf()
}

func (s *Source) NewShip(fid, sid, size, turn int, path hexagon.CoordList) overpower.ShipDat {
	sh := &Ship{
		GID:      s.GID,
		FID:      fid,
		SID:      sid,
		Size:     size,
		Launched: turn,
		Path:     path,
	}
	s.M.CreateShip(sh)
	return sh.Intf()
}

func (s *Source) NewShipView(sh overpower.ShipDat, fid, turn int, loc, dest hexagon.NullCoord, trail hexagon.CoordList) overpower.ShipViewDat {
	sv := &ShipView{
		GID:        s.GID,
		FID:        fid,
		Turn:       turn,
		Loc:        loc,
		Dest:       dest,
		Trail:      trail,
		Controller: sh.FID(),
		SID:        sh.SID(),
		Size:       sh.Size(),
	}
	s.M.CreateShipView(sv)
	return sv.Intf()
}

func (s *Source) NewLaunchRecord(turn int, o overpower.LaunchOrderDat, ship overpower.ShipDat) {
	lr := &LaunchRecord{
		GID:       s.GID,
		FID:       o.FID(),
		Turn:      turn,
		Source:    o.Source(),
		Target:    o.Target(),
		OrderSize: o.Size(),
	}
	if ship != nil {
		lr.Size = ship.Size()
	}
	s.M.CreateLaunchRecord(lr)
}
func (s *Source) NewBattleRecord(ship overpower.ShipDat, fid, turn,
	initPrimaryFac, initPrPres,
	initSecondaryFac, initSePres int,
	result overpower.PlanetDat,
	betrayals [][2]int,
) {
	btr := make([]int, 0, len(betrayals)*2)
	for _, pt := range betrayals {
		btr = append(btr, pt[0], pt[1])
	}

	br := &BattleRecord{
		GID:       s.GID,
		FID:       fid,
		Turn:      turn,
		Loc:       result.Loc(),
		Betrayals: btr,

		InitPrimaryPresence:   initPrPres,
		InitSecondaryPresence: initSePres,
		PrimaryPresence:       result.PrimaryPresence(),
		SecondaryPresence:     result.SecondaryPresence(),
	}
	if ship != nil {
		br.ShipFaction = sql.NullInt64{Valid: true, Int64: int64(ship.FID())}
		br.ShipSize = ship.Size()
	}
	if initPrimaryFac != 0 {
		br.InitPrimaryFaction = sql.NullInt64{Valid: true, Int64: int64(initPrimaryFac)}
	}
	if initSecondaryFac != 0 {
		br.InitSecondaryFaction = sql.NullInt64{Valid: true, Int64: int64(initSecondaryFac)}
	}

	if resPrFac := result.PrimaryFaction(); resPrFac != 0 {
		br.PrimaryFaction = sql.NullInt64{Valid: true, Int64: int64(resPrFac)}
	}
	if resSeFac := result.SecondaryFaction(); resSeFac != 0 {
		br.SecondaryFaction = sql.NullInt64{Valid: true, Int64: int64(resSeFac)}
	}

	s.M.CreateBattleRecord(br)
}

func (s *Source) NewPowerOrder(fid int, planet overpower.PlanetDat) overpower.PowerOrderDat {
	po := &PowerOrder{
		GID: s.GID,
		FID: fid,
		Loc: planet.Loc(),
	}
	s.M.CreatePowerOrder(po)
	return po.Intf()
}
