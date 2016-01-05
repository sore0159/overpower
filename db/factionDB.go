package db

import (
	"fmt"
	"mule/mydb"
	"mule/overpower"
)

func (d DB) DropInProgressFaction(gid, fid int) (ok bool) {
	// drop fid orders
	if !d.DropAllFidOrders(gid, fid) {
		return false
	}
	planets, ok := d.GetAllGidPlanets(gid)
	if !ok {
		return false
	}
	updateList := []mydb.Updater{}
	// set planets neutral
	for _, p := range planets {
		if p.Controller() == fid {
			p.SetController(0)
			p.SetResources(p.Resources() + p.Parts())
			p.SetParts(0)
			updateList = append(updateList, p)
		}
	}
	if !mydb.Update(d.db, updateList) {
		return false
	}
	// drop faction
	if !d.DropFaction(gid, fid) {
		return false
	}
	return true
}

func (d DB) MakeFaction(gid int, owner, name string) (ok bool) {
	f := NewFaction()
	f.gid, f.owner, f.name = gid, owner, name
	return f.Insert(d.db)
}
func (d DB) DropFaction(gid, fid int) (ok bool) {
	query := fmt.Sprintf("DELETE FROM factions WHERE gid = %d AND fid = %d", gid, fid)
	return mydb.Exec(d.db, query)
}

func (d DB) GetFidFaction(gid, fid int) (faction overpower.Faction, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM factions WHERE gid = %d AND fid = %d", FACSQLVAL, gid, fid)
	f := NewFaction()
	return f, mydb.GetOneIf(d.db, query, f)
}

func (d DB) GetOwnerFaction(gid int, owner string) (faction overpower.Faction, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM factions WHERE gid = %d AND owner = '%s'", FACSQLVAL, gid, owner)
	f := NewFaction()
	return f, mydb.GetOneIf(d.db, query, f)
}

func (d DB) GetGidFactions(gid int) (factions []overpower.Faction, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM factions WHERE gid = %d", FACSQLVAL, gid)
	return d.GetFactionsQuery(query)
}
func (d DB) GetAllOwnerFactions(owner string) (factions []overpower.Faction, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM factions WHERE owner = '%s'", FACSQLVAL, owner)
	return d.GetFactionsQuery(query)
}

func (d DB) GetFactionsQuery(query string) (factions []overpower.Faction, ok bool) {
	factionL := []*Faction{}
	maker := func() mydb.Rower {
		f := NewFaction()
		return f
	}
	if !mydb.Get(d.db, query, &factionL, maker) {
		return nil, false
	}
	factions = make([]overpower.Faction, len(factionL))
	for i, f := range factionL {
		factions[i] = f
	}
	return factions, true
}
