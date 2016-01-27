package main

import (
	"mule/hexagon"
	"mule/myweb"
	"mule/overpower"
	"mule/overpower/db"
	"mule/overpower/mapping"
	"net/http"
)

// /overpower/command/GAMEID/TURN/ACTIONNAME/ARGS
func (h *Handler) Recenter(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction) (isOk bool) {
	x, okX := h.IntAt(6)
	y, okY := h.IntAt(7)
	if !okY || !okX {
		http.Error(w, "BAD COORD ARGS IN RECENTER PATH", http.StatusBadRequest)
		return
	}
	err := OPDB.UpdateMapView("center", hexagon.Coord{x, y}, "WHERE", "gid", g.Gid(), "fid", f.Fid())
	if my, bad := Check(err, "recenter failure", "coord", hexagon.Coord{x, y}, "gid", g.Gid(), "fid", f.Fid()); bad {
		Log(my)
		http.Error(w, my.Error(), http.StatusInternalServerError)
		return false

	}
	return true
}

func (h *Handler) MapClick(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction) (isOk bool) {
	if f == nil {
		http.Error(w, "USER HAS NO FACTION IN THIS GAME", http.StatusBadRequest)
		return
	}
	if g.Turn() < 1 {
		http.Error(w, "GAME NOT IN PROGRESS", http.StatusBadRequest)
		return
	}
	ints, ok := GetInts(r, "turn", "button", "clickx", "clicky")
	if !ok {
		http.Error(w, "BAD FORM DATA", http.StatusBadRequest)
		return
	}
	turn, button, clickx, clicky := ints[0], ints[1], ints[2], ints[3]
	if turn != g.Turn() {
		http.Error(w, "FORM TURN DOES NOT MATCH GAME TURN", http.StatusBadRequest)
		return
	}
	gid, fid := f.Gid(), f.Fid()
	mv, err := OPDB.GetMapView("gid", gid, "fid", fid)
	if my, bad := Check(err, "mapclick failure", "gid", gid, "fid", fid); bad {
		return Bail(w, my)
	}
	shift := r.FormValue("shift") == "true"
	if button == 3 { // wheelspin
		if clicky == 0 {
			http.Error(w, "BAD CLICKY FORM DATA", http.StatusBadRequest)
			return
		} else if clicky > 0 {
			if zoom := mv.Zoom(); zoom < 50 {
				var amount int
				if shift {
					if zoom+10 > 50 {
						amount = 50
					} else {
						amount = zoom + 10
					}
				} else {
					amount = zoom + 1
				}
				mv.SetZoom(amount)
				err = OPDB.UpdateMapViews(mv)
				if my, bad := Check(err, "mapclick update failure", "mv", mv); bad {
					return Bail(w, my)
				}
				return true
			} else {
				return true
			}
		} else {
			if zoom := mv.Zoom(); zoom > 1 {
				var amount int
				if shift {
					if zoom-10 < 1 {
						amount = 1
					} else {
						amount = zoom - 10
					}
				} else {
					amount = zoom - 1
				}
				mv.SetZoom(amount)
				err = OPDB.UpdateMapViews(mv)
				if my, bad := Check(err, "mapclick update failure", "mv", mv); bad {
					return Bail(w, my)
				}
				return true
			} else {
				return true
			}
		}
		return true
	}
	vp := mapping.GetVP(mv)
	c := vp.HexContaining(hexagon.Pixel{float64(clickx), float64(clicky)})
	var pvList []overpower.PlanetView
	var found bool
	pvList, err = OPDB.GetPlanetViews("gid", gid, "fid", fid)
	if my, bad := Check(err, "mapclick getplanets failure", "gid", gid, "fid", fid); bad {
		return Bail(w, my)
	}
	var foundNear bool
	var near hexagon.Coord
	if mv.Zoom() < 15 {
		for _, pv := range pvList {
			if loc := pv.Loc(); loc == c {
				found = true
				break
			} else if !foundNear && loc.StepsTo(c) < 2 {
				near = loc
				foundNear = true
			}
		}
		if !found && foundNear {
			c = near
			found = true
		}
	} else {
		for _, pv := range pvList {
			if loc := pv.Loc(); loc == c {
				found = true
				break
			}
		}
	}
	switch button {
	case 0: // left
		if shift {
			mv.SetCenter(c)
		} else {
			mv.SetTarget1(c)
		}
	case 1: // wheelclick
		mv.SetCenter(c)
	case 2: // right
		if shift {
			mv.DropTarget2()
		} else {
			mv.SetTarget2(c)
		}
	case 3: // wheelspin
	// SEE ABOVE
	default:
		http.Error(w, "UNKNOWN BUTTON VALUE", http.StatusBadRequest)
		return
	}
	err = OPDB.UpdateMapViews(mv)
	if my, bad := Check(err, "mapclick update failure", "mv", mv); bad {
		return Bail(w, my)
	}
	return true
}

