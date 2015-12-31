package db

import (
	"mule/hexagon"
	"mule/overpower"
)

type Source struct {
	db              DB
	MadePlanets     []*Planet
	MadePlanetViews map[[2]int]*PlanetView
	MadeShipViews   []*ShipView
	MadeMapViews    []*MapView
	MadeReports     []*Report
	DroppedShips    []overpower.Ship
}

func (d DB) NewSource() *Source {
	return &Source{
		db:              d,
		MadePlanets:     []*Planet{},
		MadePlanetViews: map[[2]int]*PlanetView{},
		MadeShipViews:   []*ShipView{},
		MadeMapViews:    []*MapView{},
		MadeReports:     []*Report{},
		DroppedShips:    []overpower.Ship{},
	}
}

func (s *Source) NewReport(gid, fid, turn int) (report overpower.Report, isOk bool) {
	r := NewReport()
	r.gid, r.fid, r.turn = gid, fid, turn
	s.MadeReports = append(s.MadeReports, r)
	return r, true
}
func (s *Source) NewMapView(gid, fid int, center hexagon.Coord) (mapview overpower.MapView, isOk bool) {
	mv := NewMapView()
	mv.gid, mv.fid, mv.center = gid, fid, center
	mv.zoom = overpower.DEFAULTZOOM
	mv.focusValid = false
	s.MadeMapViews = append(s.MadeMapViews, mv)
	return mv, true
}

func (s *Source) NewShipView(gid, fid, turn, sid, controller, size int, loc hexagon.Coord, locValid bool, trail []hexagon.Coord) (shipview overpower.ShipView, isOk bool) {
	sv := NewShipView()
	sv.gid, sv.fid, sv.turn, sv.sid, sv.controller, sv.size = gid, fid, turn, sid, controller, size
	if locValid {
		sv.locValid = true
		sv.loc = loc
	}
	sv.trail = trail
	s.MadeShipViews = append(s.MadeShipViews, sv)
	return sv, true
}

func (s *Source) NewShip(gid, fid, size, turn int, path []hexagon.Coord) (ship overpower.Ship, ok bool) {
	sh := NewShip()
	sh.gid, sh.fid, sh.size, sh.launched, sh.path = gid, fid, size, turn, path
	if !sh.Insert(s.db.db) {
		return nil, false
	}
	return sh, true
}

func (s *Source) DropShip(ship overpower.Ship) (ok bool) {
	s.DroppedShips = append(s.DroppedShips, ship)
	return true
}

func (s *Source) NewPlanet(name string, gid, pid, controller, inhab, res, parts int, loc hexagon.Coord) (planet overpower.Planet, ok bool) {
	p := NewPlanet()
	p.name = name
	p.gid = gid
	p.pid = pid
	p.controller = controller
	p.inhabitants = inhab
	p.resources = res
	p.parts = parts
	p.loc = loc
	s.MadePlanets = append(s.MadePlanets, p)
	return p, true
}

func (s *Source) NewPlanetView(fid, turn int, pl overpower.Planet) (planetview overpower.PlanetView, ok bool) {
	pv := NewPlanetView()
	pv.fid = fid
	pv.turn = turn
	if turn == 0 {
		pv.gid = pl.Gid()
		pv.name = pl.Name()
		pv.pid = pl.Pid()
		pv.loc = pl.Loc()
	} else {
		pv.Copy(pl)
	}
	s.MadePlanetViews[[2]int{pv.pid, pv.fid}] = pv
	return pv, true
}
