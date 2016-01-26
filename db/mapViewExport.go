package db

import (
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) UpdateMapView(vals ...interface{}) error {
	set := make([]interface{}, 0)
	where := make([]interface{}, 0)
	var flag bool
	for _, item := range vals {
		if flag {
			where = append(where, item)
			continue
		}
		if str, ok := item.(string); ok && str == "WHERE" {
			flag = true
		} else {
			set = append(set, item)
		}
	}
	return d.updateItem("mapviews", set, where)
}

func (d DB) MakeMapView(gid, fid int, center hexagon.Coord) (err error) {
	item := &MapView{gid: gid, fid: fid, center: center}
	group := &MapViewGroup{[]*MapView{item}}
	return d.makeGroup(group)
}

func (d DB) DropMapViews(conditions ...interface{}) error {
	return d.dropItems("mapviews", conditions)
}

func (d DB) UpdateMapViews(list ...overpower.MapView) error {
	mylist, err := convertMapViews2DB(list...)
	if my, bad := Check(err, "update MapViews conversion failure"); bad {
		return my
	}
	return d.updateGroup(&MapViewGroup{mylist})
}

func (d DB) GetMapView(conditions ...interface{}) (overpower.MapView, error) {
	list, err := d.GetMapViews(conditions...)
	if my, bad := Check(err, "get MapView failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetMapViews(conditions ...interface{}) ([]overpower.MapView, error) {
	group := NewMapViewGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get MapViews failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertMapViews2OP(list...)
	return converted, nil
}
