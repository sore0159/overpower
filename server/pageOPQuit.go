package main

import (
	"fmt"
	"net/http"
)

var (
	TPQUIT = MixTemp("frame", "titlebar", "quit")
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
		http.Error(w, "NOT LOGGED IN", http.StatusBadRequest)
		return
	}
	gid, ok := h.IntAt(3)
	if !ok {
		http.Error(w, "BAD GID IN URL", http.StatusBadRequest)
		return
	}
	g, err := OPDB.GetGame("gid", gid)
	if my, bad := Check(err, "resource failure", "page", "opquit", "resource", "game", "gid", gid); bad {
		Bail(w, my)
		return
	}
	if g.Turn() < 1 {
		http.Error(w, "GAME NOT IN PROGRESS", http.StatusNotFound)
		return
	}
	f, err := OPDB.GetFaction("gid", gid, "owner", h.User.String())
	if err == ErrNoneFound {
		http.Error(w, "NO FACTION FOUND FOR THIS USER FOR THIS GAME", http.StatusBadRequest)
		return
	} else if my, bad := Check(err, "resource failure", "page", "opquit", "resource", "faction", "gid", gid, "user", h.User); bad {
		Bail(w, my)
		return
	}
	if r.Method == "POST" {
		confirm := r.FormValue("confirm")
		if confirm == "true" {
			err := OPDB.DropFactions("gid", g.Gid(), "fid", f.Fid())
			if my, bad := Check(err, "data update failure", "data", f); bad {
				Bail(w, my)
				return
			}
		}
		http.Redirect(w, r, fmt.Sprintf("/overpower/view/%d", g.Gid()), http.StatusFound)
		return
	}
	m := h.DefaultApp()
	m["game"] = g
	m["faction"] = f
	h.Apply(TPQUIT, w)
}
