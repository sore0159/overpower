package json

import (
	"mule/overpower"
)

type Game struct {
	Turn       int     `json:"turn"`
	Gid        int     `json:"gid"`
	Name       string  `json:"name"`
	Owner      string  `json:"owner"`
	HasPW      bool    `json:"haspw"`
	AutoDays   [7]bool `json:"autodays"`
	WinPercent int     `json:"winpercent"`
	HighScore  int     `json:"highscore"`
	FreeAutos  int     `json:"freeautos"`
}

func LoadGame(g overpower.Game) *Game {
	return &Game{
		Turn:       g.Turn(),
		Gid:        g.Gid(),
		Name:       g.Name(),
		Owner:      g.Owner(),
		HasPW:      g.HasPW(),
		AutoDays:   g.AutoDays(),
		FreeAutos:  g.FreeAutos(),
		WinPercent: g.WinPercent(),
		HighScore:  g.HighScore(),
	}
}

func LoadGames(games []overpower.Game) []*Game {
	gs := make([]*Game, len(games))
	for i, g := range games {
		gs[i] = LoadGame(g)
	}
	return gs
}
