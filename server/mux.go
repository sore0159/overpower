package main

import (
	"fmt"
	"net/http"
)

func SetupMux() {
	http.HandleFunc("/", pageIndex)
	http.HandleFunc("/auth/", muxAuth)
	http.HandleFunc("/overpower/", pageOPIndex)
	http.HandleFunc("/overpower/view/", muxView)
	http.HandleFunc("/overpower/command/", muxCommand)
	http.HandleFunc("/img/", muxImg)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))
}

func muxAuth(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/auth/" {
		pageAuthIndex(w, r)
		return
	}
	h := MakeHandler(w, r)
	lastFull := h.LastFull()
	if lastFull > 2 {
		http.Redirect(w, r, "/auth/", http.StatusFound)
		return
	}
	switch h.Path[2] {
	case "create":
		h.pageAuthCreate(w, r)
	case "login":
		fmt.Fprint(w, "AUTH LOGIN")
	case "logout":
		fmt.Fprint(w, "AUTH LOGOUT")
	case "createcommand":
		fmt.Fprint(w, "AUTH CREATECOMMAND")
	default:
		http.Redirect(w, r, "/", http.StatusFound)
	}
	return
}

func muxView(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	_ = h
	fmt.Fprint(w, "MUX VIEW")
}

func muxCommand(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	_ = h
	fmt.Fprint(w, "MUX COMMAND")
}

func muxImg(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	_ = h
	fmt.Fprint(w, "MUX IMG")
}
