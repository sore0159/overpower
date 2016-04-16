package main

import (
	"fmt"
	"mule/mybad"
	"net/http"
)

var (
	TPOPERROR = MixTemp("frame", "titlebar", "operror")
)

func (h *Handler) pageOPError(w http.ResponseWriter, r *http.Request, errS, errU error) {
	m := h.DefaultApp()
	m["ServerError"] = errS
	m["UserError"] = errU
	m["Back"] = r.URL.Path
	h.Apply(TPOPERROR, w)
}

func (h *Handler) HandleServerError(w http.ResponseWriter, r *http.Request, my *mybad.MuleError) {
	Log(my)
	h.pageOPError(w, r, my, nil)
	//http.Error(w, my.Error(), http.StatusInternalServerError)
}

func (h *Handler) HandleUserError(w http.ResponseWriter, r *http.Request, msg string, args ...interface{}) {
	h.pageOPError(w, r, nil, fmt.Errorf(msg, args...))
	//http.Error(w, fmt.Sprintf(msg, args...), http.StatusBadRequest)
}
