package main

import (
	"net/http"
)

func Login(userName, password string, w http.ResponseWriter) bool {
	if !ValidLogin(userName, password) {
		Log("Bad attempt at login for user", userName, ":", password)
		return false
	}
	//Log(userName, "logging in")
	cookie := &http.Cookie{
		Name:  "username",
		Value: userName,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	return true
}

func Logout(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func CookieUserName(r *http.Request) string {
	cookie, err := r.Cookie("username")
	if err == http.ErrNoCookie {
		return ""
	} else if err != nil {
		Log(err)
		return ""
	}
	return cookie.Value
}
