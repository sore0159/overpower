package db

import (
	"fmt"
	"mule/mydb"
	"mule/overpower"
	"strings"
)

func (d DB) GetReport(gid, fid, turn int) (overpower.Report, bool) {
	query := fmt.Sprintf("SELECT contents FROM reports WHERE gid = %d AND fid = %d AND turn = %d", gid, fid, turn)
	r := &Report{gid, fid, turn, []string{}}
	return r, mydb.GetOneIf(d.db, query, r)
}

func (r *Report) Insert(db mydb.SQLer) (ok bool) {
	if len(r.contents) == 0 {
		return true
	}
	parts := make([]string, len(r.contents))
	for i, x := range r.contents {
		parts[i] = fmt.Sprintf("'%s'", strings.Replace(x, "'", "''", -1))
	}
	query := fmt.Sprintf("INSERT INTO contents (gid, fid, turn, contents) VALUES (", r.gid, r.fid, r.turn)
	query += strings.Join(parts, ", ") + ")"
	_ = query
	return true
}

func (r *Report) RowScan(s mydb.Scanner) error {
	return nil
}
