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

type LaunchOrder struct {
	GID    int           `json:"gid"`
	FID    int           `json:"fid"`
	Source hexagon.Coord `json:"source"`
	Target hexagon.Coord `json:"target"`
	Size   int           `json:"size"`
	sql    gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewLaunchOrder() *LaunchOrder {
	return &LaunchOrder{
	//
	}
}

type LaunchOrderIntf struct {
	item *LaunchOrder
}

func (item *LaunchOrder) Intf() overpower.LaunchOrderDat {
	return &LaunchOrderIntf{item}
}

func (i LaunchOrderIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *LaunchOrder) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "fid":
		return item.FID
	case "sourcex":
		return item.Source[0]
	case "sourcey":
		return item.Source[1]
	case "targetx":
		return item.Target[0]
	case "targety":
		return item.Target[1]
	case "size":
		return item.Size
	}
	return nil
}

func (item *LaunchOrder) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "fid":
		return &item.FID
	case "sourcex":
		return &item.Source[0]
	case "sourcey":
		return &item.Source[1]
	case "targetx":
		return &item.Target[0]
	case "targety":
		return &item.Target[1]
	case "size":
		return &item.Size
	}
	return nil
}
func (item *LaunchOrder) SQLTable() string {
	return "launchorder"
}

func (i LaunchOrderIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i LaunchOrderIntf) UnmarshalJSON(data []byte) error {
	i.item = &LaunchOrder{}
	return json.Unmarshal(data, i.item)
}

func (i LaunchOrderIntf) GID() int {
	return i.item.GID
}

func (i LaunchOrderIntf) FID() int {
	return i.item.FID
}

func (i LaunchOrderIntf) Source() hexagon.Coord {
	return i.item.Source
}

func (i LaunchOrderIntf) Target() hexagon.Coord {
	return i.item.Target
}

func (i LaunchOrderIntf) Size() int {
	return i.item.Size
}

func (i LaunchOrderIntf) SetSize(x int) {
	if i.item.Size == x {
		return
	}
	i.item.Size = x
	i.item.sql.UPDATE = true
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type LaunchOrderGroup struct {
	List []*LaunchOrder
}

func NewLaunchOrderGroup() *LaunchOrderGroup {
	return &LaunchOrderGroup{
		List: []*LaunchOrder{},
	}
}

func (item *LaunchOrder) SQLGroup() gp.SQLGrouper {
	return NewLaunchOrderGroup()
}

func (group *LaunchOrderGroup) New() gp.SQLer {
	item := NewLaunchOrder()
	group.List = append(group.List, item)
	return item
}

func (group *LaunchOrderGroup) UpdateList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.UPDATE && !item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *LaunchOrderGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *LaunchOrderGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *LaunchOrderGroup) SQLTable() string {
	return "launchorder"
}

func (group *LaunchOrderGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
	}
}

func (group *LaunchOrderGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
		"size",
	}
}

func (group *LaunchOrderGroup) InsertScanCols() []string {
	return []string{}
}

func (group *LaunchOrderGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
		"size",
	}
}

func (group *LaunchOrderGroup) UpdateCols() []string {
	return []string{
		"size",
	}
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type LaunchOrderSession struct {
	*LaunchOrderGroup
	*gp.Session
}

func NewLaunchOrderSession(d db.DBer) *LaunchOrderSession {
	group := NewLaunchOrderGroup()
	return &LaunchOrderSession{
		LaunchOrderGroup: group,
		Session:          gp.NewSession(group, d),
	}
}

func (s *LaunchOrderSession) Select(conditions ...interface{}) ([]overpower.LaunchOrderDat, error) {
	cur := len(s.LaunchOrderGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "LaunchOrder select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertLaunchOrder2Intf(s.LaunchOrderGroup.List[cur:]...), nil
}

func (s *LaunchOrderSession) SelectWhere(where sq.Condition) ([]overpower.LaunchOrderDat, error) {
	cur := len(s.LaunchOrderGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "LaunchOrder SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertLaunchOrder2Intf(s.LaunchOrderGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertLaunchOrder2Struct(list ...overpower.LaunchOrderDat) ([]*LaunchOrder, error) {
	mylist := make([]*LaunchOrder, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(LaunchOrderIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad LaunchOrder struct type for conversion")
		}
	}
	return mylist, nil
}

func convertLaunchOrder2Intf(list ...*LaunchOrder) []overpower.LaunchOrderDat {
	converted := make([]overpower.LaunchOrderDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func LaunchOrderTableCreate(d db.DBer) error {
	query := `create table launchorder(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid integer NOT NULL REFERENCES faction ON DELETE CASCADE,
	sourcex integer NOT NULL,
	sourcey integer NOT NULL,
	targetx integer NOT NULL,
	targety integer NOT NULL,
	size integer NOT NULL,
	FOREIGN KEY(gid, sourcex, sourcey) REFERENCES planet ON DELETE CASCADE,
	FOREIGN KEY(gid, targetx, targety) REFERENCES planet ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, sourcex, sourcey, targetx, targety)
);`

	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed LaunchOrder table creation", "query", query); bad {
		return my
	}
	return nil
}

func LaunchOrderTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS launchorder CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed LaunchOrder table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
