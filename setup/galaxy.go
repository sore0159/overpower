package setup

import (
	"mule/hexagon"
)

type Planet struct {
	Name      string
	Pid       int
	Fid       int
	Inhab     int
	Resources int
	Parts     int
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
	names := GetNames(numP + len(fids))
	planets := []Planet{{"Planet Borion", 999, 0, 10, 30, 0, [2]int{0, 0}}}
	places := map[hexagon.Coord]bool{hexagon.Coord{0, 0}: true}
	pids := map[int]bool{0: true, 999: true}
	for i := 0; i < numP; i++ {
		p := Planet{Name: names[i]}
		if i < bigN {
			p.Resources = pick(6) + 4
			p.Inhab = pick(4)
		} else {
			p.Resources = pick(10)
		}
		for pids[p.Pid] {
			p.Pid = 100 + pick(898)
		}
		pids[p.Pid] = true
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
		p.Loc = [2]int{test[0], test[1]}
		planets = append(planets, p)
	}
	homeDist := 25
	for i, fid := range fids {
		var test hexagon.Coord
		for {
			spot := hexagon.Polar{homeDist + pick(3), 0}
			spot[1] = spot[0]*i + spot[0]/4 + pick(spot[0]/2) - 1
			test = spot.Coord()
			if !places[test] {
				break
			}
		}
		p := Planet{Name: names[numP+i], Fid: fid, Inhab: 5, Resources: 15, Parts: 5, Loc: [2]int{test[0], test[1]}}
		for pids[p.Pid] {
			p.Pid = 100 + pick(898)
		}
		pids[p.Pid] = true
		planets = append(planets, p)
	}
	return planets
}
