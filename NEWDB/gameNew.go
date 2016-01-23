package db

import (
	"database/sql"
)

func NewGame() *Game {
	return &Game{}
}

type Game struct {
	gid        int
	owner      string
	name       string
	turn       int
	autoturn   int
	freeautos  int
	password   sql.NullString
	winpercent int
	highscore  int
	//
	modified bool
}

func (g *Game) SetTurn(t int) {
	if g.turn == t {
		return
	}
	g.turn = t
	g.modified = true
}
func (g *Game) IncTurn() {
	g.turn++
	g.modified = true
}

func (g *Game) Turn() int {
	return g.turn
}
func (g *Game) Gid() int {
	return g.gid
}
func (g *Game) Owner() string {
	return g.owner
}
func (g *Game) Name() string {
	return g.name
}
func (g *Game) IsPwd(test string) bool {
	if !g.password.Valid {
		return true
	}
	return g.password.String == test
}
func (g *Game) HasPW() bool {
	return g.password.Valid
}

func (g *Game) AutoTurn() int {
	return g.autoturn
}
func (g *Game) SetAutoTurn(x int) {
	if x == g.autoturn {
		return
	}
	g.autoturn = x
	g.modified = true
}

func (g *Game) WinPercent() int {
	return g.winpercent
}
func (g *Game) HighScore() int {
	return g.highscore
}
func (g *Game) SetHighScore(x int) {
	if x == g.highscore {
		return
	}
	g.highscore = x
	g.modified = true
}

func (g *Game) FreeAutos() int {
	return g.freeautos
}
func (g *Game) SetFreeAutos(x int) {
	if x == g.freeautos {
		return
	}
	g.freeautos = x
	g.modified = true
}

func (g *Game) AutoDays() (days [7]bool) {
	sum := g.autoturn
	for i := 0; i < 7; i++ {
		if sum%2 == 1 {
			days[i] = true
		}
		sum = sum / 2
	}
	return
}

func (g *Game) SetAutoDays(days [7]bool) {
	var sum int
	for i, b := range days {
		if b {
			sum += 1 << uint32(i)
		}
	}
	g.SetAutoTurn(sum)
}

func (g *Game) SQLVal(name string) interface{} {
	switch name {
	case "gid":
		return g.gid
	case "owner":
		return g.owner
	case "name":
		return g.name
	case "turn":
		return g.turn
	case "autoturn":
		return g.autoturn
	case "freeautos":
		return g.freeautos
	case "password":
		return g.password
	case "winpercent":
		return g.winpercent
	case "highscore":
		return g.highscore
	}
	return nil
}

func (g *Game) SQLPtr(name string) interface{} {
	switch name {
	case "gid":
		return &g.gid
	case "owner":
		return &g.owner
	case "name":
		return &g.name
	case "turn":
		return &g.turn
	case "autoturn":
		return &g.autoturn
	case "freeautos":
		return &g.freeautos
	case "password":
		return &g.password
	case "winpercent":
		return &g.winpercent
	case "highscore":
		return &g.highscore
	}
	return nil
}

func (g *Game) SQLTable() string {
	return "games"
}

/*
type GameGroup struct {
	Games []*Game
}

func (gg *GameGroup) Current() []mydb.SQLer {
	list := gg.Games
	if list == nil {
		return nil
	}
	sqlers := make([]mydb.SQLer, len(list))
	for i, item := range list {
		sqlers[i] = item
	}
	return sqlers
}

func (gg *GameGroup) New() mydb.SQLer {
	g := NewGame()
	gg.Games = append(gg.Games, g)
	return g
}

func (gg *GameGroup) SQLTable() string {
	return "games"
}
func (gg *GameGroup) SQLPK() []string {
	return []string{"gid"}
}
func (gg *GameGroup) SQLSerials() []string {
	return []string{"gid"}
}
func (gg *GameGroup) SQLCols() []string {
	return []string{
		"gid",
		"owner",
		"name",
		"turn",
		"autoturn",
		"freeautos",
		"password",
		"winpercent",
		"highscore",
	}
}
*/

// query := "UPDATE games SET turn = $1, autoturn = $2, freeautos = $3 WHERE gid = $4"
// args :=  []string{"turn", "autoturn", "freeautos", "gid"}
