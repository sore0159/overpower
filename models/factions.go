package models

import "database/sql"

/*
create table factions (
	fid SERIAL,
	gid integer NOT NULL REFERENCES games(id) ON DELETE CASCADE,
	owner varchar(20) NOT NULL,
	name varchar(20) NOT NULL,
	done bool NOT NULL,
	UNIQUE(gid, name)
	PRIMARY KEY(gid, fid)
);
*/

type Faction interface {
	ID() [2]int // GID FID
	Owner() string
	Name() string
}

type FactionDB struct {
	db  *sql.DB
	GID int
	FID int
}

func (f FactionDB) ID() [2]int {
	return [2]int{f.GID, f.FID}
}

func (f FactionDB) Owner() (owner string) {
	err := f.db.QueryRow("SELECT owner FROM factions WHERE gid = $1 AND fid = $2", f.GID, f.FID).Scan(&owner)
	if err != nil {
		Log(err)
		return "ERROR"
	}
	return
}

func (f FactionDB) Name() (name string) {
	err := f.db.QueryRow("SELECT name FROM factions WHERE gid = $1 AND fid = $2", f.GID, f.FID).Scan(&name)
	if err != nil {
		Log(err)
		return "ERROR"
	}
	return
}

func (g GameDB) GetFaction(fid int) Faction {
	return FactionDB{g.db, g.GID, fid}
}
