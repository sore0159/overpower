package json

import (
	"mule/overpower"
)

func LoadOrder(item overpower.Order) *Order {
	return &Order{
		Gid:    item.Gid(),
		Fid:    item.Fid(),
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
	Gid    int `json:"gid"`
	Fid    int `json:"fid"`
	Source int `json:"source"`
	Target int `json:"target"`
	Size   int `json:"size"`
}
