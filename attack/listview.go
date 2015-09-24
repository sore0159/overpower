package attack

type TextView struct {
	Center        [2]int
	PolarC        [3]int
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
	Polar  [3]int
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
	tv.PolarC = HexPolar(center)
	for loc, plv := range sector.PlanetGrid {
		c := NewCoordView()
		c.Coord = loc
		c.Polar = HexPolar(loc)
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
		c.Polar = HexPolar(loc)
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
		c.Polar = HexPolar(loc)
		c.Dist = HexDist(center, loc)
		c.Trails = list
		tv.OrderedCoords = append(tv.OrderedCoords, *c)
	}
	tv.SortCoords()
	return tv
}
