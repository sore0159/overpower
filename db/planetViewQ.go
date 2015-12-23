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

func (p *PlanetView) Insert(db mydb.SQLer) (ok bool) {
	return mydb.Insert(db, p)
}

func (p *PlanetView) UpdateQ() (query string) {
	// return mydb.ModderQ(p)
	return fmt.Sprintf(`UPDATE planetviews SET
		turn = %d, controller = %d, 
		inhabitants = %d, resources = %d, parts = %d
	WHERE gid = %d AND fid = %d AND pid = %d`,
		p.turn, p.controller,
		p.inhabitants, p.resources, p.parts,
		p.gid, p.fid, p.pid,
	)
}

func (p *PlanetView) InsertScan(row *sql.Row) error {
	return nil
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
	)`,
		PVSQLVAL,
		p.pid, p.gid, p.fid,
		p.turn, p.name, locVal,
		contVal, p.inhabitants, p.resources, p.parts,
	), false
}

func (p *PlanetView) TableName() string {
	return "planetviews"
}
func (p *PlanetView) SQLID() []string {
	return []string{"gid", strconv.Itoa(p.gid), "fid", strconv.Itoa(p.fid), "pid", strconv.Itoa(p.pid)}
}
