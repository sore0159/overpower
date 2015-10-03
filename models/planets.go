package models

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
type Planet struct {
	GID  int
	PID  int
	Name string
	Loc  [2]int
}

func NewPlanet() *Planet {
	return &Planet{
	//
	}
}

/*
create table planetviews (
	fid integer NOT NULL REFERENCES factions ON DELETE CASCADE,
	pid integer NOT NULL,
	gid integer NOT NULL,
	name varchar(20) NOT NULL,
	loc point NOT NULL,
	FOREIGN KEY(gid, pid) REFERENCES planets ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, pid)
);
*/
type PlanetView struct {
	GID  int
	PID  int
	FID  int
	Name string
	Loc  [2]int
}

func NewPlanetView() *PlanetView {
	return &PlanetView{
	//
	}
}
