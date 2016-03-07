package models

import (
	"encoding/json"
	"errors"
	"mule/hexagon"
	"mule/mydb/db"
	gp "mule/mydb/group"
	sq "mule/mydb/sql"
	"mule/overpower"
)

type LaunchRecord struct {
	GID       int           `json:"gid"`
	FID       int           `json:"fid"`
	Turn      int           `json:"turn"`
	Source    hexagon.Coord `json:"source"`
	Target    hexagon.Coord `json:"target"`
	OrderSize int           `json:"ordersize"`
	Size      int           `json:"size"`
	sql       gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewLaunchRecord() *LaunchRecord {
	return &LaunchRecord{
	//
	}
}

type LaunchRecordIntf struct {
	item *LaunchRecord
}

func (item *LaunchRecord) Intf() overpower.LaunchRecordDat {
	return &LaunchRecordIntf{item}
}

func (i LaunchRecordIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *LaunchRecord) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "fid":
		return item.FID
	case "turn":
		return item.Turn
	case "sourcex":
		return item.Source[0]
	case "sourcey":
		return item.Source[1]
	case "targetx":
		return item.Target[0]
	case "targety":
		return item.Target[1]
	case "ordersize":
		return item.OrderSize
	case "size":
		return item.Size
	}
	return nil
}

func (item *LaunchRecord) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "fid":
		return &item.FID
	case "turn":
		return &item.Turn
	case "sourcex":
		return &item.Source[0]
	case "sourcey":
		return &item.Source[1]
	case "targetx":
		return &item.Target[0]
	case "targety":
		return &item.Target[1]
	case "ordersize":
		return &item.OrderSize
	case "size":
		return &item.Size
	}
	return nil
}
func (item *LaunchRecord) SQLTable() string {
	return "launchrecord"
}

func (i LaunchRecordIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i LaunchRecordIntf) UnmarshalJSON(data []byte) error {
	i.item = &LaunchRecord{}
	return json.Unmarshal(data, i.item)
}

func (i LaunchRecordIntf) GID() int {
	return i.item.GID
}

func (i LaunchRecordIntf) FID() int {
	return i.item.FID
}

func (i LaunchRecordIntf) Turn() int {
	return i.item.Turn
}

func (i LaunchRecordIntf) Source() hexagon.Coord {
	return i.item.Source
}

func (i LaunchRecordIntf) Target() hexagon.Coord {
	return i.item.Target
}

func (i LaunchRecordIntf) OrderSize() int {
	return i.item.OrderSize
}

func (i LaunchRecordIntf) Size() int {
	return i.item.Size
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type LaunchRecordGroup struct {
	List []*LaunchRecord
}

func NewLaunchRecordGroup() *LaunchRecordGroup {
	return &LaunchRecordGroup{
		List: []*LaunchRecord{},
	}
}

func (item *LaunchRecord) SQLGroup() gp.SQLGrouper {
	return NewLaunchRecordGroup()
}

func (group *LaunchRecordGroup) New() gp.SQLer {
	item := NewLaunchRecord()
	group.List = append(group.List, item)
	return item
}

func (group *LaunchRecordGroup) UpdateList() []gp.SQLer {
	return nil
}

func (group *LaunchRecordGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *LaunchRecordGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *LaunchRecordGroup) SQLTable() string {
	return "launchrecord"
}

func (group *LaunchRecordGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
	}
}

func (group *LaunchRecordGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
		"ordersize",
		"size",
	}
}

func (group *LaunchRecordGroup) InsertScanCols() []string {
	return []string{}
}

func (group *LaunchRecordGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
		"ordersize",
		"size",
	}
}

func (group *LaunchRecordGroup) UpdateCols() []string {
	return nil
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type LaunchRecordSession struct {
	*LaunchRecordGroup
	*gp.Session
}

func NewLaunchRecordSession(d db.DBer) *LaunchRecordSession {
	group := NewLaunchRecordGroup()
	return &LaunchRecordSession{
		LaunchRecordGroup: group,
		Session:           gp.NewSession(group, d),
	}
}

func (s *LaunchRecordSession) Select(conditions ...interface{}) ([]overpower.LaunchRecordDat, error) {
	cur := len(s.LaunchRecordGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "LaunchRecord select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertLaunchRecord2Intf(s.LaunchRecordGroup.List[cur:]...), nil
}

func (s *LaunchRecordSession) SelectWhere(where sq.Condition) ([]overpower.LaunchRecordDat, error) {
	cur := len(s.LaunchRecordGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "LaunchRecord SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertLaunchRecord2Intf(s.LaunchRecordGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertLaunchRecord2Struct(list ...overpower.LaunchRecordDat) ([]*LaunchRecord, error) {
	mylist := make([]*LaunchRecord, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(LaunchRecordIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad LaunchRecord struct type for conversion")
		}
	}
	return mylist, nil
}

func convertLaunchRecord2Intf(list ...*LaunchRecord) []overpower.LaunchRecordDat {
	converted := make([]overpower.LaunchRecordDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func LaunchRecordTableCreate(d db.DBer) error {
	query := `create table launchrecord(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid integer NOT NULL,
	turn int NOT NULL,
	sourcex integer NOT NULL,
	sourcey integer NOT NULL,
	targetx integer NOT NULL,
	targety integer NOT NULL,
	ordersize integer NOT NULL,
	size integer NOT NULL,
	FOREIGN KEY(gid, fid) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, sourcex, sourcey) REFERENCES planet ON DELETE CASCADE,
	FOREIGN KEY(gid, targetx, targety) REFERENCES planet ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, turn, sourcex, sourcey, targetx, targety)
);`
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed LaunchRecord table creation", "query", query); bad {
		return my
	}
	return nil
}

func LaunchRecordTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS launchrecord CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed LaunchRecord table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