func (h *Handler) SetMapView(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction) (isOk bool) {
	mp := myweb.GetIntsIf(r, "turn", "zoom", "hexx", "hexy")
	if t, ok := mp["turn"]; !ok || t != g.Turn() {
		http.Error(w, "FORM TURN DOES NOT MATCH GAME TURN", http.StatusBadRequest)
		return
	}
	gid, fid := f.Gid(), f.Fid()
	subtype := r.FormValue("subtype")
	zoom, zOk := mp["zoom"]
	hexX, xOk := mp["hexx"]
	hexY, yOk := mp["hexy"]
	isOk = true
	mv, err := OPDB.GetMapView("gid", gid, "fid", fid)
	if my, bad := Check(err, "set mapview get item failure", "gid", gid, "fid", fid); bad {
		return Bail(w, my)
	}
	switch subtype {
	case "swap":
		tar1, tar2 := mv.Target1(), mv.Target2()
		if !tar1.Valid {
			mv.DropTarget2()
		} else {
			mv.SetTarget2(tar1.Coord)
		}
		if !tar2.Valid {
			mv.DropTarget1()
		} else {
			mv.SetTarget1(tar2.Coord)
		}
	case "zoom":
		if !zOk {
			http.Error(w, "BAD DATA GIVEN FOR MAPVIEW CHANGE", http.StatusBadRequest)
			return false
		}
		if zoom < 1 {
			zoom = 1
		} else if zoom > 50 {
			zoom = 50
		}
		mv.SetZoom(zoom)
	case "target1":
		if !xOk || !yOk {
			http.Error(w, "BAD DATA GIVEN FOR MAPVIEW CHANGE", http.StatusBadRequest)
			return false
		}
		mv.SetTarget1(hexagon.Coord{hexX, hexY})
	case "target2":
		if !xOk || !yOk {
			http.Error(w, "BAD DATA GIVEN FOR MAPVIEW CHANGE", http.StatusBadRequest)
			return false
		}
		mv.SetTarget2(hexagon.Coord{hexX, hexY})
	case "center":
		if !xOk || !yOk {
			http.Error(w, "BAD DATA GIVEN FOR MAPVIEW CHANGE", http.StatusBadRequest)
			return false
		}
		mv.SetCenter(hexagon.Coord{hexX, hexY})
	default:
		http.Error(w, "UNKNOWN ACTION TYPE", http.StatusBadRequest)
		return false
	}

	if !isOk {
		http.Error(w, "DATABASE ERROR EXECUTING COMMAND", http.StatusInternalServerError)
		return false
	}
	err = OPDB.UpdateMapViews(mv)
	if my, bad := Check(err, "mapview change update failure", "mv", mv); bad {
		return Bail(w, my)
	}
	return true
}

