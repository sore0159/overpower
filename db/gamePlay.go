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
	op := overpower.NewTotallyOP()
	op.Game = g
	op.Factions, ok = d.GetGidFactions(g.Gid())
	if !ok {
		return false
	}
	list := []mydb.Updater{g}
	for _, f := range op.Factions {
		f.SetDone(false)
		list = append(list, f)
	}
	g.IncTurn()
	return mydb.Update(d.db, list)
}
