package overpower

import (
	"fmt"
	"mule/hexagon"
)

func (op *TotallyOP) RunGameTurn() (ok bool) {
	planets := make(map[hexagon.Coord]Planet, len(op.Planets))
	plids := make(map[int]Planet, len(op.Planets))
	radar := make(map[int][]hexagon.Coord, len(op.Factions))
	for _, p := range op.Planets {
		planets[p.Loc()] = p
		plids[p.Pid()] = p
		if fid := p.Controller(); fid != 0 {
			loc := p.Loc()
			if list, ok := radar[fid]; ok {
				radar[fid] = append(list, loc)
			} else {
				radar[fid] = []hexagon.Coord{loc}
			}
		}
	}
	gid, turn := op.Game.Gid(), op.Game.Turn()
	names := make(map[int]string, len(op.Factions))
	op.Reports = make(map[int]Report, len(op.Factions))
	for _, f := range op.Factions {
		f.SetDone(false)
		fid := f.Fid()
		names[fid] = "faction " + f.Name()
		rp, ok := op.Source.NewReport(gid, fid, turn)
		if !ok {
			return false
		}
		op.Reports[fid] = rp
	}
	// ---- SHIPS LAUNCH ---- //
	for _, o := range op.Orders {
		tar, ok1 := plids[o.Target()]
		src, ok2 := plids[o.Source()]
		if !(ok1 && ok2) {
			// Log("BAD ORDER, PLANETS NOT FOUND:", o)
			continue
		}
		size := o.Size()
		if size < 1 {
			continue
		}
		if cont := src.Controller(); cont == 0 || src.Parts() < size || cont != o.Fid() {
			// Log("BAD ORDER: PLANET NOT SUITABLE FOR LAUNCH", src, o)
			continue
		}
		src.SetParts(src.Parts() - size)
		path := src.Loc().PathTo(tar.Loc())
		sh, ok := op.Source.NewShip(gid, src.Controller(), size, turn, path)
		rStr := fmt.Sprintf("%s launched size %d ship toward %s", src.Name(), size, tar.Name())
		op.AddReport(src.Controller(), rStr)
		if !ok {
			return false
		}
		op.Ships = append(op.Ships, sh)
	}
	// ---- SHIPS MOVE ---- //
	// dist, ship index
	landings := map[int][]int{}
	for i, sh := range op.Ships {
		travelled, land := Travelled(sh, turn)
		if len(travelled) < 1 {
			// Log("BAD SHIP: TRAVELLED NO DIST", sh)
			return false
		}
		at := travelled[len(travelled)-1]
		// ----- SHIP MOVEMENT IS SEEN ------ //
		for fid, rList := range radar {
			var destValid, spottedShip bool
			var spotted []hexagon.Coord
			if fid == sh.Fid() {
				spotted, spottedShip = travelled, true
				destValid = true
			} else {
				spotted, spottedShip = RadarCheck(rList, travelled)
			}
			if len(spotted) > 0 {
				var trail []hexagon.Coord
				var loc, dest hexagon.Coord
				locValid := spottedShip && !land
				if locValid {
					loc = at
					trail = spotted[:len(spotted)-1]
				} else {
					trail = spotted
				}
				if destValid {
					path := sh.Path()
					dest = path[len(path)-1]
				}
				sv, ok := op.Source.NewShipView(gid, fid, turn, sh.Sid(), sh.Fid(), sh.Size(), loc, locValid, dest, destValid, trail)
				if !ok {
					return false
				}
				op.ShipViews = append(op.ShipViews, sv)
			}
		}
		// ---- LANDINGS TAGGED FOR LATER ------ //
		if land {
			dist := len(travelled) - 1
			if list, ok := landings[dist]; ok {
				landings[dist] = append(list, i)
			} else {
				landings[dist] = []int{i}
			}
		}
	}
	//
	// ---- SHIPS LAND ---- //
	// plid, amount
	arrivals := map[int]int{}
	for i := 1; i < SHIPSPEED+1; i++ {
		ships, ok := landings[i]
		if !ok {
			continue
		}
		for _, sI := range shuffleInts(ships) {
			sh := op.Ships[sI]
			path := sh.Path()
			loc := path[len(path)-1]
			p, ok := planets[loc]
			if !ok {
				// Log("TRIED TO LAND AT NONEXISTANT PLANET", sh.Loc()
				return false
			}
			if !op.PlanetaryLanding(p, sh, turn, arrivals, names) {
				return false
			}
			if !op.Source.DropShip(sh) {
				return false
			}
		}
		delete(landings, i)
	}
	if len(landings) > 0 {
		// Log("LANDINGS WITH IMPROPER DISTANCE:", landings)
		return false
	}
	//
	// ---- TURN STARTS ---- //
	op.Game.IncTurn()
	turn = op.Game.Turn()
	// ---- PLANETS PRODUCE ---- //
	for _, pl := range op.Planets {
		cont := pl.Controller()
		if cont == 0 {
			continue
		}
		if inh, res := pl.Inhabitants(), pl.Resources(); inh > 0 && res > 0 {
			var prod int
			switch {
			case inh > 5 && res > 5:
				prod = 5
			case inh > res:
				prod = res
			default:
				prod = inh
			}
			pl.SetResources(res - prod)
			pl.SetParts(pl.Parts() + prod)
		}
		// ---- ARRIVALS ARRIVE ---- //
		if x := arrivals[pl.Pid()]; x > 0 {
			pl.SetInhabitants(pl.Inhabitants() + x)
		}
		// ---- PLANETS ARE SEEN ---- //
		pv, ok := op.Source.NewPlanetView(cont, turn, pl)
		if !ok {
			return false
		}
		op.SetPV(pv)
	}
	return true
}
