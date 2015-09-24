package main

import (
	"mule/planetattack/attack"
	"net/http"
)

var TPLIST = MixTemp("frame", "listview", "listcoord")

func (g *Game) FactionView(w http.ResponseWriter, r *http.Request, f *attack.Faction) {
	if r.Method == "POST" {
		action := r.FormValue("action")
		var err error
		switch action {
		case "recenter":
			err = UserRecenter(r, f)
		case "setlorder":
			err = UserSetLaunchOrder(r, f)
		default:
			http.Error(w, "Action "+action+" not found", http.StatusInternalServerError)
			return
		}
		if err == nil {
			//g.Save()
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	}
	Apply(TPLIST, w, f)
}

func UserRecenter(r *http.Request, f *attack.Faction) error {
	m, err := GetInts(r, "adjX", "adjY")
	if err != nil {
		return err
	}
	f.CenterTV([2]int{m["adjX"], m["adjY"]})
	return nil
}

func UserSetLaunchOrder(r *http.Request, f *attack.Faction) error {
	m, err := GetInts(r, "amount", "tarX", "tarY", "sourceX", "sourceY")
	if err != nil {
		return err
	}
	amount := m["amount"]
	sCoord := [2]int{m["sourceX"], m["sourceY"]}
	if source, ok := f.View.PlanetGrid[sCoord]; !ok || !source.Yours || f.NumAvail(sCoord) < amount {
		return makeE("Bad launch order", m)
	}
	tCoord := [2]int{m["tarX"], m["tarY"]}
	if _, ok := f.View.PlanetGrid[tCoord]; !ok {
		return makeE("No planet found at target coord", tCoord)
	}
	f.SetOrder(amount, sCoord, tCoord)
	return nil
}
