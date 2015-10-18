package attack

import (
	"database/sql"
	//	"fmt"
)

type Game struct {
	Db    *sql.DB
	Gid   int
	Owner string
	Name  string
	Size  int
	Turn  int
	//
	CacheFactions map[int]*Faction
	CachePlanets  map[Point]*Planet
	CacheShips    map[Point][]*Ship
}

func (g *Game) Insert() error {
	query := "INSERT INTO games (name, owner, size, turn) VALUES($1, $2, $3, $4) RETURNING gid"
	err := g.Db.QueryRow(query, g.Name, g.Owner, g.Size, g.Turn).Scan(&(g.Gid))
	if err != nil {
		return Log(err)
	}
	return nil
}

func (g *Game) Select() error {
	var size sql.NullInt64
	query := "SELECT owner, name, size, turn FROM games WHERE gid = $1"
	err := g.Db.QueryRow(query, g.Gid).Scan(&(g.Owner), &(g.Name), &size, &(g.Turn))
	if err != nil {
		return Log(err)
	}
	if size.Valid {
		g.Size = int(size.Int64)
	}
	return nil
}

func (g *Game) IncTurn() error {
	g.Turn += 1
	query := "UPDATE games SET turn = turn + 1 WHERE gid = $1"
	res, err := g.Db.Exec(query, g.Gid)
	if err != nil {
		return Log("failed to inc game", g.Gid, "turn:", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to inc game", g.Gid, "turn: 0 rows affected")
	}
	return nil
}

func (g *Game) Factions() map[int]*Faction {
	if g.CacheFactions == nil {
		g.CacheFactions = map[int]*Faction{}
		query := "SELECT fid, owner, name, done FROM factions WHERE gid = $1"
		rows, err := g.Db.Query(query, g.Gid)
		if err != nil {
			Log(err)
			g.CacheFactions = nil
			return nil
		}
		defer rows.Close()
		for rows.Next() {
			f := &Faction{Db: g.Db, Gid: g.Gid}
			err = rows.Scan(&(f.Fid), &(f.Owner), &(f.Name), &(f.Done))
			if err != nil {
				Log(err)
				g.CacheFactions = nil
				return nil
			}
			g.CacheFactions[f.Fid] = f
		}
		if err = rows.Err(); err != nil {
			Log(err)
			g.CacheFactions = nil
			return nil
		}
	}
	return g.CacheFactions
}

func (g *Game) Planets() map[Point]*Planet {
	if g.CachePlanets == nil {
		g.CachePlanets = map[Point]*Planet{}
		query := "SELECT pid, name, loc, controller, inhabitants, resources, parts FROM planets WHERE gid = $1"
		rows, err := g.Db.Query(query, g.Gid)
		if err != nil {
			Log(err)
			g.CachePlanets = nil
			return nil
		}
		defer rows.Close()
		for rows.Next() {
			p := &Planet{Db: g.Db, Gid: g.Gid}
			var controller sql.NullInt64
			err = rows.Scan(&(p.Pid), &(p.Name), &(p.Loc), &controller, &(p.Inhabitants), &(p.Resources), &(p.Parts))
			if err != nil {
				Log(err)
				g.CachePlanets = nil
				return nil
			}
			if controller.Valid {
				p.Controller = int(controller.Int64)
			}
			g.CachePlanets[p.Loc] = p
		}
		if err = rows.Err(); err != nil {
			Log(err)
			g.CachePlanets = nil
			return nil
		}
	}
	return g.CachePlanets
}

func (g *Game) Ships() map[Point][]*Ship {
	if g.CacheShips == nil {
		g.CacheShips = map[Point][]*Ship{}
		query := "SELECT fid, sid, size, loc, path WHERE gid = $1"
		rows, err := g.Db.Query(query, g.Gid)
		if err != nil {
			Log(err)
			g.CacheShips = nil
			return nil
		}
		defer rows.Close()
		for rows.Next() {
			s := &Ship{Db: g.Db, Gid: g.Gid}
			err = rows.Scan(&(s.Fid), &(s.Sid), &(s.Size), &(s.Loc), &(s.Path))
			if err != nil {
				Log("Ship scan problem: ", err)
				g.CacheShips = nil
				return nil
			}
			if list, ok := g.CacheShips[s.Loc]; ok {
				g.CacheShips[s.Loc] = append(list, s)
			} else {
				g.CacheShips[s.Loc] = []*Ship{s}
			}
		}
		if err = rows.Err(); err != nil {
			Log(err)
			g.CacheShips = nil
			return nil
		}
	}
	return g.CacheShips
}
