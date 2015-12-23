package db

import (
	"database/sql"
	"fmt"
	"mule/mydb"
	"strconv"
)

const PLSQLVAL = `pid, gid, name, loc, controller, inhabitants, resources, parts`

func (p *Planet) RowScan(row mydb.Scanner) error {
	var controller sql.NullInt64
	err := row.Scan(&(p.pid), &(p.gid), &(p.name), &(p.loc), &controller, &(p.inhabitants), &(p.resources), &(p.parts))
	if controller.Valid {
		p.controller = int(controller.Int64)
	}
	return err
}

func (p *Planet) Insert(db mydb.SQLer) (ok bool) {
	return mydb.Insert(db, p)
}

func (p *Planet) UpdateQ() (query string) {
	return mydb.ModderQ(p)
}

func (p *Planet) InsertScan(row *sql.Row) error {
	return nil
}
func (p *Planet) InsertQ() (query string, scan bool) {
	locVal, _ := p.loc.Value()
	var contVal string
	if p.controller == 0 {
		contVal = "NULL"
	} else {
		contVal = strconv.Itoa(p.controller)
	}
	return fmt.Sprintf(`INSERT INTO planets (%s) VALUES(
		%d, %d, '%s', %s,
		%s, %d, %d, %d
	)`,
		PLSQLVAL, p.pid, p.gid, p.name, locVal,
		contVal, p.inhabitants, p.resources, p.parts,
	), false
}

func (p *Planet) TableName() string {
	return "planets"
}
func (p *Planet) SQLID() []string {
	return []string{"gid", strconv.Itoa(p.gid), "pid", strconv.Itoa(p.pid)}
}
