package json

import (
	"mule/hexagon"
	"mule/overpower"
)

func LoadLaunchRecord(item overpower.LaunchRecord) *LaunchRecord {
	return &LaunchRecord{
		Gid:    item.Gid(),
		Fid:    item.Fid(),
		Turn:   item.Turn(),
		Source: item.Source(),
		Target: item.Target(),
		Size:   item.Size(),
	}
}

func LoadLaunchRecords(list []overpower.LaunchRecord) []*LaunchRecord {
	jList := make([]*LaunchRecord, len(list))
	for i, item := range list {
		jList[i] = LoadLaunchRecord(item)
	}
	return jList
}

type LaunchRecord struct {
	Gid    int           `json:"gid"`
	Fid    int           `json:"fid"`
	Turn   int           `json:"turn"`
	Source hexagon.Coord `json:"source"`
	Target hexagon.Coord `json:"target"`
	Size   int           `json:"size"`
}
