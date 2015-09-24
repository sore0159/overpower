package main

import (
	//	"fmt"
	"mule/planetattack/attack"
	"net/http"
	//"strconv"
	//"strings"
)

var TPGAMEV = MixTemp("frame", "gameview")

func (g *Game) GameView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		action := r.FormValue("action")
		var gErr error
		if action == "toggledone" {
			m, err := GetInts(r, "fID")
			if err != nil {
				Log("Failed to set faction done in GameView")
				return
			}
			if f, ok := g.Factions[m["fID"]]; ok {
				if AuthFaction(r, f) {

					gErr = UserToggleDone(r, g, f)
				} else {
					err := Log("Failed to toggle faction", f.Name, "done in GameView, user not authorized")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				err := Log("Failed to set faction", m["fID"], "done in GameView, faction not found")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
		if gErr == nil {
			//g.Save()
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	}
	Apply(TPGAMEV, w, g)
}

func UserToggleDone(r *http.Request, g *Game, f *attack.Faction) error {
	g.UserToggleDone(f.FactionID)
	return nil
}
