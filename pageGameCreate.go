package main

import (
	"fmt"
	"net/http"
)

var TPGAMECR = MixTemp("frame", "titlebar", "gamecreate")

func (g *Game) GameCreate(w http.ResponseWriter, r *http.Request, v *View) {
	if r.Method == "POST" {
		switch r.FormValue("action") {
		case "startgame":
			if !Admin(v.User) {
				v.SetError("Only admins allowed here!")
				goto DONE
			}
			g.MakeGame()
			g.Save()
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		case "signup":
			if v.User == "" {
				v.SetError("Must be logged in!")
				goto DONE
			}
			for _, uName := range g.Signup {
				if uName == v.User {
					v.SetError("Already signed up!")
					goto DONE
				}
			}
			facName := r.FormValue("factionname")
			if !ValidUserName(facName) {
				v.SetError("Invalid faction name!")
				goto DONE
			}
			if err := g.UserSignUp(v.User, facName); err == nil {
				g.Save()
				http.Redirect(w, r, r.URL.Path, http.StatusFound)
				return
			} else {
				fmt.Println(err)
				v.SetError("Problem signing up:", err)
				goto DONE
			}
		}
	}
DONE:
	app := map[string]interface{}{}
	app["signups"] = g.Signup
	if v.User != "" {
		app["user"] = v.User
	}
	if Admin(v.User) {
		app["admin"] = true
	}
	for _, uName := range g.Signup {
		if uName == v.User {
			app["signed"] = true
			break
		}
	}
	v.SetApp(app)
	v.Apply(TPGAMECR, w)
}

func (g *Game) UserSignUp(user string, facName string) error {
	g.Signup[facName] = user
	return nil
}
