package attack

type PlanetView struct {
	Name        string
	Location    [2]int
	Yours       bool
	LastSeen    int
	Inhabitants [2]int
	Resources   int
	Launchers   int
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

func (s *Sector) MakeView(f *Faction) *SectorView {
	sv := NewSectorView()
	factionID := f.FactionID
	sv.Faction = factionID
	sv.Turn = s.Turn
	if f.View.PlanetGrid != nil {
		sv.PlanetGrid = f.View.PlanetGrid
	}
	for _, pl := range s.PlanetGrid {
		sv.UpdatePlanet(pl)
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

func (sv *SectorView) AddPlanet(pl *Planet) {
	pv := NewPlanetView()
	pv.Name = pl.Name
	pv.Location = pl.Location
	if pl.Faction() == sv.Faction {
		pv.Yours = true
		pv.LastSeen = sv.Turn
		pv.Inhabitants = pl.Inhabitants
		pv.Resources = pl.Resources
		pv.Launchers = pl.Launchers
	}
	sv.PlanetGrid[pl.Location] = *pv
}

func (sv *SectorView) UpdatePlanet(pl *Planet) {
	_, ok := sv.PlanetGrid[pl.Location]
	if !ok {
		sv.AddPlanet(pl)
		return
	} else if pl.Faction() == sv.Faction {
		delete(sv.PlanetGrid, pl.Location)
		sv.AddPlanet(pl)
	}
}

func (sv *SectorView) ViewPlanet(pl *Planet) {
	pv := NewPlanetView()
	pv.Name = pl.Name
	pv.Location = pl.Location
	pv.Yours = pl.Faction() == sv.Faction
	pv.LastSeen = sv.Turn
	pv.Inhabitants = pl.Inhabitants
	pv.Resources = pl.Resources
	delete(sv.PlanetGrid, pl.Location)
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
