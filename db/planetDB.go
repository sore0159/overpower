package db

import (
	"fmt"
	"mule/mydb"
	"mule/overpower"
	"strings"
)

func (d DB) GetPidPlanets(gid int, plids ...int) (planets []overpower.Planet, ok bool) {
	if len(plids) < 1 {
		return nil, true
	}
	query := fmt.Sprintf("SELECT %s FROM planets WHERE gid = %d AND (", PLSQLVAL, gid)
	parts := make([]string, len(plids))
	for i, _ := range parts {
		parts[i] = fmt.Sprintf("pid = %d", plids[i])
	}
	query += strings.Join(parts, " OR ") + ")"
	return d.GetPlanetsQuery(query)
}

func (d DB) GetAllGidPlanets(gid int) (planets []overpower.Planet, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM planets WHERE gid = %d", PLSQLVAL, gid)
	return d.GetPlanetsQuery(query)
}

func (d DB) GetPlanetsQuery(query string) (planets []overpower.Planet, ok bool) {
	pvs := []*Planet{}
	maker := func() mydb.Rower {
		pv := NewPlanet()
		return pv
	}
	if !mydb.Get(d.db, query, &pvs, maker) {
		return nil, false
	}
	planets = make([]overpower.Planet, len(pvs))
	for i, pv := range pvs {
		planets[i] = pv
	}
	return planets, true
}
