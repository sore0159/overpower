package attack

import (
	"database/sql"
)

type PlanetView struct {
	Db          *sql.DB
	Gid         int
	Fid         int
	Pid         int
	Name        string
	Loc         Point
	Turn        int
	Controller  int
	Inhabitants int
	Resources   int
	Parts       int
}
