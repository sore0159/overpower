package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"mule/hexagon"
	"mule/mydb/db"
	gp "mule/mydb/group"
	sq "mule/mydb/sql"
	"mule/overpower"
)

type BattleRecord struct {
	GID                   int           `json:"gid"`
	FID                   int           `json:"fid"`
	Loc                   hexagon.Coord `json:"loc"`
	Turn                  int           `json:"turn"`
	Index                 int           `json:"index"`
	Name                  string        `json:"name"`
	PrimaryFaction        sql.NullInt64 `json:"primaryfaction"`
	PrimaryPresence       int           `json:"primarypresence"`
	PrimaryPower          int           `json:"primarypower"`
	SecondaryFaction      sql.NullInt64 `json:"secondaryfaction"`
	SecondaryPresence     int           `json:"secondarypresence"`
	SecondaryPower        int           `json:"secondarypower"`
	Antimatter            int           `json:"antimatter"`
	Tachyons              int           `json:"tachyons"`
	ShipFaction           sql.NullInt64 `json:"shipfaction"`
	ShipSize              int           `json:"shipsize"`
	InitPrimaryFaction    sql.NullInt64 `json:"initprimaryfaction"`
	InitPrimaryPresence   int           `json:"initprimarypresence"`
	InitSecondaryFaction  sql.NullInt64 `json:"initsecondaryfaction"`
	InitSecondaryPresence int           `json:"initsecondarypresence"`
	Betrayals             db.IntList    `json:"betrayals"`
	sql                   gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewBattleRecord() *BattleRecord {
	return &BattleRecord{
	//
	}
}

type BattleRecordIntf struct {
	item *BattleRecord
}

func (item *BattleRecord) Intf() overpower.BattleRecordDat {
	return &BattleRecordIntf{item}
}

func (i BattleRecordIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *BattleRecord) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "fid":
		return item.FID
	case "locx":
		return item.Loc[0]
	case "locy":
		return item.Loc[1]
	case "turn":
		return item.Turn
	case "index":
		return item.Index
	case "name":
		return item.Name
	case "primaryfaction":
		return item.PrimaryFaction
	case "primarypresence":
		return item.PrimaryPresence
	case "primarypower":
		return item.PrimaryPower
	case "secondaryfaction":
		return item.SecondaryFaction
	case "secondarypresence":
		return item.SecondaryPresence
	case "secondarypower":
		return item.SecondaryPower
	case "antimatter":
		return item.Antimatter
	case "tachyons":
		return item.Tachyons
	case "shipfaction":
		return item.ShipFaction
	case "shipsize":
		return item.ShipSize
	case "initprimaryfaction":
		return item.InitPrimaryFaction
	case "initprimarypresence":
		return item.InitPrimaryPresence
	case "initsecondaryfaction":
		return item.InitSecondaryFaction
	case "initsecondarypresence":
		return item.InitSecondaryPresence
	case "betrayals":
		return item.Betrayals
	}
	return nil
}

func (item *BattleRecord) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "fid":
		return &item.FID
	case "locx":
		return &item.Loc[0]
	case "locy":
		return &item.Loc[1]
	case "turn":
		return &item.Turn
	case "index":
		return &item.Index
	case "name":
		return &item.Name
	case "primaryfaction":
		return &item.PrimaryFaction
	case "primarypresence":
		return &item.PrimaryPresence
	case "primarypower":
		return &item.PrimaryPower
	case "secondaryfaction":
		return &item.SecondaryFaction
	case "secondarypresence":
		return &item.SecondaryPresence
	case "secondarypower":
		return &item.SecondaryPower
	case "antimatter":
		return &item.Antimatter
	case "tachyons":
		return &item.Tachyons
	case "shipfaction":
		return &item.ShipFaction
	case "shipsize":
		return &item.ShipSize
	case "initprimaryfaction":
		return &item.InitPrimaryFaction
	case "initprimarypresence":
		return &item.InitPrimaryPresence
	case "initsecondaryfaction":
		return &item.InitSecondaryFaction
	case "initsecondarypresence":
		return &item.InitSecondaryPresence
	case "betrayals":
		return &item.Betrayals
	}
	return nil
}
func (item *BattleRecord) SQLTable() string {
	return "battlerecord"
}

func (i BattleRecordIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i BattleRecordIntf) UnmarshalJSON(data []byte) error {
	i.item = &BattleRecord{}
	return json.Unmarshal(data, i.item)
}

func (i BattleRecordIntf) GID() int {
	return i.item.GID
}

func (i BattleRecordIntf) FID() int {
	return i.item.FID
}

func (i BattleRecordIntf) Loc() hexagon.Coord {
	return i.item.Loc
}

func (i BattleRecordIntf) Turn() int {
	return i.item.Turn
}

