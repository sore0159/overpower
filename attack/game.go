package attack

type Game struct {
	Sector   *Sector
	Factions map[int]*Faction
}

func NewGame() *Game {
	return &Game{
		Sector:   NewSector(),
		Factions: map[int]*Faction{},
	}
}

func (g *Game) FactionName(id int) string {
	return g.Factions[id].Name
}

func (g *Game) RunTurn() {
	orders := map[int][]Order{}
	for fID, faction := range g.Factions {
		orders[fID] = faction.BuildOrders
		faction.BuildOrders = []Order{}
	}
	g.Sector.AdvanceTurn(orders)
	for fID, f := range g.Factions {
		f.View = *g.Sector.MakeView(fID)
	}
}
