package models

import (
	"encoding/json"
	"errors"
	"mule/mydb/db"
	gp "mule/mydb/group"
	sq "mule/mydb/sql"
	"mule/overpower"
)

type Faction struct {
	GID        int    `json:"gid"`
	FID        int    `json:"fid"`
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	DoneBuffer int    `json:"donebuffer"`
	Score      int    `json:"score"`
	sql        gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewFaction() *Faction {
	return &Faction{
	//
	}
}

type FactionIntf struct {
	item *Faction
}

func (item *Faction) Intf() overpower.FactionDat {
	return &FactionIntf{item}
}

func (i FactionIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *Faction) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "fid":
		return item.FID
	case "owner":
		return item.Owner
	case "name":
		return item.Name
	case "donebuffer":
		return item.DoneBuffer
	case "score":
		return item.Score
	}
	return nil
}

func (item *Faction) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "fid":
		return &item.FID
	case "owner":
		return &item.Owner
	case "name":
		return &item.Name
	case "donebuffer":
		return &item.DoneBuffer
	case "score":
		return &item.Score
	}
	return nil
}
func (item *Faction) SQLTable() string {
	return "faction"
}

func (i FactionIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i FactionIntf) MarshalPublicJSON() ([]byte, error) {
	return json.Marshal(struct {
		GID      int    `json:"gid"`
		FID      int    `json:"fid"`
		Owner    string `json:"owner"`
		Name     string `json:"name"`
		TurnDone bool   `json:"turndone"`
	}{
		GID:      i.GID(),
		FID:      i.FID(),
		Owner:    i.Owner(),
		Name:     i.Name(),
		TurnDone: i.DoneBuffer() != 0,
	})
}
func (i FactionIntf) UnmarshalJSON(data []byte) error {
	i.item = &Faction{}
	return json.Unmarshal(data, i.item)
}

func (i FactionIntf) GID() int {
	return i.item.GID
}

func (i FactionIntf) FID() int {
	return i.item.FID
}

func (i FactionIntf) Owner() string {
	return i.item.Owner
}

func (i FactionIntf) Name() string {
	return i.item.Name
}

func (i FactionIntf) IsDone() bool {
	return i.item.DoneBuffer != 0
}

func (i FactionIntf) DoneBuffer() int {
	return i.item.DoneBuffer
}

func (i FactionIntf) SetDoneBuffer(x int) {
	if i.item.DoneBuffer == x {
		return
	}
	i.item.DoneBuffer = x
	i.item.sql.UPDATE = true
}

func (i FactionIntf) Score() int {
	return i.item.Score
}

func (i FactionIntf) SetScore(x int) {
	if i.item.Score == x {
		return
	}
	i.item.Score = x
	i.item.sql.UPDATE = true
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type FactionGroup struct {
	List []*Faction
}

func NewFactionGroup() *FactionGroup {
	return &FactionGroup{
		List: []*Faction{},
	}
}

func (item *Faction) SQLGroup() gp.SQLGrouper {
	return NewFactionGroup()
}

func (group *FactionGroup) New() gp.SQLer {
	item := NewFaction()
	group.List = append(group.List, item)
	return item
}

func (group *FactionGroup) UpdateList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.UPDATE && !item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *FactionGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *FactionGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *FactionGroup) SQLTable() string {
	return "faction"
}

func (group *FactionGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
	}
}

func (group *FactionGroup) InsertCols() []string {
	return []string{
		"gid",
		//		"fid",
		"owner",
		"name",
		"donebuffer",
		"score",
	}
}

func (group *FactionGroup) InsertScanCols() []string {
	return []string{
		"fid",
	}
}

func (group *FactionGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"owner",
		"name",
		"donebuffer",
		"score",
	}
}

func (group *FactionGroup) UpdateCols() []string {
	return []string{
		"donebuffer",
		"score",
	}
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type FactionSession struct {
	*FactionGroup
	*gp.Session
}

func NewFactionSession(d db.DBer) *FactionSession {
	group := NewFactionGroup()
	return &FactionSession{
		FactionGroup: group,
		Session:      gp.NewSession(group, d),
	}
}

func (s *FactionSession) Select(conditions ...interface{}) ([]overpower.FactionDat, error) {
	cur := len(s.FactionGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "Faction select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertFaction2Intf(s.FactionGroup.List[cur:]...), nil
}

func (s *FactionSession) SelectWhere(where sq.Condition) ([]overpower.FactionDat, error) {
	cur := len(s.FactionGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "Faction SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertFaction2Intf(s.FactionGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertFaction2Struct(list ...overpower.FactionDat) ([]*Faction, error) {
	mylist := make([]*Faction, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(FactionIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad Faction struct type for conversion")
		}
	}
	return mylist, nil
}

func convertFaction2Intf(list ...*Faction) []overpower.FactionDat {
	converted := make([]overpower.FactionDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func FactionTableCreate(d db.DBer) error {
	query := `create table faction(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid SERIAL NOT NULL,
	owner varchar(20) NOT NULL,
	name varchar(20) NOT NULL,
	donebuffer int NOT NULL DEFAULT 0,
	score int NOT NULL DEFAULT 0,
	UNIQUE(gid, owner),
	PRIMARY KEY(gid, fid)
);`
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Faction table creation", "query", query); bad {
		return my
	}
	return nil
}

func FactionTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS faction CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Faction table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
