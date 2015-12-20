package db

import (
	"mule/hexagon"
	"mule/mydb"
	"mule/overpower"
)

var plIntfTest overpower.Planet = NewPlanet()

type Planet struct {
	*mydb.SQLHandler
	gid         int
	pid         int
	name        string
	loc         hexagon.Coord
	controller  int
	inhabitants int
	resources   int
	parts       int
}

func NewPlanet() *Planet {
	return &Planet{
		SQLHandler: mydb.NewSQLHandler(),
	}
}

func (p *Planet) Gid() int {
	return p.gid
}
func (p *Planet) Pid() int {
	return p.pid
}
func (p *Planet) Name() string {
	return p.name
}
func (p *Planet) Loc() hexagon.Coord {
	return p.loc
}
func (p *Planet) Controller() int {
	return p.controller
}
func (p *Planet) SetController(x int) {
	if p.controller == x {
		return
	}
	p.controller = x
	p.SetInt("controller", x)
}
func (p *Planet) Inhabitants() int {
	return p.inhabitants
}
func (p *Planet) SetInhabitants(x int) {
	if p.inhabitants == x {
		return
	}
	p.inhabitants = x
	p.SetInt("inhabitants", x)
}
func (p *Planet) Resources() int {
	return p.resources
}
func (p *Planet) SetResources(x int) {
	if p.resources == x {
		return
	}
	p.resources = x
	p.SetInt("resources", x)
}
func (p *Planet) Parts() int {
	return p.parts
}
func (p *Planet) SetParts(x int) {
	if p.parts == x {
		return
	}
	p.parts = x
	p.SetInt("parts", x)
}
