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

type Planet struct {
	GID               int           `json:"gid"`
	Loc               hexagon.Coord `json:"loc"`
	Name              string        `json:"name"`
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

func NewPlanet() *Planet {
	return &Planet{
	//
	}
}

type PlanetIntf struct {
	item *Planet
}

func (item *Planet) Intf() overpower.PlanetDat {
	return &PlanetIntf{item}
}

func (i PlanetIntf) DELETE() {
	i.item.sql.DELETE = true
}

func (item *Planet) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.GID
	case "locx":
		return item.Loc[0]
	case "locy":
		return item.Loc[1]
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
	}
	return nil
}

func (item *Planet) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.GID
	case "locx":
		return &item.Loc[0]
	case "locy":
		return &item.Loc[1]
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
	}
	return nil
}
func (item *Planet) SQLTable() string {
	return "planet"
}

func (i PlanetIntf) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.item)
}
func (i PlanetIntf) UnmarshalJSON(data []byte) error {
	i.item = &Planet{}
	return json.Unmarshal(data, i.item)
}

func (i PlanetIntf) GID() int {
	return i.item.GID
}

func (i PlanetIntf) Loc() hexagon.Coord {
	return i.item.Loc
}

func (i PlanetIntf) Name() string {
	return i.item.Name
}

func (i PlanetIntf) PrimaryFaction() int {
	if !i.item.PrimaryFaction.Valid {
		return 0
	}
	return int(i.item.PrimaryFaction.Int64)
}

