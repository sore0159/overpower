package main

import (
	"net/http"
)

var TPINDEX = MixTemp("frame", "titlebar", "index")

func IndexPage(w http.ResponseWriter, r *http.Request) {
	v := MakeView(r)
	m := map[string]string{}
	v.SetApp(m)
	v.Apply(TPINDEX, w)
}
