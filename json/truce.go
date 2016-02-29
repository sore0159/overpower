package json

import (
	"mule/hexagon"
	"mule/overpower"
)

func LoadTruce(item overpower.Truce) *Truce {
	return &Truce{
		Gid:    item.Gid(),
		Fid:    item.Fid(),
		Loc:    item.Loc(),
		Trucee: item.Trucee(),
	}
}

func LoadTruces(list []overpower.Truce) []*Truce {
	jList := make([]*Truce, len(list))
	for i, item := range list {
		jList[i] = LoadTruce(item)
	}
	return jList
}

type Truce struct {
	Gid    int           `json:"gid"`
	Fid    int           `json:"fid"`
	Loc    hexagon.Coord `json:"loc"`
	Trucee int           `json:"trucee"`
}
