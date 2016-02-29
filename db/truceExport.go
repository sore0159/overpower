package db

import (
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) MakeTruce(gid, fid, trucee int, loc hexagon.Coord) (err error) {
	item := &Truce{gid: gid, fid: fid, loc: loc, trucee: trucee}
	group := &TruceGroup{[]*Truce{item}}
	return d.makeGroup(group)
}

func (d DB) DropOPTruces(list ...overpower.Truce) error {
	for _, tr := range list {
		conds := []interface{}{"gid", tr.Gid(), "fid", tr.Fid(), "locx", tr.Loc()[0], "locy", tr.Loc()[1], "trucee", tr.Trucee()}
		err := d.dropItems("truces", conds)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d DB) DropTruces(conditions ...interface{}) error {
	return d.dropItems("truces", conditions)
}

/*
func (d DB) UpdateTruces(list ...overpower.Truce) error {
	mylist, err := convertTruces2DB(list...)
	if my, bad := Check(err, "update Truces conversion failure"); bad {
		return my
	}
	return d.updateGroup(&TruceGroup{mylist})
}
*/

func (d DB) GetTrucesByLoc(gid, fid int, loc hexagon.Coord) ([]overpower.Truce, error) {
	return d.GetTruces("gid", gid, "fid", fid, "locx", loc[0], "locy", loc[1])
}
func (d DB) GetTruce(conditions ...interface{}) (overpower.Truce, error) {
	list, err := d.GetTruces(conditions...)
	if my, bad := Check(err, "get Truce failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetTruces(conditions ...interface{}) ([]overpower.Truce, error) {
	group := NewTruceGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get Truces failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertTruces2OP(list...)
	return converted, nil
}
