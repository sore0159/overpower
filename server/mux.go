package main

import (
	"fmt"
	"mule/overpower"
	"net/http"
)

func SetupMux() {
	http.HandleFunc("/", pageMainIndex)
	http.HandleFunc("/auth/", pageAuthIndex)
	http.HandleFunc("/auth/logout", pageAuthLogout)
	http.HandleFunc("/auth/login", pageAuthLogin)
	http.HandleFunc("/auth/create", pageAuthCreate)
	http.HandleFunc("/overpower/", pageOPIndex)
	http.HandleFunc("/overpower/home", pageOPHome)
	http.HandleFunc("/overpower/view/", muxView)
	http.HandleFunc("/overpower/command/", muxCommand)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))
}

// /overpower/view/GAMEID/
func muxView(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	lastFull := h.LastFull()
	if lastFull < 3 {
		h.pageOPViewIndex(w, r)
		return
	}
	if lastFull < 5 {
		gid, ok := h.IntAt(3)
		if !ok {
			http.Error(w, "INVALID GAMEID", http.StatusBadRequest)
			return
		}
		g, ok := OPDB.GetGame(gid)
		if !ok {
			http.Error(w, fmt.Sprintf("GAME %d NOT FOUND", gid), http.StatusNotFound)
			return
		}
		if lastFull == 3 {
			h.pageOPViewGame(w, r, g)
			return
		}
		switch h.Path[4] {
		case "play":
			if !h.LoggedIn {
				http.Error(w, "NOT LOGGED IN", http.StatusBadRequest)
				return
			}
			facs, ok := OPDB.GetGidFactions(gid)
			if !ok {
				http.Error(w, "DATABASE ERROR LOADING FACTIONS", http.StatusInternalServerError)
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
			h.pageOPPlayGame(w, r, g, f, facs)
			return
		default:
			http.Redirect(w, r, h.NewPath(4), http.StatusFound)
			return
		}
	}
	h.pageOPViewIndex(w, r)
	return
}

// /overpower/command/GAMEID/TURN/ACTIONNAME/ARGS
func muxCommand(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	_ = h
	fmt.Fprint(w, "MUX COMMAND")
}
