package db

import (
	"fmt"
	"mule/mydb"
	"mule/overpower"
)

func (d DB) GetAllFactionPlanetViews(gid, fid int) (planetviews []overpower.PlanetView, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM planetviews WHERE gid = %d AND fid = %d", PVSQLVAL, gid, fid)
	return d.GetPlanetViewsQuery(query)
}

func (d DB) GetPlanetViewsQuery(query string) (planetviews []overpower.PlanetView, ok bool) {
	pvs := []*PlanetView{}
	maker := func() mydb.Rower {
		pv := NewPlanetView()
		return pv
	}
	if !mydb.Get(d.db, query, &pvs, maker) {
		return nil, false
	}
	planetviews = make([]overpower.PlanetView, len(pvs))
	for i, pv := range pvs {
		planetviews[i] = pv
	}
	return planetviews, true
}
