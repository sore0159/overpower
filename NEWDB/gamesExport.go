package db

import (
	"mule/mydb"
	"mule/overpower"
)

func (d DB) UpdateGames(games ...overpower.Game) error {
	mygames, err := d.convertGames2DB(games...)
	if my, bad := Check(err, "update games conversion failure"); bad {
		return my
	}
	return d.updateGames(nil, mygames...)
}

func (d DB) DropGames(conditions []interface{}) error {
	table := "games"
	query, args, err := mydb.DeleteQA(table, conditions)
	if my, bad := Check(err, "dropgames failure", "table", table, "conditions", conditions); bad {
		return my
	}

	return d.mustExec(query, args...)
}

func (d DB) GetGame(conditions []interface{}) (overpower.Game, error) {
	list, err := d.GetGames(conditions)
	if my, bad := Check(err, "getgame failure"); bad {
		return nil, my
	}
	switch len(list) {
	case 0:
		return nil, ErrNoneFound
	case 1:
		return list[0], nil
	default:
		return nil, ErrNotUnique
	}
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetGames(conditions []interface{}) ([]overpower.Game, error) {
	table := "games"
	cols := []string{
		"gid",
		"owner",
		"name",
		"turn",
		"autoturn",
		"freeautos",
		"password",
		"winpercent",
		"highscore",
	}
	query, args, err := mydb.SelectQA(table, cols, conditions)
	if my, bad := Check(err, "getgames failure", "table", table, "conditions", conditions, "cols", cols); bad {
		return nil, my
	}
	list, err := d.getGames(query, args...)
	if my, bad := Check(err, "getgames failure", "query", query, "args", args); bad {
		return nil, my
	}
	converted := d.convertGames2OP(list...)
	return converted, nil
}

func (d DB) MakeGame(owner, name, password string, winpercent int) (err error) {
	table := "games"
	var cols []string
	var args []interface{}
	if password == "" {
		cols = []string{"owner", "name", "winpercent"}
		args = []interface{}{owner, name, winpercent}
	} else {
		cols = []string{"owner", "name", "password", "winpercent"}
		args = []interface{}{owner, name, password, winpercent}
	}
	query := mydb.InsertQ(table, cols, nil)
	err = d.mustExec(query, args...)
	if my, bad := Check(err, "makegame failure", "query", query, "args", args); bad {
		return my
	}
	return nil
}
