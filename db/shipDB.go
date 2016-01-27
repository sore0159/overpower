package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type ShipGroup struct {
	List []*Ship
}

func NewShipGroup() *ShipGroup {
	return &ShipGroup{
		List: []*Ship{},
	}
}

func (group *ShipGroup) New() mydb.SQLer {
	item := NewShip()
	group.List = append(group.List, item)
	return item
}

func (group *ShipGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func (group *ShipGroup) DeleteList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if !item.justmade {
			list = append(list, item)
		}
	}
	return list
}

func (group *ShipGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertShips2DB(list ...overpower.Ship) ([]*Ship, error) {
	mylist := make([]*Ship, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*Ship); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad Ship struct type for op/db")
		}
	}
	return mylist, nil
}

func convertShips2OP(list ...*Ship) []overpower.Ship {
	converted := make([]overpower.Ship, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
