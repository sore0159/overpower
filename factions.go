package planetattack

import (
	"database/sql"
	"mule/hexagon"
)

type Faction struct {
	Db    *sql.DB
	Gid   int
	Fid   int
	Owner string
	Name  string
	Done  bool
	//
	CachePlanetViews map[hexagon.Coord]*PlanetView
	CacheShipViews   map[hexagon.Coord]*ShipView
}

func (f *Faction) Insert() error {
	query := "INSERT INTO factions (gid, owner, name, done) VALUES ($1, $2, $3, false) RETURNING fid"
	err := f.Db.QueryRow(query, f.Gid, f.Owner, f.Name).Scan(&(f.Fid))
	if err != nil {
		return Log(err)
	}
	return nil
}

func (f *Faction) Select() error {
	query := "SELECT owner, name, done FROM factions WHERE gid = $1 AND fid = $2"
	err := f.Db.QueryRow(query, f.Gid, f.Fid).Scan(&(f.Owner), &(f.Name), &(f.Done))
	if err != nil {
		return Log(err)
	}
	return nil
}

func AllFactions(db *sql.DB, owner string) []*Faction {
	r := make([]*Faction, 0)
	query := "SELECT gid, fid, name, done FROM factions WHERE owner = $1"
	rows, err := db.Query(query, owner)
	if err != nil {
		Log(err)
		return []*Faction{}
	}
	defer rows.Close()
	for rows.Next() {
		f := &Faction{Db: db, Owner: owner}
		err = rows.Scan(&(f.Gid), &(f.Fid), &(f.Name), &(f.Done))
		if err != nil {
			Log(err)
			return []*Faction{}
		}
		r = append(r, f)
	}
	if err = rows.Err(); err != nil {
		Log(err)
		return []*Faction{}
	}
	return r
}

func (f *Faction) ToggleDone() error {
	f.Done = !f.Done
	query := "UPDATE factions SET done = !done WHERE gid = $1 and fid = $2"
	res, err := f.Db.Exec(query, f.Gid, f.Fid)
	if err != nil {
		return Log("failed to toggle done", f.Gid, f.Fid, ":", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to toggle done", f.Gid, f.Fid, ": no rows affected")
	}
	return nil
}

func (f *Faction) PlanetViews() map[hexagon.Coord]*PlanetView {
	if f.CachePlanetViews == nil {
		f.CachePlanetViews = map[hexagon.Coord]*PlanetView{}
		query := "SELECT pid, name, loc, turn, controller, inhabitants, resources, parts FROM planetviews WHERE gid = $1 AND fid = $2"
		rows, err := f.Db.Query(query, f.Gid, f.Fid)
		if err != nil {
			Log(err)
			f.CachePlanetViews = nil
			return nil
		}
		defer rows.Close()
		for rows.Next() {
			p := &PlanetView{Db: f.Db, Gid: f.Gid, Fid: f.Fid}
			var turn, controller, resources, parts, inhabitants sql.NullInt64
			err = rows.Scan(&(p.Pid), &(p.Name), &(p.Loc), &turn, &controller, &inhabitants, &resources, &parts)
			if err != nil {
				Log(err)
				f.CachePlanetViews = nil
				return nil
			}
			x := []sql.NullInt64{turn, controller, resources, parts, inhabitants}
			y := []*int{&(p.Turn), &(p.Controller), &(p.Resources), &(p.Parts), &(p.Inhabitants)}
			for i, test := range x {
				if test.Valid {
					*(y[i]) = int(test.Int64)
				}
			}
			f.CachePlanetViews[p.Loc] = p
		}
		if err = rows.Err(); err != nil {
			Log(err)
			f.CachePlanetViews = nil
			return nil
		}
	}
	return f.CachePlanetViews
}
