package db

import (
	"mule/hexagon"
)

type MapView struct {
	gid    int
	fid    int
	center hexagon.Coord
}

func NewMapView() *MapView {
	return &MapView{
	//
	}
}

func (mv *MapView) Gid() int {
	return mv.gid
}
func (mv *MapView) Fid() int {
	return mv.fid
}
func (mv *MapView) Center() hexagon.Coord {
	return mv.center
}
func (mv *MapView) SetCenter(x hexagon.Coord) {
	mv.center = x
}

func (mv *MapView) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return mv.gid
	case "fid":
		return mv.fid
	case "center":
		return mv.center
	}
	return nil
}
func (mv *MapView) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &mv.gid
	case "fid":
		return &mv.fid
	case "center":
		return &mv.center
	}
	return nil
}

func (mv *MapView) SQLTable() string {
	return "mapviews"
}

func (group *MapViewGroup) SQLTable() string {
	return "mapviews"
}

func (group *MapViewGroup) SelectCols() []string {
	return []string{
		"gid",
		"fid",
		"center",
	}
}

func (group *MapViewGroup) UpdateCols() []string {
	return []string{
		"center",
	}
}

func (group *MapViewGroup) PKCols() []string {
	return []string{"gid", "fid"}
}

func (group *MapViewGroup) InsertCols() []string {
	return []string{
		"gid",
		"fid",
		"center",
	}
}

func (group *MapViewGroup) InsertScanCols() []string {
	return nil
}
