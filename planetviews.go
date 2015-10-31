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

func (pv *PlanetView) DistTo(c hexagon.Coord) int {
	return c.StepsTo(pv.Loc)
}

type SortedPV []*PlanetView

func (spv SortedPV) Len() int {
	return len(spv)
}

func (spv SortedPV) Swap(i, j int) {
	spv[i], spv[j] = spv[j], spv[i]
}

func (spv SortedPV) Less(i, j int) bool {
	if i == 0 {
		return true
	}
	if j == 0 {
		return false
	}
	return spv[0].Loc.StepsTo(spv[i].Loc) < spv[0].Loc.StepsTo(spv[j].Loc)
}

func (f *Faction) SortedPlanetViews(c hexagon.Coord) []*PlanetView {
	pvs := f.PlanetViews()
	list := make([]*PlanetView, len(pvs)+1)
	list[0] = &PlanetView{Loc: c}
	i := 1
	for _, pv := range pvs {
		list[i] = pv
		i++
	}
	sort.Sort(SortedPV(list))
	return list[1:]
}
