package attack

func (s *Sector) MakeView(factionID int) *SectorView {
	sv := NewSectorView()
	sv.Faction = factionID
	sv.Turn = s.Turn
	for _, pl := range s.PlanetGrid {
		sv.AddPlanet(pl)
	}
	for loc, list := range s.ShipGrid {
		if s.AreaVisible(factionID, loc) {
			for _, cl := range list {
				sv.AddShip(cl)
			}
		}
	}
	for loc, list := range s.TrailGrid {
		if s.AreaVisible(factionID, loc) {
			for _, trail := range list {
				sv.AddShipTrail(loc, trail)
			}
		}
	}
	return sv
}

func (s *Sector) AreaVisible(factionID int, area [2]int) bool {
	for loc, pl := range s.PlanetGrid {
		if pl.Faction() == factionID && HexDist(area, loc) < VISRANGE {
			return true
		}
	}
	return false
}

type SectorView struct {
	Faction    int
	Turn       int
	PlanetGrid map[[2]int]PlanetView
	ShipGrid   map[[2]int][]ShipView
	TrailGrid  map[[2]int][]ShipTrail
}

func NewSectorView() *SectorView {
	return &SectorView{
		PlanetGrid: map[[2]int]PlanetView{},
		ShipGrid:   map[[2]int][]ShipView{},
		TrailGrid:  map[[2]int][]ShipTrail{},
	}
}

type PlanetView struct {
	Name        string
	Location    [2]int
	Yours       bool
	Inhabitants [2]int
	Resources   int
}

func NewPlanetView() *PlanetView {
	return &PlanetView{}
}

type ShipView struct {
	FactionID int
	Size      int
	Location  [2]int
	Yours     bool
	Target    [2]int
}

func NewShipView() *ShipView {
	return &ShipView{}
}

func (sv *SectorView) AddPlanet(pl *Planet) {
	pv := NewPlanetView()
	pv.Name = pl.Name
	pv.Location = pl.Location
	if pl.Faction() == sv.Faction {
		pv.Yours = true
		pv.Inhabitants = pl.Inhabitants
		pv.Resources = pl.Resources
	}
	sv.PlanetGrid[pl.Location] = *pv
}

func (sv *SectorView) AddShip(cl *Ship) {
	cv := NewShipView()
	cv.FactionID = cl.FactionID
	cv.Size = cl.Size
	loc := cl.Location()
	cv.Location = loc
	if cv.FactionID == sv.Faction {
		cv.Yours = true
		cv.Target = cl.Target()
	}
	if list, ok := sv.ShipGrid[loc]; ok {
		sv.ShipGrid[loc] = append(list, *cv)
	} else {
		sv.ShipGrid[loc] = []ShipView{*cv}
	}
}

func (sv *SectorView) AddShipTrail(loc [2]int, trail ShipTrail) {
	if list, ok := sv.TrailGrid[loc]; ok {
		sv.TrailGrid[loc] = append(list, trail)
	} else {
		sv.TrailGrid[loc] = []ShipTrail{trail}
	}
}

type TextView struct {
	Center        [2]int
	OrderedCoords []CoordView
}

func NewTextView() *TextView {
	return &TextView{
		OrderedCoords: []CoordView{},
	}
}

type CoordView struct {
	Dist   int
	Coord  [2]int
	Planet PlanetView
	Ships  []ShipView
	Trails []ShipTrail
}

func NewCoordView() *CoordView {
	return &CoordView{
		Ships:  []ShipView{},
		Trails: []ShipTrail{},
	}
}

func MakeTextView(center [2]int, sector *SectorView) *TextView {
	tv := NewTextView()
	tv.Center = center
	for loc, plv := range sector.PlanetGrid {
		c := NewCoordView()
		c.Coord = loc
		c.Dist = HexDist(center, loc)
		c.Planet = plv
		if list, ok := sector.ShipGrid[loc]; ok {
			c.Ships = list
		}
		if list, ok := sector.TrailGrid[loc]; ok {
			c.Trails = list
		}
		tv.OrderedCoords = append(tv.OrderedCoords, *c)
	}
	for loc, list := range sector.ShipGrid {
		if _, ok := sector.PlanetGrid[loc]; ok {
			continue
		}
		c := NewCoordView()
		c.Coord = loc
		c.Dist = HexDist(center, loc)
		c.Ships = list
		if list, ok := sector.TrailGrid[loc]; ok {
			c.Trails = list
		}
		tv.OrderedCoords = append(tv.OrderedCoords, *c)
	}
	for loc, list := range sector.TrailGrid {
		if _, ok := sector.PlanetGrid[loc]; ok {
			continue
		}
		if _, ok := sector.ShipGrid[loc]; ok {
			continue
		}
		c := NewCoordView()
		c.Coord = loc
		c.Dist = HexDist(center, loc)
		c.Trails = list
		tv.OrderedCoords = append(tv.OrderedCoords, *c)
	}
	tv.SortCoords()
	return tv
}
