package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type FactionGroup struct {
	List []*Faction
}

func NewFactionGroup() *FactionGroup {
	return &FactionGroup{
		List: []*Faction{},
	}
}

func (group *FactionGroup) New() mydb.SQLer {
	item := NewFaction()
	group.List = append(group.List, item)
	return item
}

func (group *FactionGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.modified {
			list = append(list, item)
		}
	}
	return list
}

func (group *FactionGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertFactions2DB(list ...overpower.Faction) ([]*Faction, error) {
	mylist := make([]*Faction, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*Faction); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad Faction struct type for op/db")
		}
	}
	return mylist, nil
}

func convertFactions2OP(list ...*Faction) []overpower.Faction {
	converted := make([]overpower.Faction, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
