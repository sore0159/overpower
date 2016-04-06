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
		http.Error(w, "INVALID PATH", http.StatusBadRequest)
		return
	}
	gid, ok := h.IntAt(3)
	if !ok {
		http.Error(w, "UNPARSABLE GAMEID", http.StatusBadRequest)
		return
	}
	games, err := h.M.Game().SelectWhere(h.GID(gid))
	if my, bad := Check(err, "resource failure on OP view page", "resource", "game", "gid", gid); bad {
		Bail(w, my)
		return
	}
	if len(games) == 0 {
		http.Error(w, "NO GAME MATCHING GAMEID", http.StatusBadRequest)
		return
	}
	g := games[0]

	m := h.DefaultApp()
	facs, err := h.M.Faction().SelectWhere(h.GID(gid))
	if my, bad := Check(err, "resource failure", "page", "opview", "resource", "factions", "gid", gid); bad {
		Bail(w, my)
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
			http.Error(w, "NOT LOGGED IN", http.StatusBadRequest)
			return
		}
		action := r.FormValue("action")
		switch action {
		case "setdone":
			//if ok := h.SetTurnDone(w, r, g, ownedF); !ok {
			//return
			//}
		case "dropfac":
			/*
				if g.Turn() > 0 {
					http.Error(w, "GAME IN PROGRESS", http.StatusBadRequest)
					return
				}
				if ownedF == nil {
					http.Error(w, "USER HAS NO FACTION FOR THIS GAME", http.StatusBadRequest)
					return
				}
				err = OPDB.DropFactions("gid", g.Gid(), "fid", ownedF.Fid())
				if my, bad := Check(err, "data update failure", "data", ownedF, "page", "opviewgame"); bad {
					Bail(w, my)
					return
				}
			*/
		case "newfac":
			/*
				if g.Turn() > 0 {
					http.Error(w, "GAME IN PROGRESS", http.StatusBadRequest)
					return
				}
				if ownedF != nil {
					http.Error(w, "USER ALREADY HAS FACTION FOR THIS GAME", http.StatusBadRequest)
					return
				}
				if g.HasPW() {
					passwd := r.FormValue("password")
					if !ValidText(passwd) || !g.IsPwd(passwd) {
						http.Error(w, "BAD PASSWORD", http.StatusBadRequest)
						return
					}
				}
				facName := r.FormValue("facname")
				if !ValidText(facName) {
					http.Error(w, "BAD FACTION NAME", http.StatusBadRequest)
					return
				}
				lwFName := strings.ToLower(facName)
				for _, f := range facs {
					if strings.ToLower(f.Name()) == lwFName {
						http.Error(w, "FACTION NAME ALREADY IN USE FOR THIS GAME", http.StatusBadRequest)
						return
					}
				}
				err = OPDB.MakeFaction(g.Gid(), h.User.String(), facName)
				if my, bad := Check(err, "data creation error", "type", "faction", "gid", g.Gid(), "user", h.User, "facname", facName); bad {
					Bail(w, my)
					return
				}
			*/
		default:
			http.Error(w, "UNKNOWN ACTION TYPE", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
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
