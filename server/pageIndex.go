package main

import (
	"net/http"
)

var (
	TPINDEX     = MixTemp("frame", "titlebar", "index")
	TPOPINDEX   = MixTemp("frame", "titlebar", "opindex")
	TPAUTHINDEX = MixTemp("frame", "titlebar", "authindex")
)

func pageIndex(w http.ResponseWriter, r *http.Request) {
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
	h.Apply(TPOPINDEX, w)
}

func pageAuthIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth/" {
		http.Redirect(w, r, "/auth/", http.StatusFound)
		return
	}
	h := MakeHandler(w, r)
	h.Apply(TPAUTHINDEX, w)
}
