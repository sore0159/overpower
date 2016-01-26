package db

import (
	"mule/overpower"
)

func (d DB) MakeReport(gid, fid, turn int, contents []string) error {
	item := &Report{gid: gid, fid: fid, turn: turn, contents: contents}
	group := &ReportGroup{[]*Report{item}}
	return d.makeGroup(group)
}

func (d DB) DropReports(conditions ...interface{}) error {
	return d.dropItems("reports", conditions)
}

/*
func (d DB) UpdateReports(list ...overpower.Report) error {
	mylist, err := convertReports2DB(list...)
	if my, bad := Check(err, "update Reports conversion failure"); bad {
		return my
	}
	return d.updateGroup(&ReportGroup{mylist})
}
*/

func (d DB) GetReport(conditions ...interface{}) (overpower.Report, error) {
	list, err := d.GetReports(conditions...)
	if my, bad := Check(err, "get Report failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetReports(conditions ...interface{}) ([]overpower.Report, error) {
	group := NewReportGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get Reports failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertReports2OP(list...)
	return converted, nil
}
