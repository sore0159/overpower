package main

import (
	"mule/planetattack/attack"
)

type Game struct {
	*attack.Game
}

func NewGame() *Game {
	return &Game{
		Game: attack.NewGame(),
	}
}

func MakeGame(names []string) *Game {
	return &Game{attack.MakeGame(names)}
}
