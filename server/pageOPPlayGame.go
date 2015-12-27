package main

import (
	"mule/overpower"
	"net/http"
)

var (
	TPOPPLAY = MixTemp("frame", "titlebar", "play")
)

func (h *Handler) pageOPPlayGame(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction, facs []overpower.Faction) {
	turn := g.Turn()
	if turn < 1 {
		http.Redirect(w, r, h.NewPath(4), http.StatusFound)
		return
	}
	if r.Method == "POST" {
		if !h.LoggedIn {
			http.Error(w, "NOT LOGGED IN", http.StatusBadRequest)
			return
		}
		action := r.FormValue("action")
		switch action {
		case "setorder":
			if !h.SetOrder(w, r, g, f) {
				return
			}
		case "setdone":
			if !h.SetTurnDone(w, r, g, f) {
				return
			}
		case "setmap":
			if !h.SetMapView(w, r, g, f) {
				return
			}
		default:
			http.Error(w, "UNKNOWN ACTION TYPE", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	}
	gid, fid := g.Gid(), f.Fid()
	pvList, ok := OPDB.GetAllFactionPlanetViews(gid, fid)
	if !ok {
		http.Error(w, "DATABASE ERROR FETCHING PLANETVIEWS", http.StatusInternalServerError)
		return
	}
	orders, ok := OPDB.GetAllGidOrders(gid)
	if !ok {
		http.Error(w, "DATABASE ERROR FETCHING ORDERS", http.StatusInternalServerError)
		return
	}
	shipViews, ok := OPDB.GetFidTurnShipViews(gid, fid, turn-1)
	if !ok {
		http.Error(w, "DATABASE ERROR FETCHING SHIPVIEWS", http.StatusInternalServerError)
		return
	}
	mapView, ok := OPDB.GetFidMapView(gid, fid)
	if !ok {
		http.Error(w, "DATABASE ERROR FETCHING MAPVIEW", http.StatusInternalServerError)
		return
	}

	oMap := map[int][]overpower.Order{}
	for _, o := range orders {
		pid := o.Source()
		if list, ok := oMap[pid]; ok {
			oMap[pid] = append(list, o)
		} else {
			oMap[pid] = []overpower.Order{o}
		}
	}
	availMap := map[int]int{}
	plNames := make(map[int]string, len(pvList))
	for _, pv := range pvList {
		plNames[pv.Pid()] = pv.Name()
		if pv.Controller() != f.Fid() {
			delete(oMap, pv.Pid())
		} else {
			sum := pv.Parts()
			for _, o := range oMap[pv.Pid()] {
				sum -= o.Size()
			}
			availMap[pv.Pid()] = sum
		}
	}
	m := h.DefaultApp()
	names := map[int]string{0: "No Faction"}
	for _, fac := range facs {
		if fac.Fid() == f.Fid() {
			names[fac.Fid()] = "Your Faction"
		} else {
			names[fac.Fid()] = "Faction " + fac.Name()
		}
	}
	m["availparts"] = availMap
	m["plnames"] = plNames
	m["names"] = names
	m["orders"] = oMap
	m["game"] = g
	m["fac"] = f
	m["pvs"] = pvList
	m["svs"] = shipViews
	m["mapview"] = mapView
	h.Apply(TPOPPLAY, w)
}
