package attack

import (
	"database/sql"
	"fmt"
)

type Ship struct {
	Db   *sql.DB
	Gid  int
	Fid  int
	Sid  int
	Size int
	Loc  Point
	Path []Point
}

func (s *Ship) InsertValStr() string {
	return fmt.Sprintf("")

}

func (p *Planet) SpawnShip(size int, target *Planet) *Ship {
	s := &Ship{Db: p.Db, Gid: p.Gid, Fid: p.Controller, Size: size}
	s.Path = Pathfind(p, target)
	query := "UPDATE planets SET parts = parts - $1 WHERE gid = $2 AND pid = $3"
	res, err := p.Db.Exec(query, size, p.Gid, p.Pid)
	if err != nil {
		Log(err)
		return nil
	}
	if aff, err := res.RowsAffected(); err != nil || aff < 1 {
		Log("failed to spawn ship", p.Gid, p.Pid, "pt1: 0 rows affected")
		return nil
	}
	return nil
}

func Pathfind(begin, end *Planet) []Point {
	return []Point{begin.Loc, end.Loc}
}
