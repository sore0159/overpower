package db

import (
	"database/sql"
	"fmt"
	"mule/mydb"
	"strconv"
)

const GAMESQLVAL = `gid, owner, name, turn, password`

func (g *Game) RowScan(row mydb.Scanner) error {
	return row.Scan(&(g.gid), &(g.owner), &(g.name), &(g.turn), &(g.password))
}

func (g *Game) Insert(db *sql.DB) (ok bool) {
	return mydb.Insert(db, g)
}

func (g *Game) UpdateQ() (query string) {
	return mydb.ModderQ(g)
}

func (g *Game) Commit() {
	g.turnMod = false
}

func (g *Game) InsertScan(row *sql.Row) error {
	return row.Scan(&(g.gid))
}
func (g *Game) InsertQ() (query string, scan bool) {
	if g.gid != 0 {
		return fmt.Sprintf("INSERT INTO games (%s) VALUES(%d, '%s', '%s', %d, '%s') RETURNING gid", GAMESQLVAL, g.gid, g.owner, g.name, g.turn, g.password), true
	} else {
		return fmt.Sprintf("INSERT INTO games (%s) VALUES('%s', '%s', %d, '%s') RETURNING gid", GAMESQLVAL[5:], g.owner, g.name, g.turn, g.password), true
	}
}

func (g *Game) GetMods() map[string]string {
	if g.turnMod {
		return map[string]string{"turn": strconv.Itoa(g.turn)}
	}
	return nil
}
func (g *Game) TableName() string {
	return "games"
}
func (g *Game) SQLID() []string {
	return []string{"gid", strconv.Itoa(g.gid)}
}
