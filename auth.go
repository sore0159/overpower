package main

import (
	"mule/planetattack/attack"
)

func Admin(username string) bool {
	return username == "mule"
}

func (g *Game) AuthFaction(uName, fName string) (*attack.Faction, bool) {
	fID, ok := g.Users[uName]
	if !ok {
		return nil, false
	}
	f, ok := g.Factions[fID]
	if !ok {
		return nil, false
	}
	if f.Name != fName {
		return nil, false
	}
	return f, true
}
