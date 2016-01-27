package db

import (
	"mule/hexagon"
)

type Ship struct {
	justmade bool

	gid      int
	fid      int
	sid      int
	size     int
	launched int
	path     hexagon.CoordList
}

func NewShip() *Ship {
	return &Ship{
	//
	}
}

func (s *Ship) Gid() int {
	return s.gid
}
func (s *Ship) Fid() int {
	return s.fid
}
func (s *Ship) Sid() int {
	return s.sid
}
func (s *Ship) Size() int {
	return s.size
}
func (s *Ship) Launched() int {
	return s.launched
}
func (s *Ship) Path() hexagon.CoordList {
	return s.path
}

func (item *Ship) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	case "fid":
		return item.fid
	case "sid":
		return item.sid
	case "size":
		return item.size
	case "launched":
		return item.launched
	case "path":
		return item.path
	}
	return nil
}

func (item *Ship) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	case "fid":
		return &item.fid
	case "sid":
		return &item.sid
	case "size":
		return &item.size
	case "launched":
		return &item.launched
	case "path":
		return &item.path
	}
	return nil
}

func (item *Ship) SQLTable() string {
	return "ships"
}

func (group *ShipGroup) SQLTable() string {
	return "ships"
}

func (group *ShipGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"sid",
		"size",
		"launched",
		"path",
	}
}

func (group *ShipGroup) UpdateCols() []string {
	return []string{
	//
	}
}

func (group *ShipGroup) PKCols() []string {
	return []string{"gid", "fid", "sid"}
}

func (group *ShipGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"size",
		"launched",
		"path",
	}
}

func (group *ShipGroup) InsertScanCols() []string {
	return []string{"sid"}
}
