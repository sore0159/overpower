package main

import (
	"mule/hexagon"
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
		var ok bool
		switch action {
		case "mapclick":
			ok = h.MapClick(w, r, g, f)
		case "setorder":
			ok = h.SetOrder(w, r, g, f)
		case "setdone":
			ok = h.SetTurnDone(w, r, g, f)
		case "setmap":
			ok = h.SetMapView(w, r, g, f)
		default:
			http.Error(w, "UNKNOWN ACTION TYPE", http.StatusBadRequest)
			return
		}
		if ok {
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
		}
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
	m := h.DefaultApp()
	_ = hexagon.Coord{}
	// ---------------- DATA PROCESSING --------------- //
	center, target1, target2 := mapView.Center(), mapView.Target1(), mapView.Target2()
	// --------------- FACTION LOOP -------------- //
	fNames := make(map[int]string, len(facs)+1)
	fNames[0] = "Hostile Natives"
	for _, fac := range facs {
		testFid := fac.Fid()
		if fid == testFid {
			fNames[fid] = "your faction"
		} else {
			fNames[fid] = "Faction " + fac.Name()
		}
	}
	availMap := make(map[int]int, len(pvList))
	// --------------- PLANET LOOP -------------- //
	for _, pv := range pvList {
		loc := pv.Loc()
		if loc == center {
			m["centerpv"] = pv
		}
		if target1.IsCoord(loc) {
			m["tar1pv"] = pv
		} else if target2.IsCoord(loc) {
			m["tar2pv"] = pv
		}
		if pv.Controller() == fid {
			availMap[pv.Pid()] = pv.Parts()
		}
	}
	// --------------- ORDER LOOP -------------- //
	for _, o := range orders {
		_ = o
	}
	// --------------- SHIP LOOP -------------- //
	for _, sv := range shipViews {
		_ = sv
	}
	// ---------- BEGIN APP -------------- //
	m["availMap"] = availMap
	m["fNames"] = fNames
	m["game"] = g
	m["fac"] = f
	m["mapview"] = mapView
	m["target1"] = target1
	m["target2"] = target2
	if target1.Eq(target2) {
		m["targetsSame"] = true
	} else {
		if target1.Valid || target2.Valid {
			m["canswap"] = true
		}
		if target1.Valid && target2.Valid {
			m["tarDist"] = target1.Coord.StepsTo(target2.Coord)
		}
	}
	if target1.Valid && !target1.IsCoord(center) {
		m["cancenter1"] = true
	}
	if target2.Valid && !target2.IsCoord(center) {
		m["cancenter2"] = true
	}
	zoom := mapView.Zoom()
	if zoom > 1 {
		m["zoomout"] = zoom - 1
	}
	if zoom > 10 {
		m["bigzoomout"] = zoom - 10
	}
	if zoom < 50 {
		m["zoomin"] = zoom + 1
	}
	if zoom < 41 {
		m["bigzoomin"] = zoom + 10
	}
	h.SetCommand(g)
	h.Apply(TPOPPLAY, w)
}

/*

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
	tar1, tar2 := mapView.Target1(), target2()
	center := mapView.Center()
	locNames := make(map[hexagon.Coord]string, len(pvList))
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
		locNames[loc] = pv.Name()
	}
	/*
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
		shipVLoc := make(map[int]hexagon.Coord, len(shipViews))
		shipVLocV := make(map[int]bool, len(shipViews))
		shipVDest := make(map[int]hexagon.Coord, len(shipViews))
		shipVDestN := make(map[int]string, len(shipViews))
		shipVDestV := make(map[int]bool, len(shipViews))
		for _, sv := range shipViews {
			sid := sv.Sid()
			if test, ok := sv.Dest(); ok {
				shipVDest[sid] = test
				shipVDestN[sid] = locNames[test]
				shipVDestV[sid] = true
			}
			if test, ok := sv.Loc(); ok {
				shipVLoc[sid] = test
				shipVLocV[sid] = true
				if test == center {
					centerShips = append(centerShips, sv)
					continue
				}
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
		m["pvs"] = pvList
		m["svs"] = shipViews
		m["svsL"] = shipVLoc
		m["svsLV"] = shipVLocV
		m["svsD"] = shipVDest
		m["svsDN"] = shipVDestN
		m["svsDV"] = shipVDestV
*/