func (i BattleRecordIntf) Index() int {
	return i.item.Index
}

func (i BattleRecordIntf) Name() string {
	return i.item.Name
}

func (i BattleRecordIntf) PrimaryFaction() int {
	if !i.item.PrimaryFaction.Valid {
		return 0
	}
	return int(i.item.PrimaryFaction.Int64)
}

func (i BattleRecordIntf) SetPrimaryFaction(x int) {
	if x == 0 {
		if !i.item.PrimaryFaction.Valid {
			return
		}
		i.item.PrimaryFaction.Valid = false
		i.item.PrimaryFaction.Int64 = 0
		return
	}
	x64 := int64(x)
	if i.item.PrimaryFaction.Valid && i.item.PrimaryFaction.Int64 == x64 {
		return
	}
	i.item.PrimaryFaction.Int64 = x64
	i.item.PrimaryFaction.Valid = true
	i.item.sql.UPDATE = true
}

func (i BattleRecordIntf) PrimaryPresence() int {
	return i.item.PrimaryPresence
}

func (i BattleRecordIntf) PrimaryPower() int {
	return i.item.PrimaryPower
}

func (i BattleRecordIntf) SecondaryFaction() int {
	if !i.item.SecondaryFaction.Valid {
		return 0
	}
	return int(i.item.SecondaryFaction.Int64)
}

func (i BattleRecordIntf) SetSecondaryFaction(x int) {
	if x == 0 {
		if !i.item.SecondaryFaction.Valid {
			return
		}
		i.item.SecondaryFaction.Valid = false
		i.item.SecondaryFaction.Int64 = 0
		return
	}
	x64 := int64(x)
	if i.item.SecondaryFaction.Valid && i.item.SecondaryFaction.Int64 == x64 {
		return
	}
	i.item.SecondaryFaction.Int64 = x64
	i.item.SecondaryFaction.Valid = true
	i.item.sql.UPDATE = true
}

func (i BattleRecordIntf) SecondaryPresence() int {
	return i.item.SecondaryPresence
}

func (i BattleRecordIntf) SecondaryPower() int {
	return i.item.SecondaryPower
}

func (i BattleRecordIntf) Antimatter() int {
	return i.item.Antimatter
}

func (i BattleRecordIntf) Tachyons() int {
	return i.item.Tachyons
}

func (i BattleRecordIntf) ShipFaction() int {
	if !i.item.ShipFaction.Valid {
		return 0
	}
	return int(i.item.ShipFaction.Int64)
}

func (i BattleRecordIntf) SetShipFaction(x int) {
	if x == 0 {
		if !i.item.ShipFaction.Valid {
			return
		}
		i.item.ShipFaction.Valid = false
		i.item.ShipFaction.Int64 = 0
		return
	}
	x64 := int64(x)
	if i.item.ShipFaction.Valid && i.item.ShipFaction.Int64 == x64 {
		return
	}
	i.item.ShipFaction.Int64 = x64
	i.item.ShipFaction.Valid = true
	i.item.sql.UPDATE = true
}

func (i BattleRecordIntf) ShipSize() int {
	return i.item.ShipSize
}

func (i BattleRecordIntf) InitPrimaryFaction() int {
	if !i.item.InitPrimaryFaction.Valid {
		return 0
	}
	return int(i.item.InitPrimaryFaction.Int64)
}

func (i BattleRecordIntf) SetInitPrimaryFaction(x int) {
	if x == 0 {
		if !i.item.InitPrimaryFaction.Valid {
			return
		}
		i.item.InitPrimaryFaction.Valid = false
		i.item.InitPrimaryFaction.Int64 = 0
		return
	}
	x64 := int64(x)
	if i.item.InitPrimaryFaction.Valid && i.item.InitPrimaryFaction.Int64 == x64 {
		return
	}
	i.item.InitPrimaryFaction.Int64 = x64
	i.item.InitPrimaryFaction.Valid = true
	i.item.sql.UPDATE = true
}

func (i BattleRecordIntf) InitPrimaryPresence() int {
	return i.item.InitPrimaryPresence
}

func (i BattleRecordIntf) InitSecondaryFaction() int {
	if !i.item.InitSecondaryFaction.Valid {
		return 0
	}
	return int(i.item.InitSecondaryFaction.Int64)
}

func (i BattleRecordIntf) SetInitSecondaryFaction(x int) {
	if x == 0 {
		if !i.item.InitSecondaryFaction.Valid {
			return
		}
		i.item.InitSecondaryFaction.Valid = false
		i.item.InitSecondaryFaction.Int64 = 0
		return
	}
	x64 := int64(x)
	if i.item.InitSecondaryFaction.Valid && i.item.InitSecondaryFaction.Int64 == x64 {
		return
	}
	i.item.InitSecondaryFaction.Int64 = x64
	i.item.InitSecondaryFaction.Valid = true
	i.item.sql.UPDATE = true
}

