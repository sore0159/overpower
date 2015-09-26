package main

import (
	//	"fmt"
	"mule/planetattack/attack"
	"net/http"
	//"strconv"
	//"strings"
)

var TPGAMEV = MixTemp("frame", "titlebar", "gameview")

func (g *Game) GameView(w http.ResponseWriter, r *http.Request, v *View) {
	if r.Method == "POST" {
		var gErr error
		switch r.FormValue("action") {
		case "runturn":
			if Admin(v.User) {
				gErr = UserRunTurn(g)
			}
		case "toggledone":
			m, err := GetInts(r, "fID")
			if err != nil {
				Log("Failed to set faction done in GameView")
				v.SetError("Failed to set faction done in GameView")
				goto DONE
			}
			if f, ok := g.Factions[m["fID"]]; ok {
				if _, ok := g.AuthFaction(v.User, f.Name); ok {
					gErr = UserToggleDone(g, f)
				} else {
					err := Log("Failed to toggle faction", f.Name, "done in GameView, user not authorized")
					v.SetError(err.Error())
					goto DONE
				}
			} else {
				err := Log("Failed to set faction", m["fID"], "done in GameView, faction not found")
				v.SetError(err.Error())
				goto DONE
			}

		}
		if gErr == nil {
			g.Save()
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		}
	}
DONE:
	m := map[string]interface{}{}
	m["turn"] = g.Sector.Turn
	fac := map[string]bool{}
	m["facs"] = fac
	for fID, f := range g.Factions {
		fac[f.Name] = f.TurnDone
		if fID == g.Users[v.User] {
			m["userfac"] = f.Name
			m["userFID"] = fID
		}
	}
	if Admin(v.User) {
		m["admin"] = true
	}
	v.SetApp(m)
	v.Apply(TPGAMEV, w)
}

func UserToggleDone(g *Game, f *attack.Faction) error {
	g.UserToggleDone(f.FactionID)
	return nil
}

func UserRunTurn(g *Game) error {
	g.RunTurn()
	return nil
}
