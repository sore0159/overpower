package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type LaunchRecordGroup struct {
	List []*LaunchRecord
}

func NewLaunchRecordGroup() *LaunchRecordGroup {
	return &LaunchRecordGroup{
		List: []*LaunchRecord{},
	}
}

func (group *LaunchRecordGroup) New() mydb.SQLer {
	item := NewLaunchRecord()
	group.List = append(group.List, item)
	return item
}

func (group *LaunchRecordGroup) UpdateList() []mydb.SQLer {
	return nil
}

func (group *LaunchRecordGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertLaunchRecords2DB(list ...overpower.LaunchRecord) ([]*LaunchRecord, error) {
	mylist := make([]*LaunchRecord, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*LaunchRecord); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad LaunchRecord struct type for op/db")
		}
	}
	return mylist, nil
}

func convertLaunchRecords2OP(list ...*LaunchRecord) []overpower.LaunchRecord {
	converted := make([]overpower.LaunchRecord, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
