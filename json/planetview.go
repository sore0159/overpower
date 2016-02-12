package json

import (
	"mule/hexagon"
	"mule/overpower"
)

func LoadPlanetView(item overpower.PlanetView) *PlanetView {
	return &PlanetView{
		Gid:         item.Gid(),
		Fid:         item.Fid(),
		Loc:         item.Loc(),
		Name:        item.Name(),
		Turn:        item.Turn(),
		Controller:  item.Controller(),
		Inhabitants: item.Inhabitants(),
		Resources:   item.Resources(),
		Parts:       item.Parts(),
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
	Gid         int           `json:"gid"`
	Fid         int           `json:"fid"`
	Loc         hexagon.Coord `json:"loc"`
	Name        string        `json:"name"`
	Turn        int           `json:"turn"`
	Controller  int           `json:"controller"`
	Inhabitants int           `json:"inhabitants"`
	Resources   int           `json:"resources"`
	Parts       int           `json:"parts"`
}
