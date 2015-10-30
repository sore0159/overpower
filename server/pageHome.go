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
		case "approverq":
			if m["mygame"] == nil {
				v.SetError("USER", userN, "HAS NO GAME TO APPROVE REQUESTS")
				goto GET
			}
			facOwner := r.FormValue("owner")
			if !UserExists(facOwner) {
				v.SetError("USER", facOwner, "NOT FOUND")
				goto GET
			}
			rq, ok := g.GetRequest(facOwner)
			if !ok {
				v.SetError("REQUEST FOR USER", facOwner, "NOT FOUND")
				goto GET
			}
			if _, err := rq.Approve(); err != nil {
				v.SetError("APRROVAL ERROR FOR USER", facOwner, ":", err)
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
			_, err := planetattack.MakeGame(ATTACKDB, gName, userN)
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
	reqs := AllRequests(userN)
	for _, rq := range reqs {
		gids = append(gids, rq.Gid)
	}
	games := GetGames(gids)
	fGames := make([]*planetattack.Game, len(myfacs))
	rqGames := make([]*planetattack.Game, len(reqs))
	for i, gm := range games {
		if i >= len(myfacs) {
			rqGames[i-len(myfacs)] = gm
		} else {
			fGames[i] = gm
		}
	}
	if len(myfacs) > 0 {
		m["facs"] = myfacs
		m["fgames"] = fGames
	}
	if len(reqs) > 0 {
		m["requests"] = reqs
		m["rqgames"] = rqGames
	}
	if m["mygame"] != nil {
		m["mygamefactions"] = g.Factions()
		if g.Turn == 0 {
			m["mygamerequests"] = g.AllRequests()
		}
	}
	v.Apply(TPHOME, w)
}
