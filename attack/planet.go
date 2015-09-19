package attack

import "log"

type Planet struct {
	Name        string
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
	if pl.Resources == 0 || pl.Inhabitants[1] == 0 || pl.Inhabitants[0] == 0 {
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

func (pl *Planet) FireLaunchers(size int, target [2]int) *Ship {
	if pl.Launchers < size {
		log.Fatal("Bad firelaunchers size", size, "for planet", pl.Name, "with launchers", pl.Launchers)
	}
	pl.Launchers -= size
	c := NewShip()
	c.FactionID = pl.Faction()
	c.Size = size
	c.SetPath(pl.Location, target)
	return c
}
