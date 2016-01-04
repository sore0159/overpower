package db

func NewGame() *Game {
	return &Game{}
}

type Game struct {
	turnMod bool
	autoMod bool
	freeMod bool
	//
	gid       int
	owner     string
	name      string
	turn      int
	autoturn  int
	freeautos int
	password  string
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

func (g *Game) AutoTurn() int {
	return g.autoturn
}
func (g *Game) SetAutoTurn(x int) {
	if x == g.autoturn {
		return
	}
	g.autoturn = x
	g.autoMod = true
}

func (g *Game) FreeAutos() int {
	return g.freeautos
}
func (g *Game) SetFreeAutos(x int) {
	if x == g.freeautos {
		return
	}
	g.freeautos = x
	g.freeMod = true
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
