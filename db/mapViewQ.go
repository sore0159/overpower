package db

import (
	"database/sql"
	"fmt"
	"mule/hexagon"
	"mule/mydb"
	"strconv"
)

const MVSQLVAL = `gid, fid, center, zoom, focus`

func (mv *MapView) RowScan(row mydb.Scanner) error {
	var focusBytes []byte
	err := row.Scan(&(mv.gid), &(mv.fid), &(mv.center), &(mv.zoom), &focusBytes)
	if err != nil {
		return err
	}
	if !mydb.CheckNull(focusBytes) {
		mv.focusValid = true
		c := hexagon.Coord{}
		(&c).Scan(focusBytes)
		mv.focus = c
	}
	return nil
}

func (mv *MapView) Insert(db mydb.SQLer) (ok bool) {
	return mydb.Insert(db, mv)
}

func (mv *MapView) UpdateQ() (query string) {
	return mydb.ModderQ(mv)
}

func (mv *MapView) InsertScan(row *sql.Row) error {
	return nil
}
func (mv *MapView) InsertQ() (query string, scan bool) {
	return fmt.Sprintf(`INSERT INTO mapviews (%s) VALUES(
		%d, %d,
		%s, %d,
		NULL
	)`,
		MVSQLVAL,
		mv.gid, mv.fid,
		mv.center.SQLStr(), mv.zoom,
	), false
}

func (mv *MapView) TableName() string {
	return "mapviews"
}
func (mv *MapView) SQLID() []string {
	return []string{"gid", strconv.Itoa(mv.gid), "fid", strconv.Itoa(mv.fid)}
}
