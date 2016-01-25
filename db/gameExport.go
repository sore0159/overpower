package db

import (
	"mule/overpower"
)

func (d DB) MakeGame(owner, name, password string, winpercent int) (err error) {
	item := &Game{owner: owner, name: name, winpercent: winpercent}
	if password != "" {
		item.password.Valid = true
		item.password.String = password
	}
	group := &GameGroup{[]*Game{item}}
	return d.makeGroup(group)
}

func (d DB) DropGames(conditions []interface{}) error {
	return d.dropItems("games", conditions)
}

func (d DB) UpdateGames(list ...overpower.Game) error {
	mylist, err := convertGames2DB(list...)
	if my, bad := Check(err, "update Games conversion failure"); bad {
		return my
	}
	return d.updateGroup(&GameGroup{mylist})
}

func (d DB) GetGame(conditions []interface{}) (overpower.Game, error) {
	list, err := d.GetGames(conditions)
	if my, bad := Check(err, "get Game failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetGames(conditions []interface{}) ([]overpower.Game, error) {
	group := NewGameGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get Games failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertGames2OP(list...)
	return converted, nil
}
