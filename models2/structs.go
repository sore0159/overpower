package attack

import (
	"database/sql"
	//	"fmt"
)

type Game struct {
	db    *sql.DB
	gid   int
	owner string
	name  string
	size  int
	turn  int
}

func (g *Game) Insert() error {
	query := "INSERT INTO games (name, owner, size, turn) VALUES($1, $2, $3, $4) RETURNING gid"
	err := g.db.QueryRow(query, g.name, g.owner, g.size, 1).Scan(&(g.gid))
	if err != nil {
		return Log(err)
	}
	return nil
}

func (g *Game) Select() error {
	query := "SELECT owner, name, size, turn FROM games WHERE gid = $1"
	err := g.db.QueryRow(query, g.gid).Scan(&(g.owner), &(g.name), &(g.size), &(g.turn))
	if err != nil {
		return Log(err)
	}
	return nil
}

func (g *Game) IncTurn() error {
	g.turn += 1
	query := "UPDATE games SET turn = turn + 1 WHERE gid = $1"
	res, err := g.db.Exec(query, g.gid)
	if err != nil {
		return Log("failed to inc game", g.gid, "turn:", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to inc game", g.gid, "turn: 0 rows affected")
	}
	return nil
}

type Faction struct {
	db      *sql.DB
	gid     int
	fid     int
	owner   string
	name    string
	done    bool
	reports [][]string
}

func (f *Faction) InsertStmt() *sql.Stmt {
	query := "INSERT INTO factions (gid, owner, name, done) VALUES ($1, $2, $3, false) RETURNING fid"
	stmt, err := f.db.Prepare(query)
	if err != nil {
		Log(err)
		return nil
	}
	return stmt
}

func (f *Faction) Insert(stmt *sql.Stmt) error {
	if stmt == nil {
		return Log("Nil stmt to faction Make()")
	}
	err := stmt.QueryRow(f.gid, f.owner, f.name).Scan(&(f.fid))
	if err != nil {
		return Log(err)
	}
	return nil
}

func (f *Faction) Select() error {
	query := "SELECT owner, name, done, reports FROM factions WHERE gid = $1 AND fid = $2"
	err := f.db.QueryRow(query, f.gid, f.fid).Scan(&(f.owner), &(f.name), &(f.done), &(f.reports))
	if err != nil {
		return Log(err)
	}
	return nil
}

func (f *Faction) ToggleDone() error {
	f.done = !f.done
	query := "UPDATE factions SET done = !done WHERE gid = $1 and fid = $2"
	res, err := f.db.Exec(query, f.gid, f.fid)
	if err != nil {
		return Log("failed to toggle done", f.gid, f.fid, ":", err)
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		return Log("failed to toggle done", f.gid, f.fid, ": no rows affected")
	}
	return nil
}

type Planet struct {
	db          *sql.DB
	gid         int
	pid         int
	name        string
	loc         [2]int
	controller  int
	inhabitants int
	resources   int
	parts       int
}

func (p *Planet) Insert(tx *sql.Tx) {
}

type Ship struct {
	db   *sql.DB
	gid  int
	fid  int
	sid  int
	size int
	loc  [2]int
	path [][2]int
}

type PlanetView struct {
	db          *sql.DB
	gid         int
	fid         int
	pid         int
	name        string
	loc         [2]int
	turn        int
	controller  int
	inhabitants int
	resources   int
	parts       int
}

type ShipView struct {
	db         *sql.DB
	gid        int
	fid        int
	sid        int
	controller int
	loc        [2]int
	trail      [][2]int
	size       int
}

type Order struct {
	db     *sql.DB
	gid    int
	fid    int
	source [2]int
	target [2]int
	size   [2]int
}
