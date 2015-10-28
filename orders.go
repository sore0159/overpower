package planetattack

import (
	"database/sql"
)

type Order struct {
	Db     *sql.DB
	Gid    int
	Fid    int
	Source Point
	Target Point
	Size   Point
}
