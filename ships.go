package planetattack

import (
	"database/sql"
	"mule/hexagon"
	//"strings"
)

type Ship struct {
	db   *sql.DB
	Gid  int
	Fid  int
	Sid  int
	Size int
	Loc  int
	Path []hexagon.Coord
}

func (sh *Ship) LandOn(pl *Planet, turn int) (views []*PlanetView, reports []string) {
	views = []*PlanetView{}

	if pl.Controller == sh.Fid {
		pl.Arrivals += sh.Size
		return append(views, pl.MakeView(turn, pl.Controller)), reports
	}
	if pl.Controller == 0 {
		if pl.Inhabitants == 0 {
			pl.Arrivals += sh.Size
		} else {
			pl.SetInhabitants(pl.Inhabitants - sh.Size)
			if pl.Inhabitants < 0 {
				pl.Arrivals = pl.Inhabitants * -1
				pl.SetInhabitants(0)
				pl.SetController(sh.Fid)
			}
		}
		return append(views, pl.MakeView(turn, pl.Controller)), reports
	}
	prev := pl.Controller
	pl.Arrivals -= sh.Size
	if pl.Arrivals < 0 {
		pl.SetInhabitants(pl.Inhabitants + pl.Arrivals)
		if pl.Inhabitants < 0 {
			pl.Arrivals = pl.Inhabitants * -1
			pl.SetInhabitants(0)
			pl.SetController(sh.Fid)
		}
	}
	return []*PlanetView{pl.MakeView(turn, prev), pl.MakeView(turn, sh.Fid)}, reports
}

func (sh *Ship) Travel() (atTarget bool, dist int) {
	dist = len(sh.Path) - sh.Loc
	sh.Loc += SHIPSPEED
	return sh.Loc > len(sh.Path)-2, dist
}

func (source *Planet) LaunchShip(target *Planet, fid, size int) *Ship {
	source.SetParts(source.Parts - size)
	return &Ship{
		db:   source.db,
		Gid:  source.Gid,
		Fid:  fid,
		Size: size,
		Loc:  0,
		Path: source.Loc.PathTo(target.Loc),
	}
}

/*
func (s *Ship) InsertViewQVals(fid int, viewpoints []hexagon.Coord) string {
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
*/

/*func Visible(pt1, pt2 hexagon.Coord) bool {
	return Dist(pt1, pt2) < 26
}

func Dist(pt1, pt2 hexagon.Coord) int {
	return (pt1[0]-pt2[0])*(pt1[0]-pt2[0]) + (pt1[1]-pt2[1])*(pt1[1]-pt2[1])
}*/

func (s *Ship) JustTravelled() []hexagon.Coord {
	if s.Loc >= len(s.Path) {
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

func (s *Ship) CurLoc() (hexagon.Coord, bool) {
	if s.Loc >= len(s.Path)-1 {
		return hexagon.Coord{}, false
	}
	return s.Path[s.Loc], true
}
