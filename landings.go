package overpower

import (
	"fmt"
)

func (op *TotallyOP) PlanetaryLanding(pl Planet, sh Ship, turn int, arrivals map[int]int, names map[int]string) (ok bool) {
	fmt.Println("LANDING ON", pl, "BY", sh)
	defer op.BothSee(pl, pl.Controller(), sh.Fid(), turn, arrivals)
	plid := pl.Pid()
	def := arrivals[plid]
	atk := sh.Size()
	aSum := atk
	dSum := def + pl.Inhabitants()
	_, _ = aSum, dSum
	if def >= atk {
		if def == atk {
			delete(arrivals, plid)
		} else {
			arrivals[plid] = def - atk
		}
		fmt.Println("NATIVE ARRIVALS DEFEAT INVADERS", aSum, dSum)
		return true
	}
	delete(arrivals, plid)
	atk -= def
	def = pl.Inhabitants()
	if def >= atk {
		pl.SetInhabitants(def - atk)
		fmt.Println("NATIVES DEFEAT INVADERS", aSum, dSum)
		return true
	}
	pl.SetController(sh.Fid())
	pl.SetInhabitants(0)
	atk -= def
	arrivals[plid] = atk
	fmt.Println("INVADERS DEFEAT NATIVES", aSum, dSum)
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
