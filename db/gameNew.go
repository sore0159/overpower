package db

func NewGame() *Game {
	return &Game{}
}

type Game struct {
	turnMod bool
	//
	gid      int
	owner    string
	name     string
	turn     int
	password string
}

func (g *Game) SetTurn(t int) {
	if g.turn == t {
		return
	}
	g.turn = t
	g.turnMod = true
}
func (g *Game) IncTurn() {
	g.turn++
	g.turnMod = true
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
	return g.password == test
}
func (g *Game) HasPW() bool {
	return g.password != ""
}
