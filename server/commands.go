package main

import (
	"mule/hexagon"
	"mule/myweb"
	"mule/overpower"
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
	if !OPDB.UpdateMapViewCenter(g.Gid(), f.Fid(), hexagon.Coord{x, y}) {
		http.Error(w, "DATABASE ERROR UPDATING MAPVIEW", http.StatusInternalServerError)
		return
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
	mv, ok := OPDB.GetFidMapView(gid, fid)
	if !ok {
		http.Error(w, "DATABASE ERROR RETRIEVING MAPVIEW", http.StatusInternalServerError)
		return
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
				if !OPDB.UpdateMapViewZoom(gid, fid, amount) {
					http.Error(w, "DATABASE ERROR UPDATING MAPVIEW", http.StatusInternalServerError)
					return false
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

				if !OPDB.UpdateMapViewZoom(gid, fid, amount) {
					http.Error(w, "DATABASE ERROR UPDATING MAPVIEW", http.StatusInternalServerError)
					return false
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
	pvList, ok = OPDB.GetAllFactionPlanetViews(gid, fid)
	if !ok {
		http.Error(w, "DATABASE ERROR RETRIEVING PLANETVIEWS", http.StatusInternalServerError)
		return
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
			isOk = OPDB.UpdateMapViewCenter(gid, fid, c)
		} else {
			isOk = OPDB.UpdateMapViewTarget(gid, fid, true, hexagon.NullCoord{c, true})
		}
	case 1: // wheelclick
		isOk = OPDB.UpdateMapViewCenter(gid, fid, c)
	case 2: // right
		if shift {
			isOk = OPDB.UpdateMapViewTarget(gid, fid, false, hexagon.NullCoord{})
		} else {
			isOk = OPDB.UpdateMapViewTarget(gid, fid, false, hexagon.NullCoord{c, true})
		}
	case 3: // wheelspin
	// SEE ABOVE
	default:
		http.Error(w, "UNKNOWN BUTTON VALUE", http.StatusBadRequest)
		return
	}
	if !isOk {
		http.Error(w, "DATABASE ERROR UPDATING MAPVIEW", http.StatusInternalServerError)
		return false
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
	switch subtype {
	case "swap":
		mv, ok := OPDB.GetFidMapView(gid, fid)
		if !ok {
			http.Error(w, "DATABASE ERROR GETTING MAPVIEW", http.StatusInternalServerError)
			return false
		}
		tar1, tar2 := mv.Target1(), mv.Target2()
		if !OPDB.UpdateMapViewBothTargets(gid, fid, tar2, tar1) {
			http.Error(w, "DATABASE ERROR SETTING MAPVIEW", http.StatusInternalServerError)
			return false
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
		isOk = OPDB.UpdateMapViewZoom(gid, fid, zoom)
	case "target1":
		if !xOk || !yOk {
			http.Error(w, "BAD DATA GIVEN FOR MAPVIEW CHANGE", http.StatusBadRequest)
			return false
		}
		isOk = OPDB.UpdateMapViewTarget(gid, fid, true, hexagon.NullCoord{hexagon.Coord{hexX, hexY}, true})
	case "target2":
		if !xOk || !yOk {
			http.Error(w, "BAD DATA GIVEN FOR MAPVIEW CHANGE", http.StatusBadRequest)
			return false
		}
		isOk = OPDB.UpdateMapViewTarget(gid, fid, false, hexagon.NullCoord{hexagon.Coord{hexX, hexY}, true})
	case "center":
		if !xOk || !yOk {
			http.Error(w, "BAD DATA GIVEN FOR MAPVIEW CHANGE", http.StatusBadRequest)
			return false
		}
		isOk = OPDB.UpdateMapViewCenter(gid, fid, hexagon.Coord{hexX, hexY})
	default:
		http.Error(w, "UNKNOWN ACTION TYPE", http.StatusBadRequest)
		return false
	}
	if !isOk {
		http.Error(w, "DATABASE ERROR EXECUTING COMMAND", http.StatusInternalServerError)
		return false
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
	planets, ok := OPDB.GetPidPlanets(g.Gid(), source, target)
	if !ok {

		http.Error(w, "DATABASE ERROR FETCHING PLANETS", http.StatusInternalServerError)
		return
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
	if size > 0 {
		have := sPl.Parts()
		if size > have {
			http.Error(w, "PLANET HAS INSUFFICIENT PARTS FOR ORDER", http.StatusBadRequest)
			return
		}
		orders, ok := OPDB.GetAllSourceOrders(g.Gid(), source)
		if !ok {
			http.Error(w, "DATABASE ERROR FETCHING SOURCE ORDERS", http.StatusInternalServerError)
			return
		}
		want := size
		for _, o := range orders {
			if o.Target() != target {
				want += o.Size()
			}
		}
		if want > have {
			http.Error(w, "PLANET HAS INSUFFICIENT PARTS FOR ORDER", http.StatusBadRequest)
			return
		}
	}
	if !OPDB.SetOrder(g.Gid(), f.Fid(), source, target, size) {
		http.Error(w, "DATABASE ERROR SETTING ORDER", http.StatusInternalServerError)
		return
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
	facs, ok := OPDB.GetGidFactions(g.Gid())
	if !ok {
		http.Error(w, "DATABASE ERROR RETRIEVING FACTIONS", http.StatusInternalServerError)
		return false
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
		if !OPDB.RunGameTurn(g) {
			http.Error(w, "DATABASE ERROR RUNNING GAME TURN", http.StatusInternalServerError)
			return
		}
	} else {
		if !OPDB.SetTurnDone(f, done) {
			http.Error(w, "DATABASE ERROR SETTING FACTION TURN DONE", http.StatusInternalServerError)
			return
		}
	}
	return true
}
