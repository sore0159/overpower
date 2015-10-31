package main

import (
	"mule/planetattack"
	"net/http"
	"strconv"
)

var (
	TPGAME    = MixTemp("frame", "titlebar", "game")
	TPPREGAME = MixTemp("frame", "titlebar", "pregame")
)

type GameHandler struct {
	g *Game
	f *planetattack.Faction
	v *View
	u string
	w http.ResponseWriter
	r *http.Request
}

func GamePage(w http.ResponseWriter, r *http.Request) {
	v := MakeView(r)
	userN := CookieUserName(r)
	g := NewGame()
	var err error
	g.Gid, err = strconv.Atoi(v.path[2])
	if err != nil || g.Gid == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if !g.Select() {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if len(v.path) > 3 {
		if userN == "" || g.Turn < 1 {
			http.Redirect(w, r, v.newpath(3), http.StatusFound)
			return
		}
		f, ok := g.GetFaction(userN)
		if !ok {
			http.Redirect(w, r, v.newpath(3), http.StatusFound)
			return
		}
		turnStr := r.FormValue("turn")
		if turnStr != "" {
			turn, err := strconv.Atoi(turnStr)
			if err != nil {
				Log("Turn form data error:", err)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			if turn != g.Turn {
				v.SetError("FORM SUBMISSION FROM TURN", turn, ":CURRENT TURN", g.Turn)
				r.Method = ""
			}
		}
		h := &GameHandler{g: g, f: f, v: v, w: w, r: r, u: userN}
		switch v.path[3] {
		case "play":
			h.PlayPage()
		default:
			http.Redirect(w, r, v.newpath(3), http.StatusFound)
		}
		return
	}
	m := map[string]interface{}{"userN": userN}
	v.SetApp(m)
	m["game"] = g
	if r.Method == "POST" {
		if userN == "" {
			v.SetError("NOT LOGGED IN")
			goto GET
		}
		action := r.FormValue("action")
		switch action {
		case "toggledone":
			if g.Turn < 1 {
				v.SetError("GAME HAS NOT YET BEGUN")
				goto GET
			}
			ints, err := GetInts(r, "turn")
			if err != nil {
				v.SetError("ERROR IN FORM DATA:", err)
				goto GET
			}
			if g.Turn != ints["turn"] {
				v.SetError("FORM SUBMISSION FROM TURN", ints["turn"], ":CURRENT TURN", g.Turn)
				goto GET
			}
			f, ok := g.GetFaction(userN)
			if !ok {
				v.SetError("USER", userN, "HAS NO FACTION FOR GAME", g.Name)
				goto GET
			}
			var done bool
			switch r.FormValue("done") {
			case "true":
				done = true
			case "false":
			default:
				v.SetError("ERROR IN FORM DATA: INVALID done VALUE", r.FormValue("done"))
				goto GET
			}
			err = g.SetDone(f.Fid, done)
			if err != nil {
				v.SetError("FACTION TOGGLE DONE ERROR:", err)
				goto GET
			}
		case "leavegame":
			if g.Turn > 0 {
				v.SetError("CANNOT DELETE FACTION AFTER GAME BEGINS")
				goto GET
			}
			f, ok := g.GetFaction(userN)
			if !ok {
				v.SetError("USER", userN, "HAS NO FACTION FOR GAME", g.Name)
				goto GET
			}
			err := f.Delete()
			if err != nil {
				v.SetError("FACTION DELETION ERROR:", err)
				goto GET
			}

		case "joingame":
			if g.Password != "" {
				password := r.FormValue("password")
				if password != g.Password {
					v.SetError("INVALID PASSWORD TO JOIN GAME", g.Name)
					goto GET
				}
			}
			fName := r.FormValue("facname")
			if !ValidText(fName) {
				v.SetError("INVALID FACTION NAME")
				goto GET
			}
			f, ok := g.GetFaction(userN)
			if ok {
				v.SetError("USER", userN, "ALREADY HAS FACTION", f.Name, "FOR GAME", g.Name)
				goto GET
			}
			_, err := g.MakeFaction(userN, fName)
			if err != nil {
				v.SetError("FACTION CREATION ERROR:", err)
				goto GET
			}
		default:
			Log("invalid action for gamepage:", action)
			v.SetError("INVALID ACTION:", action)
			goto GET
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	}
GET:
	facs := g.Factions()
	m["factions"] = facs
	if userN != "" {
		for i, f := range facs {
			if f.Owner == userN {
				m["myfac"] = f
				facs[i] = nil
			}
		}
	}
	if g.Turn > 0 {
		v.Apply(TPGAME, w)
	} else {
		v.Apply(TPPREGAME, w)
	}
}
