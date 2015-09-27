package attack

import "strconv"

func (c *Ship) SetPath(start [2]int, stop [2]int) {
	c.Path = HexPath(start, stop)
}

func (c *Ship) Move() (dist int, trailGrid map[[2]int]ShipTrail) {
	trailGrid = map[[2]int]ShipTrail{}
	pathLen := len(c.Path) - 1
	for i := 0; i < SHIPSPEED; i++ {
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

func (g *Game) LandShip(c *Ship, pl *Planet) (r [2]string, fIDs [2]int) {
	fIDs = [2]int{c.FactionID, pl.Faction()}
	r[0] = "Ship Landed at " + pl.Name + ": "
	if c.FactionID == pl.Inhabitants[0] {
		pl.Arrivals += c.Size
		r[0] += "Your population there increases to " + strconv.Itoa(pl.Arrivals+c.Size)
		return
	}
	r[1] = pl.Name + " landed on by " + g.Factions[c.FactionID].Name
	shipWin := c.Size > pl.Arrivals+pl.Inhabitants[1]
	shipFight := pl.Inhabitants[0] != 0 || pl.Inhabitants[1] > 0
	pl.Arrivals -= c.Size
	if pl.Arrivals < 0 {
		pl.Inhabitants[1] += pl.Arrivals
		pl.Arrivals = 0
	}
	if !shipWin {
		r[0] += "But the crew is killed fighting "
		if pl.Inhabitants[0] == 0 {
			r[0] += "hostile natives"
		} else {
			r[0] += "colonists of " + g.Factions[pl.Inhabitants[0]].Name
			r[1] += " but manages to fight off the crew, losing " + strconv.Itoa(c.Size) + "population"
		}
		return
	}
	if shipFight {
		r[0] += "The crew encountered "
		if pl.Inhabitants[0] == 0 {
			r[0] += "hostile natives"
		} else {
			r[0] += "colonists of " + g.Factions[pl.Inhabitants[0]].Name
			r[1] += "; your colonists there are wiped out"
		}
	} else {
		r[0] += "The crew have begun colonization"
	}

	pl.Inhabitants[0] = c.FactionID
	pl.Arrivals = -1 * pl.Inhabitants[1]
	pl.Inhabitants[1] = 0
	if shipFight {
		r[0] += " but manage to defend themselves, " + strconv.Itoa(pl.Arrivals) + " of them surviving"
	}
	return

}
