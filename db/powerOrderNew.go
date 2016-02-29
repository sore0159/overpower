package db

import (
	"mule/hexagon"
)

type PowerOrder struct {
	gid     int
	fid     int
	loc     hexagon.Coord
	uppower bool
}

func NewPowerOrder() *PowerOrder {
	return &PowerOrder{
	//
	}
}

func (o *PowerOrder) Gid() int {
	return o.gid
}
func (o *PowerOrder) Fid() int {
	return o.fid
}
func (o *PowerOrder) Loc() hexagon.Coord {
	return o.loc
}
func (o *PowerOrder) SetLoc(x hexagon.Coord) {
	o.loc = x
}
func (o *PowerOrder) UpPower() bool {
	return o.uppower
}
func (o *PowerOrder) SetUpPower(x bool) {
	o.uppower = x
}
func (item *PowerOrder) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	case "fid":
		return item.fid
	case "locx":
		return item.loc[0]
	case "locy":
		return item.loc[1]
	case "uppower":
		return item.uppower
	}
	return nil
}

func (item *PowerOrder) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	case "fid":
		return &item.fid
	case "locx":
		return &item.loc[0]
	case "locy":
		return &item.loc[1]
	case "uppower":
		return &item.uppower
	}
	return nil
}

func (item *PowerOrder) SQLTable() string {
	return "powerorders"
}

func (group *PowerOrderGroup) SQLTable() string {
	return "powerorders"
}

func (group *PowerOrderGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"uppower",
	}
}

func (group *PowerOrderGroup) UpdateCols() []string {
	return []string{
		"locx",
		"locy",
		"uppower",
	}
}

func (group *PowerOrderGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
	}
}

func (group *PowerOrderGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"uppower",
	}
}

func (group *PowerOrderGroup) InsertScanCols() []string {
	return nil
}
