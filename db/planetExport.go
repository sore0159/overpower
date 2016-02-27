package db

import (
	"database/sql"
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) MakePlanet(gid, prFac, prPres, secFac, secPres, antiM, tach int, name string, loc hexagon.Coord) (err error) {
	prFacN := sql.NullInt64{}
	if prFac != 0 {
		prFacN.Valid = true
		prFacN.Int64 = int64(prFac)
	}
	secFacN := sql.NullInt64{}
	if secFac != 0 {
		secFacN.Valid = true
		secFacN.Int64 = int64(secFac)
	}
	item := &Planet{&MiniPlanet{
		gid: gid, name: name, loc: loc,
		primaryfaction:    prFacN,
		primarypresence:   prPres,
		secondaryfaction:  secFacN,
		secondarypresence: secPres,
		antimatter:        antiM, tachyons: tach,
	},
	}
	group := &PlanetGroup{[]*Planet{item}}
	return d.makeGroup(group)
}

func (d DB) DropPlanets(conditions ...interface{}) error {
	return d.dropItems("planets", conditions)
}

func (d DB) UpdatePlanets(list ...overpower.Planet) error {
	mylist, err := convertPlanets2DB(list...)
	if my, bad := Check(err, "update Planets conversion failure"); bad {
		return my
	}
	return d.updateGroup(&PlanetGroup{mylist})
}

func (d DB) GetPlanetsByLoc(gid int, locs ...hexagon.Coord) ([]overpower.Planet, error) {
	list, err := d.getPlanetsByLoc(gid, locs...)
	if my, bad := Check(err, "getplanets by locs fail", "gid", gid, "locs", locs); bad {
		return nil, my
	}
	return convertPlanets2OP(list...), nil
}

func (d DB) GetPlanet(conditions ...interface{}) (overpower.Planet, error) {
	list, err := d.GetPlanets(conditions...)
	if my, bad := Check(err, "get Planet failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetPlanets(conditions ...interface{}) ([]overpower.Planet, error) {
	group := NewPlanetGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get Planets failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertPlanets2OP(list...)
	return converted, nil
}
