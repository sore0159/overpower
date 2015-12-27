package db

import (
	"mule/hexagon"
	"mule/mydb"
)

type MapView struct {
	*mydb.SQLHandler
	gid    int
	fid    int
	center hexagon.Coord
	zoom   int
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
