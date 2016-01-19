package overpower

import (
	"math"
	"math/rand"
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
	// TESTING //
	//fids = []int{fids[0], 0, 0, 0, 0,
	//0, 0, 0, 0, 0}
	// TESTING //
	gid := op.Game.Gid()
	fids = shuffleInts(fids)
	s := op.Source
	// -------- GAME ---------- //
	op.Game.SetTurn(1)
	for _, f := range op.Factions {
		f.SetDone(false)
	}
	// -------- PLANETS -------- //
	homes := len(fids)
	bigPerPlayer := 3
	littlePerPlayer := 12
	bigN := homes * bigPerPlayer
	littleN := homes * littlePerPlayer
	names := GetNames(bigN + littleN + homes)
	// ------- //
	// ---- BORION ---- //
	borion, ok := s.NewPlanet("Planet Borion", gid, 9999, 0, 10, 30, 0, hexagon.Coord{0, 0})
	if !ok {
		return false
	}
	planets := []Planet{borion}
	places := map[hexagon.Coord]bool{hexagon.Coord{0, 0}: true}
	pids := map[int]bool{0: true, 9999: true}
	nameCount := 0
	// ---- INNER ---- //
	//area := bigN * HexArea(10)
	area := bigN * hexagon.HexArea(15)
	bigRadius := 2

	for ; hexagon.HexArea(bigRadius) < area; bigRadius += 1 {
	}
	for i := 0; i < bigN; i++ {
		name := names[nameCount]
		nameCount += 1
		res := pick(6) + 4
		inhab := pick(4)
		var pid int
		for pids[pid] {
			pid = 1000 + pick(8998)
		}
		pids[pid] = true
		var spot hexagon.Coord
		for {
			testP := hexagon.Polar{pick(bigRadius), 0}
			testP[1] = rand.Intn(testP[0] * 6)
			test := testP.Coord()
			if _, ok := places[test]; !ok {
				spot = test
				break
			}
		}
		places[spot] = true
		for _, pt := range spot.Ring(1) {
			places[pt] = true
		}
		p, ok := s.NewPlanet(name, gid, pid, 0, inhab, res, 0, spot)
		if !ok {
			return false
		}
		planets = append(planets, p)
	}
	// ---- OUTER ---- //
	//allArea := (littleN * HexArea(20)) + area
	allArea := area + (littleN * hexagon.HexArea(20))
	maxRadius := bigRadius + 1
	for ; hexagon.HexArea(maxRadius) < allArea; maxRadius += 1 {
	}
	for i := 0; i < homes; i++ {
		for j := 0; j < littlePerPlayer; j++ {
			name := names[nameCount]
			nameCount += 1
			res := pick(10)
			inhab := 0
			var pid int
			for pids[pid] {
				pid = 1000 + pick(8998)
			}
			pids[pid] = true
			var spot hexagon.Coord
			for {
				dist := bigRadius + pick(maxRadius-bigRadius)
				hexRange := float64(dist*6) / float64(homes)
				hexFloat := RandF(hexRange) + hexRange*float64(i)
				hexInt := int(math.Floor(hexFloat))
				testP := hexagon.Polar{dist, hexInt}
				test := testP.Coord()
				if _, ok := places[test]; !ok {
					spot = test
					break
				}
			}
			places[spot] = true
			for _, pt := range spot.Ring(1) {
				places[pt] = true
			}
			p, ok := s.NewPlanet(name, gid, pid, 0, inhab, res, 0, spot)
			if !ok {
				return false
			}
			planets = append(planets, p)
		}
	}
	// ---- HOME PLANETS ---- //
	homeRadiusStart := (3 * maxRadius) / 5
	homeRadiusEnd := (4 * maxRadius) / 5
	for i, fid := range fids {
		name := names[nameCount]
		nameCount += 1
		var pid int
		for pids[pid] {
			pid = 1000 + pick(8998)
		}
		pids[pid] = true
		var spot hexagon.Coord
		for {
			dist := homeRadiusStart + pick(homeRadiusEnd-homeRadiusStart)
			hexRange := float64(dist*6) / float64(homes)
			hexFloat := RandF(hexRange*.5) + hexRange*(float64(i)+.25)
			hexInt := int(math.Floor(hexFloat))
			testP := hexagon.Polar{dist, hexInt}
			test := testP.Coord()
			if _, ok := places[test]; !ok {
				spot = test
				break
			}
		}
		places[spot] = true
		for _, pt := range spot.Ring(1) {
			places[pt] = true
		}
		p, ok := s.NewPlanet(name, gid, pid, fid, 5, 15, 5, spot)
		if !ok {
			return false
		}
		planets = append(planets, p)
		// TESTING //
		if fid == 0 {
			continue
		}
		// TESTING //
		s.NewMapView(gid, fid, spot)
	}
	op.Planets = planets
	// -------- VIEWS --------- //
	for _, fid := range fids {
		// TESTING //
		if fid == 0 {
			continue
		}
		// TESTING //
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
