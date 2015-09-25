package attack

type Sector struct {
	Turn       int
	PlanetIDs  map[int][2]int
	PlanetGrid map[[2]int]*Planet
	ShipGrid   map[[2]int][]*Ship
	TrailGrid  map[[2]int][]ShipTrail
}

func NewSector() *Sector {
	return &Sector{
		PlanetIDs:  map[int][2]int{},
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

func (s *Sector) AddPlanet(pl *Planet, loc [2]int) {
	id := 99 + pick(899)
	for _, ok := s.PlanetIDs[id]; ok; {
		id = 99 + pick(899)
		_, ok = s.PlanetIDs[id]
	}
	pl.ID = id
	s.PlanetIDs[id] = loc
	s.PlanetGrid[loc] = pl
	pl.Location = loc
}

func (s *Sector) GetPlanet(id int) (pl *Planet, ok bool) {
	if loc, ok := s.PlanetIDs[id]; !ok {
		return nil, false
	} else {
		pl, ok = s.PlanetGrid[loc]
		return pl, ok
	}
}
