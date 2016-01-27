package overpower

import (
	"math"
	"math/rand"
	"mule/hexagon"
)

func MakeGalaxy(source Source, exodus bool) error {
	factions, err := source.Factions()
	if my, bad := Check(err, "make galaxy resource failure"); bad {
		return my
	}
	if len(factions) < 1 {
		return ErrBadArgs
	}
	fids := make([]int, len(factions))
	for i, f := range factions {
		fids[i] = f.Fid()
		f.SetDone(false)
	}
	// TESTING //
	//fids = []int{fids[0], 0, 0, 0, 0,
	//0, 0, 0, 0, 0}
	// TESTING //
	fids = shuffleInts(fids)
	// -------- GAME ---------- //
	game, err := source.Game()
	if my, bad := Check(err, "make galaxy resource failure"); bad {
		return my
	}
	game.SetTurn(1)
	// -------- PLANETS -------- //
	homes := len(fids)
	bigPerPlayer := 3
	littlePerPlayer := 12
	bigN := homes * bigPerPlayer
	littleN := homes * littlePerPlayer
	names := GetNames(bigN + littleN + homes)
	// ------- //
	// ---- BORION ---- //
	borion := source.NewPlanet("Planet Borion", 9999, 0, 10, 30, 0, hexagon.Coord{0, 0})
	planets := []Planet{borion}
	places := map[hexagon.Coord]bool{hexagon.Coord{0, 0}: true}
	for _, pt := range borion.Loc().Ring(1) {
		places[pt] = true
	}
	pids := map[int]bool{0: true, 9999: true}
	nameCount := 0
	// ------------------ INNER PLANETS ------------------ //
	//area := bigN * HexArea(10)
	area := bigN * hexagon.HexArea(15)
	bigRadius := 2

	for ; hexagon.HexArea(bigRadius) < area; bigRadius += 1 {
	}
	//Ping("INNER RAD:", bigRadius)
	for i := 0; i < bigN; i++ {
		name := names[nameCount]
		nameCount += 1
		extra := pick(5)
		res := 4 + extra
		inhab := pick(extra)
		var pid int
		for pids[pid] {
			pid = 1000 + pick(8998)
		}
		pids[pid] = true
		var spot hexagon.Coord
		for {

			testP := hexagon.Polar{pick(bigRadius), 0}
			//Ping("TEST SPOT RAD:", testP[0])
			//Ping("THETA OUT OF:", testP[0]*6)
			testP[1] = rand.Intn(testP[0] * 6)
			//Ping("THETA", testP[1])
			test := testP.Coord()
			if _, ok := places[test]; !ok {
				//Ping("AVAIL!")
				spot = test
				break
			}
		}
		places[spot] = true
		for _, pt := range spot.Ring(1) {
			places[pt] = true
		}
		p := source.NewPlanet(name, pid, 0, inhab, res, 0, spot)
		planets = append(planets, p)
	}
	// ------------------ OUTER PLANETS ------------------ //
	//allArea := (littleN * HexArea(20)) + area
	allArea := area + (littleN * hexagon.HexArea(20))
	maxRadius := bigRadius + 1
	for ; hexagon.HexArea(maxRadius) < allArea; maxRadius += 1 {
	}
	//Ping("OUTER RAD", maxRadius)
	//Ping("PICKING", bigRadius, "PLUS", maxRadius-bigRadius)
	for i := 0; i < homes; i++ {
		for j := 0; j < littlePerPlayer; j++ {
			name := names[nameCount]
			nameCount += 1
			res := 1 + pick(5)
			inhab := 0
			var pid int
			for pids[pid] {
				pid = 1000 + pick(8998)
			}
			pids[pid] = true
			var spot hexagon.Coord
			for {
				dist := bigRadius + pick(maxRadius-bigRadius)
				//Ping("TEST DIST", dist)
				hexRange := float64(dist*6) / float64(homes)
				//Ping("HEX RANGE:", hexRange)
				//Ping("PICKING HEX RANGE:", hexRange*float64(i), "PLUS", hexRange)
				hexFloat := RandF(hexRange) + hexRange*float64(i)
				hexInt := int(math.Floor(hexFloat))
				//Ping("PICKED", hexFloat, hexInt)
				testP := hexagon.Polar{dist, hexInt}
				test := testP.Coord()
				if _, ok := places[test]; !ok {
					//Ping("AVAIL")
					spot = test
					break
				}
			}
			places[spot] = true
			for _, pt := range spot.Ring(1) {
				places[pt] = true
			}
			p := source.NewPlanet(name, pid, 0, inhab, res, 0, spot)
			planets = append(planets, p)
		}
	}
	// ------------------ HOME PLANETS ------------------ //
	var homeRadiusStart int
	var homeRadiusEnd int

	if exodus {
		homeRadiusStart = bigRadius / 5
		homeRadiusEnd = (2 * bigRadius) / 5
	} else {
		homeRadiusStart = (3 * maxRadius) / 5
		homeRadiusEnd = (4 * maxRadius) / 5
	}

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
		p := source.NewPlanet(name, pid, fid, 5, 15, 5, spot)
		planets = append(planets, p)
		// TESTING //
		if fid == 0 {
			continue
		}
		// TESTING //
		source.NewMapView(fid, spot)
	}
	// -------- VIEWS --------- //
	for _, fid := range fids {
		// TESTING //
		if fid == 0 {
			continue
		}
		// TESTING //
		for _, p := range planets {
			source.NewPlanetView(fid, p, exodus)
		}
	}
	return nil
}
