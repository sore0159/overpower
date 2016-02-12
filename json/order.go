package json

import (
	"mule/hexagon"
	"mule/overpower"
)

func LoadOrder(item overpower.Order) *Order {
	return &Order{
		Gid:    item.Gid(),
		Fid:    item.Fid(),
		Turn:   item.Turn(),
		Source: item.Source(),
		Target: item.Target(),
		Size:   item.Size(),
	}
}

func LoadOrders(list []overpower.Order) []*Order {
	jList := make([]*Order, len(list))
	for i, item := range list {
		jList[i] = LoadOrder(item)
	}
	return jList
}

type Order struct {
	Gid    int           `json:"gid"`
	Fid    int           `json:"fid"`
	Turn   int           `json:"turn"`
	Source hexagon.Coord `json:"source"`
	Target hexagon.Coord `json:"target"`
	Size   int           `json:"size"`
}
