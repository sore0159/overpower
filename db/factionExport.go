package db

import (
	"mule/overpower"
)

func (d DB) MakeFaction(gid int, owner, name string) (err error) {
	item := &Faction{gid: gid, owner: owner, name: name}
	group := &FactionGroup{[]*Faction{item}}
	return d.makeGroup(group)
}

func (d DB) DropFactions(conditions ...interface{}) error {
	return d.dropItems("factions", conditions)
}

func (d DB) UpdateFactions(list ...overpower.Faction) error {
	mylist, err := convertFactions2DB(list...)
	if my, bad := Check(err, "update Factions conversion failure"); bad {
		return my
	}
	return d.updateGroup(&FactionGroup{mylist})
}

func (d DB) GetFaction(conditions ...interface{}) (overpower.Faction, error) {
	list, err := d.GetFactions(conditions...)
	if my, bad := Check(err, "get Faction failure"); bad {
		return nil, my
	}
	if err = IsUnique(len(list)); err != nil {
		return nil, err
	}
	return list[0], nil
}

// Example conditions:  C{"gid",1} C{"owner","mule"}, nil
func (d DB) GetFactions(conditions ...interface{}) ([]overpower.Faction, error) {
	group := NewFactionGroup()
	err := d.getGroup(group, conditions)
	if my, bad := Check(err, "get Factions failure", "conditions", conditions); bad {
		return nil, my
	}
	list := group.List
	converted := convertFactions2OP(list...)
	return converted, nil
}
