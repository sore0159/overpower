package main

import (
	"mule/overpower"
	"net/http"
)

func SetupMux() {
	http.HandleFunc("/", pageMainIndex)
	http.HandleFunc("/favicon.ico", imgFavIcon)
	http.HandleFunc("/auth/", pageAuthIndex)
	http.HandleFunc("/auth/logout", pageAuthLogout)
	http.HandleFunc("/auth/login", pageAuthLogin)
	http.HandleFunc("/auth/create", pageAuthCreate)
	http.HandleFunc("/overpower/", pageOPIndex)
	http.HandleFunc("/overpower/home", pageOPHome)
	http.HandleFunc("/overpower/quit/", pageOPQuit)
	http.HandleFunc("/overpower/view/", muxView)
	http.HandleFunc("/overpower/json/", apiJson)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("STATIC/"))))
}

// /overpower/view/GAMEID/
func muxView(w http.ResponseWriter, r *http.Request) {
	if DBLOCK {
		http.Error(w, "GAME DOWN FOR DAYLY MAINT: 10-20MIN", http.StatusInternalServerError)
		return
	}
	h := MakeHandler(w, r)
	lastFull := h.LastFull()
	if lastFull < 3 {
		h.pageOPViewIndex(w, r)
		return
	}
	gid, ok := h.IntAt(3)
	if !ok {
		http.Error(w, "INVALID GAMEID", http.StatusBadRequest)
		return
	}
	g, err := OPDB.GetGame("gid", gid)
	if err == ErrNoneFound {
		http.Error(w, "GAMEID NOT FOUND", http.StatusNotFound)
		return
	} else if my, bad := Check(err, "resource aquisition error", "gid", gid); bad {
		Bail(w, my)
		return
	}
	if lastFull == 3 {
		h.pageOPViewGame(w, r, g)
		return
	}
	switch h.Path[4] {
	case "play", "reports":
		if !h.LoggedIn {
			http.Error(w, "NOT LOGGED IN", http.StatusBadRequest)
			return
		}
		facs, err := OPDB.GetFactions("gid", gid)
		if my, bad := Check(err, "resource aquisition error", "gid", gid); bad {
			Bail(w, my)
			return
		}
		var f overpower.Faction
		for _, test := range facs {
			if test.Owner() == h.User.String() {
				f = test
				break
			}
		}
		if f == nil {
			http.Error(w, "USER HAS NO FACTION IN THIS GAME", http.StatusBadRequest)
			return
		}
		switch h.Path[4] {
		case "play":
			h.pageCanvas(w, r, g, f, facs)
		case "reports":
			h.pageOPReports(w, r, g, f, facs)
		}
		return
	default:
		http.Redirect(w, r, h.NewPath(4), http.StatusFound)
		return
	}
	h.pageOPViewIndex(w, r)
	return
}

func imgFavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "STATIC/img/favicon.ico")
}
