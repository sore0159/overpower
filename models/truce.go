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

type Truce struct {
	GID    int           `json:"gid"`
	FID    int           `json:"fid"`
	Loc    hexagon.Coord `json:"loc"`
	Trucee int           `json:"trucee"`
	sql    gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewTruce() *Truce {
	return &Truce{
	//
	}
}

type TruceIntf struct {
	item *Truce
}

func (item *Truce) Intf() overpower.TruceDat {
	return &TruceIntf{item}
}

func (i TruceIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *Truce) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "fid":
		return item.FID
	case "locx":
		return item.Loc[0]
	case "locy":
		return item.Loc[1]
	case "trucee":
		return item.Trucee
	}
	return nil
}

func (item *Truce) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "fid":
		return &item.FID
	case "locx":
		return &item.Loc[0]
	case "locy":
		return &item.Loc[1]
	case "trucee":
		return &item.Trucee
	}
	return nil
}
func (item *Truce) SQLTable() string {
	return "truce"
}

func (i TruceIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i TruceIntf) UnmarshalJSON(data []byte) error {
	i.item = &Truce{}
	return json.Unmarshal(data, i.item)
}

func (i TruceIntf) GID() int {
	return i.item.GID
}

func (i TruceIntf) FID() int {
	return i.item.FID
}

func (i TruceIntf) Loc() hexagon.Coord {
	return i.item.Loc
}

func (i TruceIntf) Trucee() int {
	return i.item.Trucee
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type TruceGroup struct {
	List []*Truce
}

func NewTruceGroup() *TruceGroup {
	return &TruceGroup{
		List: []*Truce{},
	}
}

func (item *Truce) SQLGroup() gp.SQLGrouper {
	return NewTruceGroup()
}

func (group *TruceGroup) New() gp.SQLer {
	item := NewTruce()
	group.List = append(group.List, item)
	return item
}

func (group *TruceGroup) UpdateList() []gp.SQLer {
	return nil
}

func (group *TruceGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *TruceGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *TruceGroup) SQLTable() string {
	return "truce"
}

func (group *TruceGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"trucee",
	}
}

func (group *TruceGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"trucee",
	}
}

func (group *TruceGroup) InsertScanCols() []string {
	return []string{}
}

func (group *TruceGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"trucee",
	}
}

func (group *TruceGroup) UpdateCols() []string {
	return nil
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type TruceSession struct {
	*TruceGroup
	*gp.Session
}

func NewTruceSession(d db.DBer) *TruceSession {
	group := NewTruceGroup()
	return &TruceSession{
		TruceGroup: group,
		Session:    gp.NewSession(group, d),
	}
}

func (s *TruceSession) Select(conditions ...interface{}) ([]overpower.TruceDat, error) {
	cur := len(s.TruceGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "Truce select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertTruce2Intf(s.TruceGroup.List[cur:]...), nil
}

func (s *TruceSession) SelectWhere(where sq.Condition) ([]overpower.TruceDat, error) {
	cur := len(s.TruceGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "Truce SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertTruce2Intf(s.TruceGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertTruce2Struct(list ...overpower.TruceDat) ([]*Truce, error) {
	mylist := make([]*Truce, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(TruceIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad Truce struct type for conversion")
		}
	}
	return mylist, nil
}

func convertTruce2Intf(list ...*Truce) []overpower.TruceDat {
	converted := make([]overpower.TruceDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func TruceTableCreate(d db.DBer) error {
	query := `create table truce(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid integer NOT NULL,
	locx int NOT NULL,
	locy int NOT NULL,
	trucee int NOT NULL,
	FOREIGN KEY(gid, fid) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, trucee) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, locx, locy) REFERENCES planet ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, locx, locy, trucee)
);`
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Truce table creation", "query", query); bad {
		return my
	}
	return nil
}

func TruceTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS truce CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Truce table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
