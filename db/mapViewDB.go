package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type MapViewGroup struct {
	List []*MapView
}

func NewMapViewGroup() *MapViewGroup {
	return &MapViewGroup{
		List: []*MapView{},
	}
}

func (group *MapViewGroup) New() mydb.SQLer {
	item := NewMapView()
	group.List = append(group.List, item)
	return item
}

func (group *MapViewGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func (group *MapViewGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertMapViews2DB(list ...overpower.MapView) ([]*MapView, error) {
	mylist := make([]*MapView, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*MapView); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad MapView struct type for op/db")
		}
	}
	return mylist, nil
}

func convertMapViews2OP(list ...*MapView) []overpower.MapView {
	converted := make([]overpower.MapView, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
