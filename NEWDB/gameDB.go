package db

import (
	"database/sql"
	"errors"
	"mule/mydb"
	"mule/overpower"
)

func (d DB) updateGames(tx *sql.Tx, games ...*Game) error {
	if len(games) == 0 {
		return nil
	}
	table := "games"
	cols := []string{
		"turn",
		"autoturn",
		"freeautos",
		"highscore",
	}
	condCols := []string{"gid"}
	query := mydb.UpdateQ(table, cols, condCols)
	allCols := append(cols, condCols...)
	sqlers := make([]mydb.SQLer, len(games))
	for i, g := range games {
		sqlers[i] = g
	}
	if tx == nil {
		return mydb.Update(d.db, true, query, allCols, sqlers...)
	} else {
		return mydb.Update(tx, true, query, allCols, sqlers...)
	}
}

func (d DB) mustExec(query string, args ...interface{}) error {
	return mydb.ExecCheck(d.db.Exec(query, args...))
}

func (d DB) getGames(query string, args ...interface{}) ([]*Game, error) {
	gg := NewGameGroup()
	err := mydb.Get(d.db, gg, query, args...)
	return gg.Games, err
}

type GameGroup struct {
	Games []*Game
}

func NewGameGroup() *GameGroup {
	return &GameGroup{
		Games: []*Game{},
	}
}

func (gg *GameGroup) New() mydb.SQLer {
	g := NewGame()
	gg.Games = append(gg.Games, g)
	return g
}

func (d DB) convertGames2DB(list ...overpower.Game) ([]*Game, error) {
	mylist := make([]*Game, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*Game); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad Game struct type for op/db")
		}
	}
	return mylist, nil
}

func (d DB) convertGames2OP(list ...*Game) []overpower.Game {
	converted := make([]overpower.Game, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}
