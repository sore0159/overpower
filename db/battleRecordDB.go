package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type BattleRecordGroup struct {
	List []*BattleRecord
}

func NewBattleRecordGroup() *BattleRecordGroup {
	return &BattleRecordGroup{
		List: []*BattleRecord{},
	}
}

func (group *BattleRecordGroup) New() mydb.SQLer {
	item := NewBattleRecord()
	group.List = append(group.List, item)
	return item
}

func (group *BattleRecordGroup) UpdateList() []mydb.SQLer {
	return nil
}

func (group *BattleRecordGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertBattleRecords2DB(list ...overpower.BattleRecord) ([]*BattleRecord, error) {
	mylist := make([]*BattleRecord, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*BattleRecord); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad BattleRecord struct type for op/db")
		}
	}
	return mylist, nil
}

func convertBattleRecords2OP(list ...*BattleRecord) []overpower.BattleRecord {
	converted := make([]overpower.BattleRecord, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
