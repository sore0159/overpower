package planetattack

import (
	"database/sql"
	"mule/hexagon"
	"sort"
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
	Dist        int
}

func (f *Faction) PlanetViewList(center hexagon.Coord) []*PlanetView {
	r := []*PlanetView{}
	query := "SELECT pid, name, loc, turn, controller, inhabitants, resources, parts FROM planetviews WHERE gid = $1 AND fid = $2"
	rows, err := f.db.Query(query, f.Gid, f.Fid)
	if err != nil {
		Log(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		p := &PlanetView{db: f.db, Gid: f.Gid, Fid: f.Fid}
		var turn, controller, resources, parts, inhabitants sql.NullInt64
		err = rows.Scan(&(p.Pid), &(p.Name), &(p.Loc), &turn, &controller, &inhabitants, &resources, &parts)
		if err != nil {
			Log(err)
			return nil
		}
		x := []sql.NullInt64{turn, controller, resources, parts, inhabitants}
		y := []*int{&(p.Turn), &(p.Controller), &(p.Resources), &(p.Parts), &(p.Inhabitants)}
		for i, test := range x {
			if test.Valid {
				*(y[i]) = int(test.Int64)
			}
		}
		p.Dist = center.StepsTo(p.Loc)
		r = append(r, p)
	}
	if err = rows.Err(); err != nil {
		Log(err)
		return nil
	}
	return r
}

func (f *Faction) PlanetViews(center hexagon.Coord) map[hexagon.Coord]*PlanetView {
	if f.CachePlanetViews == nil {
		f.CachePlanetViews = map[hexagon.Coord]*PlanetView{}
		for _, plv := range f.PlanetViewList(center) {
			f.CachePlanetViews[plv.Loc] = plv
		}
	}
	return f.CachePlanetViews
}

func (pv *PlanetView) DistTo(c hexagon.Coord) int {
	return c.StepsTo(pv.Loc)
}

type ByDist []*PlanetView

func (spv ByDist) Len() int {
	return len(spv)
}

func (spv ByDist) Swap(i, j int) {
	spv[i], spv[j] = spv[j], spv[i]
}

func (spv ByDist) Less(i, j int) bool {
	return spv[i].Dist < spv[j].Dist
}

func (f *Faction) SortedPlanetViews(c hexagon.Coord) []*PlanetView {
	list := f.PlanetViewList(c)
	sort.Sort(ByDist(list))
	return list
}
