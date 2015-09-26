package main

import (
	"mule/planetattack/attack"
)

type Game struct {
	*attack.Game
	Users  map[string]int
	Signup map[string]string
}

func NewGame() *Game {
	return &Game{
		//		Game:  attack.NewGame(),
		Users:  map[string]int{},
		Signup: map[string]string{},
	}
}

func (g *Game) MakeGame() {
	names := g.Signup
	r := make([]string, len(names))
	var i int
	for facName, _ := range names {
		r[i] = facName
		i++
	}
	g.Game = attack.MakeGame(r)
	m := make(map[string]int, len(names))
	for fID, fac := range g.Factions {
		for facName, userName := range names {
			if facName == fac.Name {
				m[userName] = fID
				break
			}
		}
	}
	g.Users = m
	g.Signup = nil
}

func (g *Game) Save() {
	Save(g, SAVEFL)
}
