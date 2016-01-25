package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type ShipViewGroup struct {
	List []*ShipView
}

func NewShipViewGroup() *ShipViewGroup {
	return &ShipViewGroup{
		List: []*ShipView{},
	}
}

func (group *ShipViewGroup) New() mydb.SQLer {
	item := NewShipView()
	group.List = append(group.List, item)
	return item
}

func (group *ShipViewGroup) UpdateList() []mydb.SQLer {
	return nil
}

func (group *ShipViewGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertShipViews2DB(list ...overpower.ShipView) ([]*ShipView, error) {
	mylist := make([]*ShipView, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*ShipView); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad ShipView struct type for op/db")
		}
	}
	return mylist, nil
}

func convertShipViews2OP(list ...*ShipView) []overpower.ShipView {
	converted := make([]overpower.ShipView, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
