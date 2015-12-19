package setup

import (
	"mule/hexagon"
)

type Planet struct {
	Name      string
	Fid       int
	Inhab     int
	Resources int
	Loc       [2]int
}

func MakeGalaxy(fids []int) []Planet {
	if len(fids) < 1 {
		return nil
	}
	fids = shuffleInts(fids)
	numP := len(fids) * 16
	bigN := numP / 4
	clearR := 5
	// ------------- //
	names := GetNames(numP)
	planets := []Planet{{"Planet Borion", 0, 20, 30, [2]int{0, 0}}}
	places := map[hexagon.Coord]bool{hexagon.Coord{0, 0}: true}
	for i := 0; i < numP; i++ {
		p := Planet{Name: names[i]}
		planets = append(planets, p)
		if i < bigN {
			p.Resources = pick(6) + 4
			p.Inhab = pick(4)
		} else {
			p.Resources = pick(10)
		}
		spot := hexagon.Polar{pick(20), 0}
		spot[1] = pick(spot[0]*6) - 1
		test := spot.Coord()
		var steps int
		for blocked := true; blocked; {
			blocked = false
			for taken, _ := range places {
				if test.StepsTo(taken) <= clearR {
					blocked = true
					break
				}
			}
			if blocked {
				steps++
				spot = hexagon.Polar{pick(20) + 5*steps, 0}
				spot[1] = pick(spot[0]*6) - 1
				test = spot.Coord()
			} else {
				places[test] = true
			}
		}
	}
	return planets
}
