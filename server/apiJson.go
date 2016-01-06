package main

import (
	"fmt"
	"mule/overpower/opjson"
	"net/http"
)

// /overpower/json/
func apiJson(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiJsonGET(w, r)
	//case "PUT":
	default:
		opjson.HttpError(w, fmt.Errorf("BAD REQUEST METHOD"), http.StatusMethodNotAllowed)
		return
	}
}

func apiJsonGET(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	lastFull := h.LastFull()
	if lastFull < 3 {
		opjson.HttpError(w, fmt.Errorf("BAD LINK"), http.StatusNotFound)
		return
	}
	switch h.Path[3] {
	case "games":
		h.apiJsonGETGames(w, r)
		return
	default:
		opjson.HttpError(w, fmt.Errorf("BAD LINK"), http.StatusNotFound)
		return
	}
}

// /overpower/json/games/
func (h *Handler) apiJsonGETGames(w http.ResponseWriter, r *http.Request) {
	lastFull := h.LastFull()
	if lastFull == 3 {
		games, ok := OPDB.AllGames()
		if !ok {
			opjson.HttpError(w, fmt.Errorf("DATABASE ERROR"), http.StatusInternalServerError)
			return
		}
		jsonGames := opjson.LoadGames(games)
		opjson.HttpServe(w, jsonGames)
		return
	}
	gid, ok := h.IntAt(4)
	if !ok {
		opjson.HttpError(w, fmt.Errorf("UNREADABLE GID"), http.StatusBadRequest)
		return
	}
	if lastFull == 4 {
		g, ok := OPDB.GetGame(gid)
		if !ok {
			opjson.HttpError(w, fmt.Errorf("GID NOT FOUND"), http.StatusBadRequest)
			return
		}
		jsG := opjson.LoadGame(g)
		opjson.HttpServe(w, jsG)
		return
	}
	opjson.HttpError(w, fmt.Errorf("BAD LINK"), http.StatusNotFound)
	return
}
