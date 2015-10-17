package attack

func (tv *TextView) SortCoords() {
	tv.SortTool(false, tv.Center)
}

func (tv *TextView) SortTool(reOrder bool, center [2]int) {
	m := map[int][]CoordView{}
	var highest int
	for _, cv := range tv.OrderedCoords {
		if cv.Planet.Name == "" && len(cv.Ships) == 0 {
			continue
		}
		cv.Trails = [][]ShipTrail{}
		if reOrder {
			d := HexDist(center, cv.Coord)
			cv.Dist = d
		}
		if cv.Dist > highest {
			highest = cv.Dist
		}
		if list, ok := m[cv.Dist]; ok {
			m[cv.Dist] = append(list, cv)
		} else {
			m[cv.Dist] = []CoordView{cv}
		}
	}
	for _, trailBlob := range tv.Trails {
		closestLoc := [2]int{}
		var closest int = -1
		for _, trail := range trailBlob {
			d := HexDist(center, trail.Location)
			if closest == -1 || d < closest {
				closestLoc = trail.Location
				closest = d
			}
		}
		if cvList, ok := m[closest]; ok {
			var found bool
			for i, cv := range cvList {
				if cv.Coord == closestLoc {
					cvList[i].Trails = append(cv.Trails, trailBlob)
					found = true
					break
				}
			}
			if !found {
				c := NewCoordView()
				c.Coord = closestLoc
				c.Polar = HexPolar(c.Coord)
				c.Dist = closest
				c.Trails = [][]ShipTrail{trailBlob}
				m[closest] = append(cvList, *c)
			}
		} else {
			if closest > highest {
				highest = closest
			}
			c := NewCoordView()
			c.Coord = closestLoc
			c.Polar = HexPolar(c.Coord)
			c.Dist = closest
			c.Trails = [][]ShipTrail{trailBlob}
			m[closest] = []CoordView{*c}
		}
	}
	newOrder := []CoordView{}
	for i := 0; i <= highest; i++ {
		if list, ok := m[i]; ok {
			for _, cv := range list {
				newOrder = append(newOrder, cv)
			}
		}
	}
	tv.OrderedCoords = newOrder
}

func (tv *TextView) Recenter(center [2]int) {
	if tv.Center == center {
		return
	}
	tv.Center = center
	tv.PolarC = HexPolar(center)
	tv.SortTool(true, center)
}
