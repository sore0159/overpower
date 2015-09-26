package attack

type PlanetView struct {
	Name        string
	ID          int
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
	ShipID    int
	Size      int
	Location  [2]int
	Yours     bool
	Target    [2]int
	Left      int
	Detected  []ShipTrail
}

func NewShipView() *ShipView {
	return &ShipView{
	//		Detected: [][2]int{},
	}
}

func (s *Sector) MakeView(f *Faction) *SectorView {
	sv := NewSectorView()
	factionID := f.FactionID
	sv.Faction = factionID
	sv.Turn = s.Turn
	sv.PlanetIDs = s.PlanetIDs
	seen := map[int]*ShipView{}
	if f.View.PlanetGrid != nil {
		sv.PlanetGrid = f.View.PlanetGrid
	}
	for _, pl := range s.PlanetGrid {
		sv.UpdatePlanet(pl)
	}
	for loc, list := range s.ShipGrid {
		vis := s.AreaVisible(factionID, loc)
		for _, cl := range list {
			if vis || cl.FactionID == sv.Faction {
				seen[cl.ShipID] = sv.AddShip(cl)
			}
		}
	}
	for loc, list := range s.TrailGrid {
		if s.AreaVisible(factionID, loc) {
			for _, trail := range list {
				sv.AddShipTrail(loc, trail)
				if shv, ok := seen[trail.TrailID]; ok {
					shv.AddShipTrail(trail)
				}
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
	PlanetIDs  map[int][2]int
	PlanetGrid map[[2]int]PlanetView
	ShipGrid   map[[2]int][]ShipView
	TrailGrid  map[[2]int][]ShipTrail
}

func NewSectorView() *SectorView {
	return &SectorView{
		PlanetIDs:  map[int][2]int{},
		PlanetGrid: map[[2]int]PlanetView{},
		ShipGrid:   map[[2]int][]ShipView{},
		TrailGrid:  map[[2]int][]ShipTrail{},
	}
}

func (sv *SectorView) AddPlanet(pl *Planet) {
	pv := NewPlanetView()
	pv.Name = pl.Name
	pv.ID = pl.ID
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
	pv.ID = pl.ID
	pv.Location = pl.Location
	pv.Yours = pl.Faction() == sv.Faction
	pv.LastSeen = sv.Turn
	pv.Inhabitants = pl.Inhabitants
	pv.Resources = pl.Resources
	delete(sv.PlanetGrid, pl.Location)
	sv.PlanetGrid[pl.Location] = *pv
}

func (sv *SectorView) AddShip(cl *Ship) *ShipView {
	cv := NewShipView()
	cv.FactionID = cl.FactionID
	cv.Size = cl.Size
	cv.ShipID = cl.ShipID
	loc := cl.Location()
	cv.Location = loc
	cur := (len(cl.Path) - 1 - cl.ILocation)
	cv.Left = cur / SHIPSPEED
	if cur%SHIPSPEED != 0 {
		cv.Left++
	}
	if cv.FactionID == sv.Faction {
		cv.Yours = true
		cv.Target = cl.Target()
	}
	if list, ok := sv.ShipGrid[loc]; ok {
		sv.ShipGrid[loc] = append(list, *cv)
	} else {
		sv.ShipGrid[loc] = []ShipView{*cv}
	}
	return &(sv.ShipGrid[loc][len(sv.ShipGrid[loc])-1])
}

func (shv *ShipView) AddShipTrail(trail ShipTrail) {
	l := len(shv.Detected)
	if l == 0 {
		shv.Detected = []ShipTrail{trail}
		return
	}
	for i, test := range shv.Detected {
		if test.Count > trail.Count {
			shv.Detected = append(shv.Detected[:i], append([]ShipTrail{trail}, shv.Detected[i:]...)...)
			break
		}
		if i == len(shv.Detected)-1 {
			shv.Detected = append(shv.Detected, trail)
			break
		}
	}
}

func (sv *SectorView) AddShipTrail(loc [2]int, trail ShipTrail) {
	if list, ok := sv.TrailGrid[loc]; ok {
		sv.TrailGrid[loc] = append(list, trail)
	} else {
		sv.TrailGrid[loc] = []ShipTrail{trail}
	}
}
