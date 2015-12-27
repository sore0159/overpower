package db

import (
	"fmt"
	"mule/mydb"
	"mule/overpower"
)

func (d DB) GetFidTurnShipViews(gid, fid, turn int) (ships []overpower.ShipView, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM shipviews WHERE gid = %d AND fid = %d AND turn = %d", SVSQLVAL, gid, fid, turn)
	return d.GetShipViewsQuery(query)
}

func (d DB) GetShipViewsQuery(query string) (ships []overpower.ShipView, ok bool) {
	sList := []*ShipView{}
	maker := func() mydb.Rower {
		s := NewShipView()
		return s
	}
	if !mydb.Get(d.db, query, &sList, maker) {
		return nil, false
	}
	ships = make([]overpower.ShipView, len(sList))
	for i, s := range sList {
		ships[i] = s
	}
	return ships, true
}
