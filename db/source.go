package db

import (
	"mule/hexagon"
	"mule/overpower"
)

type Source struct {
	db              DB
	MadePlanets     []*Planet
	MadePlanetViews []*PlanetView
}

func (d DB) NewSource() *Source {
	return &Source{
		db:              d,
		MadePlanets:     []*Planet{},
		MadePlanetViews: []*PlanetView{},
	}
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
	s.MadePlanetViews = append(s.MadePlanetViews, pv)
	return pv, true
}
