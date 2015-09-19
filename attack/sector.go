package attack

import "log"

type Sector struct {
	Turn       int
	PlanetGrid map[[2]int]*Planet
	ShipGrid   map[[2]int][]*Ship
	TrailGrid  map[[2]int][]ShipTrail
}

func NewSector() *Sector {
	return &Sector{
		PlanetGrid: map[[2]int]*Planet{},
		ShipGrid:   map[[2]int][]*Ship{},
		TrailGrid:  map[[2]int][]ShipTrail{},
	}
}

func (s *Sector) AddShip(ship *Ship) {
	loc := ship.Location()
	if list, ok := s.ShipGrid[loc]; ok {
		s.ShipGrid[loc] = append(list, ship)
	} else {
		s.ShipGrid[loc] = []*Ship{ship}
	}
}

func (s *Sector) AddTrailGrid(grid map[[2]int]ShipTrail) {
	for loc, trail := range grid {
		if list, ok := s.TrailGrid[loc]; ok {
			s.TrailGrid[loc] = append(list, trail)
		} else {
			s.TrailGrid[loc] = []ShipTrail{trail}
		}
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
	for i := 1; i < 11; i++ {
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
