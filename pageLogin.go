package main

import (
	"fmt"
	"net/http"
)

var TPLOGIN = MixTemp("frame", "login")

func LoginPage(w http.ResponseWriter, r *http.Request) {
	//d := NewTPContext()
	if r.Method == "POST" {
		userName := r.FormValue("username")
		password := r.FormValue("password")
		create := r.FormValue("create")
		if create != "true" {
			if Login(userName, password, w) {
				Apply(TPLOGIN, w, fmt.Sprintf("LOGIN AS", userName, "SUCCESSFUL"))
			} else {
				Apply(TPLOGIN, w, fmt.Sprintf("LOGIN PAGE: FAILURE: INCORRECT PASSWORD"))
			}
			return
		}
		if !ValidUserName(userName) {
			Apply(TPLOGIN, w, fmt.Sprintf("LOGIN PAGE: CREATE: BAD USERNAME"))
			return
		}
		if UserNameInUse(userName) {
			Apply(TPLOGIN, w, fmt.Sprintf("LOGIN PAGE: CREATE: USERNAME ALREADY IN USE"))
			return
		}
		if !ValidPassword(password) {
			Apply(TPLOGIN, w, fmt.Sprintf("LOGIN PAGE: CREATE: BAD PASSWORD"))
			return
		}
		if err := CreateUser(userName, password); err != nil {
			Apply(TPLOGIN, w, fmt.Sprintf("LOGIN PAGE: CREATE: ERROR:", err))
			return
		}
		if !Login(userName, password, w) {
			Log(userName, "failed creation login")
			Apply(TPLOGIN, w, fmt.Sprintf("CREATION COMPLETE, LOGIN ERROR"))
			return
		}
		Apply(TPLOGIN, w, fmt.Sprintf("LOGIN PAGE: ACCOUNT", userName, "CREATED SUCCESSFULLY"))
		return
	} else {
		Apply(TPLOGIN, w, fmt.Sprintf("LOGIN PAGE: GET"))
	}
}
