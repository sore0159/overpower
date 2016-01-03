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
		parts[i] = fmt.Sprintf("'%s'", strings.Replace(strings.Replace(x, "'", "''", -1), ",", "_", -1))
	}
	query := fmt.Sprintf("INSERT INTO reports (gid, fid, turn, contents) VALUES (%d, %d, %d, ARRAY[%s])", r.gid, r.fid, r.turn, strings.Join(parts, ", "))
	return mydb.Exec(db, query)
}

func (r *Report) RowScan(row mydb.Scanner) error {
	var conBytes []byte
	err := row.Scan(&conBytes)
	if err != nil {
		return err
	}
	str := string(conBytes)
	if len(str) < 2 {
		r.contents = []string{}
		return nil
	}
	parts := strings.Split(str, ",")
	if len(parts) == 1 {
		r.contents = []string{parts[0][1 : len(parts[0])-1]}
		return nil
	}
	list := make([]string, len(parts))
	for i, str := range parts {
		switch i {
		case 0:
			list[i] = strings.Replace(str[1:], "_", ",", -1)
		case len(parts) - 1:
			list[i] = strings.Replace(str[:len(str)-1], "_", ",", -1)
		default:
			list[i] = strings.Replace(str, "_", ",", -1)
		}
	}
	r.contents = list
	return nil
}
