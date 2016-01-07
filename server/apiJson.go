package main

import (
	"mule/jsend"
	"mule/overpower/json"
	"net/http"
)

// /overpower/json/
func apiJson(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiJsonGET(w, r)
	//case "PUT":
	default:
		jsend.Fail(w, 400, map[string]string{"method": "BAD REQUEST METHOD"})
		return
	}
}

func apiJsonGET(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	lastFull := h.LastFull()
	if lastFull < 3 {
		jsend.Fail(w, 404, map[string]string{"url": "no object type given"})
		return
	}
	switch h.Path[3] {
	case "games":
		h.apiJsonGETGames(w, r)
		return
	default:
		jsend.Fail(w, 404, map[string]string{"url": "unknown object type given"})
		return
	}
}

// /overpower/json/games/
func (h *Handler) apiJsonGETGames(w http.ResponseWriter, r *http.Request) {
	lastFull := h.LastFull()
	if lastFull == 3 {
		games, ok := OPDB.AllGames()
		if !ok {
			jsend.Kirk(w)
			return
		}
		jsonGames := json.LoadGames(games)
		jsend.Success(w, jsonGames)
		return
	}
	gid, ok := h.IntAt(4)
	if !ok {
		jsend.Fail(w, 400, map[string]string{"params": "unreadable id for object"})
		return
	}
	if lastFull == 4 {
		g, ok := OPDB.GetGame(gid)
		if !ok {
			jsend.Fail(w, 400, map[string]string{"params": "id does not correspond to any existing object"})
			return
		}
		jsG := json.LoadGame(g)
		jsend.Success(w, jsG)
		return
	}
	jsend.Fail(w, 404, map[string]string{"url": "rambling path, too many parameters given"})
	return
}
