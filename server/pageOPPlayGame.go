package main

import (
	"mule/overpower"
	"net/http"
)

var (
	TPOPPLAY = MixTemp("frame", "titlebar", "play")
)

func (h *Handler) pageOPPlayGame(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction, facs []overpower.Faction) {
	pvList, ok := OPDB.GetAllFactionPlanetViews(g.Gid(), f.Fid())
	if !ok {
		http.Error(w, "DATABASE ERROR FETCHING PLANETVIEWS", http.StatusInternalServerError)
		return
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
	m["names"] = names
	m["game"] = g
	m["fac"] = f
	m["pvs"] = pvList
	h.Apply(TPOPPLAY, w)
}
