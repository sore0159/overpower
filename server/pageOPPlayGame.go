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
		case "mapclick":
			if !h.MapClick(w, r, g, f) {
				return
			}
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
	focus, hasFocus := mapView.Focus()
	center := mapView.Center()
	m := h.DefaultApp()
	if hasFocus {
		m["focus"] = focus
		if focus == center {
			m["cenfocused"] = true
		}
	}
	var fYou, cYou bool
	var fPid, cPid int
	var fPV overpower.PlanetView
	for _, pv := range pvList {
		plNames[pv.Pid()] = pv.Name()
		if pv.Controller() != fid {
			delete(oMap, pv.Pid())
		} else {
			sum := pv.Parts()
			for _, o := range oMap[pv.Pid()] {
				sum -= o.Size()
			}
			availMap[pv.Pid()] = sum
		}
		loc := pv.Loc()
		if hasFocus && focus == loc {
			m["focuspv"] = pv
			fPid = pv.Pid()
			fPV = pv
			if pv.Controller() == fid {
				fYou = true
			}
		}
		if loc == center {
			m["centerpv"] = pv
			cPid = pv.Pid()
			if pv.Controller() == fid {
				cYou = true
			}
		}
	}
	if fPid != 0 && cPid != 0 && (cYou || fYou) {
		for _, o := range orders {
			if fYou && o.Source() == fPid && o.Target() == cPid {
				m["oftoc"] = o
				fOrds := make([]overpower.Order, 0, len(oMap[cPid]))
				for _, test := range oMap[fPid] {
					if test != o {
						fOrds = append(fOrds, test)
					}
				}
				m["fords"] = fOrds
				continue
			}
			if cYou && o.Source() == cPid && o.Target() == fPid {
				m["octof"] = o
				cOrds := make([]overpower.Order, 0, len(oMap[cPid]))
				for _, test := range oMap[cPid] {
					if test != o {
						cOrds = append(cOrds, test)
					}
				}
				m["cords"] = cOrds
				continue
			}
		}
	}
	if fPid != 0 {
		m["fcdist"] = fPV.Loc().StepsTo(center)
	}
	if _, ok := m["cords"]; cYou && !ok {
		m["cords"] = oMap[cPid]
	}
	if _, ok := m["fords"]; fYou && !ok {
		m["fords"] = oMap[fPid]
	}

	m["fyou"], m["cyou"] = fYou, cYou
	names := map[int]string{0: "Hostile Natives"}
	for _, fac := range facs {
		if fac.Fid() == f.Fid() {
			names[fac.Fid()] = "Your Faction"
		} else {
			names[fac.Fid()] = "Faction " + fac.Name()
		}
	}
	centerShips := []overpower.ShipView{}
	for _, sv := range shipViews {
		if test, ok := sv.Loc(); ok && test == center {
			centerShips = append(centerShips, sv)
			continue
		}
		for _, test := range sv.Trail() {
			if test == center {
				centerShips = append(centerShips, sv)
				break
			}
		}
	}
	m["centersv"] = centerShips
	m["availparts"] = availMap
	m["plnames"] = plNames
	m["names"] = names
	m["orders"] = oMap
	m["game"] = g
	m["fac"] = f
	m["pvs"] = pvList
	m["svs"] = shipViews
	m["mapview"] = mapView
	zoom := mapView.Zoom()
	if zoom > 1 {
		m["zoomout"] = zoom - 1
	}
	if zoom > 10 {
		m["bigzoomout"] = zoom - 10
	}
	if zoom < 100 {
		m["zoomin"] = zoom + 1
	}
	if zoom < 91 {
		m["bigzoomin"] = zoom + 10
	}
	h.SetCommand(g)
	h.Apply(TPOPPLAY, w)
}
