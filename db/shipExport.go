package db

import (
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) MakeShip(gid, fid, size, launched int, path hexagon.CoordList) error {
	item := &Ship{gid: gid, fid: fid, size: size, launched: launched, path: path}
	group := &ShipGroup{[]*Ship{item}}
	return d.makeGroup(group)
}

func (d DB) DropShips(conditions ...interface{}) error {
	return d.dropItems("ships", conditions)
}
func (d DB) dropTheseShips(ships []overpower.Ship) error {
	mylist, err := convertShips2DB(ships...)
	if my, bad := Check(err, "drop Ships conversion failure"); bad {
		return my
	}
	return d.dropGroup(&ShipGroup{mylist})
}

func (d DB) UpdateShips(list ...overpower.Ship) error {
	mylist, err := convertShips2DB(list...)
	if my, bad := Check(err, "update Ships conversion failure"); bad {
		return my
	}
	return d.updateGroup(&ShipGroup{mylist})
}

func (d DB) GetShip(conditions ...interface{}) (overpower.Ship, error) {
	list, err := d.GetShips(conditions...)
	if my, bad := Check(err, "get Ship failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetShips(conditions ...interface{}) ([]overpower.Ship, error) {
	group := NewShipGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get Ships failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertShips2OP(list...)
	return converted, nil
}
