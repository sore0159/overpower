package main

import (
	"net/http"
)

var TPHOME = MixTemp("frame", "titlebar", "home")

func HomePage(w http.ResponseWriter, r *http.Request) {
	v := MakeView(r)
	m := map[string]string{}
	v.SetApp(m)
	v.Apply(TPHOME, w)
}
