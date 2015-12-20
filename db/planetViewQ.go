package db

import (
	"database/sql"
	"fmt"
	"mule/mydb"
	"strconv"
)

const PVSQLVAL = `pid, gid, fid, turn, name, loc, controller, inhabitants, resources, parts`

func (p *PlanetView) RowScan(row mydb.Scanner) error {
	var controller sql.NullInt64
	err := row.Scan(&(p.pid), &(p.gid), &(p.fid), &(p.turn), &(p.name), &(p.loc), &controller, &(p.inhabitants), &(p.resources), &(p.parts))
	if controller.Valid {
		p.controller = int(controller.Int64)
	}
	return err
}

func (p *PlanetView) Insert(db *sql.DB) (ok bool) {
	return mydb.Insert(db, p)
}

func (p *PlanetView) UpdateQ() (query string) {
	return mydb.ModderQ(p)
}

func (p *PlanetView) InsertScan(row *sql.Row) error {
	return row.Scan(&(p.pid))
}
func (p *PlanetView) InsertQ() (query string, scan bool) {
	locVal, _ := p.loc.Value()
	var contVal string
	if p.controller == 0 {
		contVal = "NULL"
	} else {
		contVal = strconv.Itoa(p.controller)
	}
	return fmt.Sprintf(`INSERT INTO planetviews (%s) VALUES(
		%d, %d, %d,
		%d, '%s', %s,
		%s, %d, %d, %d
	) RETURNING pid`,
		PVSQLVAL,
		p.pid, p.gid, p.fid,
		p.turn, p.name, locVal,
		contVal, p.inhabitants, p.resources, p.parts,
	), true
}

func (p *PlanetView) TableName() string {
	return "planetviews"
}
func (p *PlanetView) SQLID() []string {
	return []string{"gid", strconv.Itoa(p.gid), "fid", strconv.Itoa(p.fid), "pid", strconv.Itoa(p.pid)}
}
