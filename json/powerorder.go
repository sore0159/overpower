package json

import (
	"mule/hexagon"
	"mule/overpower"
)

func LoadPowerOrder(item overpower.PowerOrder) *PowerOrder {
	return &PowerOrder{
		Gid:     item.Gid(),
		Fid:     item.Fid(),
		Loc:     item.Loc(),
		UpPower: item.UpPower(),
	}
}

func LoadPowerOrders(list []overpower.PowerOrder) []*PowerOrder {
	jList := make([]*PowerOrder, len(list))
	for i, item := range list {
		jList[i] = LoadPowerOrder(item)
	}
	return jList
}

type PowerOrder struct {
	Gid     int           `json:"gid"`
	Fid     int           `json:"fid"`
	Loc     hexagon.Coord `json:"loc"`
	UpPower int           `json:"uppower"`
}
