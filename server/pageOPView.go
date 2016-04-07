package main

import (
	"mule/overpower"
	"net/http"
)

var (
	TPOPVIEW = MixTemp("frame", "titlebar", "opview")
)

// /overpower/view/GID
func pageOPView(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	if len(h.Path) != 4 {
		h.HandleUserError(w, "INVALID PATH")
		return
	}
	gid, ok := h.IntAt(3)
	if !ok {
		h.HandleUserError(w, "UNPARSABLE GAMEID")
		return
	}
	games, err := h.M.Game().SelectWhere(h.GID(gid))
	if my, bad := Check(err, "resource failure on OP view page", "resource", "game", "gid", gid); bad {
		h.HandleServerError(w, my)
		return
	}
	if len(games) == 0 {
		h.HandleUserError(w, "NO GAME MATCHING GAMEID")
		return
	}
	g := games[0]

	m := h.DefaultApp()
	facs, err := h.M.Faction().SelectWhere(h.GID(gid))
	if my, bad := Check(err, "resource failure", "page", "opview", "resource", "factions", "gid", gid); bad {
		h.HandleServerError(w, my)
		return
	}
	m["factions"] = facs
	m["active"] = g.Turn() > 0
	var ownedF overpower.FactionDat
	if h.LoggedIn {
		m["user"] = h.User.String()
		for i, f := range facs {
			if f.Owner() == h.User.String() {
				m["ownedf"] = f
				facs[i] = nil
				ownedF = f
				break
			}
		}
	}
	if len(facs) > 1 || (len(facs) > 0 && ownedF == nil) {
		m["otherf"] = true
	}
	m["game"] = g
	if r.Method == "POST" {
		if !h.LoggedIn {
			h.HandleUserError(w, "NOT LOGGED IN")
			return
		}
		var errS, errU error
		action := r.FormValue("action")
		switch action {
		case "setdone":
			errS, errU = h.CommandSetTurnDone(g, facs, ownedF)
		case "dropfac":
			errS, errU = h.CommandDropFaction(g, ownedF)
		case "newfac":
			errS, errU = h.CommandNewFaction(g, facs, ownedF, r.FormValue("password"), r.FormValue("facname"))
		default:
			errU = NewError("UNKNOWN ACTION TYPE")
		}
		if my, bad := Check(errS, "page op view action failure", "action", action, "game", g); bad {
			h.HandleServerError(w, my)
		} else if errU != nil {
			h.HandleUserError(w, errU.Error())
		} else {
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
		}
		return
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
	h.Apply(TPOPVIEW, w)
}
