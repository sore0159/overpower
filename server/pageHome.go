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
		case "startgame":
			if m["mygame"] == nil {
				v.SetError("USER", userN, "HAS NO GAME BEGIN")
				goto GET
			}
			if g.Turn > 0 {
				v.SetError("USER", userN, "GAME", g.Name, "HAS ALREADY BEGUN!")
				goto GET
			}
			facs, err := g.Factions()
			if err != nil {
				v.SetError("ERROR RETRIEVING FACTIONS:", err)
				goto GET
			}
			if len(facs) < 1 {
				v.SetError("USER", userN, "GAME", g.Name, "HAS NO PLAYERS!")
				goto GET
			}
			err = g.Start()
			if err != nil {
				v.SetError("GAME START ERROR:", err)
				goto GET
			}
		case "deletegame":
			if m["mygame"] == nil {
				v.SetError("USER", userN, "HAS NO GAME TO DELETE GAME")
				goto GET
			}
			planetattack.DelGame(ATTACKDB, g.Gid)
		case "creategame":
			// test user owns game already
			if m["mygame"] != nil {
				v.SetError("USER", userN, "ALREADY OWNS GAME", g.Name, g.Gid)
				goto GET
			}
			gName := r.FormValue("gamename")
			if !ValidGameName(gName) {
				v.SetError("INVALID GAME NAME: '", gName, "'")
				goto GET
			}
			password := r.FormValue("password")
			if password != "" && !ValidText(password) {
				v.SetError("INVALID GAME PASSWORD: '", password, "'")
				goto GET
			}
			_, err := planetattack.MakeGame(ATTACKDB, gName, userN, password)
			if err != nil {
				v.SetError("Creation error:", err)
				goto GET
			}
		default:
			Log("invalid action for homepage:", action)
			v.SetError("INVALID ACTION:", action)
			goto GET
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
GET:
	gids := []int{}
	myfacs := AllFactions(userN)
	for _, f := range myfacs {
		gids = append(gids, f.Gid)
	}
	if len(myfacs) > 0 {
		games := GetGames(gids)
		if len(games) > 0 {
			m["facs"] = myfacs
			m["fgames"] = games
		}
	}
	if m["mygame"] != nil {
		facs, _ := g.Factions()
		m["mygamefactions"] = facs
	}
	v.Apply(TPHOME, w)
}
