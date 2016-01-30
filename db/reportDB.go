package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type ReportGroup struct {
	List []*Report
}

func NewReportGroup() *ReportGroup {
	return &ReportGroup{
		List: []*Report{},
	}
}

func (group *ReportGroup) New() mydb.SQLer {
	item := NewReport()
	group.List = append(group.List, item)
	return item
}

func (group *ReportGroup) UpdateList() []mydb.SQLer {
	return nil
}

func (group *ReportGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if len(item.contents) > 0 {
			list = append(list, item)
		}
	}
	return list
}

func convertReports2DB(list ...overpower.Report) ([]*Report, error) {
	mylist := make([]*Report, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*Report); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad Report struct type for op/db")
		}
	}
	return mylist, nil
}

func convertReports2OP(list ...*Report) []overpower.Report {
	converted := make([]overpower.Report, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
