package json

import (
	"mule/hexagon"
	"mule/overpower"
)

func LoadPlanetView(item overpower.PlanetView) *PlanetView {
	return &PlanetView{
		Gid:  item.Gid(),
		Fid:  item.Fid(),
		Loc:  item.Loc(),
		Name: item.Name(),
		Turn: item.Turn(),

		PrimaryFaction:  item.PrimaryFaction(),
		PrimaryPresence: item.PrimaryPresence(),
		PrimaryPower:    item.PrimaryPower(),

		SecondaryFaction:  item.SecondaryFaction(),
		SecondaryPresence: item.SecondaryPresence(),
		SecondaryPower:    item.SecondaryPower(),

		Antimatter: item.Antimatter(),
		Tachyons:   item.Tachyons(),
	}
}

func LoadPlanetViews(list []overpower.PlanetView) []*PlanetView {
	jList := make([]*PlanetView, len(list))
	for i, item := range list {
		jList[i] = LoadPlanetView(item)
	}
	return jList
}

type PlanetView struct {
	Gid  int           `json:"gid"`
	Fid  int           `json:"fid"`
	Loc  hexagon.Coord `json:"loc"`
	Name string        `json:"name"`
	Turn int           `json:"turn"`

	PrimaryFaction  int `json:"primaryfaction"`
	PrimaryPresence int `json:"primarypresence"`
	PrimaryPower    int `json:"primarypower"`

	SecondaryFaction  int `json:"secondaryfaction"`
	SecondaryPresence int `json:"secondarypresence"`
	SecondaryPower    int `json:"secondarypower"`

	Antimatter int `json:"antimatter"`
	Tachyons   int `json:"tachyons"`
}
