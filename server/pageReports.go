package main

import (
	"fmt"
	"mule/overpower"
	"net/http"
)

var (
	TPREPORTS = MixTemp("frame", "titlebar", "reports")
)

func (h *Handler) pageOPReports(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction, facs []overpower.Faction) {
	turn, ok := h.IntAt(5)
	if !ok {
		var pathStr = fmt.Sprintf("%s/%d", h.NewPath(5), g.Turn()-1)
		http.Redirect(w, r, pathStr, http.StatusFound)
		return
	}
	if turn < 1 || turn >= g.Turn() {
		http.Error(w, "BAD TURN FOR REPORTS", http.StatusBadRequest)
		return
	}
	if g.Turn() < 1 {
		http.Error(w, "GAME NOT YET IN PROGRESS", http.StatusBadRequest)
		return
	}
	rp, err := OPDB.GetReport("gid", g.Gid(), "fid", f.Fid(), "turn", turn)

	if err == ErrNoneFound {
	} else if my, bad := Check(err, "resource error", "page", "opreports", "resource", "reports", "gid", g.Gid(), "fid", f.Fid(), "turn", turn); bad {
		Bail(w, my)
		return
	}
	m := h.DefaultApp()
	m["game"] = g
	m["f"] = f
	if ok {
		m["rp"] = rp
	}
	m["turn"] = turn
	if turn > 1 {
		m["prev"] = turn - 1
	}
	if turn < g.Turn()-1 {
		m["next"] = turn + 1
	}
	h.Apply(TPREPORTS, w)
}
