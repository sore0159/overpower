package main

import (
	"net/http"
)

var TPLOGIN = MixTemp("frame", "titlebar", "login")

func LoginPage(w http.ResponseWriter, r *http.Request) {
	v := MakeView(r)
	m := map[string]string{}
	v.SetApp(m)
	if r.Method == "POST" {
		userName := r.FormValue("username")
		password := r.FormValue("password")
		create := r.FormValue("create")
		m["username"] = userName
		m["password"] = password
		if create == "true" {
			m["create"] = create
		}
		if create != "true" {
			if Login(userName, password, w) {
				http.Redirect(w, r, "/home", http.StatusFound)
			} else {
				delete(m, "password")
				v.SetError("LOGIN PAGE: FAILURE: INCORRECT PASSWORD")
				v.Apply(TPLOGIN, w)
			}
			return
		}
		if !ValidUserName(userName) {
			delete(m, "username")
			v.SetError("LOGIN PAGE: CREATE: BAD USERNAME")
			v.Apply(TPLOGIN, w)
			return
		}
		if !UserNameAvail(userName) {
			delete(m, "username")
			v.SetError("LOGIN PAGE: CREATE: USERNAME NOT AVAILABLE")
			v.Apply(TPLOGIN, w)
			return
		}
		if !ValidPassword(password) {
			delete(m, "password")
			v.SetError("LOGIN PAGE: CREATE: BAD PASSWORD")
			v.Apply(TPLOGIN, w)
			return
		}
		if err := CreateUser(userName, password); err != nil {
			v.SetError("LOGIN PAGE: CREATE: ERROR:", err)
			v.Apply(TPLOGIN, w)
			return
		}
		if !Login(userName, password, w) {
			Log(userName, "failed creation login")
			v.SetError("CREATION COMPLETE, LOGIN ERROR")
			v.Apply(TPLOGIN, w)
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	} else {
		v.Apply(TPLOGIN, w)
	}
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	Logout(w)
	http.Redirect(w, r, "/", http.StatusFound)
}
