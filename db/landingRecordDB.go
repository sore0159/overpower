package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type LandingRecordGroup struct {
	List []*LandingRecord
}

func NewLandingRecordGroup() *LandingRecordGroup {
	return &LandingRecordGroup{
		List: []*LandingRecord{},
	}
}

func (group *LandingRecordGroup) New() mydb.SQLer {
	item := NewLandingRecord()
	group.List = append(group.List, item)
	return item
}

func (group *LandingRecordGroup) UpdateList() []mydb.SQLer {
	return nil
}

func (group *LandingRecordGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertLandingRecords2DB(list ...overpower.LandingRecord) ([]*LandingRecord, error) {
	mylist := make([]*LandingRecord, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*LandingRecord); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad LandingRecord struct type for op/db")
		}
	}
	return mylist, nil
}

func convertLandingRecords2OP(list ...*LandingRecord) []overpower.LandingRecord {
	converted := make([]overpower.LandingRecord, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
