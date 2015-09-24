package attack

func (tv *TextView) SortCoords() {
	tv.SortTool(false, [2]int{})
}

func (tv *TextView) SortTool(reOrder bool, newCenter [2]int) {
	m := map[int][]int{}
	var highest int
	for i, cv := range tv.OrderedCoords {
		if reOrder {
			d := HexDist(newCenter, cv.Coord)
			tv.OrderedCoords[i].Dist = d
			cv.Dist = d
		}
		if cv.Dist > highest {
			highest = cv.Dist
		}
		if list, ok := m[cv.Dist]; ok {
			m[cv.Dist] = append(list, i)
		} else {
			m[cv.Dist] = []int{i}
		}
	}
	newOrder := []CoordView{}
	for i := 0; i <= highest; i++ {
		if list, ok := m[i]; ok {
			for _, j := range list {
				newOrder = append(newOrder, tv.OrderedCoords[j])
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
