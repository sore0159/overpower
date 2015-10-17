package models

import "database/sql"

/*
create table planets (
	pid integer NOT NULL,
	gid integer NOT NULL REFERENCES games(id) ON DELETE CASCADE,
	name varchar(20) NOT NULL,
	loc point NOT NULL,
	UNIQUE(gid, name),
	PRIMARY KEY(gid, pid)
);
*/
type Planet interface {
	ID() [2]int // GID, PID
	/*Name() string
	Loc() [2]int*/
}

type PlanetDB struct {
	db   *sql.DB
	GID  int
	PLID int
}

func (p PlanetDB) ID() [2]int {
	return [2]int{p.GID, p.PLID}
}

func (g GameDB) GetPlanet(plid int) Planet {
	return PlanetDB{g.db, g.GID, plid}
}

/*
create table planetviews (
	fid integer NOT NULL REFERENCES factions ON DELETE CASCADE,
	pid integer NOT NULL,
	gid integer NOT NULL,
	name varchar(20) NOT NULL,
	loc point NOT NULL,
	turn int,
	FOREIGN KEY(gid, pid) REFERENCES planets ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, pid)
);
*/

type PlanetView interface {
	ID() [3]int //GID, FID, PID
	/*Name()
	Loc() [2]int*/
}
type PlanetViewDB struct {
	db   *sql.DB
	GID  int
	FID  int
	PLID int
}

func (p PlanetViewDB) ID() [3]int {
	return [3]int{p.GID, p.FID, p.PLID}
}

func (f FactionDB) GetPlanetView(plid int) PlanetView {
	return PlanetViewDB{f.db, f.GID, f.FID, plid}
}

func (g GameDB) UpdateView(f Faction, pl Planet) PlanetView {
	query := "UPDATE planetviews SET turn = $1 WHERE gid = $2 AND fid = $3 and pid = $4"
	res, err := g.db.Exec(query, g.Turn(), g.GID, f.ID()[1], pl.ID()[1])
	if err != nil {
		Log("failed to update view:", g.GID, f.ID(), pl.ID(), err)
		return nil
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		Log("failed to update view", f.ID(), pl.ID(), ": 0 rows affected")
		return nil
	}
	return PlanetViewDB{g.db, g.GID, f.ID()[1], pl.ID()[1]}
}