func (i BattleRecordIntf) InitSecondaryPresence() int {
	return i.item.InitSecondaryPresence
}

func (i BattleRecordIntf) Betrayals() [][2]int {
	bet := i.item.Betrayals
	r := make([][2]int, len(bet)/2)
	for i, _ := range r {
		r[i] = [2]int{bet[2*i], bet[2*i+1]}
	}
	return r
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type BattleRecordGroup struct {
	List []*BattleRecord
}

func NewBattleRecordGroup() *BattleRecordGroup {
	return &BattleRecordGroup{
		List: []*BattleRecord{},
	}
}

func (item *BattleRecord) SQLGroup() gp.SQLGrouper {
	return NewBattleRecordGroup()
}

func (group *BattleRecordGroup) New() gp.SQLer {
	item := NewBattleRecord()
	group.List = append(group.List, item)
	return item
}

func (group *BattleRecordGroup) UpdateList() []gp.SQLer {
	return nil
}

func (group *BattleRecordGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *BattleRecordGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *BattleRecordGroup) SQLTable() string {
	return "battlerecord"
}

func (group *BattleRecordGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"turn",
		"index",
	}
}

func (group *BattleRecordGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"turn",
		"index",
		"name",
		"primaryfaction",
		"primarypresence",
		"primarypower",
		"secondaryfaction",
		"secondarypresence",
		"secondarypower",
		"antimatter",
		"tachyons",
		"shipfaction",
		"shipsize",
		"initprimaryfaction",
		"initprimarypresence",
		"initsecondaryfaction",
		"initsecondarypresence",
		"betrayals",
	}
}

func (group *BattleRecordGroup) InsertScanCols() []string {
	return []string{}
}

func (group *BattleRecordGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"turn",
		"index",
		"name",
		"primaryfaction",
		"primarypresence",
		"primarypower",
		"secondaryfaction",
		"secondarypresence",
		"secondarypower",
		"antimatter",
		"tachyons",
		"shipfaction",
		"shipsize",
		"initprimaryfaction",
		"initprimarypresence",
		"initsecondaryfaction",
		"initsecondarypresence",
		"betrayals",
	}
}

func (group *BattleRecordGroup) UpdateCols() []string {
	return nil
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type BattleRecordSession struct {
	*BattleRecordGroup
	*gp.Session
}

func NewBattleRecordSession(d db.DBer) *BattleRecordSession {
	group := NewBattleRecordGroup()
	return &BattleRecordSession{
		BattleRecordGroup: group,
		Session:           gp.NewSession(group, d),
	}
}

func (s *BattleRecordSession) Select(conditions ...interface{}) ([]overpower.BattleRecordDat, error) {
	cur := len(s.BattleRecordGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "BattleRecord select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertBattleRecord2Intf(s.BattleRecordGroup.List[cur:]...), nil
}

func (s *BattleRecordSession) SelectWhere(where sq.Condition) ([]overpower.BattleRecordDat, error) {
	cur := len(s.BattleRecordGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "BattleRecord SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertBattleRecord2Intf(s.BattleRecordGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertBattleRecord2Struct(list ...overpower.BattleRecordDat) ([]*BattleRecord, error) {
	mylist := make([]*BattleRecord, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(BattleRecordIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad BattleRecord struct type for conversion")
		}
	}
	return mylist, nil
}

func convertBattleRecord2Intf(list ...*BattleRecord) []overpower.BattleRecordDat {
	converted := make([]overpower.BattleRecordDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func BattleRecordTableCreate(d db.DBer) error {
	query := `create table battlerecord(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid integer NOT NULL,
	turn int NOT NULL,
	index int NOT NULL,

	targetx integer NOT NULL,
	targety integer NOT NULL,

	shipfaction int,
	shipsize int NOT NULL,

	initprimaryfaction int,
	initprimarypresence int NOT NULL,
	initsecondaryfaction int,
	initsecondarypresence int NOT NULL,

	resultprimaryfaction int,
	resultprimarypresence int NOT NULL,
	resultsecondaryfaction int,
	resultsecondarypresence int NOT NULL,

	betrayals int[],

	FOREIGN KEY(gid, fid) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, shipfaction) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, initprimaryfaction) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, resultprimaryfaction) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, initprimaryfaction) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, initsecondaryfaction) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, targetx, targety) REFERENCES planet ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, turn, index)
);`
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed BattleRecord table creation", "query", query); bad {
		return my
	}
	return nil
}

func BattleRecordTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS battlerecord CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed BattleRecord table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
