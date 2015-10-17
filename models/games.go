package models

import "database/sql"

/*
create table games(
	id SERIAL PRIMARY KEY,
	owner varchar(20) NOT NULL UNIQUE,
	name varchar(20) NOT NULL,
	turn int
);
*/
type Game interface {
	ID() int // GID
	Owner() string
	Name() string
	Turn() int
	Factions() map[int]Faction
}

type GameDB struct {
	db  *sql.DB
	GID int
}

func (g GameDB) ID() int {
	return g.GID
}

func (g GameDB) Owner() (owner string) {
	err := g.db.QueryRow("SELECT owner FROM games WHERE id = $1", g.GID).Scan(&owner)
	if err != nil {
		Log(err)
		return "ERROR"
	}
	return
}

func (g GameDB) Name() (name string) {
	err := g.db.QueryRow("SELECT name FROM games WHERE id = $1", g.GID).Scan(&name)
	if err != nil {
		Log(err)
		return "ERROR"
	}
	return
}

func (g GameDB) Turn() (turn int) {
	err := g.db.QueryRow("SELECT turn FROM games WHERE id = $1", g.GID).Scan(&turn)
	if err != nil {
		Log(err)
		return -1
	}
	return
}

func (g GameDB) Factions() map[int]Faction {
	m := map[int]Faction{}
	query := "SELECT fid from factions WHERE gid = $1"
	rows, err := g.db.Query(query, g.GID)
	if err != nil {
		Log(err)
		return nil
	}
	defer rows.Close()
	var fid int
	for rows.Next() {
		err := rows.Scan(&fid)
		if err != nil {
			Log(err)
			return nil
		}
		m[fid] = g.GetFaction(fid)
	}
	err = rows.Err()
	if err != nil {
		Log(err)
		return nil
	}
	return m
}
