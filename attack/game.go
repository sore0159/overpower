package attack

import "log"

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

func (g *Game) RunTurn() {
	orders := map[int][]Order{}
	for fID, faction := range g.Factions {
		orders[fID] = faction.BuildOrders
		faction.BuildOrders = []Order{}
	}
	g.Sector.AdvanceTurn(orders)
	for fID, f := range g.Factions {
		f.View = *g.Sector.MakeView(fID)
		f.TurnDone = false
	}
}

func (s *Sector) AdvanceTurn(orders map[int][]Order) {
	// Ship Spawning //
	for fac, list := range orders {
		for _, order := range list {
			pl, ok := s.PlanetGrid[order.Location]
			if !ok {
				log.Fatal("Bad order", order, "by faction", fac, ": no planet found at", order.Location)
			}
			if pl.Faction() != fac {
				log.Fatal("Bad faction", fac, "ordering", order, "to planet", pl.Name, "belonging to", pl.Faction())
			}
			cl := pl.FireLaunchers(order.Size, order.Target)
			s.AddShip(cl)
		}
	}
	// Ship Movement //
	oldShipGrid := s.ShipGrid
	s.TrailGrid = map[[2]int][]ShipTrail{}
	s.ShipGrid = map[[2]int][]*Ship{}
	LandList := map[int][]*Ship{}
	for _, shipList := range oldShipGrid {
		for _, ship := range shipList {
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
	for i := 1; i <= CLOUDSPEED; i++ {
		if len(LandList[i]) == 0 {
			continue
		}
		list := shuffleShips(LandList[i])
		for _, ship := range list {
			pl, ok := s.PlanetGrid[ship.Location()]
			if !ok {
				log.Fatal("Bad ship landing: no planet at", ship.Location(), "for ship", ship)
			}
			ship.Land(pl)
		}
	}
	// Build Launchers //
	for _, pl := range s.PlanetGrid {
		pl.BuildLaunchers()
	}
	// The End //
	s.Turn++
}
