package overpower

import (
	"mule/hexagon"
)

func PlanetaryLanding(source Source, pl Planet, sh Ship, turn int, arrivals map[hexagon.Coord]int, names map[int]string) (logErr error) {
	shFid := sh.Fid()
	plFid := pl.Controller()
	atk := sh.Size()
	loc := pl.Loc()
	def := arrivals[loc]
	defer BothLandingReports(source, plFid, turn, sh, pl, arrivals)
	if atk < 1 {
		return
	}
	if shFid == plFid {
		arrivals[loc] += atk
		return
	}
	defer BothSee(source, pl, pl.Controller(), sh.Fid(), turn, arrivals)
	if def >= atk {
		if def == atk {
			delete(arrivals, loc)
		} else {
			arrivals[loc] = def - atk
		}
		return
	}
	delete(arrivals, loc)
	atk -= def
	def = pl.Inhabitants()
	if def >= atk {
		pl.SetInhabitants(def - atk)
		return
	}
	pl.SetController(sh.Fid())
	pl.SetInhabitants(0)
	atk -= def
	arrivals[loc] = atk
	return
}

func BothSee(source Source, pl Planet, fid1, fid2, turn int, arrivals map[hexagon.Coord]int) {
	arv := arrivals[pl.Loc()]
	for _, fid := range []int{fid1, fid2} {
		if fid == 0 || fid == pl.Controller() {
			continue
		}
		pv := source.UpdatePlanetView(fid, turn, pl)
		if arv > 0 {
			pv.SetInhabitants(pv.Inhabitants() + arv)
		}
	}
}

func BothLandingReports(source Source, firstController, turn int, lander Ship, planet Planet, arrivals map[hexagon.Coord]int) {
	arv := arrivals[planet.Loc()]
	res := [3]int{firstController, planet.Controller(), planet.Inhabitants() + arv}
	shFid := lander.Fid()
	source.NewLandingRecord(shFid, turn, lander, res)
	if firstController != 0 && firstController != shFid {
		source.NewLandingRecord(firstController, turn, lander, res)
	}
}
