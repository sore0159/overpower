package db

import (
	"database/sql"
	"fmt"
	"mule/mydb"
	"strconv"
)

var gametest mydb.Updater = &Game{}

const GAMESQLVAL = `gid, owner, name, turn, password`

func (g *Game) RowScan(row mydb.Scanner) error {
	return row.Scan(&(g.gid), &(g.owner), &(g.name), &(g.turn), &(g.password))
}

func (g *Game) Insert() (ok bool) {
	return mydb.Insert(g.db, g)
}

func GetGame(db *sql.DB, gid int) (g *Game, ok bool) {
	g = NewGame()
	g.db, g.gid = db, gid
	query := fmt.Sprintf("SELECT %s FROM games WHERE gid = %d", GAMESQLVAL, gid)
	return g, mydb.GetOne(db, query, g)
}

func DropGame(db *sql.DB, gid int) (ok bool) {
	query := fmt.Sprintf("DELETE FROM games WHERE gid = %d", gid)
	return mydb.Exec(db, query)
}

func AllGames(db *sql.DB) (games []*Game, ok bool) {
	query := fmt.Sprintf("SELECT %s FROM games", GAMESQLVAL)
	games = []*Game{}
	maker := func() mydb.Rower {
		g := NewGame()
		g.db = db
		return g
	}
	return games, mydb.Get(db, query, &games, maker)
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
