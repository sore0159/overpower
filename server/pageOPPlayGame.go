package main

import (
	"mule/hexagon"
	"mule/overpower"
	"net/http"
)

var (
	TPOPPLAY = MixTemp("frame", "titlebar", "play")
)

// /overpower/view/GID/oldplay
func (h *Handler) pageOPPlayGame(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction, facs []overpower.Faction) {
	if h.LastFull() > 4 {
		http.Redirect(w, r, h.NewPath(5), http.StatusFound)
		return
	}
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
	pvList, err := OPDB.GetPlanetViews("gid", gid, "fid", fid)
	if my, bad := Check(err, "resource failure", "page", "op play game", "resource", "planetviews", "gid", gid, "fid", fid); bad {
		Bail(w, my)
		return
	}
	orders, err := OPDB.GetOrders("gid", gid, "fid", fid)
	if my, bad := Check(err, "resource failure", "page", "op play game", "resource", "orders", "gid", gid, "fid", fid); bad {
		Bail(w, my)
		return
	}
	shipViews, err := OPDB.GetShipViews("gid", gid, "fid", fid, "turn", turn-1)
	if my, bad := Check(err, "resource failure", "page", "op play game", "resource", "shipviews", "gid", gid, "fid", fid, "turn", turn); bad {
		Bail(w, my)
		return
	}
	mapView, err := OPDB.GetMapView("gid", gid, "fid", fid)
	if my, bad := Check(err, "resource failure", "page", "op play game", "resource", "mapview", "gid", gid, "fid", fid); bad {
		Bail(w, my)
		return
	}
	rp, err := OPDB.GetReport("gid", gid, "fid", fid, "turn", turn-1)
	if err == ErrNoneFound {
	} else if my, bad := Check(err, "resource failure", "page", "op play game", "resource", "mapview", "gid", gid, "fid", fid); bad {
		Bail(w, my)
		return
	}
	m := h.DefaultApp()
	m["prevturn"] = turn - 1
	if rp != nil {
		m["reportlen"] = len(rp.Contents())
	}
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
			fNames[testFid] = "Faction " + fac.Name()
		}
	}
	// --------------- ORDER LOOP -------------- //
	oMap := make(map[int][]overpower.Order, len(orders))
	for _, o := range orders {
		src := o.Source()
		if list, ok := oMap[src]; ok {
			oMap[src] = append(list, o)
		} else {
			oMap[src] = []overpower.Order{o}
		}
	}
	// --------------- PLANET LOOP -------------- //
	availMap := make(map[int]int, len(pvList))
	var t1Ords, t2Ords []overpower.Order
	var t1Pid, t2Pid int
	pidPV := make(map[int]overpower.PlanetView, len(pvList))
	gridPV := make(map[hexagon.Coord]overpower.PlanetView, len(pvList))
	m["pidPV"] = pidPV
	for _, pv := range pvList {
		pid := pv.Pid()
		loc := pv.Loc()
		pidPV[pid] = pv
		gridPV[loc] = pv
		yours := pv.Controller() == fid
		if yours {
			sum := pv.Parts()
			if list, ok := oMap[pid]; ok {
				for _, o := range list {
					sum -= o.Size()
				}
			}
			availMap[pid] = sum
		} else {
			delete(oMap, pid)
		}
		if loc == center {
			m["centerpv"] = pv
		}
		if target1.IsCoord(loc) {
			m["tar1pv"] = pv
			t1Pid = pv.Pid()
			if list, ok := oMap[pv.Pid()]; yours && ok {
				t1Ords = list
			}
		} else if target2.IsCoord(loc) {
			m["tar2pv"] = pv
			t2Pid = pv.Pid()
			if list, ok := oMap[pv.Pid()]; yours && ok {
				t2Ords = list
			}
		}
	}
	if t1Ords != nil {
		list := make([]overpower.Order, 0, len(t1Ords))
		for _, o := range t1Ords {
			if o.Target() == t2Pid {
				m["t1t2ord"] = o
			} else {
				list = append(list, o)
			}
		}
		if len(list) > 0 {
			m["t1ords"] = list
		}
	}
	if t2Ords != nil {
		list := make([]overpower.Order, 0, len(t2Ords))
		for _, o := range t2Ords {
			if o.Target() == t1Pid {
				m["t2t1ord"] = o
			} else {
				list = append(list, o)
			}
		}
		if len(list) > 0 {
			m["t2ords"] = list
		}
	}
	// --------------- SHIP LOOP -------------- //
	svDest := map[overpower.ShipView]overpower.PlanetView{}
	svLoc := make(map[overpower.ShipView]hexagon.NullCoord, len(shipViews))
	svAtT1 := []overpower.ShipView{}
	svTrT1 := []overpower.ShipView{}
	svAtT2 := []overpower.ShipView{}
	svTrT2 := []overpower.ShipView{}
	for _, sv := range shipViews {
		loc, ok := sv.Loc()
		svLoc[sv] = hexagon.NullCoord{loc, ok}
		if ok {
			if target1.IsCoord(loc) {
				svAtT1 = append(svAtT1, sv)
			} else if target2.IsCoord(loc) {
				svAtT2 = append(svAtT2, sv)
			}
		}
		for _, c := range sv.Trail() {
			if target1.IsCoord(c) {
				svTrT1 = append(svTrT1, sv)
			} else if target2.IsCoord(c) {
				svTrT2 = append(svTrT2, sv)
			}
		}
		if dest, ok := sv.Dest(); ok {
			svDest[sv] = gridPV[dest]
		}
	}
	m["svAt1"] = svAtT1
	m["svAt2"] = svAtT2
	m["svTr1"] = svTrT1
	m["svTr2"] = svTrT2
	m["lAt1"] = len(svAtT1)
	m["lAt2"] = len(svAtT2)
	m["lTr1"] = len(svTrT1)
	m["lTr2"] = len(svTrT2)
	m["svDest"] = svDest
	m["svLoc"] = svLoc
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
