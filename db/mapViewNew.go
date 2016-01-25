package db

import (
	"mule/hexagon"
)

type MapView struct {
	modified bool
	//
	gid     int
	fid     int
	center  hexagon.Coord
	target1 hexagon.NullCoord
	target2 hexagon.NullCoord
	zoom    int
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
	if mv.center == x {
		return
	}
	mv.center = x
	mv.modified = true
}

func (mv *MapView) Zoom() int {
	return mv.zoom
}
func (mv *MapView) SetZoom(x int) {
	if mv.zoom == x {
		return
	}
	mv.zoom = x
	mv.modified = true
}

func (mv *MapView) Target1() hexagon.NullCoord {
	return mv.target1
}
func (mv *MapView) Target2() hexagon.NullCoord {
	return mv.target2
}

func (mv *MapView) SetTarget1(x hexagon.Coord) {
	if mv.target1.Valid && mv.target1.Coord == x {
		return
	}
	mv.target1 = hexagon.NullCoord{x, true}
	mv.modified = true
}
func (mv *MapView) DropTarget1() {
	if !mv.target1.Valid {
		return
	}
	mv.target1 = hexagon.NullCoord{}
	mv.modified = true
}

func (mv *MapView) SetTarget2(x hexagon.Coord) {
	if mv.target2.Valid && mv.target2.Coord == x {
		return
	}
	mv.target2 = hexagon.NullCoord{x, true}
	mv.modified = true
}
func (mv *MapView) DropTarget2() {
	if !mv.target2.Valid {
		return
	}
	mv.target2 = hexagon.NullCoord{}
	mv.modified = true
}
func (mv *MapView) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return mv.gid
	case "fid":
		return mv.fid
	case "center":
		return mv.center
	case "target1":
		return mv.target1
	case "target2":
		return mv.target2
	case "zoom":
		return mv.zoom
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
	case "target1":
		return &mv.target1
	case "target2":
		return &mv.target2
	case "zoom":
		return &mv.zoom
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
		"target1",
		"target2",
		"zoom",
	}
}

func (group *MapViewGroup) UpdateCols() []string {
	return []string{
		"center",
		"target1",
		"target2",
		"zoom",
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
