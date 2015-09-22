package main

import (
	//	"fmt"
	"mule/planetattack/attack"
	"net/http"
	"strconv"
	"strings"
)

var TEMPGAME = attack.MakeGame([]string{"P1", "P2", "P3"})
var TPTV = MixTemp("frame", "textview")

func tempGamePage(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) > 3 {
		http.Redirect(w, r, strings.Join(pathParts[:3], "/"), http.StatusFound)
		return
	}
	/*factionStr := pathParts[len(pathParts)-1]
	if factionStr == "" {
		http.Redirect(w, r, r.URL.Path+"1", http.StatusFound)
		return
	}
	factionID, err := strconv.Atoi(factionStr)
	if err != nil || (factionID != 1 && factionID != 2 && factionID != 3) {
		http.Redirect(w, r, "/view/1", http.StatusFound)
		return
	}
	//fmt.Fprintf(w, "GAME VIEW FOR FACTION"+factionStr)
	*/
	f := TEMPGAME.Factions[1]
	recenterStr := pathParts[len(pathParts)-1]
	if recenterStr != "" {
		CoordS := strings.Split(recenterStr, "x")
		if len(CoordS) != 2 {
			http.Redirect(w, r, "/view/", http.StatusFound)
			return
		}
		x, err1 := strconv.Atoi(CoordS[0])
		y, err2 := strconv.Atoi(CoordS[1])
		if err1 != nil || err2 != nil {
			http.Redirect(w, r, "/view/", http.StatusFound)
			return
		}
		f.CenterTV([2]int{x, y})
	}
	err := TPTV.ExecuteTemplate(w, "frame", f.TV.OrderedCoords)
	if err != nil {
		Log("Error executing gameview template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
