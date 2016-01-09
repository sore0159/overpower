package db

import (
	"database/sql"
	"fmt"
	"mule/mydb"
	"strconv"
)

const ODSQLVAL = `gid, fid, source, target, size`

func (o *Order) RowScan(row mydb.Scanner) error {
	return row.Scan(&(o.gid), &(o.fid), &(o.source), &(o.target), &(o.size))
}

func (o *Order) Insert(db mydb.SQLer) (ok bool) {
	return mydb.Insert(db, o)
}

func (o *Order) UpdateQ() (query string) {
	if o.size > 0 {
		return fmt.Sprintf("UPDATE orders SET size = %d WHERE gid = %d AND fid = %d AND source = %d AND target = %d", o.size, o.gid, o.fid, o.source, o.target)
	} else {
		return fmt.Sprintf("DELETE FROM orders WHERE gid = %d AND fid = %d AND source = %d AND target = %d", o.gid, o.fid, o.source, o.target)
	}
}

func (o *Order) Commit() {
}

func (o *Order) InsertScan(row *sql.Row) error {
	return nil
}
func (o *Order) InsertQ() (query string, scan bool) {
	if o.size < 1 {
		return "", false
	}
	return fmt.Sprintf(`INSERT INTO orders (%s) VALUES(
		%d, %d, 
		%d, %d, %d
	)`,
		ODSQLVAL,
		o.gid, o.fid,
		o.source, o.target, o.size,
	), false
}

func (o *Order) TableName() string {
	return "orders"
}
func (o *Order) SQLID() []string {
	return []string{"gid", strconv.Itoa(o.gid), "fid", strconv.Itoa(o.fid), "source", strconv.Itoa(o.source), "target", strconv.Itoa(o.target)}
}
