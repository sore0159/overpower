package db

import (
	"database/sql"
	"mule/hexagon"
)

type MiniPlanet struct {
	modified bool
	//
	gid  int
	fid  int
	loc  hexagon.Coord
	name string

	turn int

	primaryfaction    sql.NullInt64
	primarypresence   int
	primarypower      int
	secondaryfaction  sql.NullInt64
	secondarypresence int
	secondarypower    int

	antimatter int
	tachyons   int
}

func NewMiniPlanet() *MiniPlanet {
	return &MiniPlanet{
	//
	}
}

func (p *MiniPlanet) Gid() int {
	return p.gid
}
func (p *MiniPlanet) Fid() int {
	return p.fid
}
func (p *MiniPlanet) Loc() hexagon.Coord {
	return p.loc
}
func (p *MiniPlanet) Name() string {
	return p.name
}

func (p *MiniPlanet) Turn() int {
	return p.turn
}

func (p *MiniPlanet) SetTurn(x int) {
	if p.turn == x {
		return
	}
	p.turn = x
	p.modified = true
}

func (p *MiniPlanet) PrimaryPower() int {
	return p.primarypower
}

func (p *MiniPlanet) SetPrimaryPower(x int) {
	if p.primarypower == x {
		return
	}
	p.primarypower = x
	p.modified = true
}

func (p *MiniPlanet) SecondaryPower() int {
	return p.secondarypower
}

func (p *MiniPlanet) SetSecondaryPower(x int) {
	if p.secondarypower == x {
		return
	}
	p.secondarypower = x
	p.modified = true
}

func (p *MiniPlanet) Antimatter() int {
	return p.antimatter
}
func (p *MiniPlanet) SetAntimatter(x int) {
	if p.antimatter == x {
		return
	}
	p.antimatter = x
	p.modified = true
}

func (p *MiniPlanet) Tachyons() int {
	return p.tachyons
}
func (p *MiniPlanet) SetTachyons(x int) {
	if p.tachyons == x {
		return
	}
	p.tachyons = x
	p.modified = true
}

func (p *MiniPlanet) PrimaryPresence() int {
	return p.primarypresence
}
func (p *MiniPlanet) SetPrimaryPresence(x int) {
	if p.primarypresence == x {
		return
	}
	p.primarypresence = x
	p.modified = true
}

func (p *MiniPlanet) SecondaryPresence() int {
	return p.secondarypresence
}
func (p *MiniPlanet) SetSecondaryPresence(x int) {
	if p.secondarypresence == x {
		return
	}
	p.secondarypresence = x
	p.modified = true
}

func (p *MiniPlanet) PrimaryFaction() int {
	if p.primaryfaction.Valid {
		return int(p.primaryfaction.Int64)
	}
	return 0
}
func (p *MiniPlanet) SetPrimaryFaction(x int) {
	if x == 0 {
		if !p.primaryfaction.Valid {
			return
		}
		p.primaryfaction.Int64 = 0
		p.primaryfaction.Valid = false
		p.modified = true
	} else {
		if p.primaryfaction.Valid && int(p.primaryfaction.Int64) == x {
			return
		}
		p.primaryfaction.Int64 = int64(x)
		p.primaryfaction.Valid = true
		p.modified = true
	}
}

func (p *MiniPlanet) SecondaryFaction() int {
	if p.secondaryfaction.Valid {
		return int(p.secondaryfaction.Int64)
	}
	return 0
}
func (p *MiniPlanet) SetSecondaryFaction(x int) {
	if x == 0 {
		if !p.secondaryfaction.Valid {
			return
		}
		p.secondaryfaction.Int64 = 0
		p.secondaryfaction.Valid = false
		p.modified = true
	} else {
		if p.secondaryfaction.Valid && int(p.secondaryfaction.Int64) == x {
			return
		}
		p.secondaryfaction.Int64 = int64(x)
		p.secondaryfaction.Valid = true
		p.modified = true
	}
}

func (p *MiniPlanet) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return p.gid
	case "fid":
		return p.fid
	case "locx":
		return p.loc[0]
	case "locy":
		return p.loc[1]
	case "name":
		return p.name
	case "turn":
		return p.turn
	case "primaryfaction":
		return p.primaryfaction
	case "primarypresence":
		return p.primarypresence
	case "primarypower":
		return p.primarypower
	case "secondaryfaction":
		return p.secondaryfaction
	case "secondarypower":
		return p.secondarypower
	case "secondarypresence":
		return p.secondarypresence
	case "antimatter":
		return p.antimatter
	case "tachyons":
		return p.tachyons
	}
	return nil
}
func (p *MiniPlanet) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &p.gid
	case "locx":
		return &p.loc[0]
	case "locy":
		return &p.loc[1]
	case "name":
		return &p.name
	case "turn":
		return &p.turn
	case "primaryfaction":
		return &p.primaryfaction
	case "primarypresence":
		return &p.primarypresence
	case "primarypower":
		return &p.primarypower
	case "secondaryfaction":
		return &p.secondaryfaction
	case "secondarypresence":
		return &p.secondarypresence
	case "secondarypower":
		return &p.secondarypower
	case "antimatter":
		return &p.antimatter
	case "tachyons":
		return &p.tachyons
	}
	return nil
}
