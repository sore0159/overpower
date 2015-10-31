package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"mule/planetattack"
)

var (
	USERDB   *sql.DB
	ATTACKDB *sql.DB
)

func LoadUserData() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASS, USERDB_NAME))
	if err != nil {
		return nil, Log(err)
	}
	err = db.Ping()
	if err != nil {
		return nil, Log(err)
	}
	return db, nil
}

func LoadPlanetDB() (*sql.DB, error) {
	return planetattack.LoadDB()
}

func NewGame() *Game {
	return &Game{
		Game: *planetattack.NewGame(ATTACKDB),
	}
}

func GetGame(gid int) (*Game, bool) {
	g := NewGame()
	g.Gid = gid
	return g, g.Select()
}

func AllFactions(userN string) []*planetattack.Faction {
	return planetattack.AllFactions(ATTACKDB, userN)
}

func GetGames(gids []int) []*planetattack.Game {
	return planetattack.GetGames(ATTACKDB, gids)
}
