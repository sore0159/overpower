package db

import (
	"database/sql"
	"errors"
	"fmt"
	"mule/hexagon"
	"mule/mydb"
)

const SVSQLVAL = `sid, gid, fid, turn, size, controller, loc, dest, trail`

func (s *ShipView) RowScan(row mydb.Scanner) error {
	var trailBytes []byte
	var locBytes []byte
	var destBytes []byte
	err := row.Scan(&(s.sid), &(s.gid), &(s.fid), &(s.turn), &(s.size), &(s.controller), &(locBytes), &(destBytes), &(trailBytes))
	if err != nil {
		return err
	}
	if trailBytes == nil {
		//if mydb.CheckNull(trailBytes) {
		s.trail = []hexagon.Coord{}
	} else {
		var ok bool
		s.trail, ok = hexagon.Sql2CoordList(trailBytes)
		if !ok {
			return errors.New("BAD BYTES SCANNED FOR SHIPVIEW TRAIL")
		}
	}
	//if !mydb.CheckNull(locBytes) {
	if locBytes != nil {
		s.locValid = true
		c := hexagon.Coord{}
		(&c).Scan(locBytes)
		s.loc = c
	}
	if destBytes != nil {
		//if !mydb.CheckNull(destBytes) {
		s.destValid = true
		c := hexagon.Coord{}
		(&c).Scan(destBytes)
		s.dest = c
	}
	return nil
}

func (s *ShipView) Insert(db mydb.SQLer) (ok bool) {
	return mydb.Insert(db, s)
}
func (s *ShipView) InsertScan(row *sql.Row) error {
	return nil
}
func (s *ShipView) InsertQ() (query string, scan bool) {
	var trailStr, locStr, destStr string
	if len(s.trail) > 0 {
		trailStr = hexagon.CoordList2Sql(s.trail)
	} else {
		trailStr = "NULL"
	}
	if s.locValid {
		locStr = s.loc.SQLStr()
	} else {
		locStr = "NULL"
	}
	if s.destValid {
		destStr = s.dest.SQLStr()
	} else {
		destStr = "NULL"
	}
	return fmt.Sprintf(`INSERT INTO shipviews (%s) VALUES(
		%d, %d, %d, %d,
		%d, %d, 
		%s, %s, %s
	) RETURNING sid`,
		SVSQLVAL,
		s.sid, s.gid, s.fid, s.turn,
		s.size, s.controller,
		locStr, destStr, trailStr,
	), false
}
