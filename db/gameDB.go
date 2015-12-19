package db

import (
	"fmt"
	"mule/mydb"
	"mule/overpower"
)

func (d DB) MakeGame(owner, name, password string) (ok bool) {
	g := NewGame()
	g.owner, g.name, g.password = owner, name, password
	return g.Insert(d.db)
}

func (d DB) GetGame(gid int) (g overpower.Game, ok bool) {
	db := d.db
	game := NewGame()
	query := fmt.Sprintf("SELECT %s FROM games WHERE gid = %d", GAMESQLVAL, gid)
	return game, mydb.GetOneIf(db, query, game)
}

func (d DB) GetGameFor(owner string) (g overpower.Game, ok bool) {
	db := d.db
	game := NewGame()
	query := fmt.Sprintf("SELECT %s FROM games WHERE owner = '%s'", GAMESQLVAL, owner)
	return game, mydb.GetOneIf(db, query, game)
}

func (d DB) DropGame(gid int) (ok bool) {
	db := d.db
	query := fmt.Sprintf("DELETE FROM games WHERE gid = %d", gid)
	return mydb.Exec(db, query)
}

func (d DB) AllGames() (games []overpower.Game, ok bool) {
	db := d.db
	query := fmt.Sprintf("SELECT %s FROM games", GAMESQLVAL)
	gamesL := []*Game{}
	maker := func() mydb.Rower {
		g := NewGame()
		return g
	}
	if !mydb.Get(db, query, &gamesL, maker) {
		return nil, false
	}
	games = make([]overpower.Game, len(gamesL))
	for i, g := range gamesL {
		games[i] = g
	}
	return games, true
}
