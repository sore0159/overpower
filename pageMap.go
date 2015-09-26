package main

import (
	"image/png"
	"mule/planetattack/attack"
	"mule/planetattack/mapping"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func MapView(w http.ResponseWriter, r *http.Request, f *attack.Faction, v *View) {
	var turn int
	//1/2view/3fac/4map/5turn
	if len(v.path) == 5 {
		iStr := ""
		for i, char := range v.path[4] {
			if i == 3 {
				break
			}
			iStr += string(char)
		}
		if t, err := strconv.Atoi(iStr); err == nil && t <= f.View.Turn && t > 0 {
			turn = t
		} else {
			turn = f.View.Turn
		}
	} else {
		turn = f.View.Turn
	}
	fileName := mapping.Filename(f.Name, turn)
	list := strings.Split(fileName, "/")
	pngName := list[len(list)-1]
	if len(v.path) < 5 {
		http.Redirect(w, r, r.URL.Path+"/"+pngName, http.StatusFound)
		return
	} else if v.path[4] != pngName {
		v.path[4] = pngName
		http.Redirect(w, r, v.newpath(5), http.StatusFound)
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		http.Error(w, Log(err).Error(), http.StatusInternalServerError)
		return
	}
	img, err := png.Decode(file)
	if err != nil {
		http.Error(w, Log(err).Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "image/png")
	err = png.Encode(w, img)
	if err != nil {
		Log(err)
	}
	return
}
