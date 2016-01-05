package main

import (
	"fmt"
	"mule/overpower"
	"net/http"
)

func SetupMux() {
	http.HandleFunc("/", pageMainIndex)
	http.HandleFunc("/favicon.ico", imgFavIcon)
	http.HandleFunc("/auth/", pageAuthIndex)
	http.HandleFunc("/auth/logout", pageAuthLogout)
	http.HandleFunc("/auth/login", pageAuthLogin)
	http.HandleFunc("/auth/create", pageAuthCreate)
	http.HandleFunc("/overpower/", pageOPIndex)
	http.HandleFunc("/overpower/home", pageOPHome)
	http.HandleFunc("/overpower/quit/", pageOPQuit)
	http.HandleFunc("/overpower/view/", muxView)
	http.HandleFunc("/overpower/command/", muxCommand)
	http.HandleFunc("/overpower/img/", pageMap)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))
}

// /overpower/view/GAMEID/
func muxView(w http.ResponseWriter, r *http.Request) {
	if DBLOCK {
		http.Error(w, "GAME DOWN FOR DAYLY MAINT: 10-20MIN", http.StatusInternalServerError)
		return
	}
	h := MakeHandler(w, r)
	lastFull := h.LastFull()
	if lastFull < 3 {
		h.pageOPViewIndex(w, r)
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
	if lastFull == 3 {
		h.pageOPViewGame(w, r, g)
		return
	}
	switch h.Path[4] {
	case "play", "reports":
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
		if h.Path[4] == "play" {
			h.pageOPPlayGame(w, r, g, f, facs)
		} else {
			h.pageOPReports(w, r, g, f, facs)
		}
		return
	default:
		http.Redirect(w, r, h.NewPath(4), http.StatusFound)
		return
	}
	h.pageOPViewIndex(w, r)
	return
}

// /overpower/command/GAMEID/TURN/ACTIONNAME/ARGS
func muxCommand(w http.ResponseWriter, r *http.Request) {
	if DBLOCK {
		http.Error(w, "GAME DOWN FOR DAYLY MAINT: 10-20MIN", http.StatusInternalServerError)
		return
	}
	h := MakeHandler(w, r)
	if !h.LoggedIn {
		http.Error(w, "NOT LOGGED IN", http.StatusBadRequest)
		return
	}
	lastFull := h.LastFull()
	gid, ok := h.IntAt(3)
	if !ok {
		http.Error(w, "BAD GID IN COMMAND PATH", http.StatusBadRequest)
		return
	}
	g, ok := OPDB.GetGame(gid)
	if !ok {
		http.Error(w, "GAME NOT FOUND", http.StatusNotFound)
		return
	}
	f, ok := OPDB.GetOwnerFaction(gid, h.User.String())
	if !ok {
		http.Error(w, "NO FACTION FOUND FOR THIS USER FOR THIS GAME", http.StatusBadRequest)
		return
	}
	turn, ok := h.IntAt(4)
	if !ok || turn != g.Turn() {
		http.Error(w, "BAD TURN IN COMMAND PATH", http.StatusBadRequest)
		return
	}
	if lastFull < 5 {
		http.Error(w, "ACTION IN COMMAND PATH", http.StatusBadRequest)
		return
	}
	switch h.Path[5] {
	case "recenter":
		ok = h.Recenter(w, r, g, f)
	default:
		http.Error(w, "ACTION IN COMMAND PATH", http.StatusBadRequest)
		return
	}
	if !ok {
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/overpower/view/%d/play", gid), http.StatusFound)
}

func imgFavIcon(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PING")
	http.ServeFile(w, r, "STATIC/img/favicon.ico")
}
