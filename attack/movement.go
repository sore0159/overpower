package attack

func (c *Ship) SetPath(start [2]int, stop [2]int) {
	c.Path = HexPath(start, stop)
}

func (c *Ship) Move() (dist int, trailGrid map[[2]int]ShipTrail) {
	trailGrid = map[[2]int]ShipTrail{}
	pathLen := len(c.Path) - 1
	for i := 0; i <= CLOUDSPEED; i++ {
		if c.ILocation == pathLen {
			break
		}
		trailGrid[c.Location()] = *c.MakeShipTrail()
		c.ILocation++
		dist++
	}
	return
}

func (c *Ship) AtTarget() bool {
	return c.ILocation == len(c.Path)-1
}

func (c *Ship) Land(pl *Planet) {
	if c.FactionID == pl.Inhabitants[0] {
		pl.Arrivals += c.Size
		return
	}
	pl.Inhabitants[1] -= c.Size
	if pl.Inhabitants[1] >= 0 {
		return
	}
	pl.Inhabitants[0] = c.FactionID
	pl.Arrivals = -1 * pl.Inhabitants[1]
	pl.Inhabitants[1] = 0
}
