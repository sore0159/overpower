package attack

type TextView struct {
	Center        [2]int
	PolarC        [3]int
	PlanetIDs     map[int][2]int
	OrderedCoords []CoordView
	Trails        [][]ShipTrail
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
				mList[trail.TrailID] = []ShipTrail{trail}
			}
		}
	}
	r := make([][]ShipTrail, len(mList))
	for i, list := range mList {
		r[i] = list
	}
	tv.Trails = r
	// PLANETS //
	for loc, plv := range sector.PlanetGrid {
		c := NewCoordView()
		c.Coord = loc
		c.Polar = HexPolar(loc)
		c.Dist = HexDist(center, loc)
		c.Planet = plv
		if list, ok := sector.ShipGrid[loc]; ok {
			c.Ships = list
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
		tv.OrderedCoords = append(tv.OrderedCoords, *c)
	}
	tv.SortCoords()
	return tv
}

/*func (tv *TextView) SetTrails() {
	mDist := map[int]int{}
	mList := map[int][]ShipTrail{}
	mLoc := map[int][2]int{}
	for loc, list := range sector.TrailGrid {
		dist := HexDist(center, loc)
		for _, trail := range list {
			if list2, ok := mList[trail.TrailID]; ok {
				for i, test := range list2 {
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
				mList[trail.TrailID] = []ShipTrail{trail}
			}
			if bestD, ok := mDist[trail.TrailID]; !ok || dist < bestD {
				mDist[trail.TrailID] = dist
				mLoc[trail.TrailID] = loc
			}
		}
	}
	mTrails := map[[2]int][]int{}
	for tID, loc := range mLoc {
		if list, ok := mTrails[loc]; ok {
			mTrails[loc] = append(list, tID)
		} else {
			mTrails[loc] = []int{tID}
		}
	}
	for loc, list := range mTrails {
		c := NewCoordView()
		c.Coord = loc
		c.Polar = HexPolar(loc)
		c.Dist = mDist[list[0]]
		for _, tID := range list {
			c.Trails = append(c.Trails, mList[tID])
		}
		if list, ok := sector.ShipGrid[loc]; ok {
			c.Ships = list
		}
		if plv, ok := sector.PlanetGrid[loc]; ok {
			c.Planet = plv
		}
		tv.OrderedCoords = append(tv.OrderedCoords, *c)
	}
}*/
