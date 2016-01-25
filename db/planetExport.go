package db

import (
	"database/sql"
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) MakePlanet(gid, pid, controller, inhabitants, resources, parts int, name string, loc hexagon.Coord) (err error) {
	contN := sql.NullInt64{}
	if controller != 0 {
		contN.Valid = true
		contN.Int64 = int64(controller)
	}
	item := &Planet{gid: gid, pid: pid, controller: contN, inhabitants: inhabitants, resources: resources, parts: parts}
	group := &PlanetGroup{[]*Planet{item}}
	return d.makeGroup(group)
}

func (d DB) DropPlanets(conditions []interface{}) error {
	return d.dropItems("planets", conditions)
}

func (d DB) UpdatePlanets(list ...overpower.Planet) error {
	mylist, err := convertPlanets2DB(list...)
	if my, bad := Check(err, "update Planets conversion failure"); bad {
		return my
	}
	return d.updateGroup(&PlanetGroup{mylist})
}

func (d DB) GetPlanet(conditions []interface{}) (overpower.Planet, error) {
	list, err := d.GetPlanets(conditions)
	if my, bad := Check(err, "get Planet failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetPlanets(conditions []interface{}) ([]overpower.Planet, error) {
	group := NewPlanetGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get Planets failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertPlanets2OP(list...)
	return converted, nil
}
