package main

import (
	"mule/hexagon"
	"mule/myweb"
	"mule/overpower"
	"mule/overpower/mapping"
	"net/http"
)

func (h *Handler) MapClick(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction) (isOk bool) {
	if f == nil {
		http.Error(w, "USER HAS NO FACTION IN THIS GAME", http.StatusBadRequest)
		return
	}
	if g.Turn() < 1 {
		http.Error(w, "GAME NOT IN PROGRESS", http.StatusBadRequest)
		return
	}
	mp, ok := myweb.GetIntsQuiet(r, "turn", "button", "clickx", "clicky")
	if !ok {
		http.Error(w, "BAD FORM DATA", http.StatusBadRequest)
		return
	}
	turn, button, clickx, clicky := mp["turn"], mp["button"], mp["clickx"], mp["clicky"]
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
	vp := mapping.GetVP(mv)
	c := vp.HexContaining(hexagon.Pixel{float64(clickx), float64(clicky)})
	if button == 0 {
		return OPDB.UpdateMapViewCenter(gid, fid, c)
	}
	if button == 1 {
		if shift {
			if zoom := mv.Zoom(); zoom < 100 {
				return OPDB.UpdateMapViewZoom(gid, fid, zoom+1)
			} else {
				return true
			}
		} else {
			if zoom := mv.Zoom(); zoom > 1 {
				return OPDB.UpdateMapViewZoom(gid, fid, zoom-1)
			} else {
				return true
			}
		}
	}
	if button == 2 {
		pvList, ok := OPDB.GetAllFactionPlanetViews(gid, fid)
		if !ok {
			http.Error(w, "DATABASE ERROR RETRIEVING PLANETVIEWS", http.StatusInternalServerError)
			return
		}
		var found bool
		for _, pv := range pvList {
			if pv.Loc() == c {
				found = true
				break
			}
		}
		if found {
			return OPDB.UpdateMapViewFocus(gid, fid, true, c)
		} else if shift {
			return OPDB.UpdateMapViewFocus(gid, fid, false, hexagon.Coord{})
		} else {
			return true
		}
	}
	return true
}

func (h *Handler) SetMapView(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction) (isOk bool) {
	mp := myweb.GetIntsIf(r, "turn", "zoom", "hexx", "hexy")
	if t, ok := mp["turn"]; !ok || t != g.Turn() {
		http.Error(w, "FORM TURN DOES NOT MATCH GAME TURN", http.StatusBadRequest)
		return
	}
	zoom, zOk := mp["zoom"]
	hexX, xOk := mp["hexx"]
	hexY, yOk := mp["hexy"]
	if !zOk && !xOk && !yOk {
		http.Error(w, "NO DATA GIVEN FOR MAPVIEW CHANGE", http.StatusBadRequest)
		return false
	}
	if zOk && (xOk || yOk) {
		http.Error(w, "CANNOT CHANGE MAPVIEW ZOOM AND CENTER SIMULTANIOUSLY", http.StatusBadRequest)
		return false
	}
	if zOk {
		if zoom < 1 {
			zoom = 1
			//	http.Error(w, "BAD VALUE FOR MAPVIEW ZOOM", http.StatusBadRequest)
			//	return false
		} else if zoom > 100 {
			zoom = 100
		}
		return OPDB.UpdateMapViewZoom(f.Gid(), f.Fid(), zoom)
	}
	if !(xOk && yOk) {
		http.Error(w, "INCOMPLETE DATA FOR MAPVIEW CENTER", http.StatusBadRequest)
		return false
	}
	return OPDB.UpdateMapViewCenter(f.Gid(), f.Fid(), hexagon.Coord{hexX, hexY})
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
	mp, ok := GetInts(r, "turn", "source", "target", "size")
	if !ok {
		http.Error(w, "BAD FORM DATA", http.StatusBadRequest)
		return
	}
	turn, source, target, size := mp["turn"], mp["source"], mp["target"], mp["size"]
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
	if size > 0 {
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
	return OPDB.SetOrder(g.Gid(), f.Fid(), source, target, size)
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
	mp, ok := GetInts(r, "turn")
	if !ok {
		http.Error(w, "MALFORMED TURN DATA", http.StatusBadRequest)
		return
	}
	turn, ok := mp["turn"]
	if !ok || turn != g.Turn() {
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
