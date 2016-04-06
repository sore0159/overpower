package main

import (
	"net/http"
)

var (
	TPOPLIST = MixTemp("frame", "titlebar", "oplist")
)

func pageOPList(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/overpower/list/" {
		http.Redirect(w, r, "/overpower/list/", http.StatusFound)
		return
	}
	h := MakeHandler(w, r)
	games, err := h.M.Game().Select()
	if my, bad := Check(err, "resource failure on OP list page"); bad {
		Log(my)
		h.SetError("DATABASE ERROR")
	} else {
		h.SetApp(games)
	}
	h.Apply(TPOPLIST, w)
}
