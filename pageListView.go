package main

import (
	"fmt"
	"mule/planetattack/attack"
	"net/http"
)

var (
	TPLIST = MixTemp("frame", "titlebar", "listview", "listcoord")
)

func (g *Game) FactionView(w http.ResponseWriter, r *http.Request, f *attack.Faction, v *View) {
	temp := TPLIST
	if len(v.path) > 3 {
		if v.path[3] == "maps" {
			MapView(w, r, f, v)
			return
		} else if v.path[3] == "reports" {
			ReportView(w, r, f, v)
			return
		}
	}
	if r.Method == "POST" {
		action := r.FormValue("action")
		var err error
		switch action {
		case "togglefilter":
			err = UserToggleFilter(r, f)
		case "recenter":
			err = UserRecenter(r, f)
		case "setlorder":
			err = UserSetLaunchOrder(r, f)
		default:
			http.Error(w, "Action "+action+" not found", http.StatusInternalServerError)
			return
		}
		if err == nil {
			g.Save()
		} else {
			fmt.Println(err)
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	}
	v.SetApp(f)
	v.Apply(temp, w)
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
	m, err := GetInts(r, "amount", "tarID", "sourceID")
	if err != nil {
		return err
	}
	amount := m["amount"]
	source, ok := f.GetPlanetView(m["sourceID"])
	if !ok || !source.Yours {
		return makeE("Bad launch order:", m["sourceID"], "not found/owned by", f.Name)
	}
	target, ok := f.GetPlanetView(m["tarID"])
	if !ok {
		return makeE("Bad launch order:", m["tarID"], "not found by", f.Name)
	}
	numAvail := f.NumAvail(source.Location)
	if o, ok := f.BuildOrders[[4]int{source.Location[0], source.Location[1], target.Location[0], target.Location[1]}]; ok {
		numAvail += o.Size
	}
	if numAvail < amount {
		return makeE("Bad launch order:", amount, "not avail at", source.Location)
	}
	f.SetOrder(amount, source, target)
	return nil
}

func UserToggleFilter(r *http.Request, f *attack.Faction) error {
	f.ToggleFilter()
	return nil
}
