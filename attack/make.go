package attack

import (
	//"fmt"
	//	"math"
	"strings"
)

func MakeGame(facNames []string) *Game {
	g := NewGame()
	homePlanets := g.Sector.MakePlanets(len(facNames), 100)
	g.Sector.Turn = 1
	nameMap := map[int]string{}
	for i, name := range facNames {
		f := NewFaction()
		f.Name = name
		f.FactionID = i + 1
		g.Factions[i+1] = f
		nameMap[i+1] = name
		f.View = *g.Sector.MakeView(f)
		for coord, pl := range homePlanets {
			if pl.Faction() == f.FactionID {
				f.TV = MakeTextView(coord, &f.View)
				break
			}
		}
	}
	for _, f := range g.Factions {
		f.OtherNames = nameMap
	}
	return g
}

func (s *Sector) MakePlanets(homeworlds, total int) (homePlanetMap map[[2]int]*Planet) {
	/*num := 120
	bigN := num/4 - 1
	bigRange := 40
	littleRange := 80
	hwRange := [2]int{70, 75}
	*/
	num := 16 * homeworlds
	bigN := num/4 - 1
	littleN := num - (bigN + 1)
	bigArea := bigN * 126
	littleArea := littleN*125 + bigArea
	// a = 3*n*(n+1) + 1
	// n^2 + n = (a-1)/3
	var bigRange, littleRange int
	for i := 1; ; i++ {
		if 3*i*(i+1)+1 > bigArea {
			bigRange = i
			break
		}
	}
	for i := bigRange; ; i++ {
		if 3*i*(i+1)+1 > littleArea {
			littleRange = i
			break
		}
	}
	s.Size = littleRange
	hwRange := [2]int{bigRange + (littleRange-bigRange)/2, littleRange - (littleRange-bigRange)/4}

	orion := NewPlanet()
	orion.Name = "Planet Borion" // DO YOU WANT TO GET US ALL SUED, BOY?
	orion.Inhabitants[1] = 20
	orion.Resources = 30
	orion.ID = 999
	s.PlanetIDs[999] = [2]int{0, 0}
	s.PlanetGrid[[2]int{0, 0}] = orion
	names := shuffleWords(GetAdj())
	bigPlanets := make([]*Planet, bigN)
	littlePlanets := make([]*Planet, num-bigN-1)
	homePlanets := make([]*Planet, homeworlds)
	homePlanetMap = make(map[[2]int]*Planet, homeworlds)
	for i := 0; i < len(bigPlanets); i++ {
		pl := NewPlanet()
		pl.Name = "Planet " + strings.Title(names[i])
		pl.Resources = 10 + pick(10)
		if pick(3) == 1 {
			pl.Inhabitants[1] = pick(10)
		}
		bigPlanets[i] = pl
	}
	for i := 0; i < len(littlePlanets); i++ {
		pl := NewPlanet()
		pl.Name = "Planet " + strings.Title(names[i+bigN])
		pl.Resources = pick(10)
		littlePlanets[i] = pl
	}
	for i := 0; i < homeworlds; i++ {
		pl := NewPlanet()
		pl.Name = "Planet " + strings.Title(names[i+num])
		pl.Resources = 15
		pl.Inhabitants = [2]int{i + 1, 5}
		pl.Launchers = 5
		homePlanets[i] = pl
	}
	for i := 0; i < len(bigPlanets); {
		radius := pick(bigRange)
		circ := 6 * radius
		theta := pick(circ) - 1
		loc := HexAngle2Grid(radius, theta)
		if s.CoordClear(loc, CLEARDIST) {
			s.AddPlanet(bigPlanets[i], loc)
			i++
		}
	}
	for i := 0; i < len(littlePlanets); {
		radius := pick(littleRange-bigRange) + bigRange
		circ := 6 * radius
		theta := pick(circ) - 1
		loc := HexAngle2Grid(radius, theta)
		if s.CoordClear(loc, CLEARDIST) {
			s.AddPlanet(littlePlanets[i], loc)
			i++
		}
	}
	homeCoords := s.SplitSector(homeworlds, hwRange)
	for i, coord := range homeCoords {
		s.AddPlanet(homePlanets[i], coord)
		homePlanetMap[coord] = homePlanets[i]
	}
	return homePlanetMap
}

func (s *Sector) SplitSector(homeworlds int, hwRange [2]int) [][2]int {
	minD, maxD := hwRange[0], hwRange[1]
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
