package attack

import (
	"database/sql"
)

type Faction struct {
	Db    *sql.DB
	Gid   int
	Fid   int
	Owner string
	Name  string
	Done  bool
	//
	CachePlanetViews map[Point]*PlanetView
	CacheShipViews   map[Point]*ShipView
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
	query := "SELECT owner, name, done, FROM factions WHERE gid = $1 AND fid = $2"
	err := f.Db.QueryRow(query, f.Gid, f.Fid).Scan(&(f.Owner), &(f.Name), &(f.Done))
	if err != nil {
		return Log(err)
	}
	return nil
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

func (f *Faction) PlanetViews() map[Point]*PlanetView {
	if f.CachePlanetViews == nil {
		f.CachePlanetViews = map[Point]*PlanetView{}
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
