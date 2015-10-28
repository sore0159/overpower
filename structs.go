package planetattack

import (
	"database/sql"
)

type ShipView struct {
	Db         *sql.DB
	Gid        int
	Fid        int
	Sid        int
	Controller int
	Loc        Point
	Trail      []Point
	Size       int
}
