package db

import (
	"fmt"
	"mule/mydb"
	"mule/overpower"
	"strings"
)

func (d DB) DropSidShips(gid int, fsIds [][2]int) (ok bool) {
	if len(fsIds) == 0 {
		return true
	}
	query := fmt.Sprintf("DELETE FROM ships WHERE gid = %d AND (", gid)
	parts := make([]string, len(fsIds))
	for i, x := range fsIds {
		parts[i] = fmt.Sprintf("(fid = %d AND sid = %d)", x[0], x[1])
	}
	query += strings.Join(parts, " OR ") + ")"
	return mydb.Exec(d.db, query)
}

func (d DB) GetAllGidShips(gid int) (ships []overpower.Ship, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM ships WHERE gid = %d", SHSQLVAL, gid)
	return d.GetShipsQuery(query)
}

func (d DB) GetShipsQuery(query string) (ships []overpower.Ship, ok bool) {
	sList := []*Ship{}
	maker := func() mydb.Rower {
		s := NewShip()
		return s
	}
	if !mydb.Get(d.db, query, &sList, maker) {
		return nil, false
	}
	ships = make([]overpower.Ship, len(sList))
	for i, s := range sList {
		ships[i] = s
	}
	return ships, true
}
