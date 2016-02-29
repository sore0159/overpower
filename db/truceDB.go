package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type TruceGroup struct {
	List []*Truce
}

func NewTruceGroup() *TruceGroup {
	return &TruceGroup{
		List: []*Truce{},
	}
}

func (group *TruceGroup) New() mydb.SQLer {
	item := NewTruce()
	group.List = append(group.List, item)
	return item
}

/*
func (group *TruceGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}
*/

func (group *TruceGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertTruces2DB(list ...overpower.Truce) ([]*Truce, error) {
	mylist := make([]*Truce, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*Truce); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad Truce struct type for op/db")
		}
	}
	return mylist, nil
}

func convertTruces2OP(list ...*Truce) []overpower.Truce {
	converted := make([]overpower.Truce, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
