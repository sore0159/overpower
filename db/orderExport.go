package db

import (
	"mule/overpower"
)

func (d DB) MakeOrder(gid, fid, source, target, size int) (err error) {
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
			conditions := C{"gid", item.Gid(), "fid", item.Fid(), "source", item.Source(), "target", item.Target()}
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
