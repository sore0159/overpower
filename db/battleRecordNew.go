package db

import (
	"database/sql"
	"mule/mydb"
)

type BattleRecord struct {
	*MiniPlanet

	index int

	shipfaction           sql.NullInt64
	shipsize              int
	initprimaryfaction    sql.NullInt64
	initprimarypresence   int
	initsecondaryfaction  sql.NullInt64
	initsecondarypresence int

	betrayals mydb.IntList
}

func NewBattleRecord() *BattleRecord {
	return &BattleRecord{
		MiniPlanet: &MiniPlanet{},
	}
}

func (item *BattleRecord) Index() int {
	return item.index
}

func (p *BattleRecord) ShipFaction() int {
	if p.shipfaction.Valid {
		return int(p.shipfaction.Int64)
	}
	return 0
}
func (p *BattleRecord) ShipSize() int {
	return p.shipsize
}
func (p *BattleRecord) InitPrimaryPresence() int {
	return p.initprimarypresence
}
func (p *BattleRecord) InitSecondaryPresence() int {
	return p.initsecondarypresence
}
func (p *BattleRecord) InitPrimaryFaction() int {
	if p.initprimaryfaction.Valid {
		return int(p.initprimaryfaction.Int64)
	}
	return 0
}
func (p *BattleRecord) InitSecondaryFaction() int {
	if p.initsecondaryfaction.Valid {
		return int(p.initprimaryfaction.Int64)
	}
	return 0
}
func (p *BattleRecord) ResultSecondaryFaction() int {
	return p.SecondaryFaction()
}
func (p *BattleRecord) ResultPrimaryFaction() int {
	return p.PrimaryFaction()
}
func (p *BattleRecord) ResultSecondaryPresence() int {
	return p.SecondaryPresence()
}
func (p *BattleRecord) ResultPrimaryPresence() int {
	return p.PrimaryPresence()
}
func (p *BattleRecord) Betrayals() [][2]int {
	r := make([][2]int, len(p.betrayals)/2)
	for i, _ := range r {
		r[i] = [2]int{p.betrayals[2*i], p.betrayals[2*i+1]}
	}
	return r
}
func (p *BattleRecord) SetBetrayals(list [][2]int) {
	r := make([]int, len(list)*2)
	for i, x := range list {
		r[2*i] = x[0]
		r[2*i+1] = x[1]
	}
	p.betrayals = r
	p.modified = true
}

func (i *BattleRecord) SQLVal(name string) interface{} {
	switch name {
	case "index":
		return i.index
	case "shipfaction":
		return i.shipfaction
	case "shipsize":
		return i.shipsize
	case "initprimaryfaction":
		return i.initprimaryfaction
	case "initsecondaryfaction":
		return i.initsecondaryfaction
	case "initprimarypresence":
		return i.initprimarypresence
	case "initsecondarypresence":
		return i.initsecondarypresence
	case "resultprimaryfaction":
		return i.primaryfaction
	case "resultsecondaryfaction":
		return i.secondaryfaction
	case "resultprimarypresence":
		return i.primarypresence
	case "resultsecondarypresence":
		return i.secondarypresence
	case "betrayals":
		return i.betrayals
	default:
		return i.MiniPlanet.SQLVal(name)
	}
}

func (i *BattleRecord) SQLPtr(name string) interface{} {
	switch name {
	case "index":
		return &i.index
	case "shipfaction":
		return &i.shipfaction
	case "shipsize":
		return &i.shipsize
	case "initprimaryfaction":
		return &i.initprimaryfaction
	case "initsecondaryfaction":
		return &i.initsecondaryfaction
	case "initprimarypresence":
		return &i.initprimarypresence
	case "initsecondarypresence":
		return &i.initsecondarypresence
	case "resultprimaryfaction":
		return &i.primaryfaction
	case "resultsecondaryfaction":
		return &i.secondaryfaction
	case "resultprimarypresence":
		return &i.primarypresence
	case "resultsecondarypresence":
		return &i.secondarypresence
	case "betrayals":
		return &i.betrayals
	default:
		return i.MiniPlanet.SQLPtr(name)
	}
}

func (mv *BattleRecord) SQLTable() string {
	return "battlerecords"
}

func (group *BattleRecordGroup) SQLTable() string {
	return "battlerecords"
}

func (group *BattleRecordGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"index",
		"targetx",
		"targety",

		"shipfaction",
		"shipsize",
		"initprimaryfaction",
		"initprimarypresence",
		"initsecondaryfaction",
		"initsecondarypresence",
		"resultprimaryfaction",
		"resultprimarypresence",
		"resultsecondfaction",
		"resultsecondpresence",
		"betrayals",
	}
}

func (group *BattleRecordGroup) UpdateCols() []string {
	return nil
}

func (group *BattleRecordGroup) PKCols() []string {
	return []string{"gid", "fid", "turn", "index"}
}

func (group *BattleRecordGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"index",
		"targetx",
		"targety",

		"shipfaction",
		"shipsize",
		"initprimaryfaction",
		"initprimarypresence",
		"initsecondaryfaction",
		"initsecondarypresence",
		"resultprimaryfaction",
		"resultprimarypresence",
		"resultsecondfaction",
		"resultsecondpresence",
		"betrayals",
	}
}

func (group *BattleRecordGroup) InsertScanCols() []string {
	return nil
}
