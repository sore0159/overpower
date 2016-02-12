package db

import (
	"errors"
	"fmt"
	"mule/hexagon"
	"mule/mydb"
	"mule/overpower"
	"strings"
)

func (d DB) getPlanetsByLoc(gid int, locs ...hexagon.Coord) ([]*Planet, error) {
	parts := make([]string, 0, len(locs))
	args := make([]interface{}, 0, 2*len(locs)+1)
	args = append(args, gid)
	count := 2
	for _, item := range locs {
		parts = append(parts, fmt.Sprintf("($%d, $%d)", count, count+1))
		args = append(args, item[0], item[1])
		count += 2
	}
	locStr := strings.Join(parts, ",")
	group := NewPlanetGroup()
	query := fmt.Sprintf("SELECT %s FROM planets WHERE gid = $1 and (locx, locy) IN (%s)", strings.Join(group.SelectCols(), ","), locStr)
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
