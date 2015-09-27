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

func (g *Game) UserToggleDone(fID int) (ranturn bool) {
	f, ok := g.Factions[fID]
	if !ok {
		Log("Couldn't toggleDone faction", fID, ": not found")
		return false
	}
	f.TurnDone = !f.TurnDone
	var wait bool
	for _, f := range g.Factions {
		if !f.TurnDone {
			wait = true
			break
		}
	}
	if !wait {
		g.RunTurn()
	}
	return !wait
}

func (g *Game) RunTurn() {
	orders := map[int][]Order{}
	for fID, faction := range g.Factions {
		faction.Reports = append(faction.Reports, []string{})
		fOrders := []Order{}
		for _, o := range faction.BuildOrders {
			fOrders = append(fOrders, o)
		}
		orders[fID] = fOrders
		faction.BuildOrders = map[[4]int]Order{}
	}
	g.ExecuteOrders(orders)
	for _, f := range g.Factions {
		tv := g.Sector.MakeView(f)
		f.View = *tv
		f.TV = MakeTextView(f.TV.Center, tv)
		f.TurnDone = false
	}
}

func (g *Game) ExecuteOrders(orders map[int][]Order) {
	s := g.Sector
	// Ship Spawning //
	for fac, list := range orders {
		for _, order := range list {
			pl, ok := s.PlanetGrid[order.Location]
			if !ok {
				LogF("Bad order", order, "by faction", fac, ": no planet found at", order.Location)
				continue
			}
			if pl.Faction() != fac {
				LogF("Bad faction", fac, "ordering", order, "to planet", pl.Name, "belonging to", pl.Faction())
				continue
			}
			cl, report := g.FireLaunchers(pl, order.Size, order.Target)
			if cl != nil {
				if report != "" {
					g.Factions[cl.FactionID].AddReport(report)
				}
				s.AddShip(cl)
			}
		}
	}
	// Ship Movement //
	oldShipGrid := s.ShipGrid
	s.TrailGrid = map[[2]int][]ShipTrail{}
	s.ShipGrid = map[[2]int][]*Ship{}
	LandList := map[int][]*Ship{}
	shipCount := pick(1000)
	for _, shipList := range oldShipGrid {
		for _, ship := range shipList {
			shipCount++
			ship.ShipID = shipCount
			dist, trailGrid := ship.Move()
			s.AddTrailGrid(trailGrid)
			if ship.AtTarget() {
				if list, ok := LandList[dist]; ok {
					LandList[dist] = append(list, ship)
				} else {
					LandList[dist] = []*Ship{ship}
				}
			} else {
				s.AddShip(ship)
			}
		}
	}
	// Ship Landing //
	for i := 1; i <= SHIPSPEED; i++ {
		if len(LandList[i]) == 0 {
			continue
		}
		list := shuffleShips(LandList[i])
		for _, ship := range list {
			pl, ok := s.PlanetGrid[ship.Location()]
			if !ok {
				LogF("Bad ship landing: no planet at", ship.Location(), "for ship", ship)
			}
			reports, fIDs := g.LandShip(ship, pl)
			for i, report := range reports {
				if report == "" || fIDs[i] == 0 {
					continue
				}
				g.Factions[fIDs[i]].AddReport(report)
			}
			g.Factions[ship.FactionID].View.ViewPlanet(pl)
		}
	}
	// Build Launchers //
	for _, pl := range s.PlanetGrid {
		pl.BuildLaunchers()
	}
	// The End //
	s.Turn++
}
