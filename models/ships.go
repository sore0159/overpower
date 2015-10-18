package attack

import (
	"database/sql"
	"fmt"
	"strings"
)

type Ship struct {
	Db   *sql.DB
	Gid  int
	Fid  int
	Sid  int
	Size int
	Loc  int
	Path []Point
}

func (s *Ship) InsertViewQVals(fid int, viewpoints []Point) string {
	// q := "INSERT INTO shipviews (gid, viewer, controller, sid, size, loc, trail) VALUES "
	parts := []string{}
	for _, pt := range s.JustTravelled() {
		for _, pt2 := range viewpoints {
			if Visible(pt, pt2) {
				parts = append(parts, pt.String())
				break
			}
		}
	}
	if len(parts) < 1 {
		return ""
	}
	pathStr := "{"
	pathStr += strings.Join(parts, ", ")
	pathStr += "}"
	loc, ok := s.CurLoc()
	var hasLoc bool
	if ok {
		for _, pt2 := range viewpoints {
			if Visible(loc, pt2) {
				hasLoc = true
				break
			}
		}
	}
	if hasLoc {
		return fmt.Sprintf("(%d, %d, %d, %d, %s, %s)", s.Gid, fid, s.Fid, s.Size, loc, pathStr)
	} else {
		return fmt.Sprintf("(%d, %d, %d, %d, %s)", s.Gid, fid, s.Fid, s.Size, pathStr)
	}
}

func ShipMassInsertQ(ships []*Ship) string {
	q := "INSERT INTO ships (gid, fid, size, loc, path) VALUES "
	qStr := make([]string, len(ships))
	for i, ship := range ships {
		qStr[i] = ship.InsertQVals()
	}
	q += strings.Join(qStr, ", ")
	return q
}

func (s *Ship) InsertQVals() string {
	pathStr := "{"
	parts := make([]string, len(s.Path))
	for i, pt := range s.Path {
		parts[i] = pt.String()
	}
	pathStr += strings.Join(parts, ", ")
	pathStr += "}"
	var loc sql.NullInt64
	if s.Loc == -1 {
		loc.Valid = false
	} else {
		loc.Valid = true
		loc.Int64 = int64(s.Loc)
	}
	return fmt.Sprintf("(%d, %d, %d, %d, %s)", s.Gid, s.Fid, s.Size, loc, pathStr)
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

func Visible(pt1, pt2 Point) bool {
	return Dist(pt1, pt2) < 26
}

func Dist(pt1, pt2 Point) int {
	return (pt1[0]-pt2[0])*(pt1[0]-pt2[0]) + (pt1[1]-pt2[1])*(pt1[1]-pt2[1])
}

func (s *Ship) JustTravelled() []Point {
	if s.Loc == -1 {
		if len(s.Path) < SHIPSPEED {
			return s.Path
		} else {
			return s.Path[len(s.Path)-SHIPSPEED:]
		}
	} else if s.Loc < SHIPSPEED {
		return s.Path[:s.Loc+1]
	} else {
		return s.Path[s.Loc-SHIPSPEED : s.Loc+1]
	}
}

func (s *Ship) CurLoc() (Point, bool) {
	if s.Loc == -1 {
		return Point{}, false
	}
	return s.Path[s.Loc], true
}
