package db

import (
	"mule/hexagon"
)

type Order struct {
	gid    int
	fid    int
	turn   int
	source hexagon.Coord
	target hexagon.Coord
	size   int
}

func NewOrder() *Order {
	return &Order{
	//
	}
}

func (o *Order) Gid() int {
	return o.gid
}
func (o *Order) Fid() int {
	return o.fid
}
func (o *Order) Source() hexagon.Coord {
	return o.source
}
func (o *Order) Target() hexagon.Coord {
	return o.target
}
func (o *Order) Turn() int {
	return o.turn
}
func (o *Order) Size() int {
	return o.size
}
func (o *Order) SetSize(size int) {
	o.size = size
}
func (item *Order) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	case "fid":
		return item.fid
	case "turn":
		return item.turn
	case "sourcex":
		return item.source[0]
	case "sourcey":
		return item.source[1]
	case "targetx":
		return item.target[0]
	case "targety":
		return item.target[1]
	case "size":
		return item.size
	}
	return nil
}

func (item *Order) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	case "fid":
		return &item.fid
	case "turn":
		return &item.turn
	case "sourcex":
		return &item.source[0]
	case "sourcey":
		return &item.source[1]
	case "targetx":
		return &item.target[0]
	case "targety":
		return &item.target[1]
	case "size":
		return &item.size
	}
	return nil
}

func (item *Order) SQLTable() string {
	return "orders"
}

func (group *OrderGroup) SQLTable() string {
	return "orders"
}

func (group *OrderGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
		"size",
	}
}

func (group *OrderGroup) UpdateCols() []string {
	return []string{
		"size",
	}
}

func (group *OrderGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
	}
}

func (group *OrderGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"turn",
		"sourcex",
		"sourcey",
		"targetx",
		"targety",
		"size",
	}
}

func (group *OrderGroup) InsertScanCols() []string {
	return nil
}
