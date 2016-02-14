package json

import (
	"mule/hexagon"
	"mule/overpower"
)

func LoadLandingRecord(item overpower.LandingRecord) *LandingRecord {
	return &LandingRecord{
		Gid:               item.Gid(),
		Fid:               item.Fid(),
		Turn:              item.Turn(),
		Index:             item.Index(),
		Target:            item.Target(),
		Size:              item.Size(),
		ShipController:    item.ShipController(),
		FirstController:   item.FirstController(),
		ResultController:  item.ResultController(),
		ResultInhabitants: item.ResultInhabitants(),
	}
}

func LoadLandingRecords(list []overpower.LandingRecord) []*LandingRecord {
	jList := make([]*LandingRecord, len(list))
	for i, item := range list {
		jList[i] = LoadLandingRecord(item)
	}
	return jList
}

type LandingRecord struct {
	Gid               int           `json:"gid"`
	Fid               int           `json:"fid"`
	Turn              int           `json:"turn"`
	Index             int           `json:"index"`
	Size              int           `json:"size"`
	Target            hexagon.Coord `json:"target"`
	ShipController    int           `json:"shipcontroller"`
	FirstController   int           `json:"firstcontroller"`
	ResultController  int           `json:"resultcontroller"`
	ResultInhabitants int           `json:"resultinhabitants"`
}
