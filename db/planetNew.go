package db

import (
	"database/sql"
	"mule/hexagon"
)

type Planet struct {
	modified bool
	//
	gid         int
	loc         hexagon.Coord
	name        string
	controller  sql.NullInt64
	inhabitants int
	resources   int
	parts       int
}

func NewPlanet() *Planet {
	return &Planet{
	//
	}
}

func (p *Planet) Gid() int {
	return p.gid
}
func (p *Planet) Name() string {
	return p.name
}
func (p *Planet) Loc() hexagon.Coord {
	return p.loc
}
func (p *Planet) Controller() int {
	if p.controller.Valid {
		return int(p.controller.Int64)
	}
	return 0
}
func (p *Planet) SetController(x int) {
	if x == 0 {
		if !p.controller.Valid {
			return
		}
		p.controller.Int64 = 0
		p.controller.Valid = false
		p.modified = true
	} else {
		if p.controller.Valid && int(p.controller.Int64) == x {
			return
		}
		p.controller.Int64 = int64(x)
		p.controller.Valid = true
		p.modified = true
	}
}

func (p *Planet) Inhabitants() int {
	return p.inhabitants
}
func (p *Planet) SetInhabitants(x int) {
	if p.inhabitants == x {
		return
	}
	p.inhabitants = x
	p.modified = true
}
func (p *Planet) Resources() int {
	return p.resources
}
func (p *Planet) SetResources(x int) {
	if p.resources == x {
		return
	}
	p.resources = x
	p.modified = true
}
func (p *Planet) Parts() int {
	return p.parts
}
func (p *Planet) SetParts(x int) {
	if p.parts == x {
		return
	}
	p.parts = x
	p.modified = true
}
func (p *Planet) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return p.gid
	case "name":
		return p.name
	case "locx":
		return p.loc[0]
	case "locy":
		return p.loc[1]
	case "controller":
		return p.controller
	case "inhabitants":
		return p.inhabitants
	case "resources":
		return p.resources
	case "parts":
		return p.parts
	}
	return nil
}
func (p *Planet) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &p.gid
	case "locx":
		return &p.loc[0]
	case "locy":
		return &p.loc[1]
	case "name":
		return &p.name
	case "loc":
		return &p.loc
	case "controller":
		return &p.controller
	case "inhabitants":
		return &p.inhabitants
	case "resources":
		return &p.resources
	case "parts":
		return &p.parts
	}
	return nil
}

func (p *Planet) SQLTable() string {
	return "planets"
}

func (group *PlanetGroup) SQLTable() string {
	return "planets"
}

func (group *PlanetGroup) SelectCols() []string {
	return []string{
		"gid",
		"locx",
		"locy",
		"name",
		"controller",
		"inhabitants",
		"resources",
		"parts",
	}
}

func (group *PlanetGroup) UpdateCols() []string {
	return []string{
		"controller",
		"inhabitants",
		"resources",
		"parts",
	}
}

func (group *PlanetGroup) PKCols() []string {
	return []string{"gid", "locx", "locy"}
}

func (group *PlanetGroup) InsertCols() []string {
	return []string{
		"gid",
		"locx",
		"locy",
		"name",
		"controller",
		"inhabitants",
		"resources",
		"parts",
	}
}

func (group *PlanetGroup) InsertScanCols() []string {
	return nil
}
