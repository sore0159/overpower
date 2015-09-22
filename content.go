package main

import (
	"fmt"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	//d := NewTPContext()
	if r.Method == "POST" {
		userName := r.FormValue("username")
		password := r.FormValue("password")
		create := r.FormValue("create")
		if create != "true" {
			if Login(userName, password, w) {
				LoginTP(w, fmt.Sprintf("LOGIN AS", userName, "SUCCESSFUL"))
			} else {
				LoginTP(w, fmt.Sprintf("LOGIN PAGE: FAILURE: INCORRECT PASSWORD"))
			}
			return
		}
		if !ValidUserName(userName) {
			LoginTP(w, fmt.Sprintf("LOGIN PAGE: CREATE: BAD USERNAME"))
			return
		}
		if UserNameInUse(userName) {
			LoginTP(w, fmt.Sprintf("LOGIN PAGE: CREATE: USERNAME ALREADY IN USE"))
			return
		}
		if !ValidPassword(password) {
			LoginTP(w, fmt.Sprintf("LOGIN PAGE: CREATE: BAD PASSWORD"))
			return
		}
		if err := CreateUser(userName, password); err != nil {
			LoginTP(w, fmt.Sprintf("LOGIN PAGE: CREATE: ERROR:", err))
			return
		}
		if !Login(userName, password, w) {
			Log(userName, "failed creation login")
			LoginTP(w, fmt.Sprintf("CREATION COMPLETE, LOGIN ERROR"))
			return
		}
		LoginTP(w, fmt.Sprintf("LOGIN PAGE: ACCOUNT", userName, "CREATED SUCCESSFULLY"))
		return
	} else {
		LoginTP(w, fmt.Sprintf("LOGIN PAGE: GET"))
	}
}
