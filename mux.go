package main

import (
	"fmt"
	"net/http"
)

func gameMux(w http.ResponseWriter, r *http.Request) {
	userName := CookieUserName(r)
	if userName == "" {
		LoginPage(w, r)
		return
	}
	path := r.URL.Path
	fmt.Fprint(w, "LOGGED IN USER:", userName, "PATH:", path)
}
