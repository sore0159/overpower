package main

import (
	"mule/overpower"
	"net/http"
)

var (
	TPOPHOME = MixTemp("frame", "titlebar", "ophome")
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
	var hasG bool
	var g overpower.GameDat

	games, err := h.M.Game().Select("owner", h.User.String())
	if my, bad := Check(err, "resource failure on OP home page", "resource", "game", "owner", h.User.String()); bad {
		h.HandleServerError(w, my)
		return
	}
	if len(games) != 0 {
		g = games[0]
		hasG = true
	}

	var gFacs []overpower.FactionDat
	var gHasF bool
	if hasG {
		gFacs, err = h.M.Faction().SelectWhere(h.GID(g.GID()))
		if my, bad := Check(err, "resource error in homepage", "resource", "faction", "user", h.User, "gid", g.GID()); bad {
			h.HandleServerError(w, my)
			return
		}
		gHasF = len(gFacs) > 0
	}
	if r.Method == "POST" {
		if DBLOCK {
			http.Error(w, "GAME DOWN FOR DAYLY MAINT: 10-20MIN", http.StatusInternalServerError)
			return
		}
		action := r.FormValue("action")
		var errS, errU error
		switch action {
		case "nextturn":
			turn := r.FormValue("turn")
			errS, errU = h.CommandForceTurn(g, turn)
		case "setautos":
			dayBool := [7]bool{}
			dayBool[0] = r.FormValue("sunday") == "on"
			dayBool[1] = r.FormValue("monday") == "on"
			dayBool[2] = r.FormValue("tuesday") == "on"
			dayBool[3] = r.FormValue("wednesday") == "on"
			dayBool[4] = r.FormValue("thursday") == "on"
			dayBool[5] = r.FormValue("friday") == "on"
			dayBool[6] = r.FormValue("saturday") == "on"

			errS, errU = h.CommandSetAutos(g, dayBool)

		case "startgame":
			exodus := r.FormValue("exodus") == "on"
			errS, errU = h.CommandStartGame(g, gFacs, exodus)
		case "newgame":
			gamename, password := r.FormValue("gamename"), r.FormValue("password")
			facname, towin := r.FormValue("facname"), r.FormValue("towin")
			errS, errU = h.CommandNewGame(g, password, gamename, facname, towin)
		case "dropgame":
			errS, errU = h.CommandDropGame(g)
		default:
			errU = NewError("UNKNOWN ACTION TYPE")
		}
		if my, bad := Check(errS, "page op home action failure", "action", action, "user", h.User.String(), "game", g); bad {
			h.HandleServerError(w, my)
		} else if errU != nil {
			h.HandleUserError(w, errU.Error())
		} else {
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
		}
		return
	}
	m := h.DefaultApp()
	m["user"] = h.User.String()
	if hasG {
		m["game"] = g
		m["active"] = g.Turn() > 0
	}
	if gHasF {
		m["gfactions"] = gFacs
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
	}
	oFacs, err := h.M.Faction().Select("owner", h.User.String())
	if my, bad := Check(err, "resource error in homepage", "resource", "faction", "owner", h.User); bad {
		h.HandleServerError(w, my)
		return
	}
	oHasF := len(oFacs) > 0
	if oHasF {
		facGames := make([]overpower.GameDat, 0, len(oFacs))
		for _, f := range oFacs {
			games, err := h.M.Game().SelectWhere(h.GID(f.GID()))
			if my, bad := Check(err, "resource error in homepage", "gid", f.GID(), "fac", f, "owner", h.User); bad {
				h.HandleServerError(w, my)
				return
			} else {
				facGames = append(facGames, games...)
			}
		}
		m["ofactions"] = oFacs
		m["ogames"] = facGames
	}
	h.Apply(TPOPHOME, w)
}
