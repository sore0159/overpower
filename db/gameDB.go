package db

import (
	"errors"
	"mule/mydb"
	"mule/overpower"
)

type GameGroup struct {
	List []*Game
}

func NewGameGroup() *GameGroup {
	return &GameGroup{
		List: []*Game{},
	}
}

func (group *GameGroup) New() mydb.SQLer {
	item := NewGame()
	group.List = append(group.List, item)
	return item
}

func (group *GameGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.modified {
			list = append(list, item)
		}
	}
	return list
}

func (group *GameGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertGames2DB(list ...overpower.Game) ([]*Game, error) {
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

func convertGames2OP(list ...*Game) []overpower.Game {
	converted := make([]overpower.Game, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}

/*
func (item *Game) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	}
	return nil
}

func (item *Game) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	}
	return nil
}

func (item *Game) SQLTable() string {
	return "games"
}

func (group *GameGroup) SQLTable() string {
	return "games"
}

func (group *GameGroup) SelectCols() []string {
	return []string{
		//
	}
}

func (group *GameGroup) UpdateCols() []string {
	return []string{
		//
	}
}

func (group *GameGroup) PKCols() []string {
	return []string{
		//
	}
}

func (group *GameGroup) InsertCols() []string {
	return []string{
		//
	}
}

func (group *GameGroup) InsertScanCols() []string {
	return nil
}
*/
