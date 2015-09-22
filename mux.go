package main

import (
	"fmt"
	"net/http"
)

func gameMux(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path) > 5 && path[:6] == "/view/" {
		tempGamePage(w, r)
		return
	}
	userName := CookieUserName(r)
	if userName == "" {
		LoginPage(w, r)
		return
	}
	fmt.Fprint(w, "LOGGED IN USER:", userName, "PATH:", path)
}
