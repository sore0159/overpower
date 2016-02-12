package db

import (
	"mule/hexagon"
	"mule/overpower"
)

func (d DB) MakeOrder(gid, fid, size int, source, target hexagon.Coord) (err error) {
	if size < 1 {
		return nil
	}
	item := &Order{gid: gid, fid: fid, size: size, source: source, target: target}
	group := &OrderGroup{[]*Order{item}}
	return d.makeGroup(group)
}

func (d DB) DropOrders(conditions ...interface{}) error {
	return d.dropItems("orders", conditions)
}

func (d DB) UpdateOrders(list ...overpower.Order) error {
	if len(list) == 0 {
		return nil
	}
	mylist, err := convertOrders2DB(list...)
	if my, bad := Check(err, "update Orders conversion failure"); bad {
		return my
	}
	upList := []*Order{}
	delList := []*Order{}
	for _, item := range mylist {
		if item.Size() > 0 {
			upList = append(upList, item)
		} else {
			delList = append(delList, item)
		}
	}
	f := func(tx DB) error {
		err := tx.updateGroup(&OrderGroup{upList})
		if err != nil {
			return err
		}
		for _, item := range delList {
			source, target := item.Source(), item.Target()
			conditions := C{"gid", item.Gid(), "fid", item.Fid(), "sourcex", source[0], "sourcey", source[1], "targetx", target[0], "targety", target[1]}
			err = tx.dropItemsIf("orders", conditions)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if d.InTrans() {
		return f(d)
	}
	return d.Transact(f)
}

func (d DB) GetOrdersBySource(gid int, source hexagon.Coord) ([]overpower.Order, error) {
	return d.GetOrders("gid", gid, "sourcex", source[0], "sourcey", source[1])
}
func (d DB) GetOrdersByST(gid int, source, target hexagon.Coord) (overpower.Order, error) {
	list, err := d.GetOrders("gid", gid, "sourcex", source[0], "sourcey", source[1], "targetx", target[0], "targety", target[1])
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (d DB) GetOrder(conditions ...interface{}) (overpower.Order, error) {
	list, err := d.GetOrders(conditions...)
	if my, bad := Check(err, "get Order failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetOrders(conditions ...interface{}) ([]overpower.Order, error) {
	group := NewOrderGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get Orders failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertOrders2OP(list...)
	return converted, nil
}
