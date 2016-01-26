package main

import (
	"net/http"
)

var (
	TPINDEX     = MixTemp("frame", "titlebar", "index")
	TPOPINDEX   = MixTemp("frame", "titlebar", "opindex")
	TPVIEWINDEX = MixTemp("frame", "titlebar", "viewindex")
)

func pageMainIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	h := MakeHandler(w, r)
	h.Apply(TPINDEX, w)
}

func pageOPIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/overpower/" {
		http.Redirect(w, r, "/overpower/", http.StatusFound)
		return
	}
	h := MakeHandler(w, r)
	m := h.DefaultApp()
	m["loggedin"] = h.LoggedIn
	h.Apply(TPOPINDEX, w)
}

func (h *Handler) pageOPViewIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/overpower/view/" {
		http.Redirect(w, r, "/overpower/view/", http.StatusFound)
		return
	}
	games, err := OPDB.GetGames()
	if my, bad := Check(err, "resource failure on view index"); bad {
		Log(my)
		h.SetError("DATABASE ERROR")
	} else {
		h.SetApp(games)
	}
	h.Apply(TPVIEWINDEX, w)
}
