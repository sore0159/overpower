package main

import (
	"net/http"
)

var (
	TPLOGIN      = MixTemp("frame", "titlebar", "authlogin")
	TPAUTHCREATE = MixTemp("frame", "titlebar", "authcreate")
)

func pageAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/login" {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}
	h := MakeHandler(w, r)
	m := h.DefaultApp()
	m["loggedin"] = h.LoggedIn
	if r.Method == "POST" {
		_, ok := USERREG.Login(w, r)
		if ok {
			http.Redirect(w, r, "/overpower/home", http.StatusFound)
			return
		}
		m["username"], m["password"] = r.FormValue("username"), r.FormValue("password")
		h.SetError("INVALID LOGIN")
	}
	h.Apply(TPLOGIN, w)
}

func pageAuthLogout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/logout" || r.Method != "POST" {
		http.Redirect(w, r, "/auth/", http.StatusFound)
		return
	}
	USERREG.Logout(w, r)
	http.Redirect(w, r, "/overpower/", http.StatusFound)
}

func pageAuthCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/create" {
		http.Redirect(w, r, "/auth/create", http.StatusFound)
		return
	}
	h := MakeHandler(w, r)
	h.Apply(TPAUTHCREATE, w)
}
