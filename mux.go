package main

import (
	//	"fmt"
	//	"mule/planetattack/attack"
	"net/http"
	//	"strings"
)

const SAVEFL = "savefile.gob"

var TEMPGAME *Game

func gameMux(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch path {
	case "/login":
		LoginPage(w, r)
		return
	case "/logout":
		LogoutPage(w, r)
		return
	}
	v := MakeView(r)
	l := len(v.path)
	g := TEMPGAME
	if g.Game == nil {
		if path != "/" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		g.GameCreate(w, r, v)
		return
	}
	if l < 2 || l == 2 && v.path[1] == "" {
		g.GameView(w, r, v)
		return
	}
	if l == 2 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if v.path[1] == "view" {
		fName := v.path[2]
		if f, ok := g.AuthFaction(v.User, fName); ok {
			g.FactionView(w, r, f, v)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
