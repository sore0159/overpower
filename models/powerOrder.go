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

type PowerOrder struct {
	GID     int           `json:"gid"`
	FID     int           `json:"fid"`
	Loc     hexagon.Coord `json:"loc"`
	UpPower int           `json:"uppower"`
	sql     gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewPowerOrder() *PowerOrder {
	return &PowerOrder{
	//
	}
}

type PowerOrderIntf struct {
	item *PowerOrder
}

func (item *PowerOrder) Intf() overpower.PowerOrderDat {
	return &PowerOrderIntf{item}
}

func (i PowerOrderIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *PowerOrder) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "fid":
		return item.FID
	case "locx":
		return item.Loc[0]
	case "locy":
		return item.Loc[1]
	case "uppower":
		return item.UpPower
	}
	return nil
}

func (item *PowerOrder) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "fid":
		return &item.FID
	case "locx":
		return &item.Loc[0]
	case "locy":
		return &item.Loc[1]
	case "uppower":
		return &item.UpPower
	}
	return nil
}
func (item *PowerOrder) SQLTable() string {
	return "powerorder"
}

func (i PowerOrderIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i PowerOrderIntf) UnmarshalJSON(data []byte) error {
	i.item = &PowerOrder{}
	return json.Unmarshal(data, i.item)
}

func (i PowerOrderIntf) GID() int {
	return i.item.GID
}

func (i PowerOrderIntf) FID() int {
	return i.item.FID
}

func (i PowerOrderIntf) Loc() hexagon.Coord {
	return i.item.Loc
}

func (i PowerOrderIntf) UpPower() int {
	return i.item.UpPower
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type PowerOrderGroup struct {
	List []*PowerOrder
}

func NewPowerOrderGroup() *PowerOrderGroup {
	return &PowerOrderGroup{
		List: []*PowerOrder{},
	}
}

func (item *PowerOrder) SQLGroup() gp.SQLGrouper {
	return NewPowerOrderGroup()
}

func (group *PowerOrderGroup) New() gp.SQLer {
	item := NewPowerOrder()
	group.List = append(group.List, item)
	return item
}

func (group *PowerOrderGroup) UpdateList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.UPDATE && !item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *PowerOrderGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *PowerOrderGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *PowerOrderGroup) SQLTable() string {
	return "powerorder"
}

func (group *PowerOrderGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
	}
}

func (group *PowerOrderGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"uppower",
	}
}

func (group *PowerOrderGroup) InsertScanCols() []string {
	return []string{}
}

func (group *PowerOrderGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"uppower",
	}
}

func (group *PowerOrderGroup) UpdateCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"uppower",
	}
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type PowerOrderSession struct {
	*PowerOrderGroup
	*gp.Session
}

func NewPowerOrderSession(d db.DBer) *PowerOrderSession {
	group := NewPowerOrderGroup()
	return &PowerOrderSession{
		PowerOrderGroup: group,
		Session:         gp.NewSession(group, d),
	}
}

func (s *PowerOrderSession) Select(conditions ...interface{}) ([]overpower.PowerOrderDat, error) {
	cur := len(s.PowerOrderGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "PowerOrder select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertPowerOrder2Intf(s.PowerOrderGroup.List[cur:]...), nil
}

func (s *PowerOrderSession) SelectWhere(where sq.Condition) ([]overpower.PowerOrderDat, error) {
	cur := len(s.PowerOrderGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "PowerOrder SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertPowerOrder2Intf(s.PowerOrderGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertPowerOrder2Struct(list ...overpower.PowerOrderDat) ([]*PowerOrder, error) {
	mylist := make([]*PowerOrder, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(PowerOrderIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad PowerOrder struct type for conversion")
		}
	}
	return mylist, nil
}

func convertPowerOrder2Intf(list ...*PowerOrder) []overpower.PowerOrderDat {
	converted := make([]overpower.PowerOrderDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func PowerOrderTableCreate(d db.DBer) error {
	query := `create table powerorder(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid integer NOT NULL,
	locx int NOT NULL,
	locy int NOT NULL,
	uppower int NOT NULL,
	FOREIGN KEY(gid, fid) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, locx, locy) REFERENCES planet ON DELETE CASCADE,
	PRIMARY KEY(gid, fid)
);`

	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed PowerOrder table creation", "query", query); bad {
		return my
	}
	return nil
}

func PowerOrderTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS powerorder CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed PowerOrder table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
