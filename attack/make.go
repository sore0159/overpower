package attack

import "strings"

func MakeGame(facNames []string) *Game {
	g := NewGame()
	g.Sector.MakePlanets(len(facNames), 100)
	g.Sector.Turn = 1
	for i, name := range facNames {
		f := NewFaction()
		f.Name = name
		f.FactionID = i + 1
		g.Factions[i+1] = f
		f.View = *g.Sector.MakeView(i + 1)
	}
	return g
}

func (s *Sector) MakePlanets(homeworlds, total int) {
	num := 100
	orion := NewPlanet()
	orion.Name = "Planet Borion" // DO YOU WANT TO GET US ALL SUED, BOY?
	orion.Inhabitants[1] = 20
	orion.Resources = 30
	s.PlanetGrid[[2]int{0, 0}] = orion
	names := shuffleWords(ADJECTIVES[:])
	bigN := num/4 - 1
	bigPlanets := make([]*Planet, bigN)
	littlePlanets := make([]*Planet, num-bigN-1)
	homePlanets := make([]*Planet, homeworlds)
	for i := 0; i < len(bigPlanets); i++ {
		pl := NewPlanet()
		pl.Name = "Planet " + strings.ToTitle(names[i])
		pl.Resources = 10 + pick(10)
		if pick(3) == 1 {
			pl.Inhabitants[1] = pick(10)
		}
		bigPlanets[i] = pl
	}
	for i := 0; i < len(littlePlanets); i++ {
		pl := NewPlanet()
		pl.Name = "Planet " + strings.ToTitle(names[i+bigN])
		pl.Resources = pick(10)
		littlePlanets[i] = pl
	}
	for i := 0; i < homeworlds; i++ {
		pl := NewPlanet()
		pl.Name = "Planet " + strings.ToTitle(names[i+num])
		pl.Resources = 15
		pl.Inhabitants = [2]int{i + 1, 5}
		pl.Launchers = 5
		homePlanets[i] = pl
	}
	bigRange := 100
	for i := 0; i < len(bigPlanets); {
		radius := pick(bigRange)
		circ := 6 * radius
		theta := pick(circ) - 1
		loc := HexAngle2Grid(radius, theta)
		if s.CoordClear(loc, CLEARDIST) {
			bigPlanets[i].Location = loc
			s.PlanetGrid[loc] = bigPlanets[i]
			i++
		}
	}
	littleRange := 200
	for i := 0; i < len(littlePlanets); {
		radius := pick(littleRange-bigRange) + bigRange
		circ := 6 * radius
		theta := pick(circ) - 1
		loc := HexAngle2Grid(radius, theta)
		if s.CoordClear(loc, CLEARDIST) {
			littlePlanets[i].Location = loc
			s.PlanetGrid[loc] = littlePlanets[i]
			i++
		}
	}
	homeCoords := s.SplitSector(homeworlds, 170, 190)
	for i, coord := range homeCoords {
		s.PlanetGrid[coord] = homePlanets[i]
	}
}

func (s *Sector) SplitSector(homeworlds, minD, maxD int) [][2]int {
	coords := [][2]int{}
	for i := 0; i < homeworlds; {
		radius := minD + pick(maxD-minD)
		leg := 6 * radius / homeworlds
		theta := i*(leg) - 1 + 3*leg/8 + pick(leg/4)
		loc := HexAngle2Grid(radius, theta)
		if s.CoordClear(loc, CLEARDIST) {
			coords = append(coords, loc)
			i++
		}
	}
	return coords
}

func (s *Sector) CoordClear(loc [2]int, radius int) bool {
	if _, ok := s.PlanetGrid[loc]; ok {
		return false
	}
	for r := 1; r <= radius; r++ {
		for theta := 0; theta < r*6; theta++ {
			mod := HexAngle2Grid(r, theta)
			if _, ok := s.PlanetGrid[[2]int{loc[0] + mod[0], loc[1] + mod[1]}]; ok {
				return false
			}
		}
	}
	return true
}
