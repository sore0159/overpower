package main

import (
	"mule/planetattack/attack"
	"net/http"
	"strconv"
)

var TPSREPORT = MixTemp("frame", "titlebar", "reports")

func ReportView(w http.ResponseWriter, r *http.Request, f *attack.Faction, v *View) {
	if f.View.Turn == 1 {
		http.Redirect(w, r, v.newpath(3), http.StatusFound)
		return
	}
	var turn int
	switch len(v.path) {
	case 5:
		if v.path[4] == "" {
			turn = f.View.Turn - 1
		} else if t, err := strconv.Atoi(v.path[4]); err == nil && t < f.View.Turn && t > 0 {
			turn = t
		} else {
			http.Redirect(w, r, v.newpath(4), http.StatusFound)
			return
		}
	case 4:
		turn = f.View.Turn - 1
	default:
		http.Redirect(w, r, v.newpath(5), http.StatusFound)
		return
	}
	m := map[string]interface{}{"reports": f.Reports[turn-1], "fac": f.Name, "view": turn}
	if turn < f.View.Turn-1 {
		m["next"] = turn + 1
	}
	if turn > 1 {
		m["prev"] = turn - 1
	}
	v.SetApp(m)
	v.Apply(TPSREPORT, w)
}
