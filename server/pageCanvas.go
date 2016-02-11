package main

import (
	"mule/overpower"
	"net/http"
)

var (
	TPCANVAS = MixTemp("frame", "titlebar", "play")
)

// /overpower/view/GID/play
func (h *Handler) pageCanvas(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction, facs []overpower.Faction) {
	if len(h.Path) > 5 {
		http.Redirect(w, r, h.NewPath(4), http.StatusFound)
		return
	}
	m := h.DefaultApp()
	m["game"] = g
	m["faction"] = f
	h.Apply(TPCANVAS, w)
}
