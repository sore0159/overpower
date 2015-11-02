package planetattack

import (
	"database/sql"
	"errors"
	"fmt"
	"mule/hexagon"
	"strings"
)

type Faction struct {
	db    *sql.DB
	Gid   int
	Fid   int
	Owner string
	Name  string
	Done  bool
	//
	CachePlanetViews map[hexagon.Coord]*PlanetView
	CacheShipViews   map[hexagon.Coord]*ShipView
}

func NewFaction(db *sql.DB) *Faction {
	return &Faction{
		db: db,
	}
}

func (f *Faction) Insert() error {
	query := "INSERT INTO factions (gid, owner, name, done) VALUES ($1, $2, $3, false) RETURNING fid"
	err := f.db.QueryRow(query, f.Gid, f.Owner, f.Name).Scan(&(f.Fid))
	if err != nil {
		return Log(err)
	}
	return nil
}

func (f *Faction) Delete() error {
	query := "DELETE FROM factions where gid = $1 AND fid = $2"
	res, err := f.db.Exec(query, f.Gid, f.Fid)
	if err != nil {
		Log("failed to delete faction", f, ":", err)
		return err
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		Log("failed to delete faction", f, ": 0 rows affected")
		return err
	}
	return nil
}

func (f *Faction) Select() bool {
	var err error
	if f.Fid != 0 {
		query := "SELECT owner, name, done FROM factions WHERE gid = $1 AND fid = $2"
		err = f.db.QueryRow(query, f.Gid, f.Fid).Scan(&(f.Owner), &(f.Name), &(f.Done))
	} else if f.Owner != "" {
		query := "SELECT fid, name, done FROM factions WHERE gid = $1 AND owner = $2"
		err = f.db.QueryRow(query, f.Gid, f.Owner).Scan(&(f.Fid), &(f.Name), &(f.Done))
	} else {
		err = errors.New("tried to SELECT faction with no fid/owner")
	}
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		Log(err)
		return false
	}
	return true
}

func (g *Game) GetFaction(userN string) (f *Faction, ok bool) {
	f = &Faction{db: g.db, Gid: g.Gid, Owner: userN}
	return f, f.Select()
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
		f := &Faction{db: db, Owner: owner}
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
	return f.SetDone(!f.Done)
}

func (f *Faction) SetDone(done bool) error {
	if f.Done == done {
		return nil
	}
	query := "UPDATE factions SET done = $1 WHERE gid = $2 and fid = $3"
	res, err := f.db.Exec(query, done, f.Gid, f.Fid)
	if err != nil {
		return Log("failed to set done", f.Gid, f.Fid, ":", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to set done", f.Gid, f.Fid, ": no rows affected")
	}
	f.Done = done
	return nil
}

func (f *Faction) GetView() (center hexagon.Coord, zoom int) {
	query := "SELECT center, zoom FROM views WHERE gid = $1 AND fid = $2"
	err := f.db.QueryRow(query, f.Gid, f.Fid).Scan(&(center), &(zoom))
	if err == sql.ErrNoRows {
		Log("View not found for f:", f)
		return center, -1
	} else if err != nil {
		Log(err)
		return center, -1
	}
	return
}

func (f *Faction) SetView(center hexagon.Coord, zoom int) error {
	query := "UPDATE views SET center = $1 AND zoom = $2 WHERE gid = $3 and fid = $4"
	res, err := f.db.Exec(query, center, zoom, f.Gid, f.Fid)
	if err != nil {
		return Log("failed to set view", f.Gid, f.Fid, ":", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to set view", f.Gid, f.Fid, ": no rows affected")
	}
	return nil
}

func (f *Faction) SetViewCenter(center hexagon.Coord) error {
	query := "UPDATE views SET center = $1 WHERE gid = $2 and fid = $3"
	res, err := f.db.Exec(query, center, f.Gid, f.Fid)
	if err != nil {
		return Log("failed to set view center", f.Gid, f.Fid, ":", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to set view center", f.Gid, f.Fid, ": no rows affected")
	}
	return nil
}

func (f *Faction) SetViewZoom(zoom int) error {
	query := "UPDATE views SET zoom = $1 WHERE gid = $2 and fid = $3"
	res, err := f.db.Exec(query, zoom, f.Gid, f.Fid)
	if err != nil {
		return Log("failed to set view", f.Gid, f.Fid, ":", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to set view", f.Gid, f.Fid, ": no rows affected")
	}
	return nil
}

func ViewMassInsertQ(planets []*Planet) string {
	query := "INSERT INTO views (gid, fid, center, zoom) VALUES "
	parts := []string{}
	for _, pl := range planets {
		if pl.Controller != 0 {
			parts = append(parts, fmt.Sprintf("(%d, %d, POINT(%d,%d), 0)", pl.Gid, pl.Controller, pl.Loc[0], pl.Loc[1]))
		}
	}
	return query + strings.Join(parts, ", ")
}
