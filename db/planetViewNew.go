package db

import (
	"database/sql"
	"mule/hexagon"
)

type PlanetView struct {
	modified bool
	//
	gid         int
	pid         int
	fid         int
	turn        int
	name        string
	loc         hexagon.Coord
	controller  sql.NullInt64
	inhabitants sql.NullInt64
	resources   sql.NullInt64
	parts       sql.NullInt64
}

func NewPlanetView() *PlanetView {
	return &PlanetView{
	//
	}
}

func (p *PlanetView) Gid() int {
	return p.gid
}
func (p *PlanetView) Pid() int {
	return p.pid
}
func (p *PlanetView) Fid() int {
	return p.fid
}
func (p *PlanetView) Name() string {
	return p.name
}
func (p *PlanetView) Loc() hexagon.Coord {
	return p.loc
}
func (p *PlanetView) Controller() int {
	if p.controller.Valid {
		return int(p.controller.Int64)
	} else {
		return 0
	}
}
func (p *PlanetView) SetController(x int) {
	if x == 0 {
		if !p.controller.Valid {
			return
		}
		p.controller.Int64 = 0
		p.controller.Valid = false
		p.modified = true
	} else {
		x64 := int64(x)
		if p.controller.Valid && p.controller.Int64 == x64 {
			return
		}
		p.controller.Valid = true
		p.controller.Int64 = x64
		p.modified = true
	}
}

func (p *PlanetView) Inhabitants() int {
	if p.inhabitants.Valid {
		return int(p.inhabitants.Int64)
	}
	return 0
}
func (p *PlanetView) SetInhabitants(x int) {
	x64 := int64(x)
	if p.inhabitants.Valid && p.inhabitants.Int64 == x64 {
		return
	}
	p.inhabitants.Valid = true
	p.inhabitants.Int64 = x64
	p.modified = true
}
func (p *PlanetView) Resources() int {
	if p.resources.Valid {
		return int(p.resources.Int64)
	}
	return 0
}
func (p *PlanetView) SetResources(x int) {
	x64 := int64(x)
	if p.resources.Valid && p.resources.Int64 == x64 {
		return
	}
	p.resources.Valid = true
	p.resources.Int64 = x64
	p.modified = true
}
func (p *PlanetView) Parts() int {
	if p.parts.Valid {
		return int(p.parts.Int64)
	}
	return 0
}
func (p *PlanetView) SetParts(x int) {
	x64 := int64(x)
	if p.parts.Valid && p.parts.Int64 == x64 {
		return
	}
	p.parts.Valid = true
	p.parts.Int64 = x64
	p.modified = true
}
func (p *PlanetView) Turn() int {
	return p.turn
}
func (p *PlanetView) SetTurn(x int) {
	if p.turn == x {
		return
	}
	p.turn = x
	p.modified = true
}
func (item *PlanetView) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return item.gid
	case "pid":
		return item.pid
	case "fid":
		return item.fid
	case "turn":
		return item.turn
	case "name":
		return item.name
	case "loc":
		return item.loc
	case "controller":
		return item.controller
	case "inhabitants":
		return item.inhabitants
	case "resources":
		return item.resources
	case "parts":
		return item.parts
	}
	return nil
}

func (item *PlanetView) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &item.gid
	case "pid":
		return &item.pid
	case "fid":
		return &item.fid
	case "turn":
		return &item.turn
	case "name":
		return &item.name
	case "loc":
		return &item.loc
	case "controller":
		return &item.controller
	case "inhabitants":
		return &item.inhabitants
	case "resources":
		return &item.resources
	case "parts":
		return &item.parts
	}
	return nil
}

func (item *PlanetView) SQLTable() string {
	return "planetviews"
}

func (group *PlanetViewGroup) SQLTable() string {
	return "planetviews"
}

func (group *PlanetViewGroup) SelectCols() []string {
	return []string{
		"gid",
		"pid",
		"fid",
		"turn",
		"name",
		"loc",
		"controller",
		"inhabitants",
		"resources",
		"parts",
	}
}

func (group *PlanetViewGroup) UpdateCols() []string {
	return []string{
		"turn",
		"controller",
		"inhabitants",
		"resources",
		"parts",
	}
}

func (group *PlanetViewGroup) PKCols() []string {
	return []string{"gid", "fid", "pid"}
}

func (group *PlanetViewGroup) InsertCols() []string {
	return []string{
		"gid",
		"pid",
		"fid",
		"turn",
		"name",
		"loc",
		"controller",
		"inhabitants",
		"resources",
		"parts",
	}
}

func (group *PlanetViewGroup) InsertScanCols() []string {
	return nil
}
