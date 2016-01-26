package main

import (
	"mule/overpower"
	"mule/overpower/mapping"
	"net/http"
)

// /overpower/img/GID
func pageMap(w http.ResponseWriter, r *http.Request) {
	if DBLOCK {
		http.Error(w, "GAME DOWN FOR DAYLY MAINT: 10-20MIN", http.StatusInternalServerError)
		return
	}
	h := MakeHandler(w, r)
	if !h.LoggedIn {
		http.Error(w, "NOT LOGGED IN", http.StatusBadRequest)
		return
	}
	if h.LastFull() != 3 {
		http.Error(w, "BAD REQUEST PATH", http.StatusBadRequest)
		return
	}
	gid, ok := h.IntAt(3)
	if !ok {
		http.Error(w, "INVALID GAMEID", http.StatusBadRequest)
		return
	}
	g, err := OPDB.GetGame("gid", gid)
	if err == ErrNoneFound {
		http.Error(w, "GAME NOT FOUND", http.StatusNotFound)
		return
	} else if my, bad := Check(err, "resource failure on map page", "gid", gid); bad {
		Bail(w, my)
		return
	}
	facs, err := OPDB.GetFactions("gid", gid)
	if my, bad := Check(err, "resource failure on map page", "gid", gid); bad {
		Bail(w, my)
		return
	}
	var f overpower.Faction
	for _, test := range facs {
		if test.Owner() == h.User.String() {
			f = test
			break
		}
	}
	if f == nil {
		http.Error(w, "USER HAS NO FACTION IN THIS GAME", http.StatusBadRequest)
		return
	}
	fid := f.Fid()
	mv, err := OPDB.GetMapView("gid", gid, "fid", fid)
	if my, bad := Check(err, "resource failure on map page", "gid", gid, "fid", fid); bad {
		Bail(w, my)
		return
	}
	pvList, err := OPDB.GetPlanetViews("gid", gid, "fid", fid)
	if my, bad := Check(err, "resource failure on map page", "gid", gid, "fid", fid); bad {
		Bail(w, my)
		return
	}
	orders, err := OPDB.GetOrders("gid", gid, "fid", fid)
	if my, bad := Check(err, "resource failure on map page", "gid", gid, "fid", fid); bad {
		Bail(w, my)
		return
	}
	shipViews, err := OPDB.GetShipViews("gid", gid, "fid", fid, "turn", g.Turn()-1)
	if my, bad := Check(err, "resource failure on map page", "resource", "shipviews", "gid", gid, "fid", fid, "turn", g.Turn()-1); bad {
		Bail(w, my)
		return
	}
	mapping.ServeMap(w, mv, fid, facs, pvList, shipViews, orders)
}
