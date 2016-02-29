package db

import (
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) MakePowerOrder(gid, fid int, upPower bool, loc hexagon.Coord) (err error) {
	item := &PowerOrder{gid: gid, fid: fid, loc: loc, uppower: upPower}
	group := &PowerOrderGroup{[]*PowerOrder{item}}
	return d.makeGroup(group)
}

func (d DB) DropPowerOrders(conditions ...interface{}) error {
	return d.dropItems("powerorders", conditions)
}

func (d DB) UpdatePowerOrders(list ...overpower.PowerOrder) error {
	mylist, err := convertPowerOrders2DB(list...)
	if my, bad := Check(err, "update PowerOrders conversion failure"); bad {
		return my
	}
	return d.updateGroup(&PowerOrderGroup{mylist})
}

func (d DB) GetPowerOrdersByLoc(gid, fid int, loc hexagon.Coord) ([]overpower.PowerOrder, error) {
	return d.GetPowerOrders("gid", gid, "fid", fid, "locx", loc[0], "locy", loc[1])
}
func (d DB) GetPowerOrder(conditions ...interface{}) (overpower.PowerOrder, error) {
	list, err := d.GetPowerOrders(conditions...)
	if my, bad := Check(err, "get PowerOrder failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetPowerOrders(conditions ...interface{}) ([]overpower.PowerOrder, error) {
	group := NewPowerOrderGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get PowerOrders failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertPowerOrders2OP(list...)
	return converted, nil
}
