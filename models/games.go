package models

/*
create table games(
	id SERIAL PRIMARY KEY,
	owner varchar(20) NOT NULL UNIQUE,
	name varchar(20) NOT NULL
	turn int
);
*/
type Game struct {
	GID      int
	Owner    string
	Name     string
	Turn     int
	Factions map[int]*Faction
	Sector   *Sector
}

func NewGame() *Game {
	return &Game{
		Factions: map[int]*Faction{},
		Sector:   NewSector(),
	}
}

type Sector struct {
	Planets map[[2]int]*Planet
}

func NewSector() *Sector {
	return &Sector{
		Planets: map[[2]int]*Planet{},
	}
}
