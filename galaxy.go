package overpower

import (
	"mule/hexagon"
)

func (op *TotallyOP) MakeGalaxy() (ok bool) {
	fids := make([]int, len(op.Factions))
	for i, f := range op.Factions {
		fids[i] = f.Fid()
	}
	if len(fids) < 1 {
		return true
	}
	gid := op.Game.Gid()
	fids = shuffleInts(fids)
	s := op.Source
	// ------- //
	numP := len(fids) * 16
	bigN := numP / 4
	clearR := 5
	names := GetNames(numP + len(fids))
	// -------- GAME ---------- //
	op.Game.SetTurn(1)
	for _, f := range op.Factions {
		f.SetDone(false)
	}
	// -------- PLANETS -------- //
	borion, ok := s.NewPlanet("Planet Borion", gid, 999, 0, 10, 30, 0, hexagon.Coord{0, 0})
	if !ok {
		return false
	}
	planets := []Planet{borion}
	places := map[hexagon.Coord]bool{hexagon.Coord{0, 0}: true}
	pids := map[int]bool{0: true, 999: true}
	for i := 0; i < numP; i++ {
		name := names[i]
		var pid, res, inhab int
		if i < bigN {
			res = pick(6) + 4
			inhab = pick(4)
		} else {
			res = pick(10)
		}
		for pids[pid] {
			pid = 100 + pick(898)
		}
		pids[pid] = true
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
		p, ok := s.NewPlanet(name, gid, pid, 0, inhab, res, 0, test)
		if !ok {
			return false
		}
		planets = append(planets, p)
	}
	homeDist := 25
	for i, fid := range fids {
		var test hexagon.Coord
		var pid int
		for {
			spot := hexagon.Polar{homeDist + pick(3), 0}
			spot[1] = spot[0]*i + spot[0]/4 + pick(spot[0]/2) - 1
			test = spot.Coord()
			if !places[test] {
				break
			}
		}
		for pids[pid] {
			pid = 100 + pick(898)
		}
		pids[pid] = true
		p, ok := s.NewPlanet(names[numP+i], gid, pid, fid, 5, 15, 5, test)
		if !ok {
			return false
		}
		planets = append(planets, p)
		s.NewMapView(gid, fid, test)
	}
	op.Planets = planets
	// -------- VIEWS --------- //
	for _, fid := range fids {
		for _, p := range planets {
			var pv PlanetView
			var ok bool
			if fid == p.Controller() {
				pv, ok = s.NewPlanetView(fid, 1, p)
			} else {
				pv, ok = s.NewPlanetView(fid, 0, p)
			}
			if !ok {
				return false
			}
			op.SetPV(pv)
		}
	}
	return true
}
