package json

import (
	"mule/hexagon"
	"mule/overpower"
)

func LoadBattleRecord(item overpower.BattleRecord) *BattleRecord {
	return &BattleRecord{
		Gid:   item.Gid(),
		Fid:   item.Fid(),
		Turn:  item.Turn(),
		Index: item.Index(),
		Loc:   item.Loc(),

		ShipFaction:           item.ShipFaction(),
		ShipSize:              item.ShipSize(),
		InitPrimaryFaction:    item.InitPrimaryFaction(),
		InitPrimaryPresence:   item.InitPrimaryPresence(),
		InitSecondaryFaction:  item.InitSecondaryFaction(),
		InitSecondaryPresence: item.InitSecondaryPresence(),

		ResultPrimaryFaction:    item.ResultPrimaryFaction(),
		ResultPrimaryPresence:   item.ResultPrimaryPresence(),
		ResultSecondaryFaction:  item.ResultSecondaryFaction(),
		ResultSecondaryPresence: item.ResultSecondaryPresence(),

		Betrayals: item.Betrayals(),
	}
}

func LoadBattleRecords(list []overpower.BattleRecord) []*BattleRecord {
	jList := make([]*BattleRecord, len(list))
	for i, item := range list {
		jList[i] = LoadBattleRecord(item)
	}
	return jList
}

type BattleRecord struct {
	Gid   int           `json:"gid"`
	Fid   int           `json:"fid"`
	Turn  int           `json:"turn"`
	Index int           `json:"index"`
	Size  int           `json:"size"`
	Loc   hexagon.Coord `json:"loc"`

	ShipFaction           int `json:"shipfaction"`
	ShipSize              int `json:"shipsize"`
	InitPrimaryFaction    int `json:"initprimaryfaction"`
	InitPrimaryPresence   int `json:"initprimarypresence"`
	InitSecondaryFaction  int `json:"initsecondaryfaction"`
	InitSecondaryPresence int `json:"initsecondarypresence"`

	ResultPrimaryFaction    int `json:"resultprimaryfaction"`
	ResultPrimaryPresence   int `json:"resultprimarypresence"`
	ResultSecondaryFaction  int `json:"resultsecondaryfaction"`
	ResultSecondaryPresence int `json:"resultsecondarypresence"`

	Betrayals [][2]int `json:"betrayals"`
}
