package main

import (
	"mule/overpower"
	"net/http"
)

var (
	TPOPHOME = MixTemp("frame", "titlebar", "home")
)

func pageOPHome(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	if !h.LoggedIn {
		http.Redirect(w, r, "/auth/login", http.StatusFound)
		return
	}
	if r.URL.Path != "/overpower/home" {
		http.Redirect(w, r, "/overpower/home", http.StatusFound)
		return
	}
	g, hasG := OPDB.GetGameFor(h.User.String())
	var gFacs []overpower.Faction
	var gHasF bool
	if hasG {
		gFacs, gHasF = OPDB.GetGidFactions(g.Gid())
	}
	oFacs, oHasF := OPDB.GetOwnerFactions(h.User.String())
	if r.Method == "POST" {
		action := r.FormValue("action")
		switch action {
		case "startgame":
			if !hasG {
				h.SetError("USER HAS NO GAME TO START")
				break
			}
			if g.Turn() > 0 {
				h.SetError("GAME ALREADY IN PROGRESS")
				break
			}
			if len(gFacs) < 1 {
				h.SetError("GAME HAS NO PLAYERS")
				break
			}
			if !OPDB.StartGame(g.Gid()) {
				h.SetError("DATABASE ERROR STARTING GAME")
				break
			}
		case "newgame":
			if hasG {
				h.SetError("USER ALREADY HAS GAME IN PROGRESS")
				break
			}
			name, password := r.FormValue("gamename"), r.FormValue("password")
			if !ValidText(name) {
				h.SetError("INVALID GAME NAME")
				break
			}
			if password != "" && !ValidText(password) {
				h.SetError("INVALID GAME PASSWORD")
				break
			}
			if !OPDB.MakeGame(h.User.String(), name, password) {
				h.SetError("DATABASE ERROR IN GAME CREATION")
				break
			}
		case "dropgame":
			if !hasG {
				h.SetError("USER HAS NO GAME IN PROGRESS")
				break
			}
			OPDB.DropGame(g.Gid())
		default:
			h.SetError("UNKNOWN ACTION TYPE")
		}
		if h.Error == "" {
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		}
	}
	//
	m := h.DefaultApp()
	m["user"] = h.User
	if hasG {
		m["game"] = g
		m["active"] = g.Turn() > 0
	}
	if oHasF {
		m["ofactions"] = oFacs
	}
	if gHasF {
		m["gfactions"] = gFacs
	}
	h.Apply(TPOPHOME, w)
}
