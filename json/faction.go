package json

import (
	"mule/overpower"
)

func LoadFaction(item overpower.Faction, fid int) *Faction {
	var score int
	if fid == item.Fid() {
		score = item.Score()
	}
	return &Faction{
		Gid:   item.Gid(),
		Fid:   item.Fid(),
		Owner: item.Owner(),
		Name:  item.Name(),
		Done:  item.Done(),
		Score: score,
	}
}

func LoadFactions(list []overpower.Faction, fid int) []*Faction {
	jList := make([]*Faction, len(list))
	for i, item := range list {
		jList[i] = LoadFaction(item, fid)
	}
	return jList
}

type Faction struct {
	Gid   int    `json:"gid"`
	Fid   int    `json:"fid"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
	Done  bool   `json:"done"`
	Score int    `json:"score,omitempty"`
}
