package main

import (
	"mule/jsend"
	"mule/mybad"
	"mule/overpower/json"
	"net/http"
)

func jFail(w http.ResponseWriter, code int, mp ...interface{}) {
	args := make(map[interface{}]interface{}, len(mp)/2)
	for i := 0; i < len(mp)/2; i += 1 {
		args[mp[2*i]] = mp[2*i+1]
	}
	err := jsend.Fail(w, code, args)
	if my, bad := Check(err, "jsend fail failure", "args", args); bad {
		Log(my)
	}
}
func jError(w http.ResponseWriter, msg string) {
	err := jsend.Error(w, msg)
	if my, bad := Check(err, "jsend error failure", "msg", msg); bad {
		Log(my)
	}
}
func jSuccess(w http.ResponseWriter, obj interface{}) {
	err := jsend.Success(w, obj)
	if my, bad := Check(err, "jsend success failure", "obj", obj); bad {
		Log(my)
	}
}

func Kirk(my *mybad.MuleError, w http.ResponseWriter) {
	Log(my)
	if my, bad := Check(jsend.Kirk(w), "api kirk failure"); bad {
		Log(my)
	}
}

// /overpower/json/
func apiJson(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiJsonGET(w, r)
	//case "PUT":
	default:
		jFail(w, 400, "method", "BAD REQUEST METHOD")
		return
	}
}

func apiJsonGET(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	lastFull := h.LastFull()
	if lastFull < 3 {
		jFail(w, 400, "url", "no object type given")
		return
	}
	switch h.Path[3] {
	case "games":
		h.apiJsonGETGames(w, r)
		return
	default:
		jFail(w, 404, "url", "unknown object type given")
		return
	}
}

// /overpower/json/games/
func (h *Handler) apiJsonGETGames(w http.ResponseWriter, r *http.Request) {
	lastFull := h.LastFull()
	if lastFull == 3 {
		games, err := OPDB.GetGames()
		if my, bad := Check(err, "api get failure", "object", "games"); bad {
			Kirk(my, w)
			return
		}
		jsonGames := json.LoadGames(games)
		jSuccess(w, jsonGames)
		return
	}
	gid, ok := h.IntAt(4)
	if !ok {
		jFail(w, 400, "params", "unreadable id for object")
		return
	}
	if lastFull == 4 {
		g, err := OPDB.GetGame("gid", gid)
		if err == ErrNoneFound {
			jFail(w, 400, "params", "id does not correspond to any existing object")
			return
		} else if my, bad := Check(err, "resource failure in json api get", "resource", "games", "gid", gid); bad {
			Kirk(my, w)
			return
		}
		jsG := json.LoadGame(g)
		jSuccess(w, jsG)
		return
	}
	jFail(w, 404, "url", "too many parameters given")
	return
}
