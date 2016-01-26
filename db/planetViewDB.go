package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type PlanetViewGroup struct {
	List []*PlanetView
}

func NewPlanetViewGroup() *PlanetViewGroup {
	return &PlanetViewGroup{
		List: []*PlanetView{},
	}
}

func (group *PlanetViewGroup) New() mydb.SQLer {
	item := NewPlanetView()
	group.List = append(group.List, item)
	return item
}

func (group *PlanetViewGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func (group *PlanetViewGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertPlanetViews2DB(list ...overpower.PlanetView) ([]*PlanetView, error) {
	mylist := make([]*PlanetView, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*PlanetView); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad PlanetView struct type for op/db")
		}
	}
	return mylist, nil
}

func convertPlanetViews2OP(list ...*PlanetView) []overpower.PlanetView {
	converted := make([]overpower.PlanetView, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
