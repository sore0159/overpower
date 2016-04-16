package main

import (
	"fmt"
	"net/http"
)

var (
	TPOPQUIT = MixTemp("frame", "titlebar", "opquit")
)

// /overpower/quit/GID
func pageOPQuit(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	if DBLOCK {
		h.HandleUserError(w, r, "GAME DOWN FOR DAYLY MAINT: 10-20MIN")
		return
	}
	if h.LastFull() > 3 {
		http.Redirect(w, r, h.NewPath(4), http.StatusFound)
		return
	}
	if !h.LoggedIn {
		h.HandleUserError(w, r, "USER NOT LOGGED IN")
		return
	}
	gid, ok := h.IntAt(4)
	if !ok {
		h.HandleUserError(w, r, "NO/BAD GAME ID SPECIFIED")
		return
	}
	g, f, _, err := h.FetchBasicData(gid)
	if my, bad := Check(err, "page quit failure on resource aquisition", "gid", gid); bad {
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
	if g.Turn() < 1 {
		h.HandleUserError(w, r, "GAME HAS NOT YET BEGUN")
		return
	}
	if r.Method == "POST" {
		confirm := r.FormValue("confirm")
		var errS, errU error
		if confirm == "true" {
			errS, errU = h.CommandQuitGame(g, f, r.FormValue("turn"))
		}
		if my, bad := Check(errS, "quitgame execute failure", "game", g, "fac", f); bad {
			h.HandleServerError(w, r, my)
			return
		}
		if errU != nil {
			h.HandleUserError(w, r, errU.Error())
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/overpower/view/%d", g.GID()), http.StatusFound)
		return
	}
	m := h.DefaultApp()
	m["game"] = g
	m["faction"] = f
	h.Apply(TPOPQUIT, w)
}
