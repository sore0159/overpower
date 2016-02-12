package db

import (
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) MakePlanetView(gid, fid, turn, controller, inhabitants, resources, parts int, name string, loc hexagon.Coord) error {
	item := &PlanetView{gid: gid, fid: fid, turn: turn, name: name, loc: loc}
	if turn > 0 {
		item.controller.Valid = true
		item.controller.Int64 = int64(controller)
		item.inhabitants.Valid = true
		item.inhabitants.Int64 = int64(inhabitants)
		item.resources.Valid = true
		item.resources.Int64 = int64(resources)
		item.parts.Valid = true
		item.parts.Int64 = int64(parts)
	}
	group := &PlanetViewGroup{[]*PlanetView{item}}
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
