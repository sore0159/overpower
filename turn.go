package overpower

import (
	"fmt"
	"math"
	"mule/hexagon"
)

func RunGameTurn(source Source, auto bool) (breaker, logger error) {
	game, err := source.Game()
	if my, bad := Check(err, "run turn resource failure"); bad {
		return my, nil
	}
	planets, err := source.Planets()
	if my, bad := Check(err, "run turn resource failure"); bad {
		return my, nil
	}
	factions, err := source.Factions()
	if my, bad := Check(err, "run turn resource failure"); bad {
		return my, nil
	}
	orders, err := source.Orders()
	if my, bad := Check(err, "run turn resource failure"); bad {
		return my, nil
	}
	ships, err := source.Ships()
	if my, bad := Check(err, "run turn resource failure"); bad {
		return my, nil
	}
	// -------------------------------- //
	//source.DropShipViews()
	var errOccured bool
	loggerM, _ := Check(ErrIgnorable, "run turn problem")
	planetGrid := make(map[hexagon.Coord]Planet, len(planets))
	plids := make(map[int]Planet, len(planets))
	radar := make(map[int]hexagon.CoordList, len(factions))
	for _, p := range planets {
		planetGrid[p.Loc()] = p
		plids[p.Pid()] = p
		if fid := p.Controller(); fid != 0 {
			loc := p.Loc()
			if list, ok := radar[fid]; ok {
				radar[fid] = append(list, loc)
			} else {
				radar[fid] = hexagon.CoordList{loc}
			}
		}
	}
	turn := game.Turn()
	names := make(map[int]string, len(factions))
	reports := make(map[int]Report, len(factions))
	for _, f := range factions {
		f.SetDone(false)
		fid := f.Fid()
		names[fid] = "faction " + f.Name()
		reports[fid] = source.NewReport(fid, turn)
	}
	// ------ AUTO TURN ------- //
	if !auto {
		game.SetFreeAutos(game.FreeAutos() + 1)
	}
	// ---- SHIPS LAUNCH ---- //
	for _, o := range orders {
		tar, ok1 := plids[o.Target()]
		src, ok2 := plids[o.Source()]
		if !(ok1 && ok2) {
			errOccured = true
			loggerM.AddContext("bad order", "planets not found", "order", o)
			continue
		}
		size := o.Size()
		if size < 1 {
			errOccured = true
			loggerM.AddContext("bad order", "size <0", "order", o)
			continue
		}
		if cont := src.Controller(); cont == 0 || src.Parts() < size || cont != o.Fid() {
			errOccured = true
			loggerM.AddContext("bad order", "misc", "order", o)
			continue
		}
		src.SetParts(src.Parts() - size)
		path := src.Loc().PathTo(tar.Loc())
		sh := source.NewShip(src.Controller(), size, turn, path)
		ships = append(ships, sh)
		rStr := fmt.Sprintf("%s launched size %d ship toward %s", src.Name(), size, tar.Name())
		if !source.AddReportEvent(src.Controller(), rStr) {
			loggerM.AddContext("report problem", "couldn't find report", "fid", src.Controller(), "report", rStr)
			errOccured = true
			return
		}
	}
	source.DropOrders()
	// ---- SHIPS MOVE ---- //
	// dist, ship index
	landings := map[int][]int{}
	for i, sh := range ships {
		travelled, land := Travelled(sh, turn)
		if len(travelled) < 1 {
			errOccured = true
			loggerM.AddContext("bad ship", "no travel dist", "ship", sh)
			source.DropShip(sh)
			continue
		}
		at := travelled[len(travelled)-1]
		// ----- SHIP MOVEMENT IS SEEN ------ //
		for fid, rList := range radar {
			var destValid, spottedShip bool
			var spotted hexagon.CoordList
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
				var locNC, destNC hexagon.NullCoord
				locNC.Valid = locValid
				locNC.Coord = loc
				destNC.Valid = destValid
				destNC.Coord = dest
				source.NewShipView(sh, fid, turn, locNC, destNC, trail)
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
		shipsLandings, ok := landings[i]
		if !ok {
			continue
		}
		for _, sI := range shuffleInts(shipsLandings) {
			sh := ships[sI]
			path := sh.Path()
			loc := path[len(path)-1]
			p, ok := planetGrid[loc]
			if !ok {
				loggerM.AddContext("bad ship", "landing nonexistant", "ship", sh)
				errOccured = true
			} else {
				err := PlanetaryLanding(source, p, sh, turn, arrivals, names)
				if my, bad := Check(err, "landing report problem", "planet", p, "ship", sh); bad {
					loggerM.Grab(my)
					errOccured = true
				}
			}
			source.DropShip(sh)
		}
		delete(landings, i)
	}
	if len(landings) > 0 {
		loggerM.AddContext("bad ships", "landings with improper dist", "landings", landings)
		errOccured = true
	}
	//
	// ---- TURN STARTS ---- //
	game.IncTurn()
	turn = game.Turn()
	facScores := make(map[int]int, len(factions))
	// ---- PLANETS PRODUCE ---- //
	for _, pl := range planets {
		cont := pl.Controller()
		if cont == 0 {
			continue
		}
		facScores[cont] += 1
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
		source.UpdatePlanetView(cont, turn, pl)
	}
	var highScore int
	winPercent := game.WinPercent()
	winners := make([]Faction, 0)
	for _, f := range factions {
		score := facScores[f.Fid()]
		if score > highScore {
			highScore = score
		}
		scorePercent := int(math.Floor(100 * float64(score) / float64(len(planets))))
		f.SetScore(scorePercent)
		if scorePercent >= winPercent {
			winners = append(winners, f)
		}
	}
	percent := 100 * float64(highScore) / float64(len(planets))
	game.SetHighScore(int(math.Floor(percent)))
	if len(winners) > 0 {
		Ping("TODO: WINNING!", winners)
	}
	if errOccured {
		return nil, loggerM
	} else {
		return nil, nil
	}
}
