package attack

type TextView struct {
	Center         [2]int
	PolarC         [3]int
	PlanetIDs      map[int][2]int
	OrderedCoords  []CoordView
	Filter         bool
	FilteredCoords []CoordView
	Trails         [][]ShipTrail
}

func NewTextView() *TextView {
	return &TextView{
		PlanetIDs:     map[int][2]int{},
		OrderedCoords: []CoordView{},
		Trails:        [][]ShipTrail{},
	}
}

type CoordView struct {
	Dist   int
	Coord  [2]int
	Polar  [3]int
	Planet PlanetView
	Ships  []ShipView
	Trails [][]ShipTrail
}

func NewCoordView() *CoordView {
	return &CoordView{
		Ships:  []ShipView{},
		Trails: [][]ShipTrail{},
	}
}

func MakeTextView(center [2]int, sector *SectorView) *TextView {
	tv := NewTextView()
	tv.Center = center
	tv.PolarC = HexPolar(center)
	tv.PlanetIDs = sector.PlanetIDs
	// TRAILS //
	mList := map[int][]ShipTrail{}
	for _, list := range sector.TrailGrid {
		for _, trail := range list {
			if list2, ok := mList[trail.TrailID]; ok {
				for i, test := range list2 {
					if test.FactionID == sector.Faction {
						continue
					}
					if test.Count > trail.Count {
						mList[trail.TrailID] = append(list2[:i], append([]ShipTrail{trail}, list2[i:]...)...)
						break
					}
					if i == len(list2)-1 {
						mList[trail.TrailID] = append(list2, trail)
						break
					}
				}
			} else {
				if trail.FactionID != sector.Faction {
					mList[trail.TrailID] = []ShipTrail{trail}
				}
			}
		}
	}
	// PLANETS //
	for loc, plv := range sector.PlanetGrid {
		c := NewCoordView()
		c.Coord = loc
		c.Polar = HexPolar(loc)
		c.Dist = HexDist(center, loc)
		c.Planet = plv
		if list, ok := sector.ShipGrid[loc]; ok {
			c.Ships = list
			for _, shv := range list {
				delete(mList, shv.ShipID)
			}
		}
		tv.OrderedCoords = append(tv.OrderedCoords, *c)
	}
	for loc, list := range sector.ShipGrid {
		if _, ok := sector.PlanetGrid[loc]; ok {
			continue
		}
		c := NewCoordView()
		c.Coord = loc
		c.Polar = HexPolar(loc)
		c.Dist = HexDist(center, loc)
		c.Ships = list
		for _, shv := range list {
			delete(mList, shv.ShipID)
		}
		tv.OrderedCoords = append(tv.OrderedCoords, *c)
	}
	r := [][]ShipTrail{}
	for _, list := range mList {
		r = append(r, list)
	}
	tv.Trails = r
	tv.SortCoords()
	return tv
}
