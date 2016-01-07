package main

import (
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
		opjson.HttpError(w, 400, "BAD REQUEST METHOD")
		return
	}
}

func apiJsonGET(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	lastFull := h.LastFull()
	if lastFull < 3 {
		opjson.HttpError(w, 404, "BAD URL PATH: No object type given")
		return
	}
	switch h.Path[3] {
	case "games":
		h.apiJsonGETGames(w, r)
		return
	default:
		opjson.HttpError(w, 404, "BAD URL PATH: unknown object type")
		return
	}
}

// /overpower/json/games/
func (h *Handler) apiJsonGETGames(w http.ResponseWriter, r *http.Request) {
	lastFull := h.LastFull()
	if lastFull == 3 {
		games, ok := OPDB.AllGames()
		if !ok {
			opjson.HttpError(w, 500, "DATABASE ERROR")
			return
		}
		jsonGames := opjson.LoadGames(games)
		opjson.HttpServe(w, jsonGames)
		return
	}
	gid, ok := h.IntAt(4)
	if !ok {
		opjson.HttpError(w, 400, "Unreadable ID for object")
		return
	}
	if lastFull == 4 {
		g, ok := OPDB.GetGame(gid)
		if !ok {
			opjson.HttpError(w, 400, "ID does not correspond to any existing object")
			return
		}
		jsG := opjson.LoadGame(g)
		opjson.HttpServe(w, jsG)
		return
	}
	opjson.HttpError(w, 404, "BAD URL PATH: rambling path, don't know where to go")
	return
}
