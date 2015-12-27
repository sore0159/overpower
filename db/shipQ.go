package db

import (
	"database/sql"
	"errors"
	"fmt"
	"mule/hexagon"
	"mule/mydb"
)

const SHSQLVAL = `sid, gid, fid, size, launched, path`

func (s *Ship) RowScan(row mydb.Scanner) error {
	//var pathStr str
	var pathBytes []byte
	err := row.Scan(&(s.sid), &(s.gid), &(s.fid), &(s.size), &(s.launched), &(pathBytes))
	if err != nil {
		return err
	}
	var ok bool
	s.path, ok = hexagon.Sql2CoordList(pathBytes)
	if !ok {
		return errors.New("BAD BYTES SCANNED FOR SHIP PATH")
	}
	return nil
}

func (s *Ship) Insert(db mydb.SQLer) (ok bool) {
	return mydb.Insert(db, s)
}
func (s *Ship) InsertScan(row *sql.Row) error {
	return row.Scan(&(s.sid))
}
func (s *Ship) InsertQ() (query string, scan bool) {
	pathStr := hexagon.CoordList2Sql(s.path)
	return fmt.Sprintf(`INSERT INTO ships (%s) VALUES(
		%d, %d,
		%d, %d, %s
	) RETURNING sid`,
		SHSQLVAL[5:],
		s.gid, s.fid,
		s.size, s.launched, pathStr,
	), true
}
