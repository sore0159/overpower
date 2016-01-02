package db

import (
	"mule/hexagon"
	"mule/mydb"
)

type MapView struct {
	*mydb.SQLHandler
	gid     int
	fid     int
	center  hexagon.Coord
	target1 hexagon.NullCoord
	target2 hexagon.NullCoord
	zoom    int
}

func NewMapView() *MapView {
	return &MapView{
		SQLHandler: mydb.NewSQLHandler(),
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
	mv.SetEtc("center", x.SQLStr())
}

func (mv *MapView) Zoom() int {
	return mv.zoom
}
func (mv *MapView) SetZoom(x int) {
	if mv.zoom == x {
		return
	}
	mv.zoom = x
	mv.SetInt("zoom", x)
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
	mv.SetEtc("target1", x.SQLStr())
}
func (mv *MapView) DropTarget1() {
	if !mv.target1.Valid {
		return
	}
	mv.target1 = hexagon.NullCoord{}
	mv.SetNull("target1")
}

func (mv *MapView) SetTarget2(x hexagon.Coord) {
	if mv.target2.Valid && mv.target2.Coord == x {
		return
	}
	mv.target2 = hexagon.NullCoord{x, true}
	mv.SetEtc("target2", x.SQLStr())
}
func (mv *MapView) DropTarget2() {
	if !mv.target2.Valid {
		return
	}
	mv.target2 = hexagon.NullCoord{}
	mv.SetNull("target2")
}
