package db

import (
	"mule/mydb"
	"mule/overpower"
)

func (d DB) SetTurnDone(f overpower.Faction, done bool) (ok bool) {
	f.SetDone(done)
	return mydb.Update(d.db, []mydb.Updater{f})
}

func (d DB) RunGameTurn(g overpower.Game) (ok bool) {
	gid := g.Gid()
	facs, ok := d.GetGidFactions(gid)
	if !ok {
		return false
	}
	planets, ok := d.GetAllGidPlanets(gid)
	if !ok {
		return false
	}
	orders, ok := d.GetAllGidOrders(gid)
	if !ok {
		return false
	}
	ships, ok := d.GetAllGidShips(gid)
	if !ok {
		return false
	}
	list := []mydb.Updater{g}
	for _, x := range facs {
		list = append(list, x)
	}
	for _, x := range planets {
		list = append(list, x)
	}
	// ------------ //
	op := overpower.NewTotallyOP()
	source := d.NewSource()
	op.Source = source
	op.Game = g
	op.Planets = planets
	op.Factions = facs
	op.Orders = orders
	op.Ships = ships
	// -------- //
	if !op.RunGameTurn() {
		return false
	}
	// -------- //
	if !d.DropAllGidOrders(gid) {
		return false
	}
	fsIds := make([][2]int, len(source.DroppedShips))
	for i, sh := range source.DroppedShips {
		fsIds[i] = [2]int{sh.Fid(), sh.Sid()}
	}
	if !d.DropSidShips(gid, fsIds) {
		return false
	}
	// Ships inserted at mid-run creation
	for _, x := range source.MadePlanetViews {
		list = append(list, x)
	}
	for _, x := range source.MadeShipViews {
		if !x.Insert(d.db) {
			return false
		}
	}
	return mydb.Update(d.db, list)
}

func (d DB) StartGame(gid int) (ok bool) {
	g, ok := d.GetGame(gid)
	if !ok {
		return
	}
	facs, ok := d.GetGidFactions(gid)
	if !ok {
		return
	}
	op := overpower.NewTotallyOP()
	source := d.NewSource()
	op.Source = source
	op.Game = g
	op.Factions = facs
	if !op.MakeGalaxy() {
		return false
	}
	list := []mydb.Updater{g}
	for _, x := range facs {
		list = append(list, x)
	}
	if !mydb.Update(d.db, list) {
		return false
	}
	for _, x := range source.MadePlanets {
		if !mydb.Insert(d.db, x) {
			return false
		}
	}
	for _, x := range source.MadePlanetViews {
		if !mydb.Insert(d.db, x) {
			return false
		}
	}
	for _, x := range source.MadeMapViews {
		if !mydb.Insert(d.db, x) {
			return false
		}
	}
	return true
}
