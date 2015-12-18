package main

import (
	"fmt"
	"net/http"
)

func SetupMux() {
	http.HandleFunc("/", pageIndex)
	http.HandleFunc("/auth/", pageAuthIndex)
	http.HandleFunc("/auth/logout", pageAuthLogout)
	http.HandleFunc("/auth/login", pageAuthLogin)
	http.HandleFunc("/auth/create", pageAuthCreate)
	http.HandleFunc("/overpower/", pageOPIndex)
	http.HandleFunc("/overpower/view/", muxView)
	http.HandleFunc("/overpower/command/", muxCommand)
	http.HandleFunc("/img/", muxImg)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))
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
