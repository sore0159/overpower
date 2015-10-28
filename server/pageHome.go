package main

import (
	"mule/planetattack"
	"net/http"
)

var TPHOME = MixTemp("frame", "titlebar", "home")

func HomePage(w http.ResponseWriter, r *http.Request) {
	userN := CookieUserName(r)
	if userN == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	v := MakeView(r)
	m := map[string]interface{}{"userN": userN}
	v.SetApp(m)
	g := NewGame()
	g.Owner = userN
	if g.Select() {
		m["mygame"] = g
	}
	if r.Method == "POST" {
		action := r.FormValue("action")
		switch action {
		case "creategame":
			// test user owns game already
			if m["mygame"] != nil {
				v.SetError("USER", userN, "ALREADY OWNS GAME", g.Name, g.Gid)
				v.Apply(TPHOME, w)
				return
			}
			gName := r.FormValue("gamename")
			if !ValidGameName(gName) {
				v.SetError("INVALID GAME NAME: '", gName, "'")
				v.Apply(TPHOME, w)
				return
			}
			_, err := planetattack.MakeGame(ATTACKDB, gName, userN)
			if err != nil {
				v.SetError("Creation error:", err)
				v.Apply(TPHOME, w)
				return
			}
		default:
			Log("invalid action for homepage:", action)
			v.SetError("INVALID ACTION:", action)
			v.Apply(TPHOME, w)
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
	facs := AllFactions(userN)
	if len(facs) > 0 {
		m["facs"] = facs
	}
	v.Apply(TPHOME, w)
}
