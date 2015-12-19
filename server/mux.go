package main

import (
	"fmt"
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
	if lastFull == 3 {
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
		h.pageOPViewGame(w, r, g)
		return
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
