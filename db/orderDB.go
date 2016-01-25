package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type OrderGroup struct {
	List []*Order
}

func NewOrderGroup() *OrderGroup {
	return &OrderGroup{
		List: []*Order{},
	}
}

func (group *OrderGroup) New() mydb.SQLer {
	item := NewOrder()
	group.List = append(group.List, item)
	return item
}

func (group *OrderGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func (group *OrderGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertOrders2DB(list ...overpower.Order) ([]*Order, error) {
	mylist := make([]*Order, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*Order); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad Order struct type for op/db")
		}
	}
	return mylist, nil
}

func convertOrders2OP(list ...*Order) []overpower.Order {
	converted := make([]overpower.Order, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
