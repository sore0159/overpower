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
	if DBLOCK {
		http.Error(w, "GAME DOWN FOR DAYLY MAINT: 10-20MIN", http.StatusInternalServerError)
		return
	}
	h := MakeHandler(w, r)
	if h.LastFull() > 3 {
		http.Redirect(w, r, h.NewPath(4), http.StatusFound)
		return
	}
	if !h.LoggedIn {
		h.HandleUserError(w, "NOT LOGGED IN")
		return
	}
	gid, ok := h.IntAt(3)
	if !ok {
		h.HandleUserError(w, "BAD GID IN URL")
		return
	}
	games, err := h.M.Game().SelectWhere(h.GID(gid))
	if my, bad := Check(err, "resource failure", "page", "opquit", "resource", "game", "gid", gid); bad {
		h.HandleServerError(w, my)
		return
	}
	if len(games) == 0 {
		h.HandleUserError(w, "NO GAME FOUND FOR GID")
		return
	}
	g := games[0]
	if g.Turn() < 1 {
		h.HandleUserError(w, "GAME NOT IN PROGRESS")
		return
	}
	facs, err := h.M.Faction().Select("gid", gid, "owner", h.User.String())
	if my, bad := Check(err, "resource failure", "page", "opquit", "resource", "faction", "gid", gid, "user", h.User); bad {
		h.HandleServerError(w, my)
		return
	}
	if len(facs) == 0 {
		h.HandleUserError(w, "NO FACTION FOUND FOR THIS USER FOR THIS GAME")
		return
	}
	if r.Method == "POST" {
		confirm := r.FormValue("confirm")
		var errS, errU error
		if confirm == "true" {
			errS, errU = h.CommandQuitGame(g, facs[0])
		}
		if my, bad := Check(errS, "quitgame execute failure", "game", g, "fac", facs[0]); bad {
			h.HandleServerError(w, my)
			return
		}
		if errU != nil {
			h.HandleUserError(w, errU.Error())
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/overpower/view/%d", g.GID()), http.StatusFound)
		return
	}
	m := h.DefaultApp()
	m["game"] = g
	m["faction"] = facs[0]
	h.Apply(TPOPQUIT, w)
}
