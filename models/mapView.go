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

type MapView struct {
	GID    int           `json:"gid"`
	FID    int           `json:"fid"`
	Center hexagon.Coord `json:"center"`
	sql    gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewMapView() *MapView {
	return &MapView{
	//
	}
}

type MapViewIntf struct {
	item *MapView
}

func (item *MapView) Intf() overpower.MapViewDat {
	return &MapViewIntf{item}
}

func (i MapViewIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *MapView) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "fid":
		return item.FID
	case "center":
		return item.Center
	}
	return nil
}

func (item *MapView) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "fid":
		return &item.FID
	case "center":
		return &item.Center
	}
	return nil
}
func (item *MapView) SQLTable() string {
	return "mapview"
}

func (i MapViewIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i MapViewIntf) UnmarshalJSON(data []byte) error {
	i.item = &MapView{}
	return json.Unmarshal(data, i.item)
}

func (i MapViewIntf) GID() int {
	return i.item.GID
}

func (i MapViewIntf) FID() int {
	return i.item.FID
}

func (i MapViewIntf) Center() hexagon.Coord {
	return i.item.Center
}

func (i MapViewIntf) SetCenter(x hexagon.Coord) {
	if i.item.Center == x {
		return
	}
	i.item.Center = x
	i.item.sql.UPDATE = true
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type MapViewGroup struct {
	List []*MapView
}

func NewMapViewGroup() *MapViewGroup {
	return &MapViewGroup{
		List: []*MapView{},
	}
}

func (item *MapView) SQLGroup() gp.SQLGrouper {
	return NewMapViewGroup()
}

func (group *MapViewGroup) New() gp.SQLer {
	item := NewMapView()
	group.List = append(group.List, item)
	return item
}

func (group *MapViewGroup) UpdateList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.UPDATE && !item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *MapViewGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *MapViewGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *MapViewGroup) SQLTable() string {
	return "mapview"
}

func (group *MapViewGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
	}
}

func (group *MapViewGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"center",
	}
}

func (group *MapViewGroup) InsertScanCols() []string {
	return []string{}
}

func (group *MapViewGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"center",
	}
}

func (group *MapViewGroup) UpdateCols() []string {
	return []string{
		"center",
	}
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type MapViewSession struct {
	*MapViewGroup
	*gp.Session
}

func NewMapViewSession(d db.DBer) *MapViewSession {
	group := NewMapViewGroup()
	return &MapViewSession{
		MapViewGroup: group,
		Session:      gp.NewSession(group, d),
	}
}

func (s *MapViewSession) Select(conditions ...interface{}) ([]overpower.MapViewDat, error) {
	cur := len(s.MapViewGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "MapView select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertMapView2Intf(s.MapViewGroup.List[cur:]...), nil
}

func (s *MapViewSession) SelectWhere(where sq.Condition) ([]overpower.MapViewDat, error) {
	cur := len(s.MapViewGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "MapView SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertMapView2Intf(s.MapViewGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertMapView2Struct(list ...overpower.MapViewDat) ([]*MapView, error) {
	mylist := make([]*MapView, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(MapViewIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad MapView struct type for conversion")
		}
	}
	return mylist, nil
}

func convertMapView2Intf(list ...*MapView) []overpower.MapViewDat {
	converted := make([]overpower.MapViewDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func MapViewTableCreate(d db.DBer) error {
	query := `create table mapview(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid integer NOT NULL,
	center point NOT NULL,
	FOREIGN KEY(gid, fid) REFERENCES faction ON DELETE CASCADE,
	PRIMARY KEY (gid, fid)
);`

	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed MapView table creation", "query", query); bad {
		return my
	}
	return nil
}

func MapViewTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS mapview CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed MapView table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
