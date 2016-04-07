package main

import (
	"net/http"
)

var (
	TPAUTHINDEX  = MixTemp("frame", "titlebar", "authindex")
	TPLOGIN      = MixTemp("frame", "titlebar", "authlogin")
	TPAUTHCREATE = MixTemp("frame", "titlebar", "authcreate")
)

func pageAuthIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/" {
		http.Redirect(w, r, "/auth/", http.StatusFound)
		return
	}
	h := MakeHandler(w, r)
	m := h.DefaultApp()
	m["loggedin"] = h.LoggedIn
	h.Apply(TPAUTHINDEX, w)
}

func pageAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/login" {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}
	h := MakeHandler(w, r)
	m := h.DefaultApp()
	m["loggedin"] = h.LoggedIn
	if r.Method == "POST" {
		_, ok, err := USERREG.Login(w, r)
		if my, bad := Check(err, "login failure"); bad {
			h.HandleServerError(w, my)
			return
		}
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
	m := h.DefaultApp()
	m["loggedin"] = h.LoggedIn
	if r.Method == "POST" {
		nameOk, passOk, err := USERREG.Create(w, r)
		if nameOk && passOk && err == nil {
			http.Redirect(w, r, "/overpower/home", http.StatusFound)
			return
		}
		if my, bad := Check(err, "auth create failure"); bad {
			h.HandleServerError(w, my)
			return
			//m["username"], m["password"] = r.FormValue("username"), r.FormValue("password")
		} else {
			if !passOk && !nameOk {
				h.SetError("BAD USERNAME AND PASSWORD")
			} else if passOk {
				h.SetError("BAD USERNAME")
				m["password"] = r.FormValue("password")
			} else {
				h.SetError("BAD PASSWORD")
				m["username"] = r.FormValue("username")
			}
		}
	}
	h.Apply(TPAUTHCREATE, w)
}
