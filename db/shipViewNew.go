package db

import (
	"mule/hexagon"
	//	"mule/overpower"
)

type ShipView struct {
	gid        int
	fid        int
	sid        int
	turn       int
	controller int
	size       int
	loc        hexagon.NullCoord
	dest       hexagon.NullCoord
	trail      hexagon.CoordList
}

func NewShipView() *ShipView {
	return &ShipView{
	//
	}
}

func (s *ShipView) Gid() int {
	return s.gid
}
func (s *ShipView) Fid() int {
	return s.fid
}
func (s *ShipView) Sid() int {
	return s.sid
}
func (s *ShipView) Turn() int {
	return s.turn
}
func (s *ShipView) Size() int {
	return s.size
}
func (s *ShipView) Controller() int {
	return s.controller
}
func (s *ShipView) Loc() (hexagon.Coord, bool) {
	if s.loc.Valid {
		return s.loc.Coord, true
	}
	return hexagon.Coord{0, 0}, false
}
func (s *ShipView) Trail() hexagon.CoordList {
	return s.trail
}
func (s *ShipView) Dest() (hexagon.Coord, bool) {
	if s.dest.Valid {
		return s.dest.Coord, true
	}
	return hexagon.Coord{0, 0}, false
}

func (item *ShipView) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	case "fid":
		return item.fid
	case "sid":
		return item.sid
	case "turn":
		return item.turn
	case "size":
		return item.size
	case "controller":
		return item.controller
	case "loc":
		return item.loc
	case "dest":
		return item.dest
	case "trail":
		return item.trail
	}
	return nil
}

func (item *ShipView) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	case "fid":
		return &item.fid
	case "sid":
		return &item.sid
	case "turn":
		return &item.turn
	case "size":
		return &item.size
	case "controller":
		return &item.controller
	case "loc":
		return &item.loc
	case "dest":
		return &item.dest
	case "trail":
		return &item.trail
	}
	return nil
}

func (item *ShipView) SQLTable() string {
	return "shipviews"
}

func (group *ShipViewGroup) SQLTable() string {
	return "shipviews"
}

func (group *ShipViewGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"sid",
		"turn",
		"controller",
		"size",
		"loc",
		"dest",
		"trail",
	}
}

func (group *ShipViewGroup) UpdateCols() []string {
	return nil
}

func (group *ShipViewGroup) PKCols() []string {
	return []string{"gid", "fid", "turn", "sid"}
}

func (group *ShipViewGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"sid",
		"turn",
		"controller",
		"size",
		"loc",
		"dest",
		"trail",
	}
}

func (group *ShipViewGroup) InsertScanCols() []string {
	return nil
}
