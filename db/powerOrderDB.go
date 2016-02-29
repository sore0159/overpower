package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type PowerOrderGroup struct {
	List []*PowerOrder
}

func NewPowerOrderGroup() *PowerOrderGroup {
	return &PowerOrderGroup{
		List: []*PowerOrder{},
	}
}

func (group *PowerOrderGroup) New() mydb.SQLer {
	item := NewPowerOrder()
	group.List = append(group.List, item)
	return item
}

func (group *PowerOrderGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func (group *PowerOrderGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertPowerOrders2DB(list ...overpower.PowerOrder) ([]*PowerOrder, error) {
	mylist := make([]*PowerOrder, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*PowerOrder); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad PowerOrder struct type for op/db")
		}
	}
	return mylist, nil
}

func convertPowerOrders2OP(list ...*PowerOrder) []overpower.PowerOrder {
	converted := make([]overpower.PowerOrder, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