func (h *Handler) SetOrder(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction) (isOk bool) {
	if f == nil {
		http.Error(w, "USER HAS NO FACTION IN THIS GAME", http.StatusBadRequest)
		return
	}
	if g.Turn() < 1 {
		http.Error(w, "GAME NOT IN PROGRESS", http.StatusBadRequest)
		return
	}
	ints, ok := GetInts(r, "turn", "source", "target", "size")
	if !ok {
		http.Error(w, "BAD FORM DATA", http.StatusBadRequest)
		return
	}
	turn, source, target, size := ints[0], ints[1], ints[2], ints[3]
	if turn != g.Turn() {
		http.Error(w, "FORM TURN DOES NOT MATCH GAME TURN", http.StatusBadRequest)
		return
	}
	planets, err := OPDB.GetPlanetsByPlid(g.Gid(), source, target)
	if my, bad := Check(err, "setorder resource get failure", "gid", g.Gid(), "source", source, "target", target); bad {
		return Bail(w, my)
	}
	if len(planets) != 2 {
		http.Error(w, "SOURCE/TARGET PLANETS NOT FOUND", http.StatusBadRequest)
		return
	}
	if planets[0].Pid() == planets[1].Pid() {
		http.Error(w, "INVALID TARGET: SOURCE", http.StatusBadRequest)
		return
	}
	var sPl overpower.Planet
	if planets[0].Pid() == source {
		sPl = planets[0]
	} else if planets[1].Pid() == source {
		sPl = planets[1]
	} else {
		http.Error(w, "SOURCE PLANET NOT FOUND", http.StatusBadRequest)
		return
	}
	if sPl.Controller() != f.Fid() {
		http.Error(w, "SOURCE PLANET NOT OWNED BY FACTION", http.StatusBadRequest)
		return
	}

	orders, err := OPDB.GetOrders("gid", g.Gid(), "source", source)
	if my, bad := Check(err, "setorder resource get failure", "gid", g.Gid(), "source", source); bad {
		return Bail(w, my)
	}
	var using int
	var curOrder overpower.Order
	for _, o := range orders {
		if o.Target() == target {
			curOrder = o
		} else {
			using += o.Size()
		}
	}
	if size > 0 {
		have := sPl.Parts()
		if using+size > have {
			http.Error(w, "PLANET HAS INSUFFICIENT PARTS FOR ORDER", http.StatusBadRequest)
			return
		}
	}
	if curOrder == nil {
		err = OPDB.MakeOrder(g.Gid(), f.Fid(), source, target, size)
	} else {
		curOrder.SetSize(size)
		err = OPDB.UpdateOrders(curOrder)
	}
	if my, bad := Check(err, "setorder update fail", "source", source, "target", target, "size", size); bad {
		return Bail(w, my)
	}
	return true
}

func (h *Handler) SetTurnDone(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction) (isOk bool) {
	if f == nil {
		http.Error(w, "USER HAS NO FACTION IN THIS GAME", http.StatusBadRequest)
		return
	}
	if g.Turn() < 1 {
		http.Error(w, "GAME NOT IN PROGRESS", http.StatusBadRequest)
		return
	}
	facs, err := OPDB.GetFactions("gid", g.Gid())
	if my, bad := Check(err, "set turn done resource failure", "gid", g.Gid()); bad {
		return Bail(w, my)
	}
	ints, ok := GetInts(r, "turn")
	if !ok {
		http.Error(w, "MALFORMED TURN DATA", http.StatusBadRequest)
		return
	}
	turn := ints[0]
	if turn != g.Turn() {
		http.Error(w, "BAD TURN DATA", http.StatusBadRequest)
		return
	}
	done := r.FormValue("done") == "true"
	fid := f.Fid()
	var allDone bool
	if done {
		allDone = true
		for _, f := range facs {
			if f.Fid() != fid && !f.Done() {
				allDone = false
				break
			}
		}
	}
	if allDone {
		err = RunGameTurn(g.Gid(), false)
		if my, bad := Check(err, "update problem in faction set done (run turn)", "faction", f); bad {
			return Bail(w, my)
		}
		return true
	} else {
		f.SetDone(done)
		err = OPDB.UpdateFactions(f)
		if my, bad := Check(err, "update problem in faction set done", "faction", f); bad {
			return Bail(w, my)
		}
		return true
	}
}

func RunGameTurn(gid int, auto bool) error {
	f := func(d db.DB) error {
		source := d.NewSource(gid)
		breakE, logE := overpower.RunGameTurn(source, auto)
		if logE != nil {
			Log(logE)
		}
		if my, bad := Check(breakE, "rungameturn failed", "gid", gid); bad {
			return my
		}
		err := source.Commit()
		if my, bad := Check(err, "rungame turn commit failed", "gid", gid); bad {
			return my
		}
		return nil
	}
	return OPDB.Transact(f)
}

func StartGame(gid int, exodus bool) error {
	f := func(d db.DB) error {
		source := d.NewSource(gid)
		err := overpower.MakeGalaxy(source, exodus)
		if my, bad := Check(err, "Start game failure", "gid", gid, "exodus", exodus); bad {
			return my
		}
		err = source.Commit()
		if my, bad := Check(err, "Start game source commit failure", "gid", gid); bad {
			return my
		}
		return nil
	}
	return OPDB.Transact(f)
}