func (i PlanetIntf) SetPrimaryFaction(x int) {
	if x == 0 {
		if !i.item.PrimaryFaction.Valid {
			return
		}
		i.item.PrimaryFaction.Valid = false
		i.item.PrimaryFaction.Int64 = 0
		i.item.sql.UPDATE = true
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

func (i PlanetIntf) PrimaryPresence() int {
	return i.item.PrimaryPresence
}

func (i PlanetIntf) SetPrimaryPresence(x int) {
	if i.item.PrimaryPresence == x {
		return
	}
	i.item.PrimaryPresence = x
	i.item.sql.UPDATE = true
}

func (i PlanetIntf) PrimaryPower() int {
	return i.item.PrimaryPower
}

func (i PlanetIntf) SetPrimaryPower(x int) {
	if i.item.PrimaryPower == x {
		return
	}
	i.item.PrimaryPower = x
	i.item.sql.UPDATE = true
}

func (i PlanetIntf) SecondaryFaction() int {
	if !i.item.SecondaryFaction.Valid {
		return 0
	}
	return int(i.item.SecondaryFaction.Int64)
}

func (i PlanetIntf) SetSecondaryFaction(x int) {
	if x == 0 {
		if !i.item.SecondaryFaction.Valid {
			return
		}
		i.item.SecondaryFaction.Valid = false
		i.item.SecondaryFaction.Int64 = 0
		i.item.sql.UPDATE = true
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

func (i PlanetIntf) SecondaryPresence() int {
	return i.item.SecondaryPresence
}

func (i PlanetIntf) SetSecondaryPresence(x int) {
	if i.item.SecondaryPresence == x {
		return
	}
	i.item.SecondaryPresence = x
	i.item.sql.UPDATE = true
}

func (i PlanetIntf) SecondaryPower() int {
	return i.item.SecondaryPower
}

func (i PlanetIntf) SetSecondaryPower(x int) {
	if i.item.SecondaryPower == x {
		return
	}
	i.item.SecondaryPower = x
	i.item.sql.UPDATE = true
}

func (i PlanetIntf) Antimatter() int {
	return i.item.Antimatter
}

func (i PlanetIntf) SetAntimatter(x int) {
	if i.item.Antimatter == x {
		return
	}
	i.item.Antimatter = x
	i.item.sql.UPDATE = true
}

func (i PlanetIntf) Tachyons() int {
	return i.item.Tachyons
}

func (i PlanetIntf) SetTachyons(x int) {
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

type PlanetGroup struct {
	List []*Planet
}

func NewPlanetGroup() *PlanetGroup {
	return &PlanetGroup{
		List: []*Planet{},
	}
}

func (item *Planet) SQLGroup() gp.SQLGrouper {
	return NewPlanetGroup()
}

func (group *PlanetGroup) New() gp.SQLer {
	item := NewPlanet()
	group.List = append(group.List, item)
	return item
}

func (group *PlanetGroup) UpdateList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.UPDATE && !item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *PlanetGroup) InsertList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.INSERT && !item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *PlanetGroup) DeleteList() []gp.SQLer {
	list := make([]gp.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.sql.DELETE {
			list = append(list, item)
		}
	}
	return list
}

func (group *PlanetGroup) SQLTable() string {
	return "planet"
}

func (group *PlanetGroup) PKCols() []string {
	return []string{
		"gid",
		"locx",
		"locy",
	}
}

func (group *PlanetGroup) InsertCols() []string {
	return []string{
		"gid",
		"locx",
		"locy",
		"name",
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

func (group *PlanetGroup) InsertScanCols() []string {
	return []string{}
}

func (group *PlanetGroup) SelectCols() []string {
	return []string{
		"gid",
		"locx",
		"locy",
		"name",
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

func (group *PlanetGroup) UpdateCols() []string {
	return []string{
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
type PlanetSession struct {
	*PlanetGroup
	*gp.Session
}

func NewPlanetSession(d db.DBer) *PlanetSession {
	group := NewPlanetGroup()
	return &PlanetSession{
		PlanetGroup: group,
		Session:     gp.NewSession(group, d),
	}
}

func (s *PlanetSession) Select(conditions ...interface{}) ([]overpower.PlanetDat, error) {
	cur := len(s.PlanetGroup.List)
	err := s.Session.Select(conditions...)
	if my, bad := Check(err, "Planet select failed", "conditions", conditions); bad {
		return nil, my
	}
	return convertPlanet2Intf(s.PlanetGroup.List[cur:]...), nil
}

func (s *PlanetSession) SelectWhere(where sq.Condition) ([]overpower.PlanetDat, error) {
	cur := len(s.PlanetGroup.List)
	err := s.Session.SelectWhere(where)
	if my, bad := Check(err, "Planet SelectWhere failed", "where", where); bad {
		return nil, my
	}
	return convertPlanet2Intf(s.PlanetGroup.List[cur:]...), nil
}
func (s *PlanetSession) SelectByLocs(gid int, locations ...hexagon.Coord) ([]overpower.PlanetDat, error) {
	coordWhere := make([]sq.Condition, len(locations))
	for i, loc := range locations {
		coordWhere[i] = sq.AND(sq.EQ("locx", loc[0]), sq.EQ("locy", loc[1]))
	}
	where := sq.AND(sq.EQ("gid", gid), sq.OR(coordWhere...))
	return s.SelectWhere(where)
}

// --------- END SESSION  ------------ //
// --------- BEGIN UTILS ------------ //

func convertPlanet2Struct(list ...overpower.PlanetDat) ([]*Planet, error) {
	mylist := make([]*Planet, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(PlanetIntf); ok {
			mylist = append(mylist, t.item)
		} else {
			return nil, errors.New("bad Planet struct type for conversion")
		}
	}
	return mylist, nil
}

func convertPlanet2Intf(list ...*Planet) []overpower.PlanetDat {
	converted := make([]overpower.PlanetDat, len(list))
	for i, item := range list {
		converted[i] = item.Intf()
	}
	return converted
}

func PlanetTableCreate(d db.DBer) error {
	query := `create table planet(
	gid integer NOT NULL REFERENCES game ON DELETE CASCADE,
	locx int NOT NULL,
	locy int NOT NULL,
	name varchar(20) NOT NULL,
	primaryfaction int REFERENCES faction ON DELETE SET NULL,
	primarypresence int NOT NULL,
	primarypower int NOT NULL,
	secondaryfaction int REFERENCES faction ON DELETE SET NULL,
	secondarypresence int NOT NULL,
	secondarypower int NOT NULL,
	antimatter int NOT NULL,
	tachyons int NOT NULL,
	UNIQUE(gid, name),
	PRIMARY KEY(gid, locx, locy)
);`
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Planet table creation", "query", query); bad {
		return my
	}
	return nil
}

func PlanetTableDelete(d db.DBer) error {
	query := "DROP TABLE IF EXISTS planet CASCADE"
	err := db.Exec(d, false, query)
	if my, bad := Check(err, "failed Planet table deletion", "query", query); bad {
		return my
	}
	return nil
}

// --------- END UTILS ------------ //
