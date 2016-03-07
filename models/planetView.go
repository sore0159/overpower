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

type PlanetView struct {
	GID               int           `json:"gid"`
	FID               int           `json:"fid"`
	Loc               hexagon.Coord `json:"loc"`
	Name              string        `json:"name"`
	Turn              int           `json:"turn"`
	PrimaryFaction    sql.NullInt64 `json:"primaryfaction"`
	PrimaryPresence   int           `json:"primarypresence"`
	PrimaryPower      int           `json:"primarypower"`
	SecondaryFaction  sql.NullInt64 `json:"secondaryfaction"`
	SecondaryPresence int           `json:"secondarypresence"`
	SecondaryPower    int           `json:"secondarypower"`
	Antimatter        int           `json:"antimatter"`
	Tachyons          int           `json:"tachyons"`
	sql               gp.SQLStruct
}

// --------- BEGIN GENERIC METHODS ------------ //

func NewPlanetView() *PlanetView {
	return &PlanetView{
	//
	}
}

type PlanetViewIntf struct {
	item *PlanetView
}

func (item *PlanetView) Intf() overpower.PlanetViewDat {
	return &PlanetViewIntf{item}
}

func (i PlanetViewIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *PlanetView) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "fid":
		return item.FID
	case "locx":
		return item.Loc[0]
	case "locy":
		return item.Loc[1]
	case "name":
		return item.Name
	case "turn":
		return item.Turn
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
	}
	return nil
}

func (item *PlanetView) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "fid":
		return &item.FID
	case "locx":
		return &item.Loc[0]
	case "locy":
		return &item.Loc[1]
	case "name":
		return &item.Name
	case "turn":
		return &item.Turn
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
	}
	return nil
}
func (item *PlanetView) SQLTable() string {
	return "planetview"
}

func (i PlanetViewIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i PlanetViewIntf) UnmarshalJSON(data []byte) error {
	i.item = &PlanetView{}
	return json.Unmarshal(data, i.item)
}

func (i PlanetViewIntf) GID() int {
	return i.item.GID
}

func (i PlanetViewIntf) FID() int {
	return i.item.FID
}

func (i PlanetViewIntf) Loc() hexagon.Coord {
	return i.item.Loc
}

func (i PlanetViewIntf) Name() string {
	return i.item.Name
}

func (i PlanetViewIntf) Turn() int {
	return i.item.Turn
}

func (i PlanetViewIntf) SetTurn(x int) {
	if i.item.Turn == x {
		return
	}
	i.item.Turn = x
	i.item.sql.UPDATE = true
}

func (i PlanetViewIntf) PrimaryFaction() int {
	if !i.item.PrimaryFaction.Valid {
		return 0
	}
	return int(i.item.PrimaryFaction.Int64)
}

func (i PlanetViewIntf) SetPrimaryFaction(x int) {
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

func (i PlanetViewIntf) PrimaryPresence() int {
	return i.item.PrimaryPresence
}

func (i PlanetViewIntf) SetPrimaryPresence(x int) {
	if i.item.PrimaryPresence == x {
		return
	}
	i.item.PrimaryPresence = x
	i.item.sql.UPDATE = true
}

func (i PlanetViewIntf) PrimaryPower() int {
	return i.item.PrimaryPower
}

func (i PlanetViewIntf) SetPrimaryPower(x int) {
	if i.item.PrimaryPower == x {
		return
	}
	i.item.PrimaryPower = x
	i.item.sql.UPDATE = true
}

func (i PlanetViewIntf) SecondaryFaction() int {
	if !i.item.SecondaryFaction.Valid {
		return 0
	}
	return int(i.item.SecondaryFaction.Int64)
}

func (i PlanetViewIntf) SetSecondaryFaction(x int) {
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

func (i PlanetViewIntf) SecondaryPresence() int {
	return i.item.SecondaryPresence
}

func (i PlanetViewIntf) SetSecondaryPresence(x int) {
	if i.item.SecondaryPresence == x {
		return
	}
	i.item.SecondaryPresence = x
	i.item.sql.UPDATE = true
}

func (i PlanetViewIntf) SecondaryPower() int {
	return i.item.SecondaryPower
}

func (i PlanetViewIntf) SetSecondaryPower(x int) {
	if i.item.SecondaryPower == x {
		return
	}
	i.item.SecondaryPower = x
	i.item.sql.UPDATE = true
}

func (i PlanetViewIntf) Antimatter() int {
	return i.item.Antimatter
}

func (i PlanetViewIntf) SetAntimatter(x int) {
	if i.item.Antimatter == x {
		return
	}
	i.item.Antimatter = x
	i.item.sql.UPDATE = true
}

func (i PlanetViewIntf) Tachyons() int {
	return i.item.Tachyons
}

func (i PlanetViewIntf) SetTachyons(x int) {
	if i.item.Tachyons == x {
		return
	}
	i.item.Tachyons = x
	i.item.sql.UPDATE = true
}

// --------- END GENERIC METHODS ------------ //
// --------- BEGIN CUSTOM METHODS ------------ //

// --------- END CUSTOM METHODS ------------ //
// --------- BEGIN GROUP ------------ //

type PlanetViewGroup struct {
	List []*PlanetView
}

func NewPlanetViewGroup() *PlanetViewGroup {
	return &PlanetViewGroup{
		List: []*PlanetView{},
	}
}

func (item *PlanetView) SQLGroup() gp.SQLGrouper {
	return NewPlanetViewGroup()
}

func (group *PlanetViewGroup) New() gp.SQLer {
	item := NewPlanetView()
	group.List = append(group.List, item)
	return item
}

func (group *PlanetViewGroup) UpdateList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.UPDATE && !item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *PlanetViewGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *PlanetViewGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *PlanetViewGroup) SQLTable() string {
	return "planetview"
}

func (group *PlanetViewGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
	}
}

