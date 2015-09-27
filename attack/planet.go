package attack

import "strconv"

type Planet struct {
	Name        string
	ID          int
	Location    [2]int
	Inhabitants [2]int // FactionID, amount
	Arrivals    int
	Resources   int
	Launchers   int
}

func NewPlanet() *Planet {
	return &Planet{}
}

func (pl *Planet) Faction() int {
	return pl.Inhabitants[0]
}

func (pl *Planet) BuildLaunchers() {
	if pl.Resources == 0 || (pl.Inhabitants[1] == 0 && pl.Arrivals == 0) || pl.Inhabitants[0] == 0 {
		return
	}
	num := pl.Inhabitants[1]
	if num > pl.Resources {
		num = pl.Resources
	}
	pl.Resources -= num
	pl.Launchers += num
	pl.Inhabitants[1] += pl.Arrivals
	pl.Arrivals = 0
}

func (g *Game) FireLaunchers(pl *Planet, size int, target [2]int) (c *Ship, report string) {
	if pl.Launchers < size {
		LogF("Bad firelaunchers size", size, "for planet", pl.Name, "with launchers", pl.Launchers)
		return nil, ""
	}
	pl.Launchers -= size
	c = NewShip()
	c.FactionID = pl.Faction()
	c.Size = size
	c.SetPath(pl.Location, target)
	report = pl.Name + " launched ship sized " + strconv.Itoa(size) + " toward " + g.Sector.PlanetGrid[target].Name
	return
}
