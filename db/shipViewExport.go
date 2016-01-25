package db

import (
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) MakeShipView(gid, fid, turn, sid, size, controller int, loc, dest hexagon.NullCoord, trail hexagon.CoordList) error {
	item := &ShipView{gid: gid, fid: fid, turn: turn, sid: sid, size: size, controller: controller, loc: loc, dest: dest, trail: trail}
	group := &ShipViewGroup{[]*ShipView{item}}
	return d.makeGroup(group)
}

func (d DB) DropShipViews(conditions []interface{}) error {
	return d.dropItems("shipviews", conditions)
}

/*
func (d DB) UpdateShipViews(list ...overpower.ShipView) error {
	mylist, err := convertShipViews2DB(list...)
	if my, bad := Check(err, "update ShipViews conversion failure"); bad {
		return my
	}
	return d.updateGroup(&ShipViewGroup{mylist})
}
*/

func (d DB) GetShipView(conditions []interface{}) (overpower.ShipView, error) {
	list, err := d.GetShipViews(conditions)
	if my, bad := Check(err, "get ShipView failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetShipViews(conditions []interface{}) ([]overpower.ShipView, error) {
	group := NewShipViewGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get ShipViews failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertShipViews2OP(list...)
	return converted, nil
}
