package db

import (
	"mule/hexagon"
	"mule/mydb"
	"mule/overpower/setup"
)

/*
type Planet struct {
	Name      string
	Pid	      int
	Fid       int
	Inhab     int
	Resources int
	Parts	  int
	Loc       [2]int
}
*/

func (d DB) StartGame(gid int) (ok bool) {
	g, ok := d.GetGame(gid)
	if !ok {
		return
	}
	facs, ok := d.GetGidFactions(gid)
	if !ok {
		return
	}
	fids := make([]int, len(facs))
	for i, f := range facs {
		fids[i] = f.Fid()
	}
	planets := setup.MakeGalaxy(fids)
	pList := make([]*Planet, len(planets))
	pViewList := make([]*PlanetView, len(planets)*len(fids))
	var j int
	for i, pl := range planets {
		p := NewPlanet()
		p.pid = pl.Pid
		p.gid = g.Gid()
		p.name = pl.Name
		p.controller = pl.Fid
		p.inhabitants = pl.Inhab
		p.resources = pl.Resources
		p.parts = pl.Parts
		p.loc = hexagon.Coord(pl.Loc)
		pList[i] = p
		if !p.Insert(d.db) {
			return false
		}
		for k, fid := range fids {
			pv := NewPlanetView()
			pv.gid = p.gid
			pv.pid = p.pid
			pv.fid = fid
			pv.name = p.name
			pv.loc = p.loc
			if fid == p.controller {
				pv.controller = p.controller
				pv.inhabitants = p.inhabitants
				pv.resources = p.resources
				pv.parts = p.parts
				pv.turn = 1
			}
			pViewList[j+k] = pv
			if !pv.Insert(d.db) {
				return false
			}
		}
		j += len(fids)
	}
	g.IncTurn()
	return mydb.Update(d.db, []mydb.Updater{g})
}
