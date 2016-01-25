package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type PlanetGroup struct {
	List []*Planet
}

func NewPlanetGroup() *PlanetGroup {
	return &PlanetGroup{
		List: []*Planet{},
	}
}

func (group *PlanetGroup) New() mydb.SQLer {
	item := NewPlanet()
	group.List = append(group.List, item)
	return item
}

func (group *PlanetGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.modified {
			list = append(list, item)
		}
	}
	return list
}

func (group *PlanetGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertPlanets2DB(list ...overpower.Planet) ([]*Planet, error) {
	mylist := make([]*Planet, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*Planet); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad Planet struct type for op/db")
		}
	}
	return mylist, nil
}

func convertPlanets2OP(list ...*Planet) []overpower.Planet {
	converted := make([]overpower.Planet, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}

/*
 */
