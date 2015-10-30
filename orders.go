package planetattack

import (
	"database/sql"
	"mule/hexagon"
)

type Order struct {
	db     *sql.DB
	Gid    int
	Fid    int
	Source hexagon.Coord
	Target hexagon.Coord
	Size   hexagon.Coord
}
