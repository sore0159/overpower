package db

import (
	"mule/hexagon"
	"mule/mydb"
)

type PlanetView struct {
	*mydb.SQLHandler
	gid         int
	pid         int
	fid         int
	turn        int
	name        string
	loc         hexagon.Coord
	controller  int
	inhabitants int
	resources   int
	parts       int
}

func NewPlanetView() *PlanetView {
	return &PlanetView{
		SQLHandler: mydb.NewSQLHandler(),
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
	return p.controller
}
func (p *PlanetView) SetController(x int) {
	if p.controller == x {
		return
	}
	p.controller = x
	p.SetInt("controller", x)
}
func (p *PlanetView) Inhabitants() int {
	return p.inhabitants
}
func (p *PlanetView) SetInhabitants(x int) {
	if p.inhabitants == x {
		return
	}
	p.inhabitants = x
	p.SetInt("inhabitants", x)
}
func (p *PlanetView) Resources() int {
	return p.resources
}
func (p *PlanetView) SetResources(x int) {
	if p.resources == x {
		return
	}
	p.resources = x
	p.SetInt("resources", x)
}
func (p *PlanetView) Parts() int {
	return p.parts
}
func (p *PlanetView) SetParts(x int) {
	if p.parts == x {
		return
	}
	p.parts = x
	p.SetInt("parts", x)
}
func (p *PlanetView) Turn() int {
	return p.turn
}
func (p *PlanetView) SetTurn(x int) {
	if p.turn == x {
		return
	}
	p.turn = x
	p.SetInt("turn", x)
}
