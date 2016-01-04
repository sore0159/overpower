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
	if r.Method == "POST" {
		if DBLOCK {
			http.Error(w, "GAME DOWN FOR DAYLY MAINT: 10-20MIN", http.StatusInternalServerError)
			return
		}
		action := r.FormValue("action")
		switch action {
		case "nextturn", "setautos":
			if !hasG {
				http.Error(w, "USER HAS NO GAME TO PROGRESS", http.StatusBadRequest)
				return
			}
			mp, ok := GetInts(r, "turn")
			if !ok {
				http.Error(w, "MALFORMED TURN DATA", http.StatusBadRequest)
				return
			}
			turn, ok := mp["turn"]
			if !ok || turn != g.Turn() {
				http.Error(w, "BAD TURN DATA", http.StatusBadRequest)
				return
			}
			if action == "nextturn" {
				if g.Turn() < 1 {
					http.Error(w, "GAME NOT YET BEGUN", http.StatusBadRequest)
					return
				}
				if !OPDB.RunGameTurn(g) {
					http.Error(w, "DATABASE ERROR RUNNING GAME TURN", http.StatusInternalServerError)
					return
				}
			} else {
				dayBool := [7]bool{}
				dayBool[0] = r.FormValue("sunday") == "on"
				dayBool[1] = r.FormValue("monday") == "on"
				dayBool[2] = r.FormValue("tuesday") == "on"
				dayBool[3] = r.FormValue("wednesday") == "on"
				dayBool[4] = r.FormValue("thursday") == "on"
				dayBool[5] = r.FormValue("friday") == "on"
				dayBool[6] = r.FormValue("saturday") == "on"
				g.SetAutoDays(dayBool)
				if !OPDB.UpdateGame(g) {
					http.Error(w, "DATABASE ERROR CHANGING GAME AUTOTURNS", http.StatusInternalServerError)
					return
				}
			}
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
			facName := r.FormValue("facname")
			if facName != "" && !ValidText(facName) {
				h.SetError("INVALID FACTION NAME")
				break
			}
			if !OPDB.MakeGame(h.User.String(), name, password) {
				h.SetError("DATABASE ERROR IN GAME CREATION")
				break
			}
			if facName != "" {
				g, ok := OPDB.GetGameFor(h.User.String())
				if !ok {
					h.SetError("DATABASE ERROR IN FACTION CREATION")
					break
				}
				if !OPDB.MakeFaction(g.Gid(), h.User.String(), facName) {
					h.SetError("DATABASE ERROR IN GAME CREATION")
					break
				}
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
	m["user"] = h.User.String()
	if hasG {
		m["game"] = g
		m["active"] = g.Turn() > 0
	}
	if gHasF {
		m["gfactions"] = gFacs
	}
	days := g.AutoDays()
	var any bool
	for _, b := range days {
		if b {
			any = true
			break
		}
	}
	if !any {
		m["noauto"] = true
	}
	oFacs, oHasF := OPDB.GetAllOwnerFactions(h.User.String())
	if oHasF {
		facGames := make([]overpower.Game, len(oFacs))
		for i, f := range oFacs {
			g, ok := OPDB.GetGame(f.Gid())
			if !ok {
				http.Error(w, "DATABASE ERROR FETCHING GAMES", http.StatusInternalServerError)
				return
			}
			facGames[i] = g
		}
		m["ofactions"] = oFacs
		m["ogames"] = facGames
	}
	h.Apply(TPOPHOME, w)
}
