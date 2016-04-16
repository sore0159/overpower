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

type ShipView struct {
	GID        int               `json:"gid"`
	Turn       int               `json:"turn"`
	FID        int               `json:"fid"`
	SID        int               `json:"sid"`
	Controller int               `json:"controller"`
	Size       int               `json:"size"`
	Loc        hexagon.NullCoord `json:"loc"`
	Dest       hexagon.NullCoord `json:"dest"`
	Trail      hexagon.CoordList `json:"trail"`
	sql        gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewShipView() *ShipView {
	return &ShipView{
	//
	}
}

type ShipViewIntf struct {
	item *ShipView
}

func (item *ShipView) Intf() overpower.ShipViewDat {
	return &ShipViewIntf{item}
}

func (i ShipViewIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *ShipView) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "turn":
		return item.Turn
	case "fid":
		return item.FID
	case "sid":
		return item.SID
	case "controller":
		return item.Controller
	case "size":
		return item.Size
	case "loc":
		return item.Loc
	case "dest":
		return item.Dest
	case "trail":
		return item.Trail
	}
	return nil
}

func (item *ShipView) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "turn":
		return &item.Turn
	case "fid":
		return &item.FID
	case "sid":
		return &item.SID
	case "controller":
		return &item.Controller
	case "size":
		return &item.Size
	case "loc":
		return &item.Loc
	case "dest":
		return &item.Dest
	case "trail":
		return &item.Trail
	}
	return nil
}
func (item *ShipView) SQLTable() string {
	return "shipview"
}

func (i ShipViewIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i ShipViewIntf) UnmarshalJSON(data []byte) error {
	i.item = &ShipView{}
	return json.Unmarshal(data, i.item)
}

func (i ShipViewIntf) GID() int {
	return i.item.GID
}

func (i ShipViewIntf) Turn() int {
	return i.item.Turn
}

func (i ShipViewIntf) FID() int {
	return i.item.FID
}

func (i ShipViewIntf) SID() int {
	return i.item.SID
}

func (i ShipViewIntf) Controller() int {
	return i.item.Controller
}

func (i ShipViewIntf) Size() int {
	return i.item.Size
}

func (i ShipViewIntf) Loc() hexagon.NullCoord {
	return i.item.Loc
}

func (i ShipViewIntf) Dest() hexagon.NullCoord {
	return i.item.Dest
}

func (i ShipViewIntf) Trail() hexagon.CoordList {
	return i.item.Trail
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type ShipViewGroup struct {
	List []*ShipView
}

func NewShipViewGroup() *ShipViewGroup {
	return &ShipViewGroup{
		List: []*ShipView{},
	}
}

func (item *ShipView) SQLGroup() gp.SQLGrouper {
	return NewShipViewGroup()
}

func (group *ShipViewGroup) New() gp.SQLer {
	item := NewShipView()
	group.List = append(group.List, item)
	return item
}

func (group *ShipViewGroup) UpdateList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.UPDATE && !item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *ShipViewGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *ShipViewGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *ShipViewGroup) SQLTable() string {
	return "shipview"
}

func (group *ShipViewGroup) PKCols() []string {
	return []string{
		"gid",
		"turn",
		"fid",
		"sid",
	}
}

func (group *ShipViewGroup) InsertCols() []string {
	return []string{
		"gid",
		"turn",
		"fid",
		"sid",
		"controller",
		"size",
		"loc",
		"dest",
		"trail",
	}
}

func (group *ShipViewGroup) InsertScanCols() []string {
	return []string{}
}

func (group *ShipViewGroup) SelectCols() []string {
	return []string{
		"gid",
		"turn",
		"fid",
		"sid",
		"controller",
		"size",
		"loc",
		"dest",
		"trail",
	}
}

func (group *ShipViewGroup) UpdateCols() []string {
	return []string{
		"gid",
		"turn",
		"fid",
		"sid",
		"controller",
		"size",
		"loc",
		"dest",
		"trail",
	}
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type ShipViewSession struct {
	*ShipViewGroup
	*gp.Session
}

func NewShipViewSession(d db.DBer) *ShipViewSession {
	group := NewShipViewGroup()
	return &ShipViewSession{
		ShipViewGroup: group,
		Session:       gp.NewSession(group, d),
	}
}

func (s *ShipViewSession) Select(conditions ...interface{}) ([]overpower.ShipViewDat, error) {
	cur := len(s.ShipViewGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "ShipView select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertShipView2Intf(s.ShipViewGroup.List[cur:]...), nil
}

func (s *ShipViewSession) SelectWhere(where sq.Condition) ([]overpower.ShipViewDat, error) {
	cur := len(s.ShipViewGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "ShipView SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertShipView2Intf(s.ShipViewGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertShipView2Struct(list ...overpower.ShipViewDat) ([]*ShipView, error) {
	mylist := make([]*ShipView, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(ShipViewIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad ShipView struct type for conversion")
		}
	}
	return mylist, nil
}

func convertShipView2Intf(list ...*ShipView) []overpower.ShipViewDat {
	converted := make([]overpower.ShipViewDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func ShipViewTableCreate(d db.DBer) error {
	query := `create table shipview(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid integer NOT NULL REFERENCES faction ON DELETE CASCADE,
	controller integer NOT NULL REFERENCES faction ON DELETE CASCADE,
	sid integer NOT NULL,
	turn integer NOT NULL,
	loc point,
	dest point,
	trail point[] NOT NULL,
	size int NOT NULL,
	PRIMARY KEY(gid, fid, turn, sid)
);`
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed ShipView table creation", "query", query); bad {
		return my
	}
	return nil
}

func ShipViewTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS shipview CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed ShipView table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
