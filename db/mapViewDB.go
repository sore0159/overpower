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
