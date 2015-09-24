package main

import (
	"fmt"
	"mule/planetattack/attack"
	"net/http"
	"strings"
)

var TEMPGAME = &Game{attack.MakeGame([]string{"P1", "P2", "P3"})}

func gameMux(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathParts := strings.Split(path, "/")
	if len(pathParts) > 1 && pathParts[1] == "view" {
		if len(pathParts) == 2 || pathParts[2] == "" {
			TEMPGAME.GameView(w, r)
			return
		} else if len(pathParts) > 3 {
			http.Redirect(w, r, strings.Join(pathParts[:3], "/"), http.StatusFound)
			return
		} else {
			fName := pathParts[2]
			for _, f := range TEMPGAME.Factions {
				if f.Name == fName {
					if AuthFaction(r, f) {
						TEMPGAME.FactionView(w, r, f)
					} else {
						http.Error(w, "User not authorized for faction "+f.Name, http.StatusInternalServerError)
					}
					return

				}
			}
			http.Redirect(w, r, "/view/", http.StatusFound)
			return
		}
	}
	userName := CookieUserName(r)
	if userName == "" {
		LoginPage(w, r)
		return
	}
	fmt.Fprint(w, "LOGGED IN USER:", userName, "PATH:", path)
}

func AuthFaction(r *http.Request, f *attack.Faction) bool {
	return true
}
