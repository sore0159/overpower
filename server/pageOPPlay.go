package main

import (
	"net/http"
)

var (
	TPPLAY = MixTemp("frame", "titlebar", "opplay")
)

// /overpower/play/GID
func pageOPPlay(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	if len(h.Path) > 4 {
		http.Redirect(w, r, h.NewPath(3), http.StatusFound)
		return
	}
	if !h.LoggedIn {
		h.HandleUserError(w, r, "USER NOT LOGGED IN")
		return
	}
	gid, ok := h.IntAt(3)
	if !ok {
		h.HandleUserError(w, r, "NO/BAD GAME ID SPECIFIED")
		return
	}
	g, f, _, err := h.FetchBasicData(gid)
	if my, bad := Check(err, "page play failure on resource aquisition", "gid", gid); bad {
		h.HandleServerError(w, r, my)
		return
	}
	if g == nil {
		h.HandleUserError(w, r, "NO GAME FOUND")
		return
	}
	if f == nil {
		h.HandleUserError(w, r, "USER HAS NO FACTION FOR THIS GAME")
		return
	}
	turn := g.Turn() - 1
	if turn < 0 {
		h.HandleUserError(w, r, "GAME HAS NOT YET BEGUN")
		return
	}
	m := h.DefaultApp()
	m["game"] = g
	m["faction"] = f
	h.Apply(TPPLAY, w)
}
