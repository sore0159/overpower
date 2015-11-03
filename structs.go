package planetattack

import (
	"database/sql"
	"mule/hexagon"
)

type ShipView struct {
	db         *sql.DB
	Gid        int
	Fid        int
	Sid        int
	Controller int
	LocValid   bool
	Loc        hexagon.Coord
	Trail      []hexagon.Coord
	Size       int
}

func (sh *Ship) MakeView(viewpoints []hexagon.Coord, fid int) *ShipView {
	shv := &ShipView{
		db:         sh.db,
		Gid:        sh.Gid,
		Fid:        fid,
		Sid:        sh.Sid,
		Controller: sh.Fid,
		Size:       sh.Size,
	}
	loc, ok := sh.CurLoc()
	if ok && CanSee(viewpoints, loc) {
		shv.LocValid = true
		shv.Loc = loc
	}
	trail := []hexagon.Coord{}
	for _, pt := range sh.JustTravelled() {
		if CanSee(viewpoints, pt) {
			trail = append(trail, pt)
		}
	}
	shv.Trail = trail
	return shv
}

func CanSee(viewpoints []hexagon.Coord, pt hexagon.Coord) bool {
	for _, pt2 := range viewpoints {
		if pt.StepsTo(pt2) <= VISRANGE {
			return true
		}
	}
	return false
}

func (sv *ShipView) SeenFrom(c hexagon.Coord) (c2 hexagon.Coord, d int) {
	if sv.LocValid {
		return sv.Loc, c.StepsTo(sv.Loc)
	}
	d = -1
	for _, pt := range sv.Trail {
		if d2 := c.StepsTo(pt); d == -1 || d2 < d {
			d = d2
			c2 = pt
		}
	}
	return
}

type SortedSV []*ShipView

func (ssv SortedSV) Len() int {
	return len(ssv)
}

func (ssv SortedSV) Swap(i, j int) {
	ssv[i], ssv[j] = ssv[j], ssv[i]
}

func (ssv SortedSV) Less(i, j int) bool {
	if i == 0 {
		return true
	}
	if j == 0 {
		return false
	}
	_, di := ssv[i].SeenFrom(ssv[0].Loc)
	_, dj := ssv[j].SeenFrom(ssv[0].Loc)
	return di < dj
}
