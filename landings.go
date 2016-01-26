package overpower

import (
	"fmt"
)

func PlanetaryLanding(source Source, pl Planet, sh Ship, turn int, arrivals map[int]int, names map[int]string) (logErr error) {
	shFid := sh.Fid()
	plFid := pl.Controller()
	plid := pl.Pid()
	atk := sh.Size()
	def := arrivals[plid]
	aSum := atk
	dSum := def + pl.Inhabitants()
	if atk < 1 {
		return
	}
	if shFid == plFid {
		rStr := fmt.Sprintf("Your ship landed at %s with %d colonists, reinforcing inhabitants to %d.", pl.Name(), aSum, dSum+aSum)
		if !source.AddReportEvent(shFid, rStr) {
			return ErrBadArgs
		}
		arrivals[plid] += atk
		return
	}
	defer BothSee(source, pl, pl.Controller(), sh.Fid(), turn, arrivals)
	if def >= atk {
		if def == atk {
			delete(arrivals, plid)
		} else {
			arrivals[plid] = def - atk
		}
		var themStr string
		if plFid == 0 {
			themStr = "hostile natives"
		} else {
			themStr = names[plFid] + " inhabitants"
			plStr := fmt.Sprintf("%s was invaded by %d colonists from %s, but your inhabitants fought them off (%d inhabitants remaining)", pl.Name(), aSum, names[shFid], dSum-aSum)
			if !source.AddReportEvent(plFid, plStr) {
				return ErrBadArgs
			}
		}
		shStr := fmt.Sprintf("Your ship landed at %s with %d colonists, but were all killed by %s (%d inhabitants remaining).", pl.Name(), aSum, themStr, dSum-aSum)
		if !source.AddReportEvent(shFid, shStr) {
			return ErrBadArgs
		}
		return
	}
	delete(arrivals, plid)
	atk -= def
	def = pl.Inhabitants()
	if def >= atk {
		pl.SetInhabitants(def - atk)
		var themStr string
		if plFid == 0 {
			themStr = "hostile natives"
		} else {
			themStr = names[plFid] + " inhabitants"
			plStr := fmt.Sprintf("%s was invaded by %d colonists from %s, but your inhabitants fought them off (%d inhabitants remaining)", pl.Name(), aSum, names[shFid], dSum-aSum)
			if !source.AddReportEvent(plFid, plStr) {
				return ErrBadArgs
			}
		}
		shStr := fmt.Sprintf("Your ship landed at %s with %d colonists, but were all killed by %s (%d inhabitants remaining).", pl.Name(), aSum, themStr, dSum-aSum)
		if !source.AddReportEvent(shFid, shStr) {
			return ErrBadArgs
		}
		return
	}
	pl.SetController(sh.Fid())
	pl.SetInhabitants(0)
	atk -= def
	arrivals[plid] = atk
	if dSum == 0 {
		var shStr string
		if plFid == 0 {
			shStr = fmt.Sprintf("Your ship landed at %s with %d colonists and found no resistance from the local populace.", pl.Name(), aSum)
		} else {
			shStr = fmt.Sprintf("Your ship landed at %s with %d colonists and found no remaining %s inhabitants left to defend it", pl.Name(), aSum, names[plFid])
			plStr := fmt.Sprintf("%s was invaded by %d colonists from %s, and you had no inhabitants present to retain control of the planet.", pl.Name(), aSum, names[shFid])
			if !source.AddReportEvent(plFid, plStr) {
				return ErrBadArgs
			}
		}
		if !source.AddReportEvent(shFid, shStr) {
			return ErrBadArgs
		}
	} else {
		var themStr string
		if plFid == 0 {
			themStr = "hostile natives"
		} else {
			themStr = names[plFid] + " inhabitants"
			plStr := fmt.Sprintf("%s was invaded by %d colonists from %s, who killed your %d inhabitants there and took control (%d invaders remaining)", pl.Name(), aSum, names[shFid], dSum, aSum-dSum)
			if !source.AddReportEvent(plFid, plStr) {
				return ErrBadArgs
			}
		}
		shStr := fmt.Sprintf("Your ship landed at %s with %d colonists and were attacked by %d %s, but your colonists managed to defend themselves and take the planet (%d colonists survived).", pl.Name(), aSum, dSum, themStr, aSum-dSum)
		if !source.AddReportEvent(shFid, shStr) {
			return ErrBadArgs
		}
	}
	return
}

func BothSee(source Source, pl Planet, fid1, fid2, turn int, arrivals map[int]int) {
	for _, fid := range []int{fid1, fid2} {
		if fid == 0 || fid == pl.Controller() {
			continue
		}
		pv := source.UpdatePlanetView(fid, turn, pl)
		arv := arrivals[pl.Pid()]
		if arv > 0 {
			pv.SetInhabitants(pv.Inhabitants() + arv)
		}
	}
}
