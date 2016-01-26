package db

import (
	"errors"
	"fmt"
	"mule/mydb"
	"mule/overpower"
	"strings"
)

func (d DB) getPlanetsByPlid(gid int, plids ...int) ([]*Planet, error) {
	parts := make([]string, len(plids))
	args := make([]interface{}, len(plids)+1)
	args[0] = gid
	for i, item := range plids {
		parts[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = item
	}
	pidStr := strings.Join(parts, ",")
	group := NewPlanetGroup()
	query := fmt.Sprintf("SELECT %s FROM planets WHERE gid = $1 and pid IN (%s)", strings.Join(group.SelectCols(), ","), pidStr)
	err := mydb.Get(d.db(), group, query, args...)
	if my, bad := Check(err, "getplanets failure", "query", query, "args", args); bad {
		return nil, my
	}
	return group.List, nil
}

type PlanetGroup struct {
	List []*Planet
}

func NewPlanetGroup() *PlanetGroup {
	return &PlanetGroup{
		List: []*Planet{},
	}
}

func (group *PlanetGroup) New() mydb.SQLer {
	item := NewPlanet()
	group.List = append(group.List, item)
	return item
}

func (group *PlanetGroup) UpdateList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		if item.modified {
			list = append(list, item)
		}
	}
	return list
}

func (group *PlanetGroup) InsertList() []mydb.SQLer {
	list := make([]mydb.SQLer, 0, len(group.List))
	for _, item := range group.List {
		list = append(list, item)
	}
	return list
}

func convertPlanets2DB(list ...overpower.Planet) ([]*Planet, error) {
	mylist := make([]*Planet, 0, len(list))
	for _, test := range list {
		if test == nil {
			continue
		}
		if t, ok := test.(*Planet); ok {
			mylist = append(mylist, t)
		} else {
			return nil, errors.New("bad Planet struct type for op/db")
		}
	}
	return mylist, nil
}

func convertPlanets2OP(list ...*Planet) []overpower.Planet {
	converted := make([]overpower.Planet, len(list))
	for i, item := range list {
		converted[i] = item
	}
	return converted
}

/*
 */
