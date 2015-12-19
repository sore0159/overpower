package db

import (
	"database/sql"
	"fmt"
	"mule/mydb"
	"mule/overpower"
	"strconv"
)

var intfTest overpower.Faction = NewFaction()

const FACSQLVAL = `fid, gid, owner, name, done`

func (f *Faction) RowScan(row mydb.Scanner) error {
	return row.Scan(&(f.fid), &(f.gid), &(f.owner), &(f.name), &(f.done))
}

func (f *Faction) Insert(db *sql.DB) (ok bool) {
	return mydb.Insert(db, f)
}

func (f *Faction) UpdateQ() (query string) {
	return mydb.ModderQ(f)
}

func (f *Faction) Commit() {
	f.doneMod = false
}

func (f *Faction) InsertScan(row *sql.Row) error {
	return row.Scan(&(f.fid))
}
func (f *Faction) InsertQ() (query string, scan bool) {
	return fmt.Sprintf("INSERT INTO factions (%s) VALUES(%d, '%s', '%s', %t) RETURNING fid", FACSQLVAL[5:], f.gid, f.owner, f.name, f.done), true
}

func (f *Faction) GetMods() map[string]string {
	if f.doneMod {
		return map[string]string{"done": fmt.Sprintf("%t", f.done)}
	}
	return nil
}
func (f *Faction) TableName() string {
	return "factions"
}
func (f *Faction) SQLID() []string {
	return []string{"gid", strconv.Itoa(f.gid), "fid", strconv.Itoa(f.fid)}
}