func (group *PlanetViewGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"name",
		"turn",
		"primaryfaction",
		"primarypresence",
		"primarypower",
		"secondaryfaction",
		"secondarypresence",
		"secondarypower",
		"antimatter",
		"tachyons",
	}
}

func (group *PlanetViewGroup) InsertScanCols() []string {
	return []string{}
}

func (group *PlanetViewGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"name",
		"turn",
		"primaryfaction",
		"primarypresence",
		"primarypower",
		"secondaryfaction",
		"secondarypresence",
		"secondarypower",
		"antimatter",
		"tachyons",
	}
}

func (group *PlanetViewGroup) UpdateCols() []string {
	return []string{
		"turn",
		"primaryfaction",
		"primarypresence",
		"primarypower",
		"secondaryfaction",
		"secondarypresence",
		"secondarypower",
		"antimatter",
		"tachyons",
	}
}

// --------- END GROUP ------------ //
// --------- BEGIN SESSION ------------ //
type PlanetViewSession struct {
	*PlanetViewGroup
	*gp.Session
}

func NewPlanetViewSession(d db.DBer) *PlanetViewSession {
	group := NewPlanetViewGroup()
	return &PlanetViewSession{
		PlanetViewGroup: group,
		Session:         gp.NewSession(group, d),
	}
}

func (s *PlanetViewSession) Select(conditions ...interface{}) ([]overpower.PlanetViewDat, error) {
	cur := len(s.PlanetViewGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "PlanetView select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertPlanetView2Intf(s.PlanetViewGroup.List[cur:]...), nil
}

func (s *PlanetViewSession) SelectWhere(where sq.Condition) ([]overpower.PlanetViewDat, error) {
	cur := len(s.PlanetViewGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "PlanetView SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertPlanetView2Intf(s.PlanetViewGroup.List[cur:]...), nil
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertPlanetView2Struct(list ...overpower.PlanetViewDat) ([]*PlanetView, error) {
	mylist := make([]*PlanetView, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(PlanetViewIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad PlanetView struct type for conversion")
		}
	}
	return mylist, nil
}

func convertPlanetView2Intf(list ...*PlanetView) []overpower.PlanetViewDat {
	converted := make([]overpower.PlanetViewDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func PlanetViewTableCreate(d db.DBer) error {
	query := `create table planetview(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	fid integer NOT NULL,
	locx int NOT NULL,
	locy int NOT NULL,
	name varchar(20) NOT NULL,
	turn int NOT NULL,
	primaryfaction int,
	primarypresence int NOT NULL,
	secondaryfaction int,
	secondarypresence int NOT NULL,
	antimatter int NOT NULL,
	tachyons int NOT NULL,
	FOREIGN KEY(gid, fid) REFERENCES faction ON DELETE CASCADE,
	FOREIGN KEY(gid, locx, locy) REFERENCES planet ON DELETE CASCADE,
	PRIMARY KEY(gid, fid, locx, locy)
);`
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed PlanetView table creation", "query", query); bad {
		return my
	}
	return nil
}

func PlanetViewTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS planetview CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed PlanetView table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
