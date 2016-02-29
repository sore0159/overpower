package db

import (
	"mule/hexagon"
)

type Truce struct {
	gid    int
	fid    int
	loc    hexagon.Coord
	trucee int
}

func NewTruce() *Truce {
	return &Truce{
	//
	}
}

func (o *Truce) Gid() int {
	return o.gid
}
func (o *Truce) Fid() int {
	return o.fid
}
func (o *Truce) Loc() hexagon.Coord {
	return o.loc
}
func (o *Truce) SetLoc(x hexagon.Coord) {
	o.loc = x
}
func (o *Truce) Trucee() int {
	return o.trucee
}
func (item *Truce) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	case "fid":
		return item.fid
	case "locx":
		return item.loc[0]
	case "locy":
		return item.loc[1]
	case "trucee":
		return item.trucee
	}
	return nil
}

func (item *Truce) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	case "fid":
		return &item.fid
	case "locx":
		return &item.loc[0]
	case "locy":
		return &item.loc[1]
	case "trucee":
		return &item.trucee
	}
	return nil
}

func (item *Truce) SQLTable() string {
	return "truces"
}

func (group *TruceGroup) SQLTable() string {
	return "truces"
}

func (group *TruceGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"trucee",
	}
}

func (group *TruceGroup) UpdateCols() []string {
	return nil
}

func (group *TruceGroup) PKCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"trucee",
	}
}

func (group *TruceGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"locx",
		"locy",
		"trucee",
	}
}

func (group *TruceGroup) InsertScanCols() []string {
	return nil
}
