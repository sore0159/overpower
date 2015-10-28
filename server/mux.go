package main

import (
	//	"fmt"
	"net/http"
	//	"strings"
)

func gameMux(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch path {
	case "/":
		IndexPage(w, r)
		return
	case "/login":
		LoginPage(w, r)
		return
	case "/logout":
		LogoutPage(w, r)
		return
	case "/home":
		HomePage(w, r)
		return
	default:
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}
