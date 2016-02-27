package db

import (
	"database/sql"
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) MakePlanetView(gid, fid, turn, prFac, prPres, secFac, secPres, antiM, tach int, name string, loc hexagon.Coord) (err error) {
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
		fid: fid, turn: turn,
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

func (d DB) DropPlanetViews(conditions ...interface{}) error {
	return d.dropItems("planetviews", conditions)
}

func (d DB) UpdatePlanetViews(list ...overpower.PlanetView) error {
	mylist, err := convertPlanetViews2DB(list...)
	if my, bad := Check(err, "update PlanetViews conversion failure"); bad {
		return my
	}
	return d.updateGroup(&PlanetViewGroup{mylist})
}

func (d DB) GetPlanetView(conditions ...interface{}) (overpower.PlanetView, error) {
	list, err := d.GetPlanetViews(conditions...)
	if my, bad := Check(err, "get PlanetView failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetPlanetViews(conditions ...interface{}) ([]overpower.PlanetView, error) {
	group := NewPlanetViewGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get PlanetViews failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertPlanetViews2OP(list...)
	return converted, nil
}
