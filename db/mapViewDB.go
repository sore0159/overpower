package db

import (
	"fmt"
	"mule/hexagon"
	"mule/mydb"
	"mule/overpower"
)

func (d DB) GetFidMapView(gid, fid int) (mapview overpower.MapView, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM mapviews WHERE gid = %d AND fid = %d", MVSQLVAL, gid, fid)
	mv := NewMapView()
	return mv, mydb.GetOne(d.db, query, mv)
}

func (d DB) UpdateMapViewZoom(gid, fid, zoom int) (ok bool) {
	query := fmt.Sprintf("UPDATE mapviews SET zoom = %d WHERE gid = %d AND fid = %d", zoom, gid, fid)
	return mydb.Exec(d.db, query)
}
func (d DB) UpdateMapViewCenter(gid, fid int, center hexagon.Coord) (ok bool) {
	query := fmt.Sprintf("UPDATE mapviews SET center = %s WHERE gid = %d AND fid = %d", center.SQLStr(), gid, fid)
	return mydb.Exec(d.db, query)
}

func (d DB) UpdateMapViewTarget(gid, fid int, first bool, nc hexagon.NullCoord) (ok bool) {
	var tarStr string
	if first {
		tarStr = "target1"
	} else {
		tarStr = "target2"
	}
	query := fmt.Sprintf("UPDATE mapviews SET %s = %s WHERE gid = %d AND fid = %d", tarStr, nc.SQLStr(), gid, fid)
	return mydb.Exec(d.db, query)
}

func (d DB) UpdateMapViewBothTargets(gid, fid int, nc1, nc2 hexagon.NullCoord) (ok bool) {
	query := fmt.Sprintf("UPDATE mapviews SET target1 = %s, target2 = %s WHERE gid = %d AND fid = %d", nc1.SQLStr(), nc2.SQLStr(), gid, fid)
	return mydb.Exec(d.db, query)
}
