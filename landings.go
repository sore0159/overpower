package overpower

import (
	"fmt"
)

func (op *TotallyOP) PlanetaryLanding(pl Planet, sh Ship, turn int, arrivals map[int]int, names map[int]string) (ok bool) {
	shFid := sh.Fid()
	plFid := pl.Controller()
	plid := pl.Pid()
	atk := sh.Size()
	def := arrivals[plid]
	aSum := atk
	dSum := def + pl.Inhabitants()
	if shFid == plFid {
		rStr := fmt.Sprintf("Your ship landed at %s with %d colonists, reinforcing inhabitants to %d.", pl.Name(), aSum, dSum+aSum)
		op.AddReport(shFid, rStr)
		arrivals[plid] += atk
		return true
	}
	defer op.BothSee(pl, pl.Controller(), sh.Fid(), turn, arrivals)
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
			op.AddReport(plFid, plStr)
		}
		shStr := fmt.Sprintf("Your ship landed at %s with %d colonists, but were all killed by %s (%d inhabitants remaining).", pl.Name(), aSum, themStr, dSum-aSum)
		op.AddReport(shFid, shStr)
		return true
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
			op.AddReport(plFid, plStr)
		}
		shStr := fmt.Sprintf("Your ship landed at %s with %d colonists, but were all killed by %s (%d inhabitants remaining).", pl.Name(), aSum, themStr, dSum-aSum)
		op.AddReport(shFid, shStr)
		return true
	}
	pl.SetController(sh.Fid())
	pl.SetInhabitants(0)
	atk -= def
	arrivals[plid] = atk
	pl.SetInhabitants(def - atk)
	if dSum == 0 {
		var shStr string
		if plFid == 0 {
			shStr = fmt.Sprintf("Your ship landed at %s with %d colonists and found no resistance from the local populace.", pl.Name(), aSum)
		} else {
			shStr = fmt.Sprintf("Your ship landed at %s with %d colonists and found no remaining %s inhabitants left to defend it", pl.Name(), aSum, names[plFid])
			plStr := fmt.Sprintf("%s was invaded by %d colonists from %s, and you had no inhabitants present to retain control of the planet.", pl.Name(), aSum, names[shFid])
			op.AddReport(plFid, plStr)
		}
		op.AddReport(shFid, shStr)
	} else {
		var themStr string
		if plFid == 0 {
			themStr = "hostile natives"
		} else {
			themStr = names[plFid] + " inhabitants"
			plStr := fmt.Sprintf("%s was invaded by %d colonists from %s, who killed your %d inhabitants there and took control (%d invaders remaining)", pl.Name(), aSum, names[shFid], dSum, aSum-dSum)
			op.AddReport(plFid, plStr)
		}
		shStr := fmt.Sprintf("Your ship landed at %s with %d colonists and were attacked by %d %s, but your colonists managed to defend themselves and take the planet (%d colonists survived).", pl.Name(), aSum, dSum, themStr, aSum-dSum)
		op.AddReport(shFid, shStr)
	}
	return true
}

func (op *TotallyOP) BothSee(pl Planet, fid1, fid2, turn int, arrivals map[int]int) (ok bool) {
	for _, fid := range []int{fid1, fid2} {
		if fid == 0 || fid == pl.Controller() {
			continue
		}
		pv, ok := op.Source.NewPlanetView(fid, turn, pl)
		if !ok {
			return false
		}
		arv := arrivals[pl.Pid()]
		if arv > 0 {
			pv.SetInhabitants(pv.Inhabitants() + arv)
		}
		op.SetPV(pv)
	}
	return true
}
