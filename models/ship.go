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

type Ship struct {
	GID      int               `json:"gid"`
	FID      int               `json:"fid"`
	SID      int               `json:"sid"`
	Size     int               `json:"size"`
	Launched int               `json:"launched"`
	Path     hexagon.CoordList `json:"path"`
	sql      gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewShip() *Ship {
	return &Ship{
	//
	}
}

type ShipIntf struct {
	item *Ship
}

func (item *Ship) Intf() overpower.ShipDat {
	return &ShipIntf{item}
}

func (i ShipIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *Ship) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "fid":
		return item.FID
	case "sid":
		return item.SID
	case "size":
		return item.Size
	case "launched":
		return item.Launched
	case "path":
		return item.Path
	}
	return nil
}

func (item *Ship) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "fid":
		return &item.FID
	case "sid":
		return &item.SID
	case "size":
		return &item.Size
	case "launched":
		return &item.Launched
	case "path":
		return &item.Path
	}
	return nil
}
func (item *Ship) SQLTable() string {
	return "ship"
}

func (i ShipIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i ShipIntf) UnmarshalJSON(data []byte) error {
	i.item = &Ship{}
	return json.Unmarshal(data, i.item)
}

func (i ShipIntf) GID() int {
	return i.item.GID
}

func (i ShipIntf) FID() int {
	return i.item.FID
}

func (i ShipIntf) SID() int {
	return i.item.SID
}

func (i ShipIntf) Size() int {
	return i.item.Size
}

func (i ShipIntf) Launched() int {
	return i.item.Launched
}

func (i ShipIntf) Path() hexagon.CoordList {
	return i.item.Path
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type ShipGroup struct {
	List []*Ship
}

func NewShipGroup() *ShipGroup {
	return &ShipGroup{
		List: []*Ship{},
	}
}

func (item *Ship) SQLGroup() gp.SQLGrouper {
	return NewShipGroup()
}

func (group *ShipGroup) New() gp.SQLer {
	item := NewShip()
	group.List = append(group.List, item)
	return item
}

func (group *ShipGroup) UpdateList() []gp.SQLer {
	return nil
}

func (group *ShipGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *ShipGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *ShipGroup) SQLTable() string {
	return "ship"
}

func (group *ShipGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"sid",
	}
}

func (group *ShipGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"sid",
		"size",
		"launched",
		"path",
	}
}

func (group *ShipGroup) InsertScanCols() []string {
	return []string{}
}

func (group *ShipGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"sid",
		"size",
		"launched",
		"path",
	}
}

func (group *ShipGroup) UpdateCols() []string {
	return nil
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type ShipSession struct {
	*ShipGroup
	*gp.Session
}

func NewShipSession(d db.DBer) *ShipSession {
	group := NewShipGroup()
	return &ShipSession{
		ShipGroup: group,
		Session:   gp.NewSession(group, d),
	}
}

func (s *ShipSession) Select(conditions ...interface{}) ([]overpower.ShipDat, error) {
	cur := len(s.ShipGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "Ship select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertShip2Intf(s.ShipGroup.List[cur:]...), nil
}

func (s *ShipSession) SelectWhere(where sq.Condition) ([]overpower.ShipDat, error) {
	cur := len(s.ShipGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "Ship SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertShip2Intf(s.ShipGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertShip2Struct(list ...overpower.ShipDat) ([]*Ship, error) {
	mylist := make([]*Ship, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(ShipIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad Ship struct type for conversion")
		}
	}
	return mylist, nil
}

func convertShip2Intf(list ...*Ship) []overpower.ShipDat {
	converted := make([]overpower.ShipDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func ShipTableCreate(d db.DBer) error {
	query := `create table ship(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid int NOT NULL REFERENCES faction ON DELETE CASCADE,
	sid int NOT NULL,
	size int NOT NULL,
	launched int NOT NULL,
	path point[] NOT NULL,
	PRIMARY KEY(gid, fid, sid)
);`

	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Ship table creation", "query", query); bad {
		return my
	}
	return nil
}

func ShipTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS ship CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Ship table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
