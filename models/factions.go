package models

/*
create table factions (
	fid SERIAL PRIMARY KEY,
	gid integer NOT NULL REFERENCES games(id) ON DELETE CASCADE,
	owner varchar(20) NOT NULL,
	name varchar(20) NOT NULL,
	UNIQUE(gid, name)
);
*/
type Faction struct {
	FID   int
	GID   int
	Owner string
	Name  string
}

func NewFaction() *Faction {
	return &Faction{
	//
	}
}
