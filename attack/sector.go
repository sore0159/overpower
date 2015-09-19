package attack

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
