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
	g, ok := OPDB.GetGame(gid)
	if !ok {
		http.Error(w, "GAME NOT FOUND", http.StatusNotFound)
		return
	}
	if g.Turn() < 1 {
		http.Error(w, "GAME NOT IN PROGRESS", http.StatusNotFound)
		return
	}
	f, ok := OPDB.GetOwnerFaction(gid, h.User.String())
	if !ok {
		http.Error(w, "NO FACTION FOUND FOR THIS USER FOR THIS GAME", http.StatusBadRequest)
		return
	}
	if r.Method == "POST" {
		confirm := r.FormValue("confirm")
		if confirm == "true" {
			if !OPDB.DropInProgressFaction(g.Gid(), f.Fid()) {
				http.Error(w, "DATABASE ERROR DROPPING FACTION", http.StatusInternalServerError)
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
