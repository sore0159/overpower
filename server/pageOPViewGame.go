package main

import (
	"mule/overpower"
	"net/http"
	"strings"
)

var (
	TPVIEW = MixTemp("frame", "titlebar", "view")
)

func (h *Handler) pageOPViewGame(w http.ResponseWriter, r *http.Request, g overpower.Game) {
	m := h.DefaultApp()
	facs, ok := OPDB.GetGidFactions(g.Gid())
	if ok {
		m["factions"] = facs
	}
	m["active"] = g.Turn() > 0
	var ownedF overpower.Faction
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
	m["game"] = g
	if r.Method == "POST" {
		if !h.LoggedIn {
			http.Error(w, "NOT LOGGED IN", http.StatusBadRequest)
			return
		}
		action := r.FormValue("action")
		switch action {
		case "dropfac":
			if g.Turn() > 0 {
				http.Error(w, "GAME IN PROGRESS", http.StatusBadRequest)
				return
			}
			if ownedF == nil {
				http.Error(w, "USER HAS NO FACTION FOR THIS GAME", http.StatusBadRequest)
				return
			}
			if !OPDB.DropFaction(g.Gid(), ownedF.Fid()) {
				http.Error(w, "DATABASE ERROR DELETING FACTION", http.StatusInternalServerError)
				return
			}
		case "newfac":
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
			if !OPDB.MakeFaction(g.Gid(), h.User.String(), facName) {
				http.Error(w, "DATABASE ERROR CREATING FACTION", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "UNKNOWN ACTION TYPE", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	}
	h.Apply(TPVIEW, w)
}
