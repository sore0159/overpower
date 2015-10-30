package planetattack

import (
	"database/sql"
	"mule/hexagon"
)

type PlanetView struct {
	db          *sql.DB
	Gid         int
	Fid         int
	Pid         int
	Name        string
	Loc         hexagon.Coord
	Turn        int
	Controller  int
	Inhabitants int
	Resources   int
	Parts       int
}
