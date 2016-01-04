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
	g, ok := OPDB.GetGame(gid)
	if !ok {
		http.Error(w, "GAME NOT FOUND", http.StatusNotFound)
		return
	}
	facs, ok := OPDB.GetGidFactions(gid)
	if !ok {
		http.Error(w, "DATABASE ERROR RETRIEVING FACTIONS", http.StatusInternalServerError)
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
	mv, ok := OPDB.GetFidMapView(gid, fid)
	if !ok {
		http.Error(w, "DATABASE ERROR RETRIEVING MAPVIEW", http.StatusInternalServerError)
		return
	}
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
	shipViews, ok := OPDB.GetFidTurnShipViews(gid, fid, g.Turn()-1)
	if !ok {
		http.Error(w, "DATABASE ERROR FETCHING SHIPVIEWS", http.StatusInternalServerError)
		return
	}
	mapping.ServeMap(w, mv, fid, facs, pvList, shipViews, orders)
}
